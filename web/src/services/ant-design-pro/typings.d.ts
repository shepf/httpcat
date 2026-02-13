// @ts-ignore
/* eslint-disable */

declare namespace API {
  type CurrentUser = {
    name?: string;
    avatar?: string;
    userid?: string;
    email?: string;
    signature?: string;
    title?: string;
    group?: string;
    tags?: { key?: string; label?: string }[];
    notifyCount?: number;
    unreadCount?: number;
    country?: string;
    access?: string;
    geographic?: {
      province?: { label?: string; key?: string };
      city?: { label?: string; key?: string };
    };
    address?: string;
    phone?: string;
  };

  type LoginResult = {
    status?: string;
    type?: string;
    currentAuthority?: string;
    token?: string;
  };

  type PageParams = {
    current?: number;
    pageSize?: number;
  };

  type UploadHistoryLogsList = {
    errorCode?: number;
    msg?: string;
    data?: {
      list: UploadHistoryLogItem[];
      current: number;
      pageSize: number;
      total: number;
    };
  };

  type UploadHistoryLogItem = {
    id?: number;
    ip?: string;
    upload_time?: string;
    filename?: string;
    file_size?: string;
    file_md5?: string;
    file_created_time?: number;
    file_modified_time?: number;
    isFileExist?: boolean;
  };

  type LoginParams = {
    username?: string;
    password?: string;
    autoLogin?: boolean;
    type?: string;
  };

  type ErrorResponse = {
    errorCode?: number;
    msg?: string;
    data?: string;
  };

  // 上传token管理相关数据类型定义
  type UploadTokenLists = {
    errorCode: string;
    msg?: string;
    data?: UploadTokenItem[];
  };

  type UploadTokenItem = {
    id?: number;
    appkey?: React.Key;
    appsecret?: string;
    uploadToken?: string;
    state?: string;
    desc?: string;
    created_at?: number;
    is_sys_built?: boolean;
  };

  // 版本信息返回格式
  interface AntResponseData<T> {
    data: T;
    success: boolean;
  }

  interface Version {
    build?: string;
    ci?: string;
    commit?: string;
    version?: string;
    uptime?: string;
  }

  type UploadAvailableSpace = {
    freeSpace?: number;
    usedSpace?: number;
  };

  interface UploadStatistics {
    monthPercentage?: string;
    monthUploadCount?: number;
    lastMonthUploadCount?: number;
    todayPercentage?: string;
    todayUploadCount?: number;
    yesterdayUploadCount?: number;
    totalUploadCount?: number;
  }

  interface DownloadStatistics {
    todayDownloadCount?: number;
    yesterdayDownloadCount?: number;
    todayPercentage?: string;
    monthDownloadCount?: number;
    lastMonthDownloadCount?: number;
    monthPercentage?: string;
    totalDownloadCount?: number;
  }

  interface HttpcatConf {
    downloadDir?: string;
    fileUploadEnable?: boolean;
    uploadDir?: string;
    webDir?: string;
    workDir?: string;
    absUploadDir?: string;
    absDownloadDir?: string;
    absWebDir?: string;
  }

  // 使用泛型定义请求返回数据类型
  type MyResponse<T> = {
    errorCode?: number;
    msg?: string;
    data?: T;
  };

  interface FileInfo {
    fileName?: string;
    lastModified?: boolean;
    md5?: string;
    md5Match?: boolean;
    size?: string;
  }

  // 图片管理相关类型
  interface ImageItem {
    FileName: string;
    ThumbnailBase64?: string;
  }

  // 文件列表项
  interface FileItem {
    FileName: string;
    LastModified: string;
    Size: string;
  }

  interface ImageListResponse {
    data: ImageItem[];
    pagination: {
      page: number;
      pageSize: number;
      totalItems: number;
    };
  }

  interface ImageUploadResponse {
    data: {
      url: string;
      thumbUrl: string;
      name: string;
    };
  }
}
