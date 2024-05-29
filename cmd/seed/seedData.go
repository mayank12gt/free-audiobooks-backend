package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Res struct {
	Books []Audiobook `json:"books"`
}

type Audiobook struct {
	ID            string       `bson:"id" json:"id"`
	Title         string       `bson:"title" json:"title"`
	Description   string       `bson:"description" json:"description"`
	URLTextSource string       `bson:"url_text_source" json:"url_text_source"`
	Language      string       `bson:"language" json:"language"`
	CopyrightYear string       `bson:"copyright_year" json:"copyright_year"`
	NumSections   string       `bson:"num_sections" json:"num_sections"`
	URLRSS        string       `bson:"url_rss" json:"url_rss"`
	URLZipFile    string       `bson:"url_zip_file"  json:"url_zip_file"`
	URLProject    string       `bson:"url_project" json:"url_project"`
	URLLibrivox   string       `bson:"url_librivox" json:"url_librivox"`
	URLOther      string       `bson:"url_other" json:"url_other"`
	TotalTime     string       `bson:"totaltime" json:"totaltime"`
	TotalTimeSecs int          `bson:"totaltimesecs" json:"totaltimesecs"`
	Authors       []Author     `bson:"authors" json:"authors,omitempty"`
	Sections      []Section    `bson:"sections" json:"sections,omitempty"`
	Genres        []Genre      `bson:"genres" json:"genres"`
	Translators   []Translator `bson:"translators" json:"translators"`
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

type Meta struct {
	TotalRecords int64     `bson:"total_records"`
	LastUpdated  time.Time `bson:"last_updated"`
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Print("no env file found")
	}

	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Print("No DSN found")
	}

	db, err := openDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Client().Disconnect(context.TODO()); err != nil {
			log.Panic(err)
		}
	}()

	collection_main := db.Collection("seed_stage")
	collection_incomplete := db.Collection("incomplete")

	var limit = 500
	var offset = 2000
	var len int = 500
	var cnt int = 0
	var response Res

	for len == limit {
		response, len = getPage(limit, offset)
		if len == 0 {
			break
		}
		log.Printf("%d records fetched\n", len)
		// for _, book := range response.Books {
		// 	_, err := collection.InsertOne(context.Background(), book)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// }

		var recordsToInsertMain []interface{}
		var recordsToInsertIncomplete []interface{}

		for _, book := range response.Books {
			if book.TotalTimeSecs == 0 {
				recordsToInsertIncomplete = append(recordsToInsertIncomplete, book)

			} else {

				recordsToInsertMain = append(recordsToInsertMain, book)
			}
		}

		collection_main.InsertMany(context.Background(), recordsToInsertMain)

		collection_incomplete.InsertMany(context.Background(), recordsToInsertIncomplete)

		log.Printf("%d records inserted\n", len)

		count, err := collection_main.CountDocuments(context.Background(), bson.D{})
		if err != nil {
			log.Fatal(err.Error())
		}

		count_incomplete, err := collection_incomplete.CountDocuments(context.Background(), bson.D{})
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Printf("%d records in main collection\n", count)

		log.Printf("%d records in incomplete collection\n", count_incomplete)

		if cnt == 5 {
			//break
		}
		cnt++
		offset += limit
	}

	_, err = db.Collection("meta_data").DeleteMany(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err.Error())
	}

	count, err := collection_main.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err.Error())
	}

	count_incomplete, err := collection_incomplete.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err.Error())
	}

	meta := Meta{
		TotalRecords: count + count_incomplete,
		LastUpdated:  time.Now(),
	}
	db.Collection("meta_data").InsertOne(context.Background(), meta)

}

func openDB(dsn string) (*mongo.Database, error) {
	opts := options.Client().ApplyURI(dsn)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}
	log.Print("DB connected")
	return client.Database("audiobooksDB"), nil
}

func getPage(limit, offset int) (Res, int) {

	reqString := fmt.Sprintf("https://librivox.org/api/feed/audiobooks?limit=%d&offset=%d&format=json&extended=1", limit, offset)

	fmt.Println(reqString)

	res, err := http.Get(reqString)
	if err != nil {
		log.Fatal(err)
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(responseData))
	var response Res
	json.Unmarshal(responseData, &response)
	// fmt.Println(response.Books[1])
	fmt.Println(len(response.Books))
	return response, len(response.Books)
}
