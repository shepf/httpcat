package v1

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"httpcat/internal/common"
	"httpcat/internal/common/utils"
	"httpcat/internal/common/ylog"
	"httpcat/internal/models"
	"httpcat/internal/storage/auth"

	"github.com/gin-gonic/gin"
)

// ========== v0.7.0 新增：分片上传 + 断点续传 ==========
//
// 流程：
//   1. POST /api/v1/file/upload/init       → 创建会话，返回 uploadId
//   2. GET  /api/v1/file/upload/status     → 查询已上传分片索引（用于断点续传）
//   3. POST /api/v1/file/upload/chunk      → 上传单个分片（幂等，重复上传覆盖）
//   4. POST /api/v1/file/upload/complete   → 合并所有分片为最终文件
//   5. POST /api/v1/file/upload/abort      → 中止会话并清理分片
//
// 鉴权复用现有 TokenOrAKSKAuth（JWT 或 AK/SK 签名）；在启用 UploadToken 校验时，
// init/chunk/complete 均要求请求头 UploadToken。

const (
	// DefaultChunkSize 默认分片大小（5MB），客户端可在 init 时指定
	DefaultChunkSize = 5 * 1024 * 1024
	// MaxChunkSize 单分片最大 100MB
	MaxChunkSize = 100 * 1024 * 1024
	// MinChunkSize 单分片最小 64KB（防止分片过小导致元数据膨胀）
	MinChunkSize = 64 * 1024
	// MaxFileSize 单文件最大 100GB（安全上限，防止磁盘打满）
	MaxFileSize = 100 * 1024 * 1024 * 1024
	// DefaultSessionTTL 会话默认有效期 24 小时
	DefaultSessionTTL = 24 * time.Hour
)

// chunkWriteLock 对同一 uploadID 的状态更新串行化，避免并发分片写入后 bitmap 竞态
var chunkWriteLock sync.Map // map[uploadID]*sync.Mutex

func lockSession(uploadID string) *sync.Mutex {
	mu, _ := chunkWriteLock.LoadOrStore(uploadID, &sync.Mutex{})
	return mu.(*sync.Mutex)
}

func unlockAndCleanup(uploadID string, mu *sync.Mutex) {
	mu.Unlock()
	// 注意：不 Delete，因为可能还有并发分片；会话结束后随临时目录一起清理
}

// ---------- 请求/响应结构 ----------

type initUploadReq struct {
	FileName    string `json:"fileName" binding:"required"`
	FileSize    int64  `json:"fileSize" binding:"required"`
	ChunkSize   int64  `json:"chunkSize"` // 可选，默认 5MB
	TotalChunks int    `json:"totalChunks"` // 可选，若不传则由服务端根据 fileSize/chunkSize 计算
	FileMD5     string `json:"fileMD5"`     // 可选，整体 MD5
	Dir         string `json:"dir"`         // 可选，目标子目录
	Overwrite   bool   `json:"overwrite"`   // 若最终文件已存在，是否允许覆盖（默认 false）
}

type initUploadResp struct {
	UploadID      string `json:"uploadId"`
	ChunkSize     int64  `json:"chunkSize"`
	TotalChunks   int    `json:"totalChunks"`
	UploadedNum   int    `json:"uploadedNum"`   // 秒传场景可能直接为 totalChunks
	UploadedIdx   []int  `json:"uploadedIdx"`   // 已上传的分片索引列表（断点续传用）
	Instant       bool   `json:"instant"`       // 是否命中秒传
	ExpireAt      int64  `json:"expireAt"`      // 会话过期时间（unix 秒）
}

// ---------- 工具：UploadToken 校验（复用 UploadFile 的逻辑） ----------

