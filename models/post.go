package models

import (
	"strings"
	"time"
)

type EventPostCreated struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *EventPostCreated) ApplyTo(post *Post) {
	post.ID = t.ID
	post.Title = t.Title
	post.Text = t.Text
	post.CreatedAt = t.CreatedAt
	post.UpdatedAt = t.UpdatedAt
	post.SortableTitle = strings.ToLower(t.Title)
}

type EventPostTitleUpdated struct {
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *EventPostTitleUpdated) ApplyTo(post *Post) {
	post.Title = t.Title
	post.SortableTitle = strings.ToLower(t.Title)
	post.UpdatedAt = t.UpdatedAt
}

type EventPostTextUpdated struct {
	Text      string    `json:"text"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *EventPostTextUpdated) ApplyTo(post *Post) {
	post.Text = t.Text
	post.UpdatedAt = t.UpdatedAt
}

type EventPostTitleAndTextUpdated struct {
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *EventPostTitleAndTextUpdated) ApplyTo(post *Post) {
	post.Title = t.Title
	post.Text = t.Text
	post.SortableTitle = strings.ToLower(t.Title)
	post.UpdatedAt = t.UpdatedAt
}

type PostEvents interface {
	ApplyTo(post *Post)
}

type Post struct {
	ID            string       `json:"id,omitempty" badgerhold:"key"`
	Title         string       `json:"title,omitempty"`
	Text          string       `json:"text,omitempty"`
	CreatedAt     time.Time    `json:"created_at,omitempty"`
	UpdatedAt     time.Time    `json:"updated_at,omitempty"`
	Events        []PostEvents `json:"events,omitempty"`
	SortableTitle string       `json:"-"`
}

func (t *Post) AddAndApply(event PostEvents) {
	t.Events = append(t.Events, event)
	event.ApplyTo(t)
}
