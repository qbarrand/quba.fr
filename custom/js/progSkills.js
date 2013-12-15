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

        k = 0;
        skill = data[i].languages[j].skill;

        half = false;
        if(skill % 1 != 0) 
        {
            skill = Math.floor(skill);
            half = true;
        }

        for(;k < skill; k++)
          table += "<i class=\"fa fa-star\"> </i>";


        if(half == true)
        {
            table += "<i class=\"fa fa-star-half-o\"> </i>";
            k += 1;
        }

        for(; k < 5; k++)
          table += "<i class=\"fa fa-star-o\"> </i>";

        table += "</td>";
        table += "</tr>";

      }
    }

    table += "</table>";

    document.getElementById("tableProgSkills").innerHTML = table;
});


function loadProgSkills(callback) {  

    var xobj = new XMLHttpRequest();
    xobj.overrideMimeType("application/json");
    xobj.open('GET', 'custom/data/prog_skills.json', true);
    xobj.onreadystatechange = function () {
        if (xobj.readyState == 4 && xobj.status == "200") {
           
            // .open will NOT return a value but simply returns undefined in async mode so use a callback
            callback(xobj.responseText);
     
        }
    }
    xobj.send(null);
}
