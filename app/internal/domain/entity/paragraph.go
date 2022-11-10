package entity

type Paragraph struct {
	ID        uint64
	Num       uint64
	HasLinks  bool
	IsTable   bool
	IsNFT     bool
	Class     string
	Content   string
	ChapterID uint64
}

type ParagraphsInput struct {
	Paragraphs []Paragraph
}

type ParagraphsOutput struct {
	Message string
}
