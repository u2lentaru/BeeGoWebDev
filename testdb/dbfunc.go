package testdb

import (
	"BeeGoWebDev/models"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type explorer struct {
	Db           *mongo.Client
	DbName       string
	DbCollection string
}

func (e explorer) addPost(post models.TPost) error {
	c := e.Db.Database(e.DbName).Collection(e.DbCollection)
	_, err := c.InsertOne(context.Background(), post)

	return err
}
