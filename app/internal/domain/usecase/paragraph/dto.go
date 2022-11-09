package usecase_paragraph

import "prod_serv/internal/domain/entity"

type CreateParagraphsInput struct {
	Paragraphs []entity.Paragraph
}

type CreateParagraphsOutput struct {
	Message string
}
