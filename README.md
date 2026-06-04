# Digital Space Station (DSS)

> 个人数字空间站 —— 零成本、免运维、前后端分离的个人专属站点。
> 用于沉淀技术博客、记录碎片化灵感（Moments）、留存人生里程碑。

线上地址：**https://artemis-space.vercel.app**

---

## ✨ 功能

- **博客（Blog）** —— 支持 Markdown 长文，含标题、摘要、分类、标签、草稿/发布状态。
- **碎片（Moments）** —— 极简的「短文本 + 心情」随手记。
- **里程碑（Milestones）** —— 时间轴形式记录重要时刻。
- **单一管理员鉴权** —— 不开放注册，JWT 登录，前台直接内联编辑/发布。
- **桌面看板娘** —— 一只会漂浮、追鼠标、说话的 SVG 小机器人。

---

## 🧱 技术栈与架构

```
┌─────────────────┐      ┌──────────────────┐      ┌────────────────┐
│  space-station  │      │     backend      │      │    Supabase    │
│  Astro (前台)   │ ───► │   Go + Gin (API) │ ───► │  PostgreSQL    │
│  Vercel 托管    │      │   Render 托管    │      │   (数据库)     │
└─────────────────┘      └──────────────────┘      └────────────────┘
```

| 层 | 技术 | 托管 |
|----|------|------|
| 前台 | [Astro](https://astro.build) 4 + Tailwind CSS | Vercel |
| 后端 | Go + [Gin](https://gin-gonic.com) + [pgx](https://github.com/jackc/pgx) + JWT | Render |
| 数据库 | PostgreSQL | Supabase |
| 后台（可选） | Vue 3 + Vite（`admin/`，目前编辑直接在前台完成） | — |

**性能设计**：首页与博客列表为**静态预渲染（SSG）**，挂在 Vercel CDN 上秒开；
文章数据由浏览器**客户端异步拉取**（带骨架屏），不阻塞于后端冷启动。

---

## 📁 仓库结构（Monorepo）

```
digital_space_dev/
├── backend/            # Go API 服务
│   ├── main.go
│   ├── schema.sql      # 数据库建表脚本
│   └── internal/
│       ├── config/     # 环境变量加载
│       ├── handler/    # 路由处理（blog / moments / milestones / auth / upload）
│       ├── middleware/ # CORS、JWT 鉴权
│       ├── model/      # 数据模型
│       └── store/      # 数据库查询
├── space-station/      # Astro 前台
│   └── src/
│       ├── pages/      # 路由页面
│       ├── components/ # 组件（含看板娘 Companion）
│       └── layouts/    # 全局布局
└── admin/              # Vue 后台（可选）
```

---

## 🚀 本地开发

### 后端

```bash
cd backend
cp .env.example .env        # 填写下方环境变量
go run .                    # 默认监听 :8080
```

`.env` 必填项：

| 变量 | 说明 |
|------|------|
| `DATABASE_URL` | PostgreSQL 连接串（Supabase 提供） |
| `JWT_SECRET` | JWT 签名密钥（随机字符串） |
| `ADMIN_PASSWORD` | 管理员初始密码 |
| `ADMIN_USERNAME` | 管理员用户名（可选，默认 `admin`） |
| `FRONTEND_URL` | 前台地址，用于 CORS 白名单 |
| `PORT` | 端口（可选，默认 `8080`） |

首次启动会自动建立管理员账号（若数据库无 admin 角色用户）。
建表脚本见 `backend/schema.sql`，在 Supabase SQL Editor 中执行。

### 前台

```bash
cd space-station
npm install
npm run dev                 # 默认 http://localhost:4321
```

前台环境变量（`.env`）：

| 变量 | 说明 |
|------|------|
| `PUBLIC_API_URL` | 后端 API 地址，如 `http://localhost:8080` |

---

## ☁️ 部署

| 组件 | 平台 | 关键配置 |
|------|------|----------|
| 后端 | Render（Web Service） | Root: `backend`，Build: `go build -o app .`，Start: `./app` |
| 前台 | Vercel | Root: `space-station`，Framework: Astro，Node: 20.x |
| 数据库 | Supabase | 执行 `backend/schema.sql` 建表 |

部署时记得：
- Render 配置全部环境变量，`FRONTEND_URL` 指向 Vercel 域名。
- Vercel 配置 `PUBLIC_API_URL` 指向 Render 域名。
- 修改任一域名后，需同步更新对端配置（否则 CORS 拦截）。

> ⚠️ Render 免费实例闲置约 15 分钟会休眠，再次访问有数十秒冷启动延迟，属正常现象。
> 可用 UptimeRobot 等免费服务定时 ping `/health` 保活。

---

## 🔌 API 概览

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| POST | `/api/auth/login` | — | 登录获取 token |
| GET | `/api/public/posts` | — | 公开的已发布文章列表 |
| GET | `/api/posts` | ✅ | 本人全部文章 |
| GET | `/api/posts/:slug` | ✅ | 单篇文章 |
| POST | `/api/posts` | ✅ | 发布文章 |
| DELETE | `/api/posts/:id` | ✅ | 删除文章 |
| GET/POST/DELETE | `/api/moments` | ✅ | 碎片增删查 |
| GET/POST/DELETE | `/api/milestones` | ✅ | 里程碑增删查 |
| GET | `/health` | — | 健康检查 |

---

## 📄 License

个人项目，自用为主。
