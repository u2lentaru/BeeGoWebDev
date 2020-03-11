package controllers

import (
	"BeeGoWebDev/models"
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"

	"github.com/astaxie/beego"
)

// BlogController struct
type BlogController struct {
	beego.Controller
	Db       *sql.DB
	currBlog string
}

// Get func
func (c *BlogController) Get() {
	c.currBlog = "1"
	blog, err := getBlog(c.Db, c.currBlog)
	if err != nil {
		log.Fatal(err)
		return
	}

	c.Data["Blog"] = blog
	c.TplName = "blogs.tpl"
}

func getBlog(db *sql.DB, id string) (models.TBlog, error) {
	blog := models.TBlog{}
	//blog := make(models.TBlog, 0, 1)
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	row := db.QueryRow("select * from myblog.blogs where blogs.id = ?", id)
	err := row.Scan(&blog.ID, &blog.Name, &blog.Title)
	if err != nil {
		return models.TBlog{}, err
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("select * from posts where blogid = ?", id)
	if err != nil {
		return models.TBlog{}, err
	}
	defer rows.Close()

	for rows.Next() {
		post := models.TPost{}
		//post := make(models.TPost, 0, 1)

		err := rows.Scan(&post.ID, new(int), &post.Subj, &post.PostTime, &post.Text)
		if err != nil {
			log.Println(err)
			continue
		}

		blog.PostList = append(blog.PostList, post)
	}

	return blog, nil
}

type postRequest struct {
	Subj     string `json:"subj"`
	PostTime string `json:"posttime"`
	Text     string `json:"text"`
}

/*
	curl -vX POST -H "Content-Type: application/json"  -d'{"subj":"NewSubj",
"posttime":"02.02.2020","text":"NewText"}' http://localhost:8080/
*/

// Post func
func (c *BlogController) Post() {
	resp := new(postRequest)
	if err := readAndUnmarshall(resp, c.Ctx.Request.Body); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
		return
	}

	post := models.TPost{
		Subj:     resp.Subj,
		PostTime: resp.PostTime,
		Text:     resp.Text,
	}

	if err := createPost(c.Db, c.currBlog, post); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
		return
	}

	c.Ctx.ResponseWriter.WriteHeader(200)
	_, _ = c.Ctx.ResponseWriter.Write([]byte(`SUCCESS\n`))
}

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

func createPost(db *sql.DB, currBlog string, post models.TPost) error {
	_, err := db.Exec("insert into myblog.posts (blogid,subj,posttime,text) values (?,??,?)",
		currBlog, post.Subj, post.PostTime, post.Text)

	return err
}
