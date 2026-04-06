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
    mustChangePassword?: boolean;
  };

  type LoginResult = {
    status?: string;
    type?: string;
    currentAuthority?: string;
    token?: string;
    mustChangePassword?: boolean;
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
    fileBaseDir?: string;
    absFileBaseDir?: string;
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

  // 系统配置
  interface SysConfig {
    fileBaseDir?: string;     // 文件根目录（只读，只能通过配置文件修改）
    uploadDir?: string;       // 上传子目录
    downloadDir?: string;     // 下载子目录
    fullUploadDir?: string;   // 完整上传路径（只读）
    fullDownloadDir?: string; // 完整下载路径（只读）
    httpPort?: number;
    fileUploadEnable?: boolean;
    enableUploadToken?: boolean;
    uploadPolicyDeadline?: number;
    uploadPolicyFSizeMin?: number;
    uploadPolicyFSizeLimit?: number;
    persistentNotifyUrl?: string;
    notifyEnable?: boolean;
    thumbWidth?: number;
    thumbHeight?: number;
    logLevel?: number;
  }

  interface SysConfigUpdateResult {
    changes?: string[];
    needRestart?: boolean;
    message?: string;
  }

  // ===== 分享功能 =====
  interface ShareItem {
    id?: number;
    shareCode?: string;
    filePath?: string;
    fileName?: string;
    fileType?: string;
    createdBy?: string;
    extractCode?: string;
    expireAt?: string;
    maxDownloads?: number;
    curDownloads?: number;
    isActive?: boolean;
    createdAt?: string;
    updatedAt?: string;
  }

  interface CreateShareParams {
    filePath: string;
    fileName: string;
    fileType?: string;
    extractCode?: string;
    expireHours?: number;
    maxDownloads?: number;
  }

  interface CreateShareResult {
    shareCode?: string;
    extractCode?: string;
    shareUrl?: string;
    expireAt?: string;
  }

  interface ShareInfoResult {
    valid?: boolean;
    reason?: string;
    share?: {
      shareCode?: string;
      fileName?: string;
      fileType?: string;
      hasExtractCode?: boolean;
      expireAt?: string;
      maxDownloads?: number;
      curDownloads?: number;
      isActive?: boolean;
      createdBy?: string;
      createdAt?: string;
    };
  }

  interface ShareStats {
    totalShares?: number;
    activeShares?: number;
    expiredShares?: number;
    totalDownloads?: number;
  }

  interface ShareConfig {
    shareEnable?: boolean;
    anonymousAccess?: boolean;
  }
}
