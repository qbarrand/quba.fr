/* Get the data from JSON */
function modal_skills_loadSkills() {  

    var xobj = new XMLHttpRequest();
    xobj.overrideMimeType("application/json");
    xobj.open('GET', 'custom/data/skills.json', true);
    xobj.onreadystatechange = function () {
        if (xobj.readyState == 4 && xobj.status == "200") {
            modal_skills_printTable(xobj.responseText);
        }
    }
    xobj.send(null);
}


/* Actually print the table into the modal */
function modal_skills_printTable(response) {

    data = JSON.parse(response);

    // console.log(data)

    var table;
    table = "<center><table>";

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

        table += "<tr><td>&nbsp;</td><td></td></tr>";

    }

    table += "</table></center>";

    document.getElementById("tableProgSkills").innerHTML = table;
}
