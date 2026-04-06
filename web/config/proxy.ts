/**
 * 在生产环境 代理是无法生效的，所以这里没有生产环境的配置
 * -------------------------------
 * The agent cannot take effect in the production environment
 * so there is no configuration of the production environment
 * For details, please see
 * https://pro.ant.design/docs/deploy
 */
export default {
  dev: {
    // localhost:8000/api/** -> http://127.0.0.1:8888/api/**
    '/api/': {
      // 要代理的地址
      // target: 'https://preview.pro.ant.design',
      target: 'http://127.0.0.1:8888',
      // 配置了这个可以从 http 代理到 https
      // 依赖 origin 的功能可能需要这个，比如 cookie
      changeOrigin: true,
    },
    // 分享公开接口代理到后端（GET /s/:code, POST /s/:code/verify, GET /s/:code/download）
    '/s/': {
      target: 'http://127.0.0.1:8888',
      changeOrigin: true,
      bypass: function (req: any) {
        // /download 和 /verify 路径始终代理到后端（不管 Accept header）
        if (req.url && (req.url.indexOf('/download') !== -1 || req.url.indexOf('/verify') !== -1)) {
          return null; // null = 不 bypass，走代理
        }
        // 其他路径：只代理 AJAX 请求，页面请求返回前端 index.html
        const accept = req.headers.accept || '';
        if (accept.indexOf('text/html') !== -1) {
          return '/index.html';
        }
      },
    },
  },
  test: {
    '/api/': {
      target: 'https://proapi.azurewebsites.net',
      changeOrigin: true,
      pathRewrite: { '^': '' },
    },
  },
  pre: {
    '/api/': {
      target: 'your pre url',
      changeOrigin: true,
      pathRewrite: { '^': '' },
    },
  },
};
