package usecase_paragraph

import "read-only_master_service/internal/domain/entity"

type CreateParagraphsInput struct {
	Paragraphs []entity.Paragraph
}

type CreateParagraphsOutput struct {
	Message string
}
