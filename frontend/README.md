# Frontend

这是一个与当前 Go API 对接的原生 HTML + JavaScript 前端，开箱即用。

## 页面

- `/frontend/index.html`：前台文章列表（筛选、搜索、分页）
- `/frontend/post.html?slug=<slug>`：文章详情
- `/frontend/admin.html`：后台登录与内容管理（文章/分类/标签/上传）

## 运行

启动后端服务后直接访问：

- `http://localhost:8080/`（会重定向到前台）
- `http://localhost:8080/frontend/admin.html`

## 说明

- 默认 API 地址为同域 `/api`。
- 登录 token 存储在 `localStorage.blog_token`。
- 如需切换 API 地址，可在控制台设置：

```js
localStorage.setItem('api_base', 'http://localhost:8080/api')
```
