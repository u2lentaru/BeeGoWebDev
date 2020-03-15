package controllers

import (
	"BeeGoWebDev/models"
	"context"

	"github.com/astaxie/beego"
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
/*func (c *BlogController) Get() {
	//c.currBlog = "1"
	//blog, err := getBlog(c.Db, c.currBlog)
	blog, err := getBlog(c.Explorer.Db)
	if err != nil {
		log.Println(err)
		return
	}

	c.Data["Blog"] = blog
	c.TplName = "blogs.tpl"
}*/

/*//func getBlog(db *sql.DB, id string) (models.TBlog, error) {
func getBlog(db *mongo.Client) (models.TBlog, error) {
	blog := models.TBlog{}
	//if err := db.Ping(); err != nil {
	//	log.Println(err)
	//}

	row := db.QueryRow("select * from myblog.blogs where blogs.id = ?", id)
	err := row.Scan(&blog.ID, &blog.Name, &blog.Title)
	if err != nil {
		return models.TBlog{}, err
	}

	rows, err := db.Query("select * from posts where blogid = ?", id)
	if err != nil {
		return models.TBlog{}, err
	}
	defer rows.Close()

	for rows.Next() {
		post := models.TPost{}

		err := rows.Scan(&post.ID, new(int), &post.Subj, &post.PostTime, &post.PostText)
		if err != nil {
			log.Println(err)
			continue
		}

		blog.PostList = append(blog.PostList, post)
	}

	return blog, nil
}*/

type postRequest struct {
	Subj     string `json:"subj"`
	PostTime string `json:"posttime"`
	PostText string `json:"posttext"`
}

/*
	curl.exe -vX POST -H "Content-Type: application/json"  -d "@data.json" http://localhost:8080/
*/

// Post func
/*func (c *BlogController) Post() {
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

	//log.Printf("post %v", post)

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
/*func (c *BlogController) Put() {
	id := c.Ctx.Request.URL.Query().Get("id")

	if len(id) == 0 {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte("empty id"))
		return
	}

	resp := new(putRequest)

	if err := readAndUnmarshall(resp, c.Ctx.Request.Body); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
		return
	}

	blog := models.TBlog{
		Name:  resp.Name,
		Title: resp.Title,
	}

	if err := updateBlog(c.Db, id, blog.Name, blog.Title); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
	}

	c.Ctx.ResponseWriter.WriteHeader(200)
	_, _ = c.Ctx.ResponseWriter.Write([]byte(`SUCCESS`))
}

func updateBlog(db *sql.DB, id, name, title string) error {
	if len(name) == 0 && len(title) == 0 {
		return nil
	}

	_, err := db.Exec("UPDATE myblog.blogs SET name=?, title=? WHERE id=?",
		name, title, id)

	return err

}*/

/*
	curl.exe -vX DELETE  http://localhost:8080?id=42
*/

// Delete func
/*func (c *BlogController) Delete() {
	id := c.Ctx.Request.URL.Query().Get("id")

	err := deletePost(c.Db, id)

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		_, _ = c.Ctx.ResponseWriter.Write([]byte(err.Error()))
	}

	c.Ctx.ResponseWriter.WriteHeader(200)
	_, _ = c.Ctx.ResponseWriter.Write([]byte(`SUCCESS`))

}

func deletePost(db *sql.DB, id string) error {
	if _, err := db.Exec("DELETE FROM myblog.posts WHERE `id`=?", id); err != nil {
		return err
	}

	return nil
}
*/