// verifyUploadTokenIfNeeded 校验 UploadToken。
// 返回 (appkey, ok)；若 EnableUploadToken=false，则直接返回 ("", true)。
// 未通过时已写入响应体，调用方只需 return。
func verifyUploadTokenIfNeeded(c *gin.Context) (string, bool) {
	if !common.EnableUploadToken {
		return "", true
	}
	uploadToken := c.Request.Header.Get("UploadToken")
	if uploadToken == "" {
		common.BadRequest(c, "UploadToken is empty")
		return "", false
	}
	parts := strings.Split(uploadToken, ":")
	if len(parts) != 3 {
		common.Unauthorized(c, "Invalid UploadToken format")
		return "", false
	}
	appkey := parts[0]
	common.UploadTokenLock.RLock()
	tokenItem, ok := common.UploadTokenTable[appkey]
	common.UploadTokenLock.RUnlock()
	if !ok {
		common.Unauthorized(c, "Invalid Appkey")
		return "", false
	}
	if tokenItem.State == "closed" {
		common.Unauthorized(c, "Invalid Appkey, appkey is closed")
		return "", false
	}
	mac := auth.New(appkey, tokenItem.Appsecret)
	if !mac.VerifyUploadToken(uploadToken) {
		common.Unauthorized(c, "UploadToken is invalid")
		return "", false
	}
	return appkey, true
}

// ---------- bitmap 辅助 ----------

// bitmapSet 将 bitmap 第 i 位置为 '1'，返回新 bitmap 和是否发生变化
func bitmapSet(bitmap string, i int, total int) (string, bool) {
	if len(bitmap) != total {
		// 初始化或长度异常：重建 bitmap
		buf := make([]byte, total)
		for k := range buf {
			if k < len(bitmap) && bitmap[k] == '1' {
				buf[k] = '1'
			} else {
				buf[k] = '0'
			}
		}
		bitmap = string(buf)
	}
	if bitmap[i] == '1' {
		return bitmap, false
	}
	b := []byte(bitmap)
	b[i] = '1'
	return string(b), true
}

// bitmapCollectMissing 收集所有未上传的分片索引
func bitmapCollectMissing(bitmap string, total int) []int {
	missing := make([]int, 0)
	for i := 0; i < total; i++ {
		if i >= len(bitmap) || bitmap[i] != '1' {
			missing = append(missing, i)
		}
	}
	return missing
}

// bitmapCollectUploaded 收集已上传的分片索引
func bitmapCollectUploaded(bitmap string, total int) []int {
	uploaded := make([]int, 0)
	for i := 0; i < total && i < len(bitmap); i++ {
		if bitmap[i] == '1' {
			uploaded = append(uploaded, i)
		}
	}
	return uploaded
}

// ---------- Handler 实现 ----------

