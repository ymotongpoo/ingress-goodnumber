function postAp() {
    var ap = parseInt(document.getElementById("ap").value);
    var data = { "ap": ap };
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/", true)
    xhr.onreadystatechange = function() {
        var result = document.getElementById("result");
        if ( xhr.readyState === 4 && xhr.status === 200 ) {
            var action = JSON.parse(xhr.responseText);
            for (var x in action) {
                result.innerHTML += x + " : " + action[x] + "<br/>";
            }
        } else if ( xhr.readyState === 4 && xhr.status === 0 ) {
            result.innerHTML = "not http";
        }
    }
    xhr.send(JSON.stringify(data));
}
