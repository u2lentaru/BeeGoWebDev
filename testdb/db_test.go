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

// TExplorer - Explorer for tests
/*type TExplorer struct {
	Db           *mongo.Client
	DbName       string
	DbCollection string
}*/

// TestaddPost add post
func TestAddPost(t *testing.T) {
	/*db, err := mongo.NewClient(options.Client().ApplyURI("mongodb://192.168.0.103:27017"))
	if err != nil {
		log.Fatal(err)
	}
	*/

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

	if err := e.addPost(post); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(resultPosts, initialState) {
		t.Errorf("Expected %v, result %v", initialState, resultPosts)
	}
}

func initDb() (TExplorer, error) {
	db, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return TExplorer{}, err
	}

	if err = db.Connect(context.Background()); err != nil {
		return TExplorer{}, err
	}
	log.Print("mongo connected")

	e := TExplorer{
		Db:           db,
		DbName:       "habr",
		DbCollection: "test_collection",
	}

	if err := e.InsertDefault(); err != nil {
		return TExplorer{}, err
	}

	return e, nil
}
