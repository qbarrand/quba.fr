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


/* Actually prints the map */
function modal_travels_printMap(response)
{
    data = JSON.parse(response);

    map = new OpenLayers.Map("map");
    var mapnik = new OpenLayers.Layer.OSM();
    map.addLayer(mapnik);

    var zoom = 0;

    var markers = new OpenLayers.Layer.Markers( "Markers" );
    map.addLayer(markers);
    
    main_marker = new OpenLayers.Marker(
            new OpenLayers.LonLat(data[0].lon, data[0].lat).transform(
                new OpenLayers.Projection("EPSG:4326"),
                new OpenLayers.Projection("EPSG:900913")
    ));

    markers.addMarker(main_marker);

    for(var i = 1; i < data.length; i++)
    {
        markers.addMarker(new OpenLayers.Marker(
            new OpenLayers.LonLat(data[i].lon, data[i].lat).transform(
                new OpenLayers.Projection("EPSG:4326"),
                new OpenLayers.Projection("EPSG:900913")
        )));
    }

    map.setCenter(main_marker, zoom);
}