// InitChunkUpload 初始化分片上传会话
// POST /api/v1/file/upload/init
func InitChunkUpload(c *gin.Context) {
	if !common.FileUploadEnable {
		common.CreateResponse(c, common.ErrorCode, "File service is not enabled")
		return
	}

	appkey, ok := verifyUploadTokenIfNeeded(c)
	if !ok {
		return
	}

	var req initUploadReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequest(c, "invalid request body: "+err.Error())
		return
	}

	// 校验文件大小
	if req.FileSize <= 0 {
		common.BadRequest(c, "fileSize must be greater than 0")
		return
	}
	if req.FileSize > MaxFileSize {
		common.BadRequest(c, fmt.Sprintf("fileSize exceeds the limit of %d bytes", MaxFileSize))
		return
	}

	// 应用系统配置的单文件大小上限（若开启）
	if common.UploadPolicyFSizeLimit > 0 && req.FileSize > common.UploadPolicyFSizeLimit {
		common.BadRequest(c, fmt.Sprintf("fileSize exceeds upload policy limit %d bytes", common.UploadPolicyFSizeLimit))
		return
	}
	if common.UploadPolicyFSizeMin > 0 && req.FileSize < common.UploadPolicyFSizeMin {
		common.BadRequest(c, fmt.Sprintf("fileSize is smaller than upload policy min %d bytes", common.UploadPolicyFSizeMin))
		return
	}

	// 规范化文件名
	fileName, err := common.NormalizeSafeFileName(req.FileName)
	if err != nil {
		common.BadRequest(c, "invalid fileName")
		return
	}

	// 解析目标目录
	uploadBaseDir := common.GetUploadDir()
	relDir := strings.TrimSpace(req.Dir)
	if relDir != "" {
		resolvedDir, err := common.ResolvePathWithinBase(common.GetUploadDir(), relDir)
		if err != nil {
			common.BadRequest(c, "invalid dir")
			return
		}
		uploadBaseDir = resolvedDir
	}

	// 校验目标文件是否已存在
	finalPath, err := common.ResolvePathWithinBase(uploadBaseDir, fileName)
	if err != nil {
		common.BadRequest(c, "invalid fileName path")
		return
	}
	if info, err := os.Stat(finalPath); err == nil && !info.IsDir() && !req.Overwrite {
		// 若已存在且 MD5 匹配，可支持"秒传"
		if req.FileMD5 != "" {
			if existingMD5, err := utils.CalculateMD5(finalPath); err == nil && strings.EqualFold(existingMD5, req.FileMD5) {
				c.JSON(http.StatusOK, gin.H{
					"errorCode": common.SuccessCode,
					"msg":       "success",
					"data": initUploadResp{
						UploadID:    "instant-" + req.FileMD5,
						ChunkSize:   0,
						TotalChunks: 0,
						UploadedNum: 0,
						UploadedIdx: []int{},
						Instant:     true,
						ExpireAt:    time.Now().Add(DefaultSessionTTL).Unix(),
					},
				})
				return
			}
		}
		common.CreateResponse(c, common.ErrorCode, "target file already exists (set overwrite=true to replace)")
		return
	}

	// 计算分片参数
	chunkSize := req.ChunkSize
	if chunkSize == 0 {
		chunkSize = DefaultChunkSize
	}
	if chunkSize < MinChunkSize || chunkSize > MaxChunkSize {
		common.BadRequest(c, fmt.Sprintf("chunkSize must be between %d and %d bytes", MinChunkSize, MaxChunkSize))
		return
	}
	totalChunks := int((req.FileSize + chunkSize - 1) / chunkSize)
	if req.TotalChunks > 0 && req.TotalChunks != totalChunks {
		common.BadRequest(c, fmt.Sprintf("totalChunks mismatch: expect %d, got %d", totalChunks, req.TotalChunks))
		return
	}

	// 秒传：若上传过相同 MD5，找到已 completed 的会话
	if req.FileMD5 != "" {
		db, err := common.GetDB()
		if err == nil {
			var prior models.UploadSessionModel
			err := db.Where("file_md5 = ? AND status = ?", req.FileMD5, "completed").
				Order("updated_at DESC").
				First(&prior).Error
			if err == nil && prior.FinalPath != "" {
				if _, err := os.Stat(prior.FinalPath); err == nil {
					// 硬链接到目标位置（同分区），不占额外空间；失败则退化为拷贝
					_ = os.MkdirAll(filepath.Dir(finalPath), 0o755)
					if err := os.Link(prior.FinalPath, finalPath); err != nil {
						if err := copyFile(prior.FinalPath, finalPath); err == nil {
							ylog.Infof("InitChunkUpload", "instant upload via copy: %s", finalPath)
						}
					} else {
						ylog.Infof("InitChunkUpload", "instant upload via hardlink: %s", finalPath)
					}
					// 记录上传日志
					if common.EnableSqlite {
						go insertUploadLog(c.ClientIP(), appkey,
							time.Now().Format("2006-01-02 15:04:05"),
							fileName, utils.FormatSize(req.FileSize),
							req.FileMD5, time.Now().Unix(), time.Now().Unix())
					}
					c.JSON(http.StatusOK, gin.H{
						"errorCode": common.SuccessCode,
						"msg":       "success",
						"data": initUploadResp{
							UploadID:    "instant-" + req.FileMD5,
							ChunkSize:   chunkSize,
							TotalChunks: totalChunks,
							UploadedNum: totalChunks,
							UploadedIdx: fullIdxList(totalChunks),
							Instant:     true,
							ExpireAt:    time.Now().Add(DefaultSessionTTL).Unix(),
						},
					})
					return
				}
			}
		}
	}

	// 创建新会话
	uploadID, err := generateUploadID()
	if err != nil {
		ylog.Errorf("InitChunkUpload", "generate upload id failed: %v", err)
		common.CreateResponse(c, common.ErrorCode, "failed to create session")
		return
	}

	username, _ := c.Get("user")
	uname, _ := username.(string)

	session := models.UploadSessionModel{
		UploadID:     uploadID,
		Appkey:       appkey,
		Username:     uname,
		FileName:     fileName,
		RelDir:       relDir,
		FileSize:     req.FileSize,
		ChunkSize:    chunkSize,
		TotalChunks:  totalChunks,
		FileMD5:      strings.ToLower(req.FileMD5),
		UploadedBits: strings.Repeat("0", totalChunks),
		UploadedNum:  0,
		Status:       "active",
		FinalPath:    finalPath,
		IP:           c.ClientIP(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		ExpireAt:     time.Now().Add(DefaultSessionTTL),
	}

	db, err := common.GetDB()
	if err != nil {
		common.CreateResponse(c, common.ErrorCode, "database unavailable")
		return
	}
	if err := db.Create(&session).Error; err != nil {
		ylog.Errorf("InitChunkUpload", "create session failed: %v", err)
		common.CreateResponse(c, common.ErrorCode, "failed to create session")
		return
	}

	// 创建分片临时目录
	if err := os.MkdirAll(common.ChunkSessionDir(uploadID), 0o755); err != nil {
		ylog.Errorf("InitChunkUpload", "create chunk dir failed: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"errorCode": common.SuccessCode,
		"msg":       "success",
		"data": initUploadResp{
			UploadID:    uploadID,
			ChunkSize:   chunkSize,
			TotalChunks: totalChunks,
			UploadedNum: 0,
			UploadedIdx: []int{},
			Instant:     false,
			ExpireAt:    session.ExpireAt.Unix(),
		},
	})
}

