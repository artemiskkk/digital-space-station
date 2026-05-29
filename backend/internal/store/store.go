package store

import (
	"context"
	"fmt"
	"math"

	"github.com/artemis/dss-backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

// ── Users ─────────────────────────────────────────────────────────────────────

func (s *Store) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var u model.User
	err := s.db.QueryRow(ctx, `SELECT id, username, password, role, created_at FROM users WHERE username=$1`, username).
		Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Store) ListUsers(ctx context.Context) ([]model.User, error) {
	rows, err := s.db.Query(ctx, `SELECT id, username, role, created_at FROM users ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (s *Store) CreateUser(ctx context.Context, username, hashedPwd, role string) (*model.User, error) {
	var u model.User
	err := s.db.QueryRow(ctx, `INSERT INTO users (username, password, role) VALUES ($1,$2,$3) RETURNING id, username, role, created_at`,
		username, hashedPwd, role).Scan(&u.ID, &u.Username, &u.Role, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Store) UpdatePassword(ctx context.Context, userID int64, hashedPwd string) error {
	_, err := s.db.Exec(ctx, `UPDATE users SET password=$1 WHERE id=$2`, hashedPwd, userID)
	return err
}

func (s *Store) DeleteUser(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `DELETE FROM users WHERE id=$1 AND role != 'admin'`, id)
	return err
}

// ── Posts (user-scoped) ──────────────────────────────────────────────────────

func (s *Store) ListPostsByUser(ctx context.Context, userID int64, page, size int) (*model.PostList, error) {
	offset := (page - 1) * size
	var total int
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM posts WHERE user_id=$1`, userID).Scan(&total); err != nil {
		return nil, fmt.Errorf("count posts: %w", err)
	}
	rows, err := s.db.Query(ctx, `
		SELECT id, title, slug, excerpt, cover_url, tags, status, read_time, created_at, updated_at
		FROM posts WHERE user_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, userID, size, offset)
	if err != nil {
		return nil, fmt.Errorf("list posts: %w", err)
	}
	defer rows.Close()
	var posts []model.Post
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Excerpt, &p.CoverURL,
			&p.Tags, &p.Status, &p.ReadTime, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return &model.PostList{
		Items: posts,
		Meta:  model.ListMeta{Total: total, Page: page, Size: size, Pages: int(math.Ceil(float64(total) / float64(size)))},
	}, nil
}

// Public: admin's published posts (for guest preview)
func (s *Store) ListPublicPosts(ctx context.Context, page, size int) (*model.PostList, error) {
	offset := (page - 1) * size
	var total int
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM posts WHERE user_id=1 AND status='published'`).Scan(&total); err != nil {
		return nil, err
	}
	rows, err := s.db.Query(ctx, `
		SELECT id, title, slug, excerpt, cover_url, tags, status, read_time, created_at, updated_at
		FROM posts WHERE user_id=1 AND status='published' ORDER BY created_at DESC LIMIT $1 OFFSET $2`, size, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []model.Post
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Excerpt, &p.CoverURL,
			&p.Tags, &p.Status, &p.ReadTime, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return &model.PostList{
		Items: posts,
		Meta:  model.ListMeta{Total: total, Page: page, Size: size, Pages: int(math.Ceil(float64(total) / float64(size)))},
	}, nil
}

func (s *Store) GetPostBySlug(ctx context.Context, slug string, userID int64) (*model.Post, error) {
	var p model.Post
	err := s.db.QueryRow(ctx, `
		SELECT id, title, slug, excerpt, content, cover_url, tags, status, read_time, created_at, updated_at
		FROM posts WHERE slug=$1 AND user_id=$2`, slug, userID).
		Scan(&p.ID, &p.Title, &p.Slug, &p.Excerpt, &p.Content, &p.CoverURL,
			&p.Tags, &p.Status, &p.ReadTime, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Store) CreatePost(ctx context.Context, p *model.Post, userID int64) error {
	return s.db.QueryRow(ctx, `
		INSERT INTO posts (title, slug, excerpt, content, cover_url, tags, status, read_time, user_id)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, created_at, updated_at`,
		p.Title, p.Slug, p.Excerpt, p.Content, p.CoverURL, p.Tags, p.Status, p.ReadTime, userID).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (s *Store) UpdatePost(ctx context.Context, p *model.Post, userID int64) error {
	_, err := s.db.Exec(ctx, `
		UPDATE posts SET title=$1, slug=$2, excerpt=$3, content=$4, cover_url=$5,
		tags=$6, status=$7, read_time=$8, updated_at=NOW()
		WHERE id=$9 AND user_id=$10`,
		p.Title, p.Slug, p.Excerpt, p.Content, p.CoverURL, p.Tags, p.Status, p.ReadTime, p.ID, userID)
	return err
}

func (s *Store) DeletePost(ctx context.Context, id, userID int64) error {
	_, err := s.db.Exec(ctx, `DELETE FROM posts WHERE id=$1 AND user_id=$2`, id, userID)
	return err
}

// ── Moments (user-scoped) ────────────────────────────────────────────────────

func (s *Store) ListMoments(ctx context.Context, userID int64, page, size int) (*model.MomentList, error) {
	offset := (page - 1) * size
	var total int
	if err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM moments WHERE user_id=$1`, userID).Scan(&total); err != nil {
		return nil, err
	}
	rows, err := s.db.Query(ctx, `
		SELECT id, text, images, mood, created_at
		FROM moments WHERE user_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, userID, size, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var moments []model.Moment
	for rows.Next() {
		var m model.Moment
		if err := rows.Scan(&m.ID, &m.Text, &m.Images, &m.Mood, &m.CreatedAt); err != nil {
			return nil, err
		}
		moments = append(moments, m)
	}
	return &model.MomentList{
		Items: moments,
		Meta:  model.ListMeta{Total: total, Page: page, Size: size, Pages: int(math.Ceil(float64(total) / float64(size)))},
	}, nil
}

func (s *Store) CreateMoment(ctx context.Context, m *model.Moment, userID int64) error {
	return s.db.QueryRow(ctx, `
		INSERT INTO moments (text, images, mood, user_id) VALUES ($1,$2,$3,$4)
		RETURNING id, created_at`,
		m.Text, m.Images, m.Mood, userID).
		Scan(&m.ID, &m.CreatedAt)
}

func (s *Store) DeleteMoment(ctx context.Context, id, userID int64) error {
	_, err := s.db.Exec(ctx, `DELETE FROM moments WHERE id=$1 AND user_id=$2`, id, userID)
	return err
}

// ── Milestones (user-scoped) ─────────────────────────────────────────────────

func (s *Store) ListMilestones(ctx context.Context, userID int64) ([]model.Milestone, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, title, detail, icon, event_date, created_at
		FROM milestones WHERE user_id=$1 ORDER BY event_date DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []model.Milestone
	for rows.Next() {
		var m model.Milestone
		if err := rows.Scan(&m.ID, &m.Title, &m.Detail, &m.Icon, &m.EventDate, &m.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	return items, nil
}

func (s *Store) CreateMilestone(ctx context.Context, m *model.Milestone, userID int64) error {
	return s.db.QueryRow(ctx, `
		INSERT INTO milestones (title, detail, icon, event_date, user_id) VALUES ($1,$2,$3,$4,$5)
		RETURNING id, created_at`,
		m.Title, m.Detail, m.Icon, m.EventDate, userID).
		Scan(&m.ID, &m.CreatedAt)
}

func (s *Store) DeleteMilestone(ctx context.Context, id, userID int64) error {
	_, err := s.db.Exec(ctx, `DELETE FROM milestones WHERE id=$1 AND user_id=$2`, id, userID)
	return err
}
