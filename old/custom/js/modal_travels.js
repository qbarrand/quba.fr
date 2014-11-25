/* Get the data from the JSON file */
function modal_travels_loadMap() {  

    var xobj = new XMLHttpRequest();
    xobj.overrideMimeType("application/json");
    xobj.open('GET', 'custom/data/travels.json', true);
    xobj.onreadystatechange = function () {
        if (xobj.readyState == 4 && xobj.status == "200") {
            modal_travels_printMap(xobj.responseText);
        }
    }
    xobj.send(null);
}


function modal_travels_printMap(response) {
    data = JSON.parse(response);

    var map = L.map('map').setView([36.6, 13.0], 2);

    L.tileLayer('http://tile.openstreetmap.org/{z}/{x}/{y}.png', {
        maxZoom: 18,
        attribution: 'Map data &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>'
    }).addTo(map);

    for(var i = 0; i < data.length; i++)
    {
        var desc = "<b>" + data[i].place + "</b>";
        
        for(var j = 0; j < data[i].trips.length; j++) {
            desc += "<br />" + data[i].trips[j].when;

            if(data[i].trips[j].why != null)
            {
                desc += " - " + data[i].trips[j].why;
            }
        }

        L.marker([data[i].lat, data[i].lon]).addTo(map)
            .bindPopup(desc)
    }

    var popup = L.popup();
}