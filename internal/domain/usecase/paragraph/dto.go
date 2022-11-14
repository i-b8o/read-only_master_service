package usecase_paragraph

import "regulations_supreme_service/internal/domain/entity"

type CreateParagraphsInput struct {
	Paragraphs []entity.Paragraph
}

type CreateParagraphsOutput struct {
	Message string
}
