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
    
    $(window).on('resize scroll', updateImages);

    $.getJSON("http://localhost:8080/get?path="+window.location.pathname.slice(1), function(data){
        pageData = data;
        $("#Title1").text(pageData.title1);
        $("#Title2").text(pageData.title2);
        $("#Description").text(pageData.description);
        updateImages();
        updateImages();
        updateImages();
        updateImages();
    });
});
