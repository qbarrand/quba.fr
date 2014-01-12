/* OpenLayers */
$("#modal_travels").on('shown.bs.modal', function() {
	modal_travels_loadMap();
});

$("#modal_travels").on('hidden.bs.modal', function() {
	document.getElementById("map").innerHTML = "";
});


/* Load Hackgyver modal */
$("#modal_hackgyver").on('shown.bs.modal', function() {
	loadHackGyverModal();
});