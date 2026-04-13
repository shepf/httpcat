package v1

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"httpcat/internal/common"
	"httpcat/internal/common/ylog"
	"httpcat/internal/storage/auth"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/disk"
	"gorm.io/gorm"
)

// GenerateAppSecret 生成 AppSecret（MD5 随机数）
func GenerateAppSecret(c *gin.Context) {
	rand.Seed(time.Now().UnixNano())
	randomNumber := strconv.Itoa(rand.Intn(1000000))

	hasher := md5.New()
	hasher.Write([]byte(randomNumber))
	hash := hex.EncodeToString(hasher.Sum(nil))

	common.CreateResponse(c, common.SuccessCode, hash)
}

// SaveUploadToken 保存/更新上传凭证
func SaveUploadToken(c *gin.Context) {
	var data struct {
		Appkey    string `json:"appkey"`
		State     string `json:"state"`
		Appsecret string `json:"appsecret"`
		Desc      string `json:"desc"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		common.CreateResponse(c, common.ParamInvalidErrorCode, err)
		return
	}

	db, err := common.GetDB()
	if err != nil {
		common.CreateResponse(c, common.ErrorCode, err)
		return
	}

	var token common.UploadTokenItem
	result := db.Table("t_upload_token").Where("appkey = ?", data.Appkey).First(&token)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 插入新记录
			err = db.Table("t_upload_token").Create(&common.UploadTokenItem{
				Appkey:    data.Appkey,
				State:     data.State,
				Appsecret: data.Appsecret,
				CreatedAt: time.Now(),
			}).Error
			if err != nil {
				common.CreateResponse(c, common.ErrorCode, err)
				return
			}
		} else {
			common.CreateResponse(c, common.ErrorCode, result.Error)
			return
		}
	} else {
		// 更新记录
		err = db.Table("t_upload_token").Where("appkey = ?", data.Appkey).Updates(common.UploadTokenItem{
			State:     data.State,
			Appsecret: data.Appsecret,
			Desc:      data.Desc,
		}).Error
		if err != nil {
			common.CreateResponse(c, common.ErrorCode, err)
			return
		}
	}

	// 查询 id 返回给前端
	var item common.UploadTokenItem
	if err := db.Table("t_upload_token").Where("appkey = ?", data.Appkey).First(&item).Error; err != nil {
		common.CreateResponse(c, common.ErrorCode, err)
		return
	}

	common.CreateResponse(c, common.SuccessCode, struct {
		ID int `json:"id"`
	}{ID: item.ID})
}

// RemoveUploadToken 删除上传凭证
func RemoveUploadToken(c *gin.Context) {
	appkey := c.Query("appkey")
	if appkey == "" {
		common.CreateResponse(c, common.ParamInvalidErrorCode, errors.New("缺少必传参数 appkey"))
		return
	}

	db, err := common.GetDB()
	if err != nil {
		common.CreateResponse(c, common.ErrorCode, err)
		return
	}

	// 查询记录
	var token common.UploadTokenItem
	result := db.Table("t_upload_token").Where("appkey = ?", appkey).First(&token)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			common.CreateResponse(c, common.ErrorCode, errors.New("appkey不存在"))
			return
		}
		common.CreateResponse(c, common.ErrorCode, result.Error)
		return
	}

	// 系统内置记录不允许删除
	if token.IsSysBuilt {
		ylog.Errorf("RemoveUploadToken", "系统内置记录不允许删除")
		c.JSON(http.StatusForbidden, gin.H{
			"error": "系统内置记录不允许删除",
		})
		return
	}

	// 删除记录
	result = db.Table("t_upload_token").Where("appkey = ?", appkey).Delete(&common.UploadTokenItem{})
	if result.Error != nil {
		common.CreateResponse(c, common.ErrorCode, result.Error)
		return
	}

	common.CreateResponse(c, common.SuccessCode, "删除成功")
}

// GetUploadTokenLists 获取上传凭证列表
func GetUploadTokenLists(c *gin.Context) {
	var logs []common.UploadTokenItem
	db, err := common.GetDB()
	if err != nil {
		common.CreateResponse(c, common.ErrorCode, err)
		return
	}

	if err := db.Table("t_upload_token").Find(&logs).Error; err != nil {
		common.CreateResponse(c, common.ErrorCode, err)
		return
	}

	common.CreateResponse(c, common.SuccessCode, logs)
}

// CheckUploadToken 校验上传凭证
func CheckUploadToken(c *gin.Context) {
	uploadToken := c.Request.Header.Get("UploadToken")

	accessKey := common.AppKey
	secretKey := common.AppSecret

	mac := auth.New(accessKey, secretKey)
	if !mac.VerifyUploadToken(uploadToken) {
		common.CreateResponse(c, common.ErrorCode, "UploadToken is invalid")
		return
	}

	common.CreateResponse(c, common.SuccessCode, "UploadToken is valid")
}

// DataOverview 数据概览
func DataOverview(c *gin.Context) {
	dir := common.GetUploadDir()

	usedSpace, availableSpace, _ := getDiskUsage(dir)

	c.JSON(http.StatusOK, gin.H{
		"used_space":      usedSpace,
		"available_space": availableSpace,
	})
}

// GetUploadAvailableSpace 获取上传可用空间
func GetUploadAvailableSpace(c *gin.Context) {
	dir := common.GetUploadDir()

	usedSpace, freeSpace, _ := getDiskUsage(dir)

	c.JSON(http.StatusOK, gin.H{
		"usedSpace": usedSpace,
		"freeSpace": freeSpace,
	})
}

// getDiskUsage 获取磁盘使用情况
func getDiskUsage(path string) (uint64, uint64, error) {
	usage, err := disk.Usage(path)
	if err != nil {
		return 0, 0, err
	}
	return usage.Used, usage.Free, nil
}
