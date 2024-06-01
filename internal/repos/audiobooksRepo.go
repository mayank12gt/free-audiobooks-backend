package repos

import (
	"context"
	"log"
	"math"

	Error "github.com/mayank12gt/free-audiobooks-backend/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AudiobooksRepo struct {
	DB *mongo.Database
}

type AudiobookDTO struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	IDInt         int                `bson:"id" json:"id"`
	Title         string             `bson:"title" json:"title"`
	Description   string             `bson:"description" json:"description"`
	URLTextSource string             `bson:"url_text_source" json:"url_text_source"`
	Language      string             `bson:"language" json:"language"`
	CopyrightYear string             `bson:"copyright_year" json:"copyright_year"`
	NumSections   string             `bson:"num_sections" json:"num_sections"`
	URLRSS        string             `bson:"url_rss" json:"url_rss"`
	URLZipFile    string             `bson:"url_zip_file"  json:"url_zip_file"`
	URLProject    string             `bson:"url_project" json:"url_project"`
	URLLibrivox   string             `bson:"url_librivox" json:"url_librivox"`
	URLOther      string             `bson:"url_other" json:"url_other"`
	TotalTime     string             `bson:"totaltime" json:"totaltime"`
	TotalTimeSecs int                `bson:"totaltimesecs" json:"totaltimesecs"`
	Authors       []Author           `bson:"authors" json:"authors,omitempty"`
	Sections      []Section          `bson:"sections" json:"sections,omitempty"`
	Genres        []Genre            `bson:"genres" json:"genres"`
	Translators   []Translator       `bson:"translators" json:"translators"`
}

type GenreDTO struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name  string             `bson:"name" json:"name"`
	IDStr string             `bson:"id" json:"id"`
}

type Audiobook struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	IDStr         string             `bson:"id" json:"id"`
	Title         string             `bson:"title" json:"title"`
	Description   string             `bson:"description" json:"description"`
	URLTextSource string             `bson:"url_text_source" json:"url_text_source"`
	Language      string             `bson:"language" json:"language"`
	CopyrightYear string             `bson:"copyright_year" json:"copyright_year"`
	NumSections   string             `bson:"num_sections" json:"num_sections"`
	URLRSS        string             `bson:"url_rss" json:"url_rss"`
	URLZipFile    string             `bson:"url_zip_file"  json:"url_zip_file"`
	URLProject    string             `bson:"url_project" json:"url_project"`
	URLLibrivox   string             `bson:"url_librivox" json:"url_librivox"`
	URLOther      string             `bson:"url_other" json:"url_other"`
	TotalTime     string             `bson:"totaltime" json:"totaltime"`
	TotalTimeSecs int                `bson:"totaltimesecs" json:"totaltimesecs"`
	Authors       []Author           `bson:"authors" json:"authors,omitempty"`
	Sections      []Section          `bson:"sections" json:"sections,omitempty"`
	Genres        []Genre            `bson:"genres" json:"genres"`
	Translators   []Translator       `bson:"translators" json:"translators"`
}

type Author struct {
	ID        string `bson:"id" json:"id"`
	FirstName string `bson:"first_name" json:"first_name"`
	LastName  string `bson:"last_name" json:"last_name"`
	DOB       string `bson:"dob" json:"dob"`
	DOD       string `bson:"dod" json:"dod"`
}

type Genre struct {
	ID   string `bson:"id" json:"id"`
	Name string `bson:"name" json:"name"`
}

type Translator struct {
	ID        string `bson:"id" json:"id"`
	FirstName string `bson:"first_name" json:"first_name"`
	LastName  string `bson:"last_name" json:"last_name"`
	DOB       string `bson:"dob" json:"dob"`
	DOD       string `bson:"dod" json:"dod"`
}

type Section struct {
	ID            string `bson:"id" json:"id"`
	SectionNumber string `bson:"section_number" json:"section_number"`
	Title         string `bson:"title" json:"title"`
	ListenURL     string `bson:"listen_url" json:"listen_url"`
	Language      string `bson:"language" json:"language"`
	Playtime      string `bson:"playtime" json:"playtime"`
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func NewAudiobookRepo(db *mongo.Database) AudiobooksRepo {
	return AudiobooksRepo{
		DB: db,
	}
}

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}
	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

func (m *AudiobooksRepo) List(search string, genres []string, language string, totalTimeMin, totalTimeMax int64, page, page_size int64, sortBy string) ([]*Audiobook, Metadata, error) {

	collection := m.DB.Collection("audiobooks")

	filter := bson.D{}
	options := options.Find().SetProjection(bson.D{{"sections", 0}, {"translators", 0}}).SetSkip((page - 1) * page_size).SetLimit(page_size)

	if search != "" {
		log.Print(search)
		// filter = bson.D{{Key: "title", Value: search}}
		filter = append(filter, bson.E{Key: "$text", Value: bson.D{{Key: "$search", Value: search}}})
		log.Print(filter)
	}

	if len(genres) != 0 {
		log.Print(genres)
		filter = append(filter, bson.E{Key: "genres.id", Value: bson.M{"$in": genres}})
	}

	if language != "" {
		log.Print(language)
		filter = append(filter, bson.E{Key: "language", Value: language})
	}

	if totalTimeMax != 0 && totalTimeMin != 0 {
		log.Print(totalTimeMin)
		log.Print(totalTimeMax)

		filter = append(filter, bson.E{Key: "totaltimesecs", Value: bson.M{"$gte": totalTimeMin, "$lte": totalTimeMax}})
	}
	if sortBy != "" {
		log.Print("sort" + sortBy)

		options = options.SetSort(bson.D{{sortBy, 1}})
	} else {
		log.Print(sortBy)
		options = options.SetSort(bson.D{{"_id", 1}})
	}

	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, Metadata{}, Error.NewError().Set("server", "Internal Server Error")
	}
	if count == 0 {
		return nil, Metadata{}, Error.NewError().Set("message", "No records found")
	}

	meta := calculateMetadata(int(count), int(page), int(page_size))

	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer cursor.Close(context.Background())

	var audiobooks []*Audiobook
	if err = cursor.All(context.TODO(), &audiobooks); err != nil {
		return nil, Metadata{}, err
	}

	return audiobooks, meta, nil
}

func (m *AudiobooksRepo) Get(id string) (*Audiobook, error) {

	collection := m.DB.Collection("audiobooks")

	filter := bson.D{{"id", id}}
	//options := options.FindOne()

	var audiobook Audiobook
	err := collection.FindOne(context.TODO(), filter).Decode(&audiobook)
	if err != nil {

		return nil, Error.NewError().Set("client", "record not found")
	}

	return &audiobook, nil

}

func (m *AudiobooksRepo) GetGenres(page, page_size int64) ([]*GenreDTO, Metadata, error) {

	collection := m.DB.Collection("genres")
	filter := bson.D{}

	options := options.Find().SetSkip((page - 1) * page_size).SetLimit(page_size)

	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, Metadata{}, Error.NewError().Set("server", "Internal Server Error")
	}
	if count == 0 {
		return nil, Metadata{}, Error.NewError().Set("message", "No records found")
	}

	meta := calculateMetadata(int(count), int(page), int(page_size))

	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer cursor.Close(context.Background())

	var genres []*GenreDTO
	if err = cursor.All(context.TODO(), &genres); err != nil {
		return nil, Metadata{}, err
	}

	return genres, meta, nil

}
