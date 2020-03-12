package controllers

import (
	"BeeGoWebDev/models"
	"database/sql"
	"log"

	"github.com/astaxie/beego"
)

// PostController struct
type PostController struct {
	beego.Controller
	Db       *sql.DB
	currBlog string
}

// Get func
func (c *PostController) Get() {
	c.currBlog = "1"
	id := c.Ctx.Request.URL.Query().Get("id")

	if len(id) == 0 {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte("empty id"))
		return
	}

	post, err := getPost(c.Db, c.currBlog, id)
	if err != nil {
		log.Fatal(err)
		return
	}

	c.Data["Post"] = post
	c.TplName = "post.tpl"
}

func getPost(db *sql.DB, blogid, id string) (models.TPost, error) {
	post := models.TPost{}

	row := db.QueryRow("select * from myblog.posts where blogs.id = ?", id)
	err := row.Scan(&post.ID, &post.Subj, &post.PostTime, &post.PostText)
	if err != nil {
		return models.TPost{}, err
	}

	return post, nil
}

/*type postRequest struct {
	Subj     string `json:"subj"`
	PostTime string `json:"posttime"`
	PostText string `json:"posttext"`
}*/

/*
	curl.exe -vX POST -H "Content-Type: application/json"  -d "@data.json" http://localhost:8080/
*/

// Post func
func (c *PostController) Post() {
	c.currBlog = "1"

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

	if err := createPost(c.Db, c.currBlog, post); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
		return
	}

	c.Ctx.ResponseWriter.WriteHeader(200)
	_, _ = c.Ctx.ResponseWriter.Write([]byte(`SUCCESS\n`))
}

/*func createPost(db *sql.DB, currBlog string, post models.TPost) error {
	_, err := db.Exec("insert into myblog.posts (blogid,subj,posttime,posttext) values (?,?,?,?)",
		currBlog, post.Subj, post.PostTime, post.PostText)

	return err
}*/

/*type putRequest struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}*/

/*
	curl.exe -vX PUT -H "Content-Type: application/json"  -d"@update.json" http://localhost:8080?id=1
*/

// Put func
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

	if err := updatePost(c.Db, id, post.Subj, post.PostTime, post.PostText); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
	}

	c.Ctx.ResponseWriter.WriteHeader(200)
	_, _ = c.Ctx.ResponseWriter.Write([]byte(`SUCCESS`))
}

func updatePost(db *sql.DB, id, subj, posttime, posttext string) error {
	if len(subj) == 0 && len(posttime) == 0 && len(posttime) == 0 {
		return nil
	}

	_, err := db.Exec("UPDATE myblog.posts SET subj=?, posttime=?, posttext=? WHERE id=?",
		subj, posttime, posttext, id)

	return err
}

/*
	curl.exe -vX DELETE  http://localhost:8080?id=42
*/

// Delete func
func (c *PostController) Delete() {
	id := c.Ctx.Request.URL.Query().Get("id")

	err := deletePost(c.Db, id)

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
	}

	c.Ctx.ResponseWriter.WriteHeader(200)
	_, _ = c.Ctx.ResponseWriter.Write([]byte(`SUCCESS`))
}

/*func deletePost(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM myblog.posts WHERE `id`=?", id)

	if err != nil {
		return err
	}

	return nil
}*/
