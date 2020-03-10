package models

// TBlog - blog struct
type TBlog struct {
	ID       string
	Name     string
	Title    string
	PostList []TPost
}

// TPost - post struct
type TPost struct {
	ID       string
	Subj     string
	PostTime string
	Text     string
}
