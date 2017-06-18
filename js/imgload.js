var url = "/getimg/"; //url to load image from
var refreshInterval = 5000; //in ms
var drawDate = true; //draw date string
var img;

function init() {
    var canvas = document.getElementById("canvas");
    var context = canvas.getContext("2d");
    img = new Image();
    img.onload = function() {
        canvas.setAttribute("width", img.width)
        canvas.setAttribute("height", img.height)
        context.drawImage(this, 0, 0);
        if(drawDate) {
            var now = new Date();
            var text = now.toLocaleDateString() + " " + now.toLocaleTimeString();
            var maxWidth = 100;
            var x = img.width-10-maxWidth;
            var y = img.height-10;
            context.strokeStyle = 'black';
            context.lineWidth = 2;
            context.strokeText(text, x, y, maxWidth);
            context.fillStyle = 'white';
            context.fillText(text, x, y, maxWidth);
        }
    };
    refresh();
}
function refresh()
{
    img.src = url + "?t=" + new Date().getTime();
    setTimeout("refresh()",refreshInterval);
}
