<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Blog</title>
</head>
<body>
    <p>{{.Title}}</p>
    <ul>
        {{range .Blog}}
        <li>{{.PostTime}}&nbsp:&nbsp<a href="/post/{{.ID}}">{{.Subj}}</a>&nbsp&nbsp&nbsp<a href="/del/{{.ID}}">Del</a></li>
        {{end}}
    </ul>
    <div><a href="/new">New&nbspPost</a></div>
</body>
</html>