// GetChunkUploadStatus 查询分片上传会话状态
// GET /api/v1/file/upload/status?uploadId=xxx
func GetChunkUploadStatus(c *gin.Context) {
	uploadID := c.Query("uploadId")
	if uploadID == "" {
		common.BadRequest(c, "uploadId is required")
		return
	}

	db, err := common.GetDB()
	if err != nil {
		common.CreateResponse(c, common.ErrorCode, "database unavailable")
		return
	}

	var session models.UploadSessionModel
	if err := db.Where("upload_id = ?", uploadID).First(&session).Error; err != nil {
		common.CreateResponse(c, common.ErrorCode, "session not found")
		return
	}

	common.CreateResponse(c, common.SuccessCode, gin.H{
		"uploadId":    session.UploadID,
		"fileName":    session.FileName,
		"fileSize":    session.FileSize,
		"chunkSize":   session.ChunkSize,
		"totalChunks": session.TotalChunks,
		"uploadedNum": session.UploadedNum,
		"uploadedIdx": bitmapCollectUploaded(session.UploadedBits, session.TotalChunks),
		"missingIdx":  bitmapCollectMissing(session.UploadedBits, session.TotalChunks),
		"status":      session.Status,
		"expireAt":    session.ExpireAt.Unix(),
	})
}

