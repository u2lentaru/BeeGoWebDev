package controllers

import (
	"BeeGoWebDev/models"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"

	"github.com/astaxie/beego"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// BlogController struct
type BlogController struct {
	beego.Controller
	Explorer Explorer
}

// Explorer struct
type Explorer struct {
	Db           *mongo.Client
	DbName       string
	DbCollection string
}

// Truncate - truncate database. Cyclomatic complexity 1
func (e Explorer) Truncate() error {
	c := e.Db.Database(e.DbName).Collection(e.DbCollection)
	_, err := c.DeleteMany(context.Background(), bson.D{})

	return err
}

// InsertDefault - create posts. Cyclomatic complexity 3
func (e Explorer) InsertDefault() error {
	for _, post := range createPosts() {
		if err := e.addPost(post); err != nil {
			return err
		}
	}

	return nil
}

// createPosts. Cyclomatic complexity 1
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

// addPost. Cyclomatic complexity 1
func (e Explorer) addPost(post models.TPost) error {
	c := e.Db.Database(e.DbName).Collection(e.DbCollection)
	_, err := c.InsertOne(context.Background(), post)

	return err
}

// Get func. Cyclomatic complexity 2
func (c *BlogController) Get() {
	blog, err := c.Explorer.getBlog()
	if err != nil {
		log.Println(err)
		return
	}

	c.Data["Blog"] = blog
	c.TplName = "blogs.tpl"
}

// getBlog. Cyclomatic complexity 3
func (e Explorer) getBlog() ([]models.TPost, error) {
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

type postRequest struct {
	ID       string `json:"id"`
	Subj     string `json:"subj"`
	PostTime string `json:"posttime"`
	PostText string `json:"posttext"`
}

/*
	curl.exe -vX POST -H "Content-Type: application/json"  -d "@data.json" http://localhost:8080/
*/

// Post func. Cyclomatic complexity 3
func (c *BlogController) Post() {
	resp := new(postRequest)

	if err := readAndUnmarshall(resp, c.Ctx.Request.Body); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
		return
	}

	post := models.TPost{
		ID:       resp.ID,
		Subj:     resp.Subj,
		PostTime: resp.PostTime,
		PostText: resp.PostText,
	}

	if err := c.Explorer.addPost(post); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
		return
	}

	c.Ctx.ResponseWriter.WriteHeader(200)
	_, _ = c.Ctx.ResponseWriter.Write([]byte(`SUCCESS\n`))
}

// readAndUnmarshall. Cyclomatic complexity 3
func readAndUnmarshall(resp interface{}, body io.ReadCloser) error {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Print("empty id")
		return err
	}
	if err := json.Unmarshal(bytes, resp); err != nil {
		return err
	}
	return nil
}

/*
	curl.exe -vX PUT -H "Content-Type: application/json"  -d"@update.json" http://localhost:8080?id=1
*/

// Put func. Cyclomatic complexity 4
func (c *BlogController) Put() {
	id := c.Ctx.Request.URL.Query().Get("id")

	if len(id) == 0 {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte("empty id"))
		return
	}

	resp := new(postRequest)

	if err := readAndUnmarshall(resp, c.Ctx.Request.Body); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
		return
	}

	post := models.TPost{
		Subj:     resp.Subj,
		PostTime: resp.PostTime,
		PostText: resp.PostText,
	}

	if err := c.Explorer.editPost(&post, id); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
	}

	c.Ctx.ResponseWriter.WriteHeader(200)
	_, _ = c.Ctx.ResponseWriter.Write([]byte(`SUCCESS`))
}

// editPost. Cyclomatic complexity 1
func (e Explorer) editPost(post *models.TPost, id string) error {
	filter := bson.D{{Key: "id", Value: id}}

	update := createUpdates(*post)

	c := e.Db.Database(e.DbName).Collection(e.DbCollection)
	_, err := c.UpdateOne(context.Background(), filter, update)

	return err
}

// createUpdates. Cyclomatic complexity 4
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

/*
	curl.exe -vX DELETE  http://localhost:8080?id=2
*/

// Delete func. Cyclomatic complexity 3
func (c *BlogController) Delete() {
	id := c.Ctx.Request.URL.Query().Get("id")

	if len(id) == 0 {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte("empty id"))
		return
	}

	err := c.Explorer.deletePost(id)

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
	}

	c.Ctx.ResponseWriter.WriteHeader(200)
	_, _ = c.Ctx.ResponseWriter.Write([]byte(`SUCCESS`))

}

// deletePost. Cyclomatic complexity 1
func (e Explorer) deletePost(id string) error {
	filter := bson.D{{Key: "id", Value: id}}

	c := e.Db.Database(e.DbName).Collection(e.DbCollection)
	_, err := c.DeleteOne(context.Background(), filter)

	return err

}
