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

  type RuleListItem = {
    key?: number;
    disabled?: boolean;
    href?: string;
    avatar?: string;
    name?: string;
    owner?: string;
    desc?: string;
    callNo?: number;
    status?: number;
    updatedAt?: string;
    createdAt?: string;
    progress?: number;
  };

  type RuleList = {
    data?: RuleListItem[];
    /** 列表的内容总数 */
    total?: number;
    success?: boolean;
  };

  type UploadHistoryLogsList = {
    /** 业务约定的错误码 */
    errorCode?: number;
    /** 业务上的错误信息 */
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
    isFileExist?: boolean; // 添加 isFileExist 属性
  };


  type FakeCaptcha = {
    code?: number;
    status?: string;
  };

  type LoginParams = {
    username?: string;
    password?: string;
    autoLogin?: boolean;
    type?: string;
  };

  type ErrorResponse = {
    /** 业务约定的错误码 */
    errorCode?: number;
    /** 业务上的错误信息 */
    msg?: string;
    data?: string;

  };

  type NoticeIconList = {
    data?: NoticeIconItem[];
    /** 列表的内容总数 */
    total?: number;
    success?: boolean;
  };

  type NoticeIconItemType = 'notification' | 'message' | 'event';

  type NoticeIconItem = {
    id?: string;
    extra?: string;
    key?: string;
    read?: boolean;
    avatar?: string;
    title?: string;
    status?: string;
    datetime?: string;
    description?: string;
    type?: NoticeIconItemType;
  };

  // 上传token管理相关数据类型定义
  type UploadTokenLists = {
      /** 业务约定的错误码 */
      errorCode: string;
      /** 业务上的错误信息 */
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
  }




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
  }


  //使用泛型定义请求返回数据类型
  type MyResponse<T> = {
    /** 业务约定的错误码 */
    errorCode?: number;
    /** 业务上的错误信息 */
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

}
