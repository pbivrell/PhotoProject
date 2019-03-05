$(document).ready(function () {
    $.fn.isBelowViewport = function() {
        var elementTop = $(this).offset().top;
        var viewportTop = $(window).scrollTop();
        var viewportBottom = viewportTop + $(window).height();
        return elementTop > viewportBottom;
    };

    var max = $("#hidden").text() || 0;
    var used = 0;

    var updateImages = function(){
        for(i=0; i < 4; i++){
            var nextCount = 0;
            $('.col'+i).each(function() {
                if ($(this).isBelowViewport()) {
                    nextCount += 1;
                }
            });
            if(nextCount < 1 && used < max) {
                var small = $("<img />").attr('src', 'static/tiny.jpg').attr('class', 'col'+i).attr('id', 'tiny').appendTo('#'+i);
                $("<img />").attr('src', '/getImage?path='+window.location.pathname + "/" + used++).attr('class', 'col'+i)
                    .on('load', {replace: small}, function(e) {
                        e.data.replace.replaceWith(this);
                    });
            }
        }
    };
    
    $(window).on('resize scroll', updateImages);

    updateImages();
    updateImages();
    updateImages();
    updateImages();
    console.log("Done");

});
