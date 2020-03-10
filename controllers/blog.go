package controllers

import (
	"database/sql"
	"log"

	"github.com/astaxie/beego"
)

type BlogController struct {
	beego.Controller
	Db *sql.DB
}

func (c *BlogController) Get() {
	posts, err := getBlog(c.Db)
	if err != nil {
		log.Fatal(err)
		return
	}

	c.Data["Blogs"] = blogs
	c.TplName = "blogs.tpl"
}

func getBlog(db *sql.DB) {
	rows, err := db.Query("select * from myblog.blogs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make(models/Blogs, 0, 1)
	for rows.Next() {
		blog := models.Blogs{
			if err := rows.Scan(&blog.Id, &blog.Name, &&blog.Title); err != nil {
				log.Println(err)
				continue
			}
		}
	}

	return res, nil
}
