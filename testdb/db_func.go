package testdb

import (
	"BeeGoWebDev/models"
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// TExplorer - Explorer for tests
type TExplorer struct {
	Db           *mongo.Client
	DbName       string
	DbCollection string
}

func (e TExplorer) addPost(post models.TPost) error {
	c := e.Db.Database(e.DbName).Collection(e.DbCollection)
	_, err := c.InsertOne(context.Background(), post)

	return err
}

// InsertDefault - create posts
func (e TExplorer) InsertDefault() error {
	for _, post := range createPosts() {
		if err := e.addPost(post); err != nil {
			return err
		}
	}

	return nil
}

func createPosts() []models.TPost {
	return []models.TPost{
		{
			ID:       "1",
			Subj:     "1st subj",
			PostTime: "2020-01-01",
			PostText: "1st text",
		},
		{
			ID:       "2",
			Subj:     "2nd subj",
			PostTime: "2020-01-02",
			PostText: "2nd text",
		},
		{
			ID:       "3",
			Subj:     "3rd subj",
			PostTime: "2020-01-03",
			PostText: "3rd text",
		},
	}
}

// Truncate - truncate database
func (e TExplorer) Truncate() error {
	c := e.Db.Database(e.DbName).Collection(e.DbCollection)
	_, err := c.DeleteMany(context.Background(), bson.D{})

	return err

}

//func getBlog(db *sql.DB, id string) (models.TBlog, error) {
func (e TExplorer) getBlog() ([]models.TPost, error) {
	blog := []models.TPost{}

	c := e.Db.Database(e.DbName).Collection(e.DbCollection)

	cur, err := c.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "Find")
	}

	if err := cur.All(context.Background(), &blog); err != nil {
		return nil, errors.Wrap(err, "All")
	}

	return blog, nil

}

func (e TExplorer) getPost(id string) (models.TPost, error) {
	c := e.Db.Database(e.DbName).Collection(e.DbCollection)

	filter := bson.D{{Key: "id", Value: id}}

	res := c.FindOne(context.Background(), filter)

	post := new(models.TPost)
	if err := res.Decode(post); err != nil {
		return models.TPost{}, errors.Wrap(err, "decode")
	}

	return *post, nil

}

func (e TExplorer) editPost(post *models.TPost, id string) error {
	filter := bson.D{{Key: "id", Value: id}}

	update := createUpdates(*post)

	c := e.Db.Database(e.DbName).Collection(e.DbCollection)
	_, err := c.UpdateOne(context.Background(), filter, update)

	return err

}

func createUpdates(post models.TPost) bson.D {
	update := bson.D{}
	if len(post.Subj) != 0 {
		update = append(update, bson.E{Key: "Subj", Value: post.Subj})
	}

	if len(post.PostTime) != 0 {
		update = append(update, bson.E{Key: "PostTime", Value: post.PostTime})
	}

	if len(post.PostText) != 0 {
		update = append(update, bson.E{Key: "PostText", Value: post.PostText})
	}

	return bson.D{{Key: "$set", Value: update}}

}

func (e TExplorer) deletePost(id string) error {
	filter := bson.D{{Key: "id", Value: id}}

	c := e.Db.Database(e.DbName).Collection(e.DbCollection)
	_, err := c.DeleteOne(context.Background(), filter)

	return err

}
