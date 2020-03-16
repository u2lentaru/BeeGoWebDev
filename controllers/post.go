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
	//Db       *sql.DB
	//currBlog string
	Explorer Explorer
}

// Get func
func (c *PostController) Get() {
	//c.currBlog = "1"
	id := c.Ctx.Request.URL.Query().Get("id")

	if len(id) == 0 {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte("empty id"))
		return
	}

	//post, err := getPost(c.Db, c.currBlog, id)
	post, err := c.Explorer.getPost(id)
	if err != nil {
		log.Fatal(err)
		return
	}

	c.Data["Post"] = post
	c.TplName = "post.tpl"
}

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

/*func getPost(db *sql.DB, blogid, id string) (models.TPost, error) {
	post := models.TPost{}

	row := db.QueryRow("select * from myblog.posts where posts.id = ?", id)
	err := row.Scan(&post.ID, new(int), &post.Subj, &post.PostTime, &post.PostText)

	if err != nil {
		return models.TPost{}, err
	}

	return post, nil
}

/*
	curl.exe -vX POST -H "Content-Type: application/json"  -d "@data.json" http://localhost:8080/post
*/

// Post func
func (c *PostController) Post() {
	//c.currBlog = "1"

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

	//if err := createPost(c.Db, c.currBlog, post); err != nil {
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

// Put func
/*func (c *PostController) Put() {
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

	if err := updatePost(c.Db, id, post.Subj, post.PostTime, post.PostText); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
	}

	c.Ctx.ResponseWriter.WriteHeader(200)
	_, _ = c.Ctx.ResponseWriter.Write([]byte(`SUCCESS`))

}

func updatePost(db *sql.DB, id, subj, posttime, posttext string) error {
	if len(subj) == 0 && len(posttime) == 0 && len(posttext) == 0 {
		return nil
	}

	_, err := db.Exec("UPDATE myblog.posts SET subj=?, posttime=?, posttext=? WHERE id=?",
		subj, posttime, posttext, id)

	return err
}*/

/*
	curl.exe -vX DELETE  http://localhost:8080/post?id=46
*/

// Delete func
/*func (c *PostController) Delete() {
	id := c.Ctx.Request.URL.Query().Get("id")

	if len(id) == 0 {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte("empty id"))
		return
	}

	err := deletePost(c.Db, id)

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
	}

	c.Ctx.ResponseWriter.WriteHeader(200)
	_, _ = c.Ctx.ResponseWriter.Write([]byte(`SUCCESS`))

}*/
