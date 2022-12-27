package entity

import "time"

type Doc struct {
	Id          uint64     `json:"id,omitempty"`
	Pseudo      string     `json:"r_pseudo,omitempty"`
	Name        string     `json:"name"`
	Header      string     `json:"header"`
	Type        string     `json:"type"`
	SubType     string     `json:"sub_type"`
	Rev         string     `json:"rev"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Keywords    string     `json:"keywords"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	Chapters    []Chapter  `json:"chapters"`
}
