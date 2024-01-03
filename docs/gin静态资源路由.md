

## router.StaticFS("/" 路由冲突问题
router.StaticFS("/", http.Dir(StaticDir))  如果指定"/"，
在 /api/v1之前注册，就会存在路由冲突。

当使用 router.StaticFS("/", http.Dir(StaticDir)) 定义静态文件服务时，它会将根路径 / 下的所有请求都交给静态文件服务处理。
因此，如果你的 /api/v1 开头的路由与静态文件服务发生冲突，可以考虑以下解决方案：

思路1：将静态文件服务的路径调整为一个更具体的路径，而不是根路径 /。例如，可以将其修改为 /static 或其他非冲突的路径：
```bash
router.StaticFS("/static", http.Dir(StaticDir))
```

思路2: 将静态文件服务的定义放在具体路径的路由之前，这样具体路径的路由会优先匹配。例如：
```bash
// 先定义具体路径的路由
router.GET("/api/v1/your-route", func(c *gin.Context) {
// 处理具体路径的逻辑
})

// 再定义静态文件服务
router.StaticFS("/", http.Dir(StaticDir))
通过以上两种方式，可以避免 /api/v1 路径与静态文件服务发生冲突。
```

本人打算： 静态文件服务路径：为静态文件服务选择一个非常用的路径，避免与其他具体的路由发生冲突。
例如，使用 /static、/assets 等较为常见的路径。

如果你按照上述代码将静态文件服务的路径定义为 /static，那么前端在请求静态文件时需要添加 /static 路由前缀。

例如，如果你有一个名为 main.css 的静态文件，在前端页面中引用该文件的 URL 应为 /static/main.css。

### 如何在前端代码中添加 /static 路由前缀
官方：https://pro.ant.design/zh-CN/config/config/

在Ant Design Pro的React项目中，如果你想为静态资源添加路由前缀"/static"，你可以按照以下步骤进行操作：


publicPath
Type: publicPath
Default: /
配置 webpack 的 publicPath。当打包的时候，webpack 会在静态文件路径前面添加 publicPath 的值，
当你需要修改静态文件地址时，比如使用 CDN 部署，把 publicPath 的值设为 CDN 的值就可以。如果使用一些特殊的文件系统，比如混合开发或者 cordova 等技术，可以尝试将 publicPath 设置成 ./ 相对路径。

相对路径 ./ 有一些限制，例如不支持多层路由 /foo/bar，只支持单层路径 /foo

如果你的应用部署在域名的子路径上，例如 https://www.your-app.com/foo/，你需要设置 publicPath 为 /foo/，如果同时要兼顾开发环境正常调试，你可以这样配置：

import { defineConfig } from 'umi';

export default defineConfig({
publicPath: process.env.NODE_ENV === 'production' ? '/foo/' : '/',
});


总结：当然，如果我们使用nginx处理前端，我们应该不用设置这个，我们的静态资源放在根目录下。
前端其实是一个单页面应用，初始化时候能拿到静态资源即可，之后的跳转都是纯前端操作。

我们这里要处理，主要是因为我们的go程序即要处理静态资源返回，还要处理 api接口请求，如果静态资源直接放在跟目录下，
这块存在路由冲突问题。 

所以这里的思路就是，配置前端publicPath配置，生产环境静态资源从 static 路由获取，完成整个单页面应用初始化。


### F5 刷新日志查询界面 404 问题
测试发现，虽然配置了 配置前端publicPath，生产环境静态资源从 static 路由获取，但是F5刷新日志查询界面，还是会404。

前端会根据你的路由配置，如下，向 /list-page 发起请求，这个请求会被gin的路由处理，返回404。 实际上这个请求后台有对应的静态文件，所以前端应该请求的是
/static/list-page/index.html，这个请求才会被gin的静态文件处理，返回静态文件。

```bash
{
name: 'list.table-list',
icon: 'table',
path: '/list-page',
component: './TableList',
},
```
说明前端配置的 publicPath，对路由配置无效。

解决思路，修改前端路由配置，我们直接把static加上。
```bash
{
name: 'list.table-list',
icon: 'table',
path: '/static/list-page',
component: './TableList',
},
```
发现这样还不行，前端编译后， 会多一个static目录，实际的后端路由变成了 /static/static/list-page

最终解决办法：
因为在F5刷新页面时，浏览器会直接向后端发送请求，而不会经过前端路由。
要解决这个问题，您需要在后端服务器上进行相应的配置，以确保在刷新页面时正确处理前端路由，并返回前端应用程序的入口文件。

