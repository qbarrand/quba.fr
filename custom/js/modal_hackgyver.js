function loadHackGyverModal()
{
	document.getElementById("modalHackGyverBody").innerHTML += ' \
	           <object id="player" width="640" height="360" classid="clsid:d27cdb6e-ae6d-11cf-96b8-444553540000" codebase="http://download.macromedia.com/pub/shockwave/cabs/flash/swflash.cab#version=6,0,40,0" name="player"> \
                <param name="allowfullscreen" value="true" /> \
                <param name="allowscriptaccess" value="always" /> \
                <param name="flashvars" value="config=http%3A%2F%2Fwww.besancon.tv%2Fgetconfig.php%3Fid_prod%3D1169%26id_app%3D149%26embed%3D1" /> \
                <param name="src" value="http://www.besancon.tv/mediaplayer5.swf" /> \
                <embed id="player" width="640" height="360" type="application/x-shockwave-flash" src="http://www.besancon.tv/mediaplayer5.swf" allowfullscreen="true" allowscriptaccess="always" flashvars="config=http%3A%2F%2Fwww.besancon.tv%2Fgetconfig.php%3Fid_prod%3D1169%26id_app%3D149%26embed%3D1" name="player" /> \
              </object>';
}