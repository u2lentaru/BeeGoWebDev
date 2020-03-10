package routers

import (
	"BeeGoWebDev/controllers"
	"database/sql"
	"log"

	"github.com/astaxie/beego"

	_ "github.com/go-sql-driver/MySQL"
)

const (
	dsn = "root:qw12345@tcp(localhost:3306)/myblog?charset=utf8"
)

func init() {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//serv := Server{database: db, currBlog: "1"}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("db pinged!")

	beego.Router("/", &controllers.BlogController{
		Controller: beego.Controller{},
		Db:         db,
	})
}
