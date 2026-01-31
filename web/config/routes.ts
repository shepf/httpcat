export default [
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
        path: '/file_manage/image-manage-page',
        name: 'image-manage-page',
        icon: 'smile',
        component: './FileManage/ImageManage',
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
