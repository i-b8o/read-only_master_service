package usecase_link

import (
	"context"
	"fmt"
	"regulations_supreme_service/internal/domain/entity"

	"github.com/i-b8o/logging"
)

type LinkService interface {
	// GetAllByChapterID(ctx context.Context, chapterID uint64) ([]*entity.Link, error)
	// CreateForChapter(ctx context.Context, link entity.Link) error
	Create(ctx context.Context, link entity.Link) error
	GetAll(ctx context.Context) ([]*entity.Link, error)
}

type ChapterService interface {
	GetOrderNum(ctx context.Context, id uint64) (orderNum uint64, err error)
}

type linkUsecase struct {
	linkService    LinkService
	chapterService ChapterService
	logger         *logging.Logger
}

func NewLinkUsecase(linkService LinkService, chapterService ChapterService, logger *logging.Logger) *linkUsecase {
	return &linkUsecase{linkService: linkService, chapterService: chapterService, logger: logger}
}

func (u linkUsecase) CreateLink(ctx context.Context, link entity.Link) error {
	return u.linkService.Create(ctx, link)
}

func (u linkUsecase) GetDartAllLinks(ctx context.Context) string {
	links, err := u.linkService.GetAll(ctx)
	if err != nil {
		return ""
	}

	dartClass := `
	import 'link.dart';
	
	class AllLinks {
		static const List<Link> links = <Link>[
			%s
		];
	}
	`

	linkssDartString, _ := u.linksDart(ctx, links)

	return fmt.Sprintf(dartClass, linkssDartString)
}

func (u linkUsecase) linksDart(ctx context.Context, links []*entity.Link) (dartLinksList string, err error) {
	for _, l := range links {
		num, err := u.chapterService.GetOrderNum(ctx, l.ChapterID)
		if err != nil {
			return "", err
		}
		dartLinksList += fmt.Sprintf(`		Link(id: %d, chapterNum: %d, ParagraphNum: %d, RID: %d),
		`, l.ID, num, l.ParagraphNum, l.RID)
	}
	return dartLinksList, nil
}
