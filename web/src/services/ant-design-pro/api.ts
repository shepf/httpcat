// @ts-ignore
/* eslint-disable */
// import { request } from 'umi';
import request from 'umi-request';

// request拦截器, 改变url 或 options.
request.interceptors.request.use((url, options) => {
  let token = localStorage.getItem('token');
  if (null === token) {
      token = '';
  }
  const authHeader = { Authorization: `Bearer ${token}` };
  return {
    url: url,
    options: { ...options, interceptors: true, headers: authHeader },
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

/** 登录接口 POST /api/v1/login/user/account */
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

/** 此处后端没有提供注释 GET /api/notices */
export async function getNotices(options?: { [key: string]: any }) {
  return request<API.NoticeIconList>('/api/notices', {
    method: 'GET',
    ...(options || {}),
  });
}

/** 获取文件上传日志列表 GET /api/rule */
export async function rule(
  params: {
    // query
    /** 当前的页码 */
    current?: number;
    /** 页面的容量 */
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  return request<API.RuleList>('/api/rule', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 获取文件上传日志列表 GET /api/v1/file//uploadHistoryLogs */
export async function uploadHistoryLogs(
  params: {
    // query
    /** 当前的页码 */
    current?: number;
    /** 页面的容量 */
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  return request<API.UploadHistoryLogsList>('/api/v1/file/uploadHistoryLogs', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}


/** 新建规则 PUT /api/rule */
export async function updateRule(options?: { [key: string]: any }) {
  return request<API.RuleListItem>('/api/rule', {
    method: 'PUT',
    ...(options || {}),
  });
}

/** 新建规则 POST /api/rule */
export async function addRule(options?: { [key: string]: any }) {
  return request<API.RuleListItem>('/api/rule', {
    method: 'POST',
    ...(options || {}),
  });
}

/** 删除上传文件历史记录 DELETE /api/v1/file/uploadHistoryLogs */
export async function removeUploadHistoryLog(options?: { [id: string]: any }) {
  console.log('removeUploadHistoryLog - options:', options);

  return request<Record<string, any>>('/api/v1/file/uploadHistoryLogs', {
    method: 'DELETE',
    params: { ...options },
  });
}








// 上传token管理api

/** 获取文件上传日志列表 GET /api/v1/file/uploadTokenLists */
export async function uploadTokenLists(
  params: {
  },
  options?: { [key: string]: any },
) {
  return request<API.UploadTokenLists>('/api/v1/user/uploadTokenLists', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

export async function saveUploadToken(options?: { [key: string]: any }) {
  console.log('saveUploadToken - options:', options);

  return request<API.UploadTokenLists>('/api/v1/user/saveUploadToken', {
    method: 'POST',
    data: {...(options || {})},
  });
}

export async function removeUploadToken(
  params: {
  },
  options?: { [key: string]: any },
) {
  return request<API.UploadTokenLists>('/api/v1/user/removeUploadToken', {
    method: 'DELETE',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}


export async function createUploadToken(options?: { [key: string]: any }) {
  console.log('createUploadToken - options:', options);

  return request<API.ErrorResponse>('/api/v1/user/createUploadToken', {
    method: 'POST',
    data: {...(options || {})},
  });
}


export async function changePasswd(options?: { [key: string]: any }) {
  console.log('changePasswd - options:', options);

  return request<API.ErrorResponse>('/api/v1/user/changePasswd', {
    method: 'POST',
    data: {...(options || {})},
  });
}


//获取软件版本信息
export async function getVersion(options?: { [key: string]: any }) {
  return request<API.AntResponseData<API.Version>>('/api/v1/conf/getVersion', {
    method: 'GET',
    ...(options || {}),
  });
}

// 获取上传统计
export async function getUploadStatistics(options?: { [key: string]: any }) {
  return request<API.AntResponseData<API.UploadStatistics>>('/api/v1/statistics/getUploadStatistics', {
    method: 'GET',
    ...(options || {}),
  });
}


// 获取下载统计
export async function getDownloadStatistics(options?: { [key: string]: any }) {
  return request<API.AntResponseData<API.DownloadStatistics>>('/api/v1/statistics/getDownloadStatistics', {
    method: 'GET',
    ...(options || {}),
  });
}

export async function getUploadAvailableSpace(options?: { [key: string]: any }) {
  return request<API.UploadAvailableSpace>('/api/v1/user/getUploadAvailableSpace', {
    method: 'GET',
    ...(options || {}),
  });
}

// httpcat 配置信息获取
export async function getConf(options?: { [key: string]: any }) {
  return request<API.AntResponseData<API.HttpcatConf>>('/api/v1/conf/getConf', {
    method: 'GET',
    ...(options || {}),
  });
}

// 获取某个文件信息
export async function getFileInfo(options?: { [key: string]: any }) {
  console.log('getFileInfo - options:', options);

  return request<API.MyResponse<API.FileInfo>>('/api/v1/file/getFileInfo', {
    method: 'GET',
    params: { ...options },
  });
}