function on(n) {
    slideIndex = n;
    document.getElementById("overlay").style.display = "block";
}

function off() {
    document.getElementById("overlay").style.display = "none";
}


window.onkeyup = function(e) {
    var key = e.keyCode ? e.keyCode : e.which;

    if (key == 27) {
        off()
    }else if (key == 37) {
        plusDivs(-1, e)
    }else if (key == 39) {
        plusDivs(1, e)
    }
}

var slideIndex = 1;
showDivs(slideIndex);

function plusDivs(n,e) {
    e.stopPropagation()
        showDivs(slideIndex += n);
}

function showDivs(n) {
    var i;
    var x = document.getElementsByClassName("mySlides");
    if (n > x.length) {slideIndex = 1}    
    if (n < 1) {slideIndex = x.length}
    for (i = 0; i < x.length; i++) {
        x[i].style.display = "none";  
    }
    x[slideIndex-1].style.display = "block";  
}
