import { Settings as LayoutSettings } from '@ant-design/pro-components';

const Settings: LayoutSettings & {
  pwa?: boolean;
  logo?: string;
} = {
  navTheme: 'light', //菜单导航 light白 | dark黑
  // 拂晓蓝
  colorPrimary: '#1890ff',
  // primaryColor: '#1890ff',
  layout: 'mix', //菜单模式,side：右侧导航，top：顶部导航,mix混合
  contentWidth: 'Fluid', //内容模式,Fluid：自适应，Fixed：定宽 1200px
  fixedHeader: false, //是否固定 header 到顶部 Boolean 默认false
  fixSiderbar: true, //是否固定导航 Boolean 默认false
  colorWeak: false,
  title: '🚀HttpCat', //标签页标题与项目标题
  pwa: false,
  // logo: 'https://gw.alipayobjects.com/zos/rmsportal/KDpgvguMpGfqaHPjicRK.svg',
  // logo: '/logo.svg',
  logo: process.env.NODE_ENV === 'production' ? '/static/logo.svg' : '/logo.svg',
  iconfontUrl: '',
};

export default Settings;
