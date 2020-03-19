package testdb

import (
	"BeeGoWebDev/models"
	"context"
	"log"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestDb test database functions
func TestDb(t *testing.T) {
	e, err := initDb()
	if err != nil {
		t.Error(err)
	}

	defer func() {
		_ = e.Truncate()
		_ = e.Db.Disconnect(context.Background())
	}()

	post := models.TPost{
		ID:       "100",
		Subj:     "NewPostSubj",
		PostTime: "2020-03-04",
		PostText: "NewPostText",
	}

	initialState := createPosts()

	if err := e.addPost(post); err != nil {
		t.Error(err)
	}

	resultPosts, err := e.getBlog()
	if err != nil {
		t.Error(err)
	}

	newPosts := append(initialState, post)
	if !reflect.DeepEqual(resultPosts, newPosts) {
		t.Errorf("Expected %v, result %v", newPosts, resultPosts)
	}

	resultPost, err := e.getPost(post.ID)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(post, resultPost) {
		t.Errorf("Expected %v, result %v", post, resultPost)
	}

	post = models.TPost{
		ID:       "100",
		Subj:     "NewUpdPostSubj",
		PostTime: "2020-03-04",
		PostText: "NewUpdPostText",
	}

	if err := e.editPost(&post, post.ID); err != nil {
		t.Error(err)
	}

	resultPosts, err = e.getBlog()
	if err != nil {
		t.Error(err)
	}

	newPosts = append(initialState, post)
	if !reflect.DeepEqual(resultPosts, newPosts) {
		t.Errorf("Expected %v, result %v", newPosts, resultPosts)
	}

	if err := e.deletePost(post.ID); err != nil {
		t.Error(err)
	}

	resultPosts, err = e.getBlog()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(resultPosts, initialState) {
		t.Errorf("Expected %v, result %v", initialState, resultPosts)
	}
}

func initDb() (TExplorer, error) {
	db, err := mongo.NewClient(options.Client().ApplyURI("mongodb://192.168.0.103:27017"))
	if err != nil {
		return TExplorer{}, err
	}

	if err = db.Connect(context.Background()); err != nil {
		return TExplorer{}, err
	}
	log.Print("Connected")

	e := TExplorer{
		Db:           db,
		DbName:       "myblog",
		DbCollection: "posts",
	}

	if err := e.InsertDefault(); err != nil {
		return TExplorer{}, err
	}

	return e, nil

}
