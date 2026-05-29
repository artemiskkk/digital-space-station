package model

import "time"

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Excerpt   string    `json:"excerpt"`
	Content   string    `json:"content,omitempty"`
	CoverURL  string    `json:"cover_url"`
	Tags      []string  `json:"tags"`
	Status    string    `json:"status"` // draft | published
	ReadTime  string    `json:"read_time"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Moment struct {
	ID        int64     `json:"id"`
	Text      string    `json:"text"`
	Images    []string  `json:"images"`
	Mood      string    `json:"mood"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password,omitempty"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type Milestone struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Detail    string    `json:"detail"`
	Icon      string    `json:"icon"`
	EventDate string    `json:"event_date"`
	CreatedAt time.Time `json:"created_at"`
}

type ListMeta struct {
	Total  int `json:"total"`
	Page   int `json:"page"`
	Size   int `json:"size"`
	Pages  int `json:"pages"`
}

type PostList struct {
	Items []Post   `json:"items"`
	Meta  ListMeta `json:"meta"`
}

type MomentList struct {
	Items []Moment `json:"items"`
	Meta  ListMeta `json:"meta"`
}