// UploadChunk 上传单个分片
// POST /api/v1/file/upload/chunk
// multipart/form-data: uploadId, chunkIndex, chunk(file), [chunkMD5]
func UploadChunk(c *gin.Context) {
	if !common.FileUploadEnable {
		common.CreateResponse(c, common.ErrorCode, "File service is not enabled")
		return
	}
	if _, ok := verifyUploadTokenIfNeeded(c); !ok {
		return
	}

	uploadID := c.PostForm("uploadId")
	chunkIndexStr := c.PostForm("chunkIndex")
	chunkMD5Expect := strings.ToLower(strings.TrimSpace(c.PostForm("chunkMD5")))

	if uploadID == "" || chunkIndexStr == "" {
		common.BadRequest(c, "uploadId and chunkIndex are required")
		return
	}
	chunkIndex, err := strconv.Atoi(chunkIndexStr)
	if err != nil || chunkIndex < 0 {
		common.BadRequest(c, "invalid chunkIndex")
		return
	}

	file, header, err := c.Request.FormFile("chunk")
	if err != nil {
		common.BadRequest(c, "missing chunk file field: "+err.Error())
		return
	}
	defer file.Close()

	db, err := common.GetDB()
	if err != nil {
		common.CreateResponse(c, common.ErrorCode, "database unavailable")
		return
	}

	var session models.UploadSessionModel
	if err := db.Where("upload_id = ?", uploadID).First(&session).Error; err != nil {
		common.CreateResponse(c, common.ErrorCode, "session not found")
		return
	}
	if session.Status != "active" {
		common.CreateResponse(c, common.ErrorCode, "session is not active: "+session.Status)
		return
	}
	if chunkIndex >= session.TotalChunks {
		common.BadRequest(c, fmt.Sprintf("chunkIndex out of range [0, %d)", session.TotalChunks))
		return
	}

	// 校验分片大小：除最后一片外必须等于 ChunkSize；最后一片为 fileSize - (total-1)*chunkSize
	expectedSize := session.ChunkSize
	if chunkIndex == session.TotalChunks-1 {
		expectedSize = session.FileSize - int64(session.TotalChunks-1)*session.ChunkSize
	}
	if header.Size != expectedSize {
		common.BadRequest(c, fmt.Sprintf("chunk size mismatch: expect %d, got %d", expectedSize, header.Size))
		return
	}

	// 将分片写入临时文件（先写 .part，完成后 rename 保证幂等）
	chunkDir := common.ChunkSessionDir(uploadID)
	if err := os.MkdirAll(chunkDir, 0o755); err != nil {
		ylog.Errorf("UploadChunk", "mkdir chunk dir: %v", err)
		common.CreateResponse(c, common.ErrorCode, "failed to prepare chunk dir")
		return
	}

	chunkPath := filepath.Join(chunkDir, fmt.Sprintf("%06d", chunkIndex))
	tmpPath := chunkPath + ".part"

	out, err := os.Create(tmpPath)
	if err != nil {
		ylog.Errorf("UploadChunk", "create tmp file: %v", err)
		common.CreateResponse(c, common.ErrorCode, "failed to write chunk")
		return
	}

	hasher := md5.New()
	if _, err := io.Copy(io.MultiWriter(out, hasher), file); err != nil {
		_ = out.Close()
		_ = os.Remove(tmpPath)
		ylog.Errorf("UploadChunk", "copy chunk: %v", err)
		common.CreateResponse(c, common.ErrorCode, "failed to write chunk")
		return
	}
	_ = out.Close()

	gotMD5 := hex.EncodeToString(hasher.Sum(nil))
	if chunkMD5Expect != "" && chunkMD5Expect != gotMD5 {
		_ = os.Remove(tmpPath)
		common.BadRequest(c, fmt.Sprintf("chunk MD5 mismatch: expect %s, got %s", chunkMD5Expect, gotMD5))
		return
	}

	if err := os.Rename(tmpPath, chunkPath); err != nil {
		_ = os.Remove(tmpPath)
		ylog.Errorf("UploadChunk", "rename chunk: %v", err)
		common.CreateResponse(c, common.ErrorCode, "failed to finalize chunk")
		return
	}

	// 更新 bitmap（同一会话串行化）
	mu := lockSession(uploadID)
	mu.Lock()
	defer unlockAndCleanup(uploadID, mu)

	// 重新读一次 session 拿最新 bitmap
	if err := db.Where("upload_id = ?", uploadID).First(&session).Error; err != nil {
		common.CreateResponse(c, common.ErrorCode, "session not found")
		return
	}

	newBitmap, changed := bitmapSet(session.UploadedBits, chunkIndex, session.TotalChunks)
	if changed {
		session.UploadedBits = newBitmap
		session.UploadedNum++
		session.UpdatedAt = time.Now()
		if err := db.Model(&models.UploadSessionModel{}).
			Where("upload_id = ?", uploadID).
			Updates(map[string]interface{}{
				"uploaded_bits": newBitmap,
				"uploaded_num":  session.UploadedNum,
				"updated_at":    session.UpdatedAt,
			}).Error; err != nil {
			ylog.Errorf("UploadChunk", "update session: %v", err)
		}
	}

	common.CreateResponse(c, common.SuccessCode, gin.H{
		"uploadId":    uploadID,
		"chunkIndex":  chunkIndex,
		"chunkMD5":    gotMD5,
		"uploadedNum": session.UploadedNum,
		"totalChunks": session.TotalChunks,
	})
}

