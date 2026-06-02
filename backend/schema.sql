-- Digital Space Station · PostgreSQL Schema

CREATE TABLE IF NOT EXISTS users (
    id         BIGSERIAL PRIMARY KEY,
    username   TEXT        NOT NULL UNIQUE,
    password   TEXT        NOT NULL,
    role       TEXT        NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'user')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS posts (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title      TEXT        NOT NULL,
    slug       TEXT        NOT NULL UNIQUE,
    excerpt    TEXT        NOT NULL DEFAULT '',
    content    TEXT        NOT NULL DEFAULT '',
    cover_url  TEXT        NOT NULL DEFAULT '',
    tags       TEXT[]      NOT NULL DEFAULT '{}',
    status     TEXT        NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published')),
    read_time  TEXT        NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_posts_user ON posts (user_id);
CREATE INDEX IF NOT EXISTS idx_posts_status_created ON posts (status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_posts_tags ON posts USING GIN (tags);

CREATE TABLE IF NOT EXISTS moments (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    text       TEXT        NOT NULL,
    images     TEXT[]      NOT NULL DEFAULT '{}',
    mood       TEXT        NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_moments_user_created ON moments (user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS milestones (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title      TEXT        NOT NULL,
    detail     TEXT        NOT NULL DEFAULT '',
    icon       TEXT        NOT NULL DEFAULT '🏁',
    event_date DATE        NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_milestones_user_date ON milestones (user_id, event_date DESC);
