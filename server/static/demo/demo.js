$(document).ready(function () {

$.fn.isInViewport = function() {
  var elementTop = $(this).offset().top;
  var elementBottom = elementTop + $(this).outerHeight();

  var viewportTop = $(window).scrollTop();
  var viewportBottom = viewportTop + $(window).height();

  return elementBottom > viewportTop && elementTop < viewportBottom;
};

$(window).on('resize scroll', function() {
  $('.color').each(function() {
      var activeColor = $(this).attr('id');
    if ($(this).isInViewport()) {
      $('#fixed-' + activeColor).addClass(activeColor + '-active');
    } else {
      $('#fixed-' + activeColor).removeClass(activeColor + '-active');
    }
  });
});
});
