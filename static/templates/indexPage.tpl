<!DOCTYPE html>
<html>
<head>
</script>
</head>

<body>
<h1>Paul & Lizzies Photo Page</h1>
<p>Click the links below to follow our adventures in Austrillia and New Zealand</p>
<ul>
{{range $i, $v := .Items}}
    <li><a href="/display?id={{$v.Id}}">{{$v.Title1}}</a></li>
{{end}}
</ul>
</body>
</html>
