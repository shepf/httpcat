export default [
  {
    path: '/s/:code',
    layout: false,
    component: './SharePage',
  },
  {
    path: '/user',
    layout: false,
    routes: [
      {
        name: 'login',
        path: '/user/login',
        component: './user/Login',
      },
      {
        name: 'change-password',
        path: '/user/change-password',
        component: './user/ChangePassword',
      },
      {
        component: './404',
      },
    ],
  },
  {
    path: '/welcome',
    name: 'welcome',
    icon: 'smile',
    component: './Welcome',
  },
  {
    path: '/admin',
    name: 'admin',
    icon: 'crown',
    access: 'canAdmin',
    routes: [
      {
        path: '/admin/sysinfo-page',
        name: 'sysinfo-page',
        icon: 'smile',
        component: './sysInfo',
      },
      {
        path: '/admin/sys-config',
        name: 'sys-config',
        icon: 'setting',
        component: './SysConfig',
      },
      {
        path: '/admin/sub-page',
        name: 'sub-page',
        icon: 'smile',
        component: './uploadTokenManage',
      },
      {
        component: './404',
      },
    ],
  },
  {
    path: '/file_manage',
    name: 'file_manage',
    icon: 'file',
    access: 'canAdmin',
    routes: [
      {
        path: '/file_manage/file-list-page',
        name: 'file-list-page',
        icon: 'folder',
        component: './FileManage/FileList',
      },
      {
        path: '/file_manage/image-manage-page',
        name: 'image-manage-page',
        icon: 'picture',
        component: './FileManage/ImageManage',
      },
      {
        path: '/file_manage/share-manage',
        name: 'share-manage',
        icon: 'share',
        component: './ShareManage',
      },
      {
        component: './404',
      },
    ],
  },  
  {
    name: 'list.table-list',
    icon: 'table',
    path: '/list-page',
    component: './TableList',
  },
  {
    path: '/',
    redirect: '/welcome',
  },
  {
    component: './404',
  },
];
