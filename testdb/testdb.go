package testdb

import (
	"BeeGoWebDev/models"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestaddPost add post
func TestaddPost(t *testing.T) {
	db, err := mongo.NewClient(options.Client().ApplyURI("mongodb://192.168.0.103:27017"))
	if err != nil {
		log.Fatal(err)
	}

	e := texplorer{
		Db:           db,
		DbName:       "myblog",
		DbCollection: "posts",
	}
	post := models.TPost{
		ID:       "100",
		Subj:     "NewPostSubj",
		PostTime: "2020-03-04",
		PostText: "NewPostText",
	}

	if err := e.addPost(post); err != nil {
		t.Error(err)
	}
}
