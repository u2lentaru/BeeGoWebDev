package routers

import (
	"BeeGoWebDev/controllers"
	"context"
	"log"

	"github.com/astaxie/beego"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//_ "github.com/go-sql-driver/MySQL"
)

const (
	//dsn = "root:qw12345@tcp(localhost:3306)/myblog?charset=utf8"
	dbName         = "myblog"
	collectionName = "blogs"
)

func init() {
	//db, err := sql.Open("mysql", dsn)
	db, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	//if err := db.Ping(); err != nil {
	if err = db.Connect(context.Background()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected!")

	e := controllers.Explorer{
		Db:           db,
		DbName:       dbName,
		DbCollection: collectionName,
	}

	if err := e.Truncate(); err != nil {
		log.Fatal(err)
	}
	log.Print("truncated")

	if err := e.InsertDefault(); err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted default")

	beego.Router("/", &controllers.BlogController{
		Explorer: e,
	})

	/*beego.Router("/", &controllers.BlogController{
		Controller: beego.Controller{},
		Db:         db,
		//currBlog:   "1",
	})

	beego.Router("/post", &controllers.PostController{
		Controller: beego.Controller{},
		Db:         db,
		//currBlog:   "1",
	})*/
}
