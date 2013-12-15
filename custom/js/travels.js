// Call to function with anonymous callback
loadTravelsMap(function(response) {

    data = JSON.parse(response);

    map = new OpenLayers.Map("modalTravelsBody");
    map.addLayer(new OpenLayers.Layer.OSM());
 
    var lonLat = new OpenLayers.LonLat( -0.1279688 ,51.5077286 )
          .transform(
            new OpenLayers.Projection("EPSG:4326"), // transform from WGS 1984
            map.getProjectionObject() // to Spherical Mercator Projection
          );
 
    var zoom=0;
 
    var markers = new OpenLayers.Layer.Markers( "Markers" );
    map.addLayer(markers);
 
    for(var i in data)
    {
        markers.addMarker(new OpenLayers.Marker(new OpenLayers.LonLat(data[i].lon, data[i].lat)));
    }

    markers.addMarker(new OpenLayers.Marker(lonLat));
 
    map.setCenter (lonLat, zoom);
});


function loadTravelsMap(callback) {  

    var xobj = new XMLHttpRequest();
    xobj.overrideMimeType("application/json");
    xobj.open('GET', 'custom/data/travels.json', true);
    xobj.onreadystatechange = function () {
        if (xobj.readyState == 4 && xobj.status == "200") {
           
            // .open will NOT return a value but simply returns undefined in async mode so use a callback
            callback(xobj.responseText);
     
        }
    }
    xobj.send(null);
}
