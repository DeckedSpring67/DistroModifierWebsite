postMessage("Worker working");

function sendPOST() {
    var request = new XMLHttpRequest();
    request.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            postMessage(this.responseText);
        }
    }
    request.open("POST", "http://deckedhost.ns0.it:8080/current_percentage", true);
    request.send();
    setTimeout("sendPOST()", 500);
}
sendPOST();