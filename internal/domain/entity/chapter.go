package entity

import "time"

type Chapter struct {
	ID         uint64      `json:"id,omitempty"`
	Pseudo     string      `json:"c_pseudo,omitempty"`
	Name       string      `json:"name"`
	Num        string      `json:"num,omitempty"`
	DocID      uint64      `json:"doc_id,omitempty"`
	OrderNum   uint32      `json:"order_num"`
	Paragraphs []Paragraph `json:"paragraphs,omitempty"`
	UpdatedAt  *time.Time  `json:"updated_at,omitempty"`
}
