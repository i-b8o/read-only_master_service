package entity

type Paragraph struct {
	ID        uint64 `json:"id"`
	Num       uint64 `json:"num"`
	HasLinks  bool   `json:"has_links"`
	IsTable   bool   `json:"is_table"`
	IsNFT     bool   `json:"is_nft"`
	Class     string `json:"class,omitempty"`
	Content   string `json:"content,omitempty"`
	ChapterID uint64 `json:"chapterid,omitempty"`
}

type ParagraphsInput struct {
	Paragraphs []Paragraph
}

type ParagraphsOutput struct {
	Message string
}
