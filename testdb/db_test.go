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

type Case struct {
	Method dbMethod
	ID     string
	Post   models.TPost

	ExpectedPosts []models.TPost
	ExpectedPost  models.TPost
}

type dbMethod string

const (
	readAll dbMethod = "readAll"
	readOne dbMethod = "readOne"
	create  dbMethod = "create"
	update  dbMethod = "update"
	delete  dbMethod = "delete"
)

func TestDb(t *testing.T) {
	e, err := initDb()
	if err != nil {
		t.Error(err)
		return

	}

	defer func() {
		_ = e.Truncate()
		_ = e.Db.Disconnect(context.Background())
	}()

	for i, c := range createCases() {
		var (
			posts = make([]models.TPost, 0, 1)
			post  = models.TPost{}
			err   error
		)

		switch c.Method {
		case readAll:
			posts, err = e.getBlog()
		case readOne:
			post, err = e.getPost(c.ID)
		case create:
			err = e.addPost(c.Post)
		case update:
			err = e.editPost(&c.Post, c.ID)
		case delete:
			err = e.deletePost(c.ID)
		default:
			t.Error("unknown method")
			continue
		}

		if err != nil {
			t.Error(err)
		}

		if c.Method == readAll {
			if !reflect.DeepEqual(posts, c.ExpectedPosts) {
				t.Errorf("[%d] Expected: %v; Result: %v", i, c.ExpectedPosts, posts)
				break
			}
		} else if c.Method == readOne {
			if !reflect.DeepEqual(post, c.ExpectedPost) {
				t.Errorf("[%d] Expected: %v; Result: %v", i, c.ExpectedPost, post)
				break
			}
		}
	}
}

func createCases() []Case {
	initialPost := createPosts()
	return []Case{
		// show all
		{
			Method:        readAll,
			ExpectedPosts: initialPost, // the first our state
		},
		// show one
		{
			Method: readOne,
			ID:     "1",
			ExpectedPost: models.TPost{
				ID:       "1",
				Subj:     "1st subj",
				PostTime: "2020-01-01",
				PostText: "1st text",
			},
		},
		// add one
		{
			Method: create,
			Post: models.TPost{
				ID:       "100",
				Subj:     "NewPostSubj",
				PostTime: "2020-03-04",
				PostText: "NewPostText",
			},
		},
		// show all (check previous case)
		{
			Method: readAll,
			ExpectedPosts: append(initialPost, models.TPost{
				ID:       "100",
				Subj:     "NewPostSubj",
				PostTime: "2020-03-04",
				PostText: "NewPostText",
			}),
		},
		// edit one
		{
			Method: update,
			ID:     "100",
			Post: models.TPost{
				ID:       "100",
				Subj:     "NewUpdPostSubj",
				PostTime: "2020-03-04",
				PostText: "NewUpdPostText",
			},
		},
		// show all (check previous case)
		{
			Method: readAll,
			ID:     "100",
			ExpectedPosts: append(initialPost, models.TPost{
				ID:       "100",
				Subj:     "NewUpdPostSubj",
				PostTime: "2020-03-04",
				PostText: "NewUpdPostText",
			}),
		},
		// delete one
		{
			Method: delete,
			ID:     "100",
		},
		// show all (check previous case)
		{
			Method:        readAll,
			ID:            "100",
			ExpectedPosts: initialPost,
		},
	}
}

func initDb() (TExplorer, error) {
	db, err := mongo.NewClient(options.Client().ApplyURI("mongodb://192.168.0.102:27017"))
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
