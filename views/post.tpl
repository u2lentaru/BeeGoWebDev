<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Post.Subj}}</title>
</head>
<body>
    <p>Viewing&nbsp{{.Post.Subj}}</p>
    <p>{{.Post.PostTime}}</p>
    <p>{{.Post.PostText}}</p>
    <p><a href="/edit/{{.Post.ID}}">Edit &nbsp{{.Post.Subj}}</a></p>
</body>
</html>
