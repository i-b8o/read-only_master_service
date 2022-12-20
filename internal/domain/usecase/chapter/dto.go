package usecase_chapter

type CreateChapterInput struct {
	Name  string
	Num   string
	DocID uint64
}

type CreateChapterOutput struct {
	ChapterID uint64
}
