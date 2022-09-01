package dto

import (
	"time"

	"github.com/samber/lo"
	"github.com/txgruppi/study-server/models"
)

type Post struct {
	ID        string    `json:"id,omitempty"`
	Title     string    `json:"title,omitempty"`
	Text      string    `json:"text,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Version   int       `json:"version,omitempty"`
}

func NewPostFromModel(post *models.Post) *Post {
	return &Post{
		ID:        post.ID,
		Title:     post.Title,
		Text:      post.Text,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Version:   len(post.Events),
	}
}

type PostWithVersions struct {
	Post
	Versions []Post `json:"versions,omitempty"`
}

func NewPostWithVersionsFromModel(post *models.Post) *PostWithVersions {
	curr := models.Post{}
	versions := make([]Post, len(post.Events))
	for i, event := range post.Events {
		curr.AddAndApply(event)
		versions[i] = *NewPostFromModel(&curr)
	}

	return &PostWithVersions{
		Post:     versions[len(versions)-1],
		Versions: lo.Reverse[Post](versions[:len(versions)-1]),
	}
}
