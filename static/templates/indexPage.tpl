<!DOCTYPE html>
<html>
<head>
</script>
</head>

<body>
<h1>Paul & Lizzies Photo Page</h1>
<p>Click the links below to follow our adventures in Austrillia and New Zealand</p>
{{range $i, $v := .Items}}
    <a href="https://drive.google.com/uc?export=view&id={{$v.Id}}">{{$v.Title1}}</a>
{{end}}
</body>
</html>
