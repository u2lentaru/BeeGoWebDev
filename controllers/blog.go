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
	//Db       *sql.DB
	//currBlog string
	Explorer Explorer
}

// Explorer struct
type Explorer struct {
	Db           *mongo.Client
	DbName       string
	DbCollection string
}

// Truncate - truncate database
func (e Explorer) Truncate() error {
	c := e.Db.Database(e.DbName).Collection(e.DbCollection)
	_, err := c.DeleteMany(context.Background(), bson.D{})

	return err
}

// InsertDefault - create posts
func (e Explorer) InsertDefault() error {
	for _, post := range createPosts() {
		if err := e.addPost(post); err != nil {
			return err
		}
	}

	return nil
}

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
	}
}

func (e Explorer) addPost(post models.TPost) error {
	c := e.Db.Database(e.DbName).Collection(e.DbCollection)
	_, err := c.InsertOne(context.Background(), post)

	return err
}

// Get func
func (c *BlogController) Get() {
	//c.currBlog = "1"
	//blog, err := getBlog(c.Db, c.currBlog)
	blog, err := c.Explorer.getBlog()
	if err != nil {
		log.Println(err)
		return
	}

	c.Data["Blog"] = blog
	c.TplName = "blogs.tpl"
}

//func getBlog(db *sql.DB, id string) (models.TBlog, error) {
func (e Explorer) getBlog() ([]models.TPost, error) {
	blog := []models.TPost{}

	c := e.Db.Database(e.DbName).Collection(e.DbCollection)

	cur, err := c.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "Find")
	}

	//res := make([]models.Post, 0, 1)

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

// Post func
func (c *BlogController) Post() {
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

	//log.Printf("post %v", post)

	//if err := createPost(c.Db, c.currBlog, post); err != nil {
	if err := c.Explorer.addPost(post); err != nil {
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

/*func createPost(db *sql.DB, currBlog string, post models.TPost) error {
	_, err := db.Exec("insert into myblog.posts (blogid,subj,posttime,posttext) values (?,?,?,?)",
		currBlog, post.Subj, post.PostTime, post.PostText)

	return err
}*/

type putRequest struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

/*
	curl.exe -vX PUT -H "Content-Type: application/json"  -d"@update.json" http://localhost:8080?id=1
*/

// Put func
func (c *BlogController) Put() {
	id := c.Ctx.Request.URL.Query().Get("id")

	if len(id) == 0 {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte("empty id"))
		return
	}

	//resp := new(putRequest)
	resp := new(postRequest)
	//resp := new(models.TPost)

	if err := readAndUnmarshall(resp, c.Ctx.Request.Body); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
		return
	}

	/*blog := models.TBlog{
		Name:  resp.Name,
		Title: resp.Title,
	}*/

	post := models.TPost{
		Subj:     resp.Subj,
		PostTime: resp.PostTime,
		PostText: resp.PostText,
	}

	//log.Printf("id %v", id)
	//log.Printf("post %v", post)
	//if err := updateBlog(c.Db, id, blog.Name, blog.Title); err != nil {
	if err := c.Explorer.editPost(&post, id); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
	}

	c.Ctx.ResponseWriter.WriteHeader(200)
	_, _ = c.Ctx.ResponseWriter.Write([]byte(`SUCCESS`))
}

/*func updateBlog(db *sql.DB, id, name, title string) error {
	if len(name) == 0 && len(title) == 0 {
		return nil
	}

	_, err := db.Exec("UPDATE myblog.blogs SET name=?, title=? WHERE id=?",
		name, title, id)

	return err

}*/

func (e Explorer) editPost(post *models.TPost, id string) error {
	filter := bson.D{{Key: "id", Value: id}}

	update := createUpdates(*post)

	c := e.Db.Database(e.DbName).Collection(e.DbCollection)
	_, err := c.UpdateOne(context.Background(), filter, update)

	return err
}

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

// Delete func
func (c *BlogController) Delete() {
	id := c.Ctx.Request.URL.Query().Get("id")

	if len(id) == 0 {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte("empty id"))
		return
	}

	//err := deletePost(c.Db, id)
	err := c.Explorer.deletePost(id)

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
	}

	c.Ctx.ResponseWriter.WriteHeader(200)
	_, _ = c.Ctx.ResponseWriter.Write([]byte(`SUCCESS`))

}

/*func deletePost(db *sql.DB, id string) error {
	if _, err := db.Exec("DELETE FROM myblog.posts WHERE `id`=?", id); err != nil {
		return err
	}

	return nil
}
*/

func (e Explorer) deletePost(id string) error {
	filter := bson.D{{Key: "id", Value: id}}

	c := e.Db.Database(e.DbName).Collection(e.DbCollection)
	_, err := c.DeleteOne(context.Background(), filter)

	return err
}
