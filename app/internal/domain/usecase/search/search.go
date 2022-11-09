package usecase_search

import (
	"context"
	"prod_serv/internal/domain/entity"
	"strings"
)

type SearchService interface {
	Search(ctx context.Context, searchQuery string, params ...string) ([]entity.Search, error)
	SearchLike(ctx context.Context, searchQuery string, params ...string) ([]entity.Search, error)
	RegSearch(ctx context.Context, searchQuery string, params ...string) ([]entity.Search, error)
	ChSearch(ctx context.Context, searchQuery string, params ...string) ([]entity.Search, error)
	PSearch(ctx context.Context, searchQuery string, params ...string) ([]entity.Search, error)
}

type searchUsecase struct {
	searchService SearchService
}

func NewSearchUsecase(searchService SearchService) *searchUsecase {
	return &searchUsecase{searchService: searchService}
}

func (u searchUsecase) Search(ctx context.Context, searchQuery string, params ...string) ([]entity.Search, error) {

	searchResults, err := u.searchService.Search(ctx, searchQuery, params...)
	if err != nil {
		return nil, err
	}
	if len(searchResults) == 0 {
		searchQueryFirstWord := strings.Split(searchQuery, " ")[0]
		searchResults, err = u.searchService.SearchLike(ctx, searchQueryFirstWord, params...)
		if err != nil {
			return nil, err
		}
	}
	return searchResults, nil
}

func (u searchUsecase) RegSearch(ctx context.Context, searchQuery string, params ...string) ([]entity.Search, error) {
	searchResults, err := u.searchService.RegSearch(ctx, searchQuery, params...)
	if err != nil {
		return nil, err
	}
	return searchResults, nil
}

func (u searchUsecase) ChSearch(ctx context.Context, searchQuery string, params ...string) ([]entity.Search, error) {
	searchResults, err := u.searchService.ChSearch(ctx, searchQuery, params...)
	if err != nil {
		return nil, err
	}
	return searchResults, nil
}

func (u searchUsecase) PSearch(ctx context.Context, searchQuery string, params ...string) ([]entity.Search, error) {
	searchResults, err := u.searchService.PSearch(ctx, searchQuery, params...)
	if err != nil {
		return nil, err
	}
	return searchResults, nil
}
