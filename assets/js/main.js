/*
	Eventually by HTML5 UP
	html5up.net | @ajlkn
	Free for personal and commercial use under the CCA 3.0 license (html5up.net/license)
*/

const allImages = {
	'shenzhen_1.jpg': 		{location: 'Shenzhen, China', 				date: 'August 2014'},
	'geneva_1.jpg': 		{location: 'Geneva, Switzerland', 			date: 'June 2016'},
	'newyork_2.jpg': 		{location: 'New York, USA',					date: 'August 2015'},
	'thun_1.jpg': 			{location: 'Thun, Switzerland',				date: 'May 2016'},
	'montreux_1.jpg': 		{location: 'Montreux, Switzerland',			date: 'October 2016',},
	'dubai_1.jpg': 			{location: 'Dubai, UAE',					date: 'June 2017'},
	'kyoto_1.jpg': 			{location: 'Kyoto, Japan',					date: 'October 2017'},
	'nuggets_point_1.jpg': 	{location: 'Nuggets Point, New Zealand',	date: 'January 2019'},
	'whaikiti_beach_1.jpg': {location: 'Whaikiti Beach, New Zealand',	date: 'January 2019'},
	'lhc_1.jpg': 			{location: 'LHC, France / Switzerland',		date: 'August 2019'}
};

let currentDiv = null;
let currentFile = '';

async function printRandomBackground(wrapper) {
	const keys = Object.keys(allImages).filter(e => e != currentFile);
	currentFile = keys[Math.floor(Math.random()*keys.length)];

	console.log(currentFile);

	const image = allImages[currentFile];

	if (image.div !== undefined) {
		newDiv = image.div;
	} else {
		console.log('Fetching ' + currentFile);

		// const width = window.innerWidth;

		const response = await fetch(`images/${currentFile}?format=webp&width=${window.innerWidth}`);
		const blob = await response.blob();
		const url = URL.createObjectURL(blob);

		image.div = document.createElement('div');
		image.div.style.backgroundImage = `url("${url}")`;
		image.div.style.backgroundPosition = 'center';
		wrapper.appendChild(image.div);
	}

	if (currentDiv != null) {
		currentDiv.classList.remove('top');
	}

	image.div.classList.add('visible');
	image.div.classList.add('top');

	document.querySelector('#where').innerHTML = image.location;
	document.querySelector('#when').innerHTML = image.date;

	const oldDiv = currentDiv;

	if (oldDiv != null) {
		window.setTimeout(() => oldDiv.classList.remove('visible'), 500);
	}
	currentDiv = image.div;

	// TODO re-enable switchDiv
	// document.querySelector('#switch > i').classList.remove('fa-spin');
}

(function() {

	"use strict";

	var	$body = document.querySelector('body');

	// Methods/polyfills.

		// classList | (c) @remy | github.com/remy/polyfills | rem.mit-license.org
			!function(){function t(t){this.el=t;for(var n=t.className.replace(/^\s+|\s+$/g,"").split(/\s+/),i=0;i<n.length;i++)e.call(this,n[i])}function n(t,n,i){Object.defineProperty?Object.defineProperty(t,n,{get:i}):t.__defineGetter__(n,i)}if(!("undefined"==typeof window.Element||"classList"in document.documentElement)){var i=Array.prototype,e=i.push,s=i.splice,o=i.join;t.prototype={add:function(t){this.contains(t)||(e.call(this,t),this.el.className=this.toString())},contains:function(t){return-1!=this.el.className.indexOf(t)},item:function(t){return this[t]||null},remove:function(t){if(this.contains(t)){for(var n=0;n<this.length&&this[n]!=t;n++);s.call(this,n,1),this.el.className=this.toString()}},toString:function(){return o.call(this," ")},toggle:function(t){return this.contains(t)?this.remove(t):this.add(t),this.contains(t)}},window.DOMTokenList=t,n(Element.prototype,"classList",function(){return new t(this)})}}();

		// canUse
			window.canUse=function(p){if(!window._canUse)window._canUse=document.createElement("div");var e=window._canUse.style,up=p.charAt(0).toUpperCase()+p.slice(1);return p in e||"Moz"+up in e||"Webkit"+up in e||"O"+up in e||"ms"+up in e};

		// window.addEventListener
			(function(){if("addEventListener"in window)return;window.addEventListener=function(type,f){window.attachEvent("on"+type,f)}})();

	// Play initial animations on page load.
		window.addEventListener('load', function() {
			window.setTimeout(function() {
				$body.classList.remove('is-preload');
			}, 100);
		});

	const $wrapper = document.createElement('div');
	$wrapper.id = 'bg';
	document.querySelector('body').appendChild($wrapper);

	printRandomBackground($wrapper);

	// TODO re-enable switchDiv
	// const switchDiv = document.querySelector('#switch');

	// switchDiv.addEventListener('click', e => {
	// 	document.querySelector('#switch > i').classList.add('fa-spin');
	// 	printRandomBackground($wrapper);
	// });

	// Slideshow Background.
		// (function() {
			// document.querySelector('#switch').addEventListener('click', e => {
			// 	document.querySelector('#switch > i').classList.add('fa-spin');
			// });

			// // Vars.
			// 	var	pos = 0, lastPos = 0,
			// 		$wrapper, $bgs = [], $bg,
			// 		k, v;

			// // Create BG wrapper, BGs.
			// 	$wrapper = document.createElement('div');
			// 		$wrapper.id = 'bg';
			// 		$body.appendChild($wrapper);

			// 	for (k in settings.images) {

			// 		// Create BG.
			// 			$bg = document.('div');
			// 				$bg.style.backgroundImage = 'url("' + k + '")';
			// 				$bg.style.backgroundPosition = settings.images[k];
			// 				$wrapper.appendChild($bg);

			// 		// Add it to array.
			// 			$bgs.push($bg);

			// 	}

			// // Main loop.
			// 	$bgs[pos].classList.add('visible');
			// 	$bgs[pos].classList.add('top');

			// 	// Bail if we only have a single BG or the client doesn't support transitions.
			// 		if ($bgs.length == 1
			// 		||	!canUse('transition'))
			// 			return;

			// 	window.setInterval(function() {

			// 		lastPos = pos;
			// 		pos++;

			// 		// Wrap to beginning if necessary.
			// 			if (pos >= $bgs.length)
			// 				pos = 0;

			// 		// Swap top images.
			// 			$bgs[lastPos].classList.remove('top');
			// 			$bgs[pos].classList.add('visible');
			// 			$bgs[pos].classList.add('top');

			// 			console.log($bgs[pos])

			// 		// Hide last image after a short delay.
			// 			window.setTimeout(function() {
			// 				$bgs[lastPos].classList.remove('visible');
			// 			}, settings.delay / 2);

			// 	}, settings.delay);

		// })();
})();