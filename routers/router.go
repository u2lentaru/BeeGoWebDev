package routers

import (
	"BeeGoWebDev/controllers"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/astaxie/beego"

	_ "github.com/go-sql-driver/MySQL"
)

// Configuration - configuration structure
type Configuration struct {
	SiteName string
	DSN      string
	LogFile  string
}

// Conf app configuration
var Conf = Configuration{
	"MyBlog",
	"root:qw12345@tcp(localhost:3306)/myblog?charset=utf8",
	"myblog.log",
}

//const (
//	dsn = "root:qw12345@tcp(localhost:3306)/myblog?charset=utf8"
//)

func init() {
	log.Println(os.Getwd())
	f, err := ioutil.ReadFile("./routers/myblog.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	if err = json.Unmarshal(f, &Conf); err != nil {
		log.Println(err)
	}

	//db, err := sql.Open("mysql", dsn)
	db, err := sql.Open("mysql", Conf.DSN)

	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	lf, err := os.OpenFile(Conf.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer lf.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(lf)

	log.Println("db pinged!")

	beego.Router("/", &controllers.BlogController{
		Controller: beego.Controller{},
		Db:         db,
		//currBlog:   "1",
	})

	beego.Router("/post", &controllers.PostController{
		Controller: beego.Controller{},
		Db:         db,
		//currBlog:   "1",
	})
}
