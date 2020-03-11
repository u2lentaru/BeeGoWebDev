package controllers

import (
	"BeeGoWebDev/models"
	"database/sql"
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

	//c.Data["Blogs"] = blog
	c.Data["PostList"] = blog.PostList
	c.TplName = "blogs.tpl"
}

func getBlog(db *sql.DB, id string) (models.TBlog, error) {
	/*rows, err := db.Query("select * from myblog.posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make(models.TPost, 0, 1)
	for rows.Next() {
		post := models.Tposts{
			if err := rows.Scan(&blog.Id, &blog.Name, &&blog.Title); err != nil {
				log.Println(err)
				continue
			}
		}
	}*/
	///////////////////////

	blog := models.TBlog{}
	//blog := make(models.TBlog, 0, 1)
	log.Println(db)
	log.Println(id)
	if err := db.Ping(); err != nil {
		log.Println("2 db not pinged!")
		log.Fatal(err)
	}
	log.Println("2 db pinged!")
	row := db.QueryRow("select * from myblog.blogs where blogs.id = ?", id)
	err := row.Scan(&blog.ID, &blog.Name, &blog.Title)
	if err != nil {
		log.Println("1 query")
		return models.TBlog{}, err
	}

	if err := db.Ping(); err != nil {
		log.Println("3 db not pinged!")
		log.Fatal(err)
	}
	log.Println("3 db pinged!")

	rows, err := db.Query("select * from posts where blogid = ?", id)
	if err != nil {
		log.Println("2 query")
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
	log.Println(blog)
	return blog, nil
}
