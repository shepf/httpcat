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

// ==================== 图片管理（统一到 umi-request） ====================

/** 获取图片缩略图列表 */
export async function listThumbImages(params: { page: number; pageSize: number }) {
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

/** 获取上传Token列表（用于 Welcome 页面获取可用 token） */
export async function getFirstUploadToken(): Promise<string> {
  const result = await uploadTokenLists({});
  const tokens = result?.data || [];
  const activeToken = tokens.find((t) => t.state === 'open' && t.uploadToken);
  return activeToken?.uploadToken || '';
}
