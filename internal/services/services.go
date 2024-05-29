package services

import (
	"github.com/mayank12gt/free-audiobooks-backend/internal/repos"
	"go.mongodb.org/mongo-driver/mongo"
)

type Services struct {
	AudiobooksService AudiobookService
}

func NewService(db *mongo.Database) Services {
	return Services{
		AudiobooksService: AudiobookService{
			audiobookRepo: repos.AudiobooksRepo{
				DB: db,
			},
		},
	}
}
