package repos

import "go.mongodb.org/mongo-driver/mongo"

type Repos struct {
	audiobooksRepo AudiobooksRepo
}

func NewRepos(db *mongo.Database) Repos {
	return Repos{
		audiobooksRepo: AudiobooksRepo{
			DB: db,
		},
	}
}
