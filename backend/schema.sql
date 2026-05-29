-- Digital Space Station · PostgreSQL Schema

CREATE TABLE IF NOT EXISTS posts (
    id         BIGSERIAL PRIMARY KEY,
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

CREATE INDEX IF NOT EXISTS idx_posts_status_created ON posts (status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_posts_tags ON posts USING GIN (tags);

CREATE TABLE IF NOT EXISTS moments (
    id         BIGSERIAL PRIMARY KEY,
    text       TEXT        NOT NULL,
    images     TEXT[]      NOT NULL DEFAULT '{}',
    mood       TEXT        NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_moments_created ON moments (created_at DESC);

CREATE TABLE IF NOT EXISTS milestones (
    id         BIGSERIAL PRIMARY KEY,
    title      TEXT        NOT NULL,
    detail     TEXT        NOT NULL DEFAULT '',
    icon       TEXT        NOT NULL DEFAULT '🏁',
    event_date DATE        NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_milestones_date ON milestones (event_date DESC);
