<!DOCTYPE html>
<html>
<head>
<link rel="stylesheet" href="css/responsive-grid.css">
<link rel="stylesheet" href="css/progressive-image.min.css">
<link rel="stylesheet" href="https://www.w3schools.com/w3css/4/w3.css">
<style>
.mySlides {display:none;}
</style>

<script
src="https://code.jquery.com/jquery-3.3.1.min.js"
integrity="sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8QtmkMRdAu8="
crossorigin="anonymous">
</script>
</head>

</style>



<body>

<div class="header">
    <h1>{{.Title1}}<b class="bar">|</b>{{.Title2}}</h1>
    <p>{{.Description}}</p>
</div>

<script src="js/progressive-image.min.js"></script>

<div class="column">
{{$colContainer := create (len .BigImages)}}
{{range $i, $v := .BigImages}}
    {{if newColumn $colContainer $i}} 
<div/>
<div class="column">
    {{end}}
    <a href="{{$v}}" class="primary progressive replace blur" id="slideshow" onclick="on({{$i}});return false;">
        <img src="{{index .TinyImages $i}}" class="preview" alt=""/>
    </a>
{{end}}
</div>

<div id="overlay" onclick="off()">
    <div class="w3-content w3-display-container">
<button class="w3-button w3-black w3-display-left" onclick="plusDivs(-1,event)">&#10094;</button>
<button class="w3-button w3-black w3-display-right" onclick="plusDivs(1,event)">&#10095;</button>
{{range $i, $v := .BigImages}}<img class="mySlides" src="{{$v}}" style="heigh:100%">
{{end}}
</div>
</div>
<script src="js/slideshow.js"></script>
</body>
</html>