// CompleteChunkUpload 合并分片
// POST /api/v1/file/upload/complete
// body: { "uploadId": "..." }
func CompleteChunkUpload(c *gin.Context) {
	if !common.FileUploadEnable {
		common.CreateResponse(c, common.ErrorCode, "File service is not enabled")
		return
	}
	appkey, ok := verifyUploadTokenIfNeeded(c)
	if !ok {
		return
	}

	var req struct {
		UploadID string `json:"uploadId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequest(c, "invalid request body")
		return
	}

	db, err := common.GetDB()
	if err != nil {
		common.CreateResponse(c, common.ErrorCode, "database unavailable")
		return
	}

	var session models.UploadSessionModel
	if err := db.Where("upload_id = ?", req.UploadID).First(&session).Error; err != nil {
		common.CreateResponse(c, common.ErrorCode, "session not found")
		return
	}
	if session.Status != "active" {
		common.CreateResponse(c, common.ErrorCode, "session is not active: "+session.Status)
		return
	}

	// 加锁防并发 complete
	mu := lockSession(req.UploadID)
	mu.Lock()
	defer unlockAndCleanup(req.UploadID, mu)

	// 再读一次最新状态
	if err := db.Where("upload_id = ?", req.UploadID).First(&session).Error; err != nil {
		common.CreateResponse(c, common.ErrorCode, "session not found")
		return
	}

	// 检查所有分片都已上传
	if session.UploadedNum != session.TotalChunks {
		missing := bitmapCollectMissing(session.UploadedBits, session.TotalChunks)
		common.CreateResponse(c, common.ErrorCode, gin.H{
			"msg":        fmt.Sprintf("not all chunks uploaded: %d/%d", session.UploadedNum, session.TotalChunks),
			"missingIdx": missing,
		})
		return
	}

	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(session.FinalPath), 0o755); err != nil {
		ylog.Errorf("CompleteChunkUpload", "mkdir final dir: %v", err)
		common.CreateResponse(c, common.ErrorCode, "failed to prepare final dir")
		return
	}

	// 合并：顺序读取每个分片，写入目标文件，同时计算 MD5
	chunkDir := common.ChunkSessionDir(req.UploadID)
	tmpFinal := session.FinalPath + ".merging"
	out, err := os.Create(tmpFinal)
	if err != nil {
		ylog.Errorf("CompleteChunkUpload", "create final tmp: %v", err)
		common.CreateResponse(c, common.ErrorCode, "failed to create final file")
		return
	}

	hasher := md5.New()
	var written int64
	for i := 0; i < session.TotalChunks; i++ {
		chunkPath := filepath.Join(chunkDir, fmt.Sprintf("%06d", i))
		in, err := os.Open(chunkPath)
		if err != nil {
			_ = out.Close()
			_ = os.Remove(tmpFinal)
			ylog.Errorf("CompleteChunkUpload", "open chunk %d: %v", i, err)
			common.CreateResponse(c, common.ErrorCode, fmt.Sprintf("failed to open chunk %d", i))
			return
		}
		n, err := io.Copy(io.MultiWriter(out, hasher), in)
		_ = in.Close()
		if err != nil {
			_ = out.Close()
			_ = os.Remove(tmpFinal)
			ylog.Errorf("CompleteChunkUpload", "merge chunk %d: %v", i, err)
			common.CreateResponse(c, common.ErrorCode, fmt.Sprintf("failed to merge chunk %d", i))
			return
		}
		written += n
	}
	if err := out.Close(); err != nil {
		_ = os.Remove(tmpFinal)
		ylog.Errorf("CompleteChunkUpload", "close final: %v", err)
		common.CreateResponse(c, common.ErrorCode, "failed to close final file")
		return
	}

	// 校验总大小
	if written != session.FileSize {
		_ = os.Remove(tmpFinal)
		common.CreateResponse(c, common.ErrorCode, fmt.Sprintf("final size mismatch: expect %d, got %d", session.FileSize, written))
		return
	}

	// 校验整体 MD5（若客户端声明了）
	finalMD5 := hex.EncodeToString(hasher.Sum(nil))
	if session.FileMD5 != "" && !strings.EqualFold(session.FileMD5, finalMD5) {
		_ = os.Remove(tmpFinal)
		common.CreateResponse(c, common.ErrorCode, fmt.Sprintf("file MD5 mismatch: expect %s, got %s", session.FileMD5, finalMD5))
		return
	}

	// 原子 rename 到最终位置
	if err := os.Rename(tmpFinal, session.FinalPath); err != nil {
		_ = os.Remove(tmpFinal)
		ylog.Errorf("CompleteChunkUpload", "rename final: %v", err)
		common.CreateResponse(c, common.ErrorCode, "failed to finalize file")
		return
	}

	// 更新会话状态
	db.Model(&models.UploadSessionModel{}).
		Where("upload_id = ?", req.UploadID).
		Updates(map[string]interface{}{
			"status":     "completed",
			"file_md5":   finalMD5,
			"updated_at": time.Now(),
		})

	// 异步清理分片目录
	go func() {
		if err := os.RemoveAll(chunkDir); err != nil {
			ylog.Errorf("CompleteChunkUpload", "remove chunk dir: %v", err)
		}
	}()

	// 记录上传日志（复用 /upload 的日志格式）
	if common.EnableSqlite {
		go insertUploadLog(c.ClientIP(), appkey,
			time.Now().Format("2006-01-02 15:04:05"),
			session.FileName, utils.FormatSize(session.FileSize),
			finalMD5, time.Now().Unix(), time.Now().Unix())
	}

	// 通知 Webhook（复用普通上传的逻辑）
	if common.PersistentNotifyURL != "" {
		go utils.SendNotify(common.PersistentNotifyURL, fmt.Sprintf(
			">分片上传完成：\n- IP地址：%s\n- 文件名：%s\n- 大小：%s\n- MD5：%s",
			c.ClientIP(), session.FileName, utils.FormatSize(session.FileSize), finalMD5))
	}

	common.CreateResponse(c, common.SuccessCode, gin.H{
		"uploadId": req.UploadID,
		"fileName": session.FileName,
		"fileSize": session.FileSize,
		"fileMD5":  finalMD5,
		"path":     session.FinalPath,
	})
}

// AbortChunkUpload 中止会话
// POST /api/v1/file/upload/abort
// body: { "uploadId": "..." }
func AbortChunkUpload(c *gin.Context) {
	var req struct {
		UploadID string `json:"uploadId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequest(c, "invalid request body")
		return
	}

	db, err := common.GetDB()
	if err != nil {
		common.CreateResponse(c, common.ErrorCode, "database unavailable")
		return
	}

	var session models.UploadSessionModel
	if err := db.Where("upload_id = ?", req.UploadID).First(&session).Error; err != nil {
		common.CreateResponse(c, common.ErrorCode, "session not found")
		return
	}

	// 清理分片目录
	chunkDir := common.ChunkSessionDir(req.UploadID)
	if err := os.RemoveAll(chunkDir); err != nil {
		ylog.Errorf("AbortChunkUpload", "remove chunk dir: %v", err)
	}

	// 更新状态
	db.Model(&models.UploadSessionModel{}).
		Where("upload_id = ?", req.UploadID).
		Updates(map[string]interface{}{
			"status":     "aborted",
			"updated_at": time.Now(),
		})

	common.CreateResponse(c, common.SuccessCode, gin.H{
		"uploadId": req.UploadID,
		"status":   "aborted",
	})
}

// ---------- 辅助函数 ----------

func generateUploadID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func fullIdxList(n int) []int {
	idx := make([]int, n)
	for i := 0; i < n; i++ {
		idx[i] = i
	}
	return idx
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
