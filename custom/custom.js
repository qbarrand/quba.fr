function loadProgSkills(callback) {  

    var xobj = new XMLHttpRequest();
    xobj.overrideMimeType("application/json");
    xobj.open('GET', 'prog_skills.json', true);
    xobj.onreadystatechange = function () {
        if (xobj.readyState == 4 && xobj.status == "200") {
           
            // .open will NOT return a value but simply returns undefined in async mode so use a callback
            callback(xobj.responseText);
     
        }
    }
    xobj.send(null);
}


// Call to function with anonymous callback
loadProgSkills(function(response) {

  	data = JSON.parse(response);

    var table;
    table = "<table>";

    for(var i in data)
    {
      table += "<tr><td><strong>" + data[i].name + "</strong></td><td></td></tr>";

      for(var j in data[i].languages)
      {
        table += "<tr>";
        table += "<td>" + data[i].languages[j].name + "</td>"

        table += "<td>";

        for(var k = 0; k < 5; k++)
        {
          table += "<i class=\"fa fa-star\"> </i>";
        }

        table += "</td>";
        table += "</tr>";

      }
    }

    table += "</table>";

    document.getElementById("tableProgSkills").innerHTML = table;
});