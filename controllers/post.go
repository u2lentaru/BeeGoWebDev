package controllers

import (
	"BeeGoWebDev/models"
	"context"
	"log"

	"github.com/astaxie/beego"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

// PostController struct
type PostController struct {
	beego.Controller
	Explorer Explorer
}

// Get func. Cyclomatic complexity 3
func (c *PostController) Get() {
	id := c.Ctx.Request.URL.Query().Get("id")

	if len(id) == 0 {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte("empty id"))
		return
	}

	post, err := c.Explorer.getPost(id)
	if err != nil {
		log.Fatal(err)
		return
	}

	c.Data["Post"] = post
	c.TplName = "post.tpl"
}

// getPost. Cyclomatic complexity 2
func (e Explorer) getPost(id string) (models.TPost, error) {
	c := e.Db.Database(e.DbName).Collection(e.DbCollection)

	filter := bson.D{{Key: "id", Value: id}}

	res := c.FindOne(context.Background(), filter)

	post := new(models.TPost)
	if err := res.Decode(post); err != nil {
		return models.TPost{}, errors.Wrap(err, "decode")
	}

	return *post, nil

}

/*
	curl.exe -vX POST -H "Content-Type: application/json"  -d "@data.json" http://localhost:8080/post
*/

// Post func. Cyclomatic complexity 3
func (c *PostController) Post() {
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

/*
	curl.exe -vX PUT -H "Content-Type: application/json"  -d"@data.json" http://localhost:8080/post?id=46
*/

// Put func. Cyclomatic complexity 4
func (c *PostController) Put() {
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

/*
	curl.exe -vX DELETE  http://localhost:8080/post?id=46
*/

// Delete func. Cyclomatic complexity 3
func (c *PostController) Delete() {
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
