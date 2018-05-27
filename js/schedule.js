postMessage("Waiting for server");
sendPOST();

function sendPOST() {
    var request = new XMLHttpRequest();
    request.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            postMessage(this.responseText);
        } else {
            postMessage("The Server is down");
        }
    }
    request.open("POST", "http://deckedhost.ns0.it:8080/get_schedule", true);
    request.send();
    setTimeout("sendPOST()", 10000);
}