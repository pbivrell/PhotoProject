$(document).ready(function () {
    $.fn.isBelowViewport = function() {
        var elementTop = $(this).offset().top;
        var viewportTop = $(window).scrollTop();
        var viewportBottom = viewportTop + $(window).height();
        return elementTop > viewportBottom;
    };

    var used = 0;

    var updateImages = function(){
        for(i=0; i < 4; i++){
            var nextCount = 0;
            $('.col'+i).each(function() {
                if ($(this).isBelowViewport()) {
                    nextCount += 1;
                }
            });
            if(nextCount < 1 && used < pageData.pictures.length) {
                var small = $("<img />").attr('src', 'http://localhost:8080/static/tiny.jpg').attr('class', 'col'+i).attr('id', 'tiny').appendTo('#'+i);
                $("<img />").attr('src', '/getImage?id='+pageData.pictures[used++]).attr('class', 'col'+i)
                    .on('load', {replace: small}, function(e) {
                        e.data.replace.replaceWith(this);
                    });
            }
        }
    };

    var updatePageData = function(){
        //$("#Nav").InnerHTML("<span>
        $("#nav").html(makeURLNav());
        $("#title1").text(pageData.title1);
        $("#title2").text(pageData.title2);
        $("#description").text(pageData.description);
    };

    var makeURLNav = function(){
        res = "<span>";
        url = "http://localhost:8080/";
        path = trim(window.location.pathname, "/");
        console.log(path);
        console.log(path.split("/"));
        if(path == ""){
            return "<span></span>"
        }
        res += '<a href="' + url + '"> ' + "All Photos" + ' </a>';
        spath = path.split("/");
        spath.forEach(function(entry, index){
            console.log(entry, index, spath.length);
            if(index < spath.length - 1){
                url += entry + "/";
                res += '/<a href="' + url + '"> ' + decodeURI(entry) + ' </a>';
            }
        });
        return res + "</span>";
    }

    var trim = function(data, c){
        temp = data;
        if(temp[0] == c){
            temp = temp.slice(1);
        }
        if(temp[temp.length-1] == c){
            temp = temp.slice(0, temp.length -1);
        }
        return temp;
    }
    
    $(window).on('resize scroll', updateImages);

    //$.getJSON("http://localhost:8080/get?path="+window.location.pathname.slice(1), function(data){
    $.getJSON("http://localhost:8080/json", function(data){
        pageData = data;
        updatePageData();
        updateImages();
        updateImages();
        updateImages();
        updateImages();
    });
});
