// @ts-ignore
/* eslint-disable */
import request from 'umi-request';

// request拦截器, 自动附加 Bearer Token
request.interceptors.request.use((url, options) => {
  const token = localStorage.getItem('token') || '';
  return {
    url,
    options: {
      ...options,
      interceptors: true,
      headers: {
        ...options.headers,
        Authorization: `Bearer ${token}`,
      },
    },
  };
});

/** 获取当前的用户 GET /api/v1/user/currentUser */
export async function currentUser(options?: { [key: string]: any }) {
  return request<{
    errorCode: number;
    msg: string;
    data: API.CurrentUser;
  }>('/api/v1/user/currentUser', {
    method: 'GET',
    ...(options || {}),
  });
}

/** 退出登录接口 POST /api/v1/user/login/outLogin */
export async function outLogin(options?: { [key: string]: any }) {
  return request<Record<string, any>>('/api/v1/user/login/outLogin', {
    method: 'POST',
    ...(options || {}),
  });
}

/** 登录接口 POST /api/v1/user/login/account */
export async function login(body: API.LoginParams, options?: { [key: string]: any }) {
  return request<API.LoginResult>('/api/v1/user/login/account', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 获取文件上传日志列表 GET /api/v1/file/uploadHistoryLogs */
export async function uploadHistoryLogs(
  params: {
    current?: number;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  return request<API.UploadHistoryLogsList>('/api/v1/file/uploadHistoryLogs', {
    method: 'GET',
    params: { ...params },
    ...(options || {}),
  });
}

/** 删除上传文件历史记录 DELETE /api/v1/file/uploadHistoryLogs */
export async function removeUploadHistoryLog(options?: { [id: string]: any }) {
  return request<Record<string, any>>('/api/v1/file/uploadHistoryLogs', {
    method: 'DELETE',
    params: { ...options },
  });
}

// ==================== 上传Token管理 ====================

/** 获取上传Token列表 GET /api/v1/user/uploadTokenLists */
export async function uploadTokenLists(
  params: {},
  options?: { [key: string]: any },
) {
  return request<API.UploadTokenLists>('/api/v1/user/uploadTokenLists', {
    method: 'GET',
    params: { ...params },
    ...(options || {}),
  });
}

/** 保存上传Token */
export async function saveUploadToken(options?: { [key: string]: any }) {
  return request<API.UploadTokenLists>('/api/v1/user/saveUploadToken', {
    method: 'POST',
    data: { ...(options || {}) },
  });
}

/** 删除上传Token */
export async function removeUploadToken(
  params: {},
  options?: { [key: string]: any },
) {
  return request<API.UploadTokenLists>('/api/v1/user/removeUploadToken', {
    method: 'DELETE',
    params: { ...params },
    ...(options || {}),
  });
}

/** 生成上传Token */
export async function createUploadToken(options?: { [key: string]: any }) {
  return request<API.ErrorResponse>('/api/v1/user/createUploadToken', {
    method: 'POST',
    data: { ...(options || {}) },
  });
}

/** 修改密码 */
export async function changePasswd(options?: { [key: string]: any }) {
  return request<API.ErrorResponse>('/api/v1/user/changePasswd', {
    method: 'POST',
    data: { ...(options || {}) },
  });
}

// ==================== 系统配置 & 统计 ====================

/** 获取软件版本信息 */
export async function getVersion(options?: { [key: string]: any }) {
  return request<API.AntResponseData<API.Version>>('/api/v1/conf/getVersion', {
    method: 'GET',
    ...(options || {}),
  });
}

/** 获取上传统计 */
export async function getUploadStatistics(options?: { [key: string]: any }) {
  return request<API.AntResponseData<API.UploadStatistics>>('/api/v1/statistics/getUploadStatistics', {
    method: 'GET',
    ...(options || {}),
  });
}

/** 获取下载统计 */
export async function getDownloadStatistics(options?: { [key: string]: any }) {
  return request<API.AntResponseData<API.DownloadStatistics>>('/api/v1/statistics/getDownloadStatistics', {
    method: 'GET',
    ...(options || {}),
  });
}

/** 获取下载历史日志 GET /api/v1/statistics/downloadHistoryLogs */
export async function downloadHistoryLogs(params: {
  current?: number;
  pageSize?: number;
  filename?: string;
  file_md5?: string;
  ip?: string;
}) {
  return request<API.MyResponse<{
    list: API.DownloadHistoryLogItem[];
    current: number;
    pageSize: number;
    total: number;
  }>>('/api/v1/statistics/downloadHistoryLogs', {
    method: 'GET',
    params,
  });
}

/** 获取文件总览统计 GET /api/v1/statistics/getFileOverview */
export async function getFileOverview(options?: { [key: string]: any }) {
  return request<API.MyResponse<API.FileOverview>>('/api/v1/statistics/getFileOverview', {
    method: 'GET',
    ...(options || {}),
  });
}

/** 获取磁盘空间 */
export async function getUploadAvailableSpace(options?: { [key: string]: any }) {
  return request<API.UploadAvailableSpace>('/api/v1/user/getUploadAvailableSpace', {
    method: 'GET',
    ...(options || {}),
  });
}

/** 获取系统配置 */
export async function getConf(options?: { [key: string]: any }) {
  return request<API.AntResponseData<API.HttpcatConf>>('/api/v1/conf/getConf', {
    method: 'GET',
    ...(options || {}),
  });
}

/** 获取系统配置（完整） */
export async function getSysConfig(options?: { [key: string]: any }) {
  return request<API.MyResponse<API.SysConfig>>('/api/v1/conf/sysConfig', {
    method: 'GET',
    ...(options || {}),
  });
}

/** 更新系统配置 */
export async function updateSysConfig(data: Partial<API.SysConfig>) {
  return request<API.MyResponse<API.SysConfigUpdateResult>>('/api/v1/conf/sysConfig', {
    method: 'PUT',
    data,
  });
}

/** 重启服务（需管理员密码） */
export async function restartServer(password: string) {
  return request<API.MyResponse<{ message: string }>>('/api/v1/conf/restart', {
    method: 'POST',
    data: { password },
  });
}

/** 获取某个文件信息 */
export async function getFileInfo(options?: { [key: string]: any }) {
  return request<API.MyResponse<API.FileInfo>>('/api/v1/file/getFileInfo', {
    method: 'GET',
    params: { ...options },
  });
}

// ==================== 文件管理 ====================

/** 获取文件列表 GET /api/v1/file/listFiles */
export async function listFiles(params?: { dir?: string }) {
  return request<API.MyResponse<API.FileItem[]>>('/api/v1/file/listFiles', {
    method: 'GET',
    params,
  });
}

/** 批量删除文件 POST /api/v1/file/delete */
export async function deleteFiles(data: API.DeleteFilesParams) {
  return request<API.MyResponse<API.DeleteFilesResult>>('/api/v1/file/delete', {
    method: 'POST',
    data,
  });
}

/** 创建文件夹 POST /api/v1/file/mkdir */
export async function createFolder(data: API.CreateFolderParams) {
  return request<API.MyResponse<string>>('/api/v1/file/mkdir', {
    method: 'POST',
    data,
  });
}

/** 重命名文件/文件夹 POST /api/v1/file/rename */
export async function renameFile(data: API.RenameFileParams) {
  return request<API.MyResponse<string>>('/api/v1/file/rename', {
    method: 'POST',
    data,
  });
}

// ==================== 图片管理（统一到 umi-request） ====================

/** 获取图片缩略图列表 */
export async function listThumbImages(params: { page: number; pageSize: number; search?: string }) {
  return request<API.ImageListResponse>('/api/v1/imageManage/listThumbImages', {
    method: 'GET',
    params,
  });
}

/** 上传图片 */
export async function uploadImage(formData: FormData) {
  return request<API.ImageUploadResponse>('/api/v1/imageManage/upload', {
    method: 'POST',
    data: formData,
  });
}

/** 下载图片（返回 blob） */
export async function downloadImage(filename: string) {
  return request('/api/v1/imageManage/download', {
    method: 'GET',
    params: { filename },
    responseType: 'blob',
  });
}

/** 删除图片 */
export async function deleteImage(filename: string) {
  return request<API.ErrorResponse>('/api/v1/imageManage/delete', {
    method: 'DELETE',
    params: { filename },
  });
}

/** 清空所有图片 */
export async function clearImages() {
  return request<API.ErrorResponse>('/api/v1/imageManage/clear', {
    method: 'DELETE',
  });
}

/** 上传文件（通用，使用 UploadToken） */
export async function uploadFile(formData: FormData, uploadToken: string, onProgress?: (percent: number) => void) {
  return request('/api/v1/file/upload', {
    method: 'POST',
    data: formData,
    headers: {
      UploadToken: uploadToken,
    },
    requestType: 'form',
  });
}

/** 获取第一个可用的上传Token（用于 Welcome 页面快捷上传） */
export async function getFirstUploadToken(): Promise<string> {
  const result = await uploadTokenLists({});
  const tokens = result?.data || [];
  const activeToken = tokens.find((t) => t.state === 'open' && t.appkey && t.appsecret);
  if (!activeToken) return '';
  // 用 appkey + appsecret 生成实际的 UploadToken
  const tokenResult = await createUploadToken({
    appkey: activeToken.appkey,
    appsecret: activeToken.appsecret,
  });
  return tokenResult?.data || '';
}

// ==================== 文件分享 ====================

/** 创建分享 POST /api/v1/share */
export async function createShare(data: API.CreateShareParams) {
  return request<API.MyResponse<API.CreateShareResult>>('/api/v1/share', {
    method: 'POST',
    data,
  });
}

/** 获取分享列表 GET /api/v1/share/list */
export async function listShares(params: { current?: number; pageSize?: number }) {
  return request<API.MyResponse<{ list: API.ShareItem[]; current: number; pageSize: number; total: number }>>('/api/v1/share/list', {
    method: 'GET',
    params,
  });
}

/** 取消分享 DELETE /api/v1/share/:code */
export async function deleteShare(code: string) {
  return request<API.MyResponse<string>>(`/api/v1/share/${code}`, {
    method: 'DELETE',
  });
}

/** 获取分享信息（公开接口） GET /s/:code */
export async function getShareInfo(code: string) {
  return request<API.ShareInfoResult>(`/s/${code}`, {
    method: 'GET',
  });
}

/** 验证提取码 POST /s/:code/verify */
export async function verifyShareCode(code: string, extractCode: string) {
  return request<{ valid: boolean; reason?: string }>(`/s/${code}/verify`, {
    method: 'POST',
    data: { extractCode },
  });
}

/** 获取分享统计 GET /api/v1/share/stats */
export async function getShareStats() {
  return request<API.MyResponse<API.ShareStats>>('/api/v1/share/stats', {
    method: 'GET',
  });
}

/** 获取分享配置 GET /api/v1/share/config */
export async function getShareConfig() {
  return request<API.MyResponse<API.ShareConfig>>('/api/v1/share/config', {
    method: 'GET',
  });
}

// ==================== v0.6.0 文件预览 & 打包下载 ====================

/** 获取文件预览信息 GET /api/v1/file/previewInfo */
export async function getPreviewInfo(params: { filename: string }) {
  return request<API.MyResponse<API.PreviewInfo>>('/api/v1/file/previewInfo', {
    method: 'GET',
    params,
  });
}

/** 获取文件预览 URL */
export function getPreviewUrl(filename: string): string {
  const token = localStorage.getItem('token') || '';
  return `/api/v1/file/preview?filename=${encodeURIComponent(filename)}&token=${encodeURIComponent(token)}`;
}

/** 打包下载 POST /api/v1/file/downloadZip */
export async function downloadZip(data: API.DownloadZipParams) {
  const token = localStorage.getItem('token') || '';
  const response = await fetch('/api/v1/file/downloadZip', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    throw new Error('Download failed');
  }

  // 从 Content-Disposition 获取文件名
  const contentDisposition = response.headers.get('Content-Disposition');
  let zipFileName = 'httpcat-download.zip';
  if (contentDisposition) {
    const match = contentDisposition.match(/filename="?(.+?)"?$/);
    if (match) {
      zipFileName = match[1];
    }
  }

  const blob = await response.blob();
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = zipFileName;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  window.URL.revokeObjectURL(url);
}

/** 上传文件到指定目录（v0.6.0：支持 dir 参数） */
export async function uploadFileToDir(
  file: File,
  dir: string,
  uploadToken: string,
  onProgress?: (percent: number) => void,
) {
  const formData = new FormData();
  formData.append('f1', file);
  if (dir) {
    formData.append('dir', dir);
  }

  return request('/api/v1/file/upload', {
    method: 'POST',
    data: formData,
    headers: {
      UploadToken: uploadToken,
    },
    requestType: 'form',
  });
}

// ==================== v0.6.0 操作日志 ====================

/** 获取操作日志列表 GET /api/v1/oplog/list */
export async function getOperationLogs(params: API.OperationLogParams) {
  return request<API.MyResponse<{
    list: API.OperationLogItem[];
    current: number;
    pageSize: number;
    total: number;
  }>>('/api/v1/oplog/list', {
    method: 'GET',
    params,
  });
}

/** 获取操作日志统计 GET /api/v1/oplog/stats */
export async function getOperationStats() {
  return request<API.MyResponse<API.OperationStats>>('/api/v1/oplog/stats', {
    method: 'GET',
  });
}

// ==================== v0.7.0 分片上传 + 断点续传 ====================

/** 初始化分片上传会话 POST /api/v1/file/upload/init */
export async function initChunkUpload(data: API.InitChunkUploadParams, uploadToken: string) {
  return request<API.MyResponse<API.InitChunkUploadResp>>('/api/v1/file/upload/init', {
    method: 'POST',
    data,
    headers: {
      UploadToken: uploadToken,
    },
  });
}

/** 查询分片上传会话状态 GET /api/v1/file/upload/status?uploadId=xxx */
export async function getChunkUploadStatus(uploadId: string) {
  return request<API.MyResponse<API.ChunkUploadStatus>>('/api/v1/file/upload/status', {
    method: 'GET',
    params: { uploadId },
  });
}

/** 上传单个分片 POST /api/v1/file/upload/chunk
 *  使用原生 XHR 以支持上传进度
 */
export function uploadSingleChunk(
  uploadId: string,
  chunkIndex: number,
  chunkBlob: Blob,
  uploadToken: string,
  onProgress?: (percent: number) => void,
  chunkMD5?: string,
): Promise<API.MyResponse<any>> {
  return new Promise((resolve, reject) => {
    const fd = new FormData();
    fd.append('uploadId', uploadId);
    fd.append('chunkIndex', String(chunkIndex));
    if (chunkMD5) fd.append('chunkMD5', chunkMD5);
    fd.append('chunk', chunkBlob);

    const xhr = new XMLHttpRequest();
    xhr.open('POST', '/api/v1/file/upload/chunk', true);
    xhr.setRequestHeader('UploadToken', uploadToken);

    xhr.upload.onprogress = (e) => {
      if (e.lengthComputable && onProgress) {
        onProgress(Math.round((e.loaded * 100) / e.total));
      }
    };
    xhr.onload = () => {
      try {
        const resp = JSON.parse(xhr.responseText);
        if (xhr.status === 200 && resp.errorCode === 0) {
          resolve(resp);
        } else {
          reject(new Error(resp.msg || `chunk upload failed: ${xhr.status}`));
        }
      } catch (e) {
        reject(e);
      }
    };
    xhr.onerror = () => reject(new Error('network error'));
    xhr.onabort = () => reject(new Error('aborted'));
    xhr.send(fd);
  });
}

/** 完成分片上传（合并分片） POST /api/v1/file/upload/complete */
export async function completeChunkUpload(uploadId: string, uploadToken: string) {
  return request<API.MyResponse<API.CompleteChunkUploadResp>>('/api/v1/file/upload/complete', {
    method: 'POST',
    data: { uploadId },
    headers: {
      UploadToken: uploadToken,
    },
  });
}

/** 中止分片上传 POST /api/v1/file/upload/abort */
export async function abortChunkUpload(uploadId: string) {
  return request<API.MyResponse<any>>('/api/v1/file/upload/abort', {
    method: 'POST',
    data: { uploadId },
  });
}

/**
 * 高级分片上传：支持自动切片、并发、断点续传、进度回调、取消
 *
 * @param file       浏览器 File 对象
 * @param dir        目标子目录（相对上传根）
 * @param uploadToken 已生成的 UploadToken
 * @param opts       可选配置
 *   - chunkSize: 单片字节数（默认 5MB）
 *   - concurrent: 并发上传分片数（默认 3）
 *   - onProgress: 总进度回调 (0-100)
 *   - onChunkProgress: 单片进度回调
 *   - fileMD5: 可选，用于秒传/校验
 *   - signal: AbortSignal 用于取消
 */
export async function chunkedUpload(
  file: File,
  dir: string,
  uploadToken: string,
  opts?: {
    chunkSize?: number;
    concurrent?: number;
    fileMD5?: string;
    overwrite?: boolean;
    onProgress?: (percent: number, uploadedBytes: number, totalBytes: number) => void;
    signal?: AbortSignal;
  },
): Promise<API.CompleteChunkUploadResp> {
  const chunkSize = opts?.chunkSize || 5 * 1024 * 1024;
  const concurrent = opts?.concurrent || 3;
  const totalChunks = Math.max(1, Math.ceil(file.size / chunkSize));

  // Step 1: init
  const initResp = await initChunkUpload(
    {
      fileName: file.name,
      fileSize: file.size,
      chunkSize,
      totalChunks,
      fileMD5: opts?.fileMD5,
      dir,
      overwrite: opts?.overwrite,
    },
    uploadToken,
  );
  if (initResp.errorCode !== 0 || !initResp.data) {
    throw new Error(initResp.msg || 'init upload failed');
  }

  // 秒传命中
  if (initResp.data.instant) {
    opts?.onProgress?.(100, file.size, file.size);
    return {
      uploadId: initResp.data.uploadId,
      fileName: file.name,
      fileSize: file.size,
      fileMD5: opts?.fileMD5 || '',
      path: '',
    };
  }

  const uploadId = initResp.data.uploadId;
  const uploadedSet = new Set<number>(initResp.data.uploadedIdx || []);

  // Step 2: 计算所有待上传分片（排除已上传）
  const pending: number[] = [];
  for (let i = 0; i < totalChunks; i++) {
    if (!uploadedSet.has(i)) pending.push(i);
  }

  let uploadedChunks = uploadedSet.size;
  const reportProgress = () => {
    const uploaded = Math.min(uploadedChunks * chunkSize, file.size);
    const percent = Math.round((uploaded / file.size) * 100);
    opts?.onProgress?.(percent, uploaded, file.size);
  };
  reportProgress();

  // Step 3: 并发上传
  let cursor = 0;
  let aborted = false;
  const fail: Error[] = [];

  const checkSignal = () => {
    if (opts?.signal?.aborted) aborted = true;
  };

  const worker = async () => {
    while (!aborted && fail.length === 0) {
      checkSignal();
      if (aborted) return;

      const idx = cursor++;
      if (idx >= pending.length) return;
      const chunkIndex = pending[idx];

      const start = chunkIndex * chunkSize;
      const end = Math.min(start + chunkSize, file.size);
      const blob = file.slice(start, end);

      // 重试 3 次
      let lastErr: Error | null = null;
      for (let attempt = 0; attempt < 3; attempt++) {
        try {
          await uploadSingleChunk(uploadId, chunkIndex, blob, uploadToken);
          lastErr = null;
          break;
        } catch (e) {
          lastErr = e as Error;
          if (aborted) return;
          await new Promise((r) => setTimeout(r, 500 * (attempt + 1)));
        }
      }
      if (lastErr) {
        fail.push(lastErr);
        return;
      }
      uploadedChunks++;
      reportProgress();
    }
  };

  await Promise.all(Array.from({ length: concurrent }, () => worker()));

  if (opts?.signal?.aborted) {
    await abortChunkUpload(uploadId).catch(() => {});
    throw new Error('upload aborted');
  }
  if (fail.length > 0) {
    throw fail[0];
  }

  // Step 4: complete
  const completeResp = await completeChunkUpload(uploadId, uploadToken);
  if (completeResp.errorCode !== 0 || !completeResp.data) {
    throw new Error(completeResp.msg || 'complete upload failed');
  }
  opts?.onProgress?.(100, file.size, file.size);
  return completeResp.data;
}


