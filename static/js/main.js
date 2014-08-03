var goodnumber = {};

(function(){
    goodnumber.postAp = function() {
        var ap = parseInt(document.getElementById("ap").value);
        var data = { "ap": ap };
        var xhr = new XMLHttpRequest();
        xhr.open("POST", "/", true)
        xhr.onreadystatechange = function() {
            var result = document.getElementById("result");
            if ( xhr.readyState === 4 && xhr.status === 200 ) {
                var actionData = JSON.parse(xhr.responseText);
                table = goodnumber.createActionTable(actionData);
                result.appendChild(table);
            } else if ( xhr.readyState === 4 && xhr.status === 0 ) {
                result.innerHTML = "not http";
            }
        }
        xhr.send(JSON.stringify(data));
    }

    goodnumber.createActionTable = function(data) {
        var table = document.createElement('table');
        var headerTr = document.createElement('tr');
        var headerAction = document.createElement('th');
        headerAction.setAttribute('class', 'header');
        headerAction.setAttribute('colspan', 2);
        headerAction.innerHTML = "Action";
        headerTr.appendChild(headerAction);
        table.appendChild(headerTr);

        var tr = document.createElement('tr');
        var dataAction = document.createElement('td');
        dataAction.setAttribute('class', 'action');
        var dataCount = document.createElement('td');
        dataCount.setAttribute('class', 'count');
        tr.appendChild(dataAction);
        tr.appendChild(dataCount)

        for (k in data) {
            var row = tr.cloneNode(true);
            row.getElementsByClassName("action")[0].innerHTML = k;
            row.getElementsByClassName("count")[0].innerHTML = data[k];
            table.appendChild(row);
        }
        return table
    }
})()
