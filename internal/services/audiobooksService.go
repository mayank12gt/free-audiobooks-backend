package services

import (
	"github.com/mayank12gt/free-audiobooks-backend/internal/repos"
)

type AudiobookService struct {
	audiobookRepo repos.AudiobooksRepo
}

type Query struct {
	Search         string
	Genres         []string
	Language       string
	TotalTimeRange TimeRange
	PageSize       int
	Page           int
	Sort           string
}

type TimeRange struct {
	TotalTimeMin int64
	TotalTimeMax int64
}

func (s *AudiobookService) List(query Query) ([]*repos.Audiobook, repos.Metadata, error) {

	query.TotalTimeRange.TotalTimeMax = query.TotalTimeRange.TotalTimeMax * 60
	query.TotalTimeRange.TotalTimeMin = query.TotalTimeRange.TotalTimeMin * 60

	audiobooks, meta, err := s.audiobookRepo.List(query.Search, query.Genres, query.Language, query.TotalTimeRange.TotalTimeMin, query.TotalTimeRange.TotalTimeMax, int64(query.Page), int64(query.PageSize), query.Sort)
	if err != nil {
		return nil, meta, err
	}

	return audiobooks, meta, nil

}

func (s *AudiobookService) Get(id string) (*repos.Audiobook, error) {
	return s.audiobookRepo.Get(id)
}
