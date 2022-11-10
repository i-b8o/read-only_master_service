package entity

type Link struct {
	ID           uint64 `json:"id"`
	ChapterID    uint64 `json:"chapter_id"`
	ParagraphNum uint64 `json:"paragraph_num"`
	RID          uint64 `json:"r_id"`
}
