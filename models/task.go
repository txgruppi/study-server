package models

import "time"

type Task struct {
	ID    string     `json:"id,omitempty" badgerhold:"key"`
	Title string     `json:"title,omitempty"`
	Start time.Time  `json:"start,omitempty"`
	End   *time.Time `json:"end,omitempty"`
}
