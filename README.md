# Digital Space Station (DSS)

> A zero-cost, maintenance-free personal website with a decoupled frontend and backend.
> Built for publishing a tech blog, jotting down fragmented thoughts (Moments), and tracking life milestones.

Live site: **https://artemis-space.vercel.app**

---

## ✨ Features

- **Blog** — Markdown long-form posts with title, excerpt, category, tags, and draft/published status.
- **Moments** — a minimal "short text + mood" quick-note feed.
- **Milestones** — a timeline of meaningful moments.
- **Single-admin auth** — no public registration; JWT login with inline editing/publishing right on the frontend.
- **Companion** — a floating SVG satellite probe that drifts, tracks the cursor, and pops speech bubbles.

---

## 🧱 Tech Stack & Architecture

```
┌─────────────────┐      ┌──────────────────┐      ┌────────────────┐
│  space-station  │      │     backend      │      │    Supabase    │
│  Astro (front)  │ ───► │   Go + Gin (API) │ ───► │  PostgreSQL    │
│  Vercel hosting │      │  Render hosting  │      │   (database)   │
└─────────────────┘      └──────────────────┘      └────────────────┘
```

| Layer | Tech | Hosting |
|-------|------|---------|
| Frontend | [Astro](https://astro.build) 4 + Tailwind CSS | Vercel |
| Backend | Go + [Gin](https://gin-gonic.com) + [pgx](https://github.com/jackc/pgx) + JWT | Render |
| Database | PostgreSQL | Supabase |
| Media | Cloudinary (unsigned direct upload) | Cloudinary |
| Admin (optional) | Vue 3 + Vite (`admin/`, editing is currently done on the frontend) | — |

**Performance design:** the homepage and blog list are **statically prerendered (SSG)** and served from Vercel's CDN for instant loads; post data is **fetched client-side** (with skeleton loaders) so pages never block on backend cold starts.

---

## 📁 Repository Layout (Monorepo)

```
digital_space_dev/
├── backend/            # Go API service
│   ├── main.go
│   ├── schema.sql      # Database schema
│   └── internal/
│       ├── config/     # Env var loading
│       ├── handler/    # Route handlers (blog / moments / milestones / auth / upload)
│       ├── middleware/ # CORS, JWT auth
│       ├── model/      # Data models
│       └── store/      # Database queries
├── space-station/      # Astro frontend
│   └── src/
│       ├── pages/      # Route pages
│       ├── components/ # Components (incl. the Companion)
│       └── layouts/    # Global layout
└── admin/              # Vue admin panel (optional)
```

---

## 🚀 Local Development

### Backend

```bash
cd backend
cp .env.example .env        # fill in the vars below
go run .                    # listens on :8080 by default
```

Required `.env` values:

| Variable | Description |
|----------|-------------|
| `DATABASE_URL` | PostgreSQL connection string (from Supabase) |
| `JWT_SECRET` | JWT signing secret (random string) |
| `ADMIN_PASSWORD` | Initial admin password |
| `ADMIN_USERNAME` | Admin username (optional, defaults to `admin`) |
| `FRONTEND_URL` | Frontend origin, used for the CORS allowlist |
| `PORT` | Port (optional, defaults to `8080`) |

On first startup the backend creates the admin account automatically (if no admin-role user exists).
The schema lives in `backend/schema.sql`; run it in the Supabase SQL Editor.

### Frontend

```bash
cd space-station
npm install
npm run dev                 # http://localhost:4321 by default
```

Frontend `.env`:

| Variable | Description |
|----------|-------------|
| `PUBLIC_API_URL` | Backend API URL, e.g. `http://localhost:8080` |

---

## ☁️ Deployment

| Component | Platform | Key settings |
|-----------|----------|--------------|
| Backend | Render (Web Service) | Root: `backend`, Build: `go build -o app .`, Start: `./app` |
| Frontend | Vercel | Root: `space-station`, Framework: Astro, Node: 20.x |
| Database | Supabase | Run `backend/schema.sql` to create tables |

Deployment notes:
- Set all backend env vars on Render; point `FRONTEND_URL` at the Vercel domain.
- Set `PUBLIC_API_URL` on Vercel to the Render domain.
- After changing either domain, update the other side too (otherwise CORS blocks requests).

> ⚠️ Render's free instance sleeps after ~15 minutes of inactivity, so the next request may take tens of seconds to wake it — this is normal.
> Use a free uptime service (e.g. UptimeRobot) to ping `/health` every 5 minutes and keep it warm.

---

## 🔌 API Overview

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/auth/login` | — | Log in and get a token |
| GET | `/api/public/posts` | — | Public list of published posts |
| GET | `/api/posts` | ✅ | The user's own posts |
| GET | `/api/posts/:slug` | ✅ | A single post |
| POST | `/api/posts` | ✅ | Create a post |
| PUT | `/api/posts/:id` | ✅ | Update a post |
| DELETE | `/api/posts/:id` | ✅ | Delete a post |
| GET/POST/DELETE | `/api/moments` | ✅ | Moments CRUD |
| GET/POST/DELETE | `/api/milestones` | ✅ | Milestones CRUD |
| GET | `/health` | — | Health check |
