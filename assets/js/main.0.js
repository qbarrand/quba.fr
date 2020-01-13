/*
    Eventually by HTML5 UP
    html5up.net | @ajlkn
    Free for personal and commercial use under the CCA 3.0 license (html5up.net/license)
*/

let webpSupported = false;
let webpPromiseRan = false;

const p = new Promise((resolve, _) => {
    let img = new Image()
    img.onload = () => resolve(img.width > 0 && img.height > 0)
    img.onerror = () => resolve(false)
    img.src = 'data:image/webp;base64,UklGRiIAAABXRUJQVlA4IBYAAAAwAQCdASoBAAEADsD+JaQAA3AAAAAA'
});

class Constraint {
    constructor(direction, n) {
        this.direction = direction
        this.n = n
    }

    requiresUpdate(other) {
        return other.direction != direction || other.n > this.n
    }

    toQueryString() {
        if (this.n == Infinity) {
            return ''
        }

        return `${this.direction}=${this.n}`
    }
}

const queries = {
    // Portrait
    '(orientation: portrait) and (max-height: 480px)':                           new Constraint('height', 480),
    '(orientation: portrait) and (min-height: 481px) and (max-height: 736px)':   new Constraint('height', 736),
    '(orientation: portrait) and (min-height: 737px) and (max-height: 980px)':   new Constraint('height', 980),
    '(orientation: portrait) and (min-height: 981px) and (max-height: 1280px)':  new Constraint('height', 1280),
    '(orientation: portrait) and (min-height: 1281px) and (max-height: 1690px)': new Constraint('height', 1690),
    '(orientation: portrait) and (min-height: 1691px)':                          new Constraint('height', Infinity),

    // Landscape
    '(orientation: landscape) and (max-width: 480px)':                          new Constraint('width', 480),
    '(orientation: landscape) and (min-width: 481px) and (max-width: 736px)':   new Constraint('width', 736),
    '(orientation: landscape) and (min-width: 737px) and (max-width: 980px)':   new Constraint('width', 980),
    '(orientation: landscape) and (min-width: 981px) and (max-width: 1280px)':  new Constraint('width', 1280),
    '(orientation: landscape) and (min-width: 1281px) and (max-width: 1690px)': new Constraint('width', 1690),
    '(orientation: landscape) and (min-width: 1691px) and (max-width: 1920px)': new Constraint('width', 1920),
    '(orientation: landscape) and (min-width: 1921px) and (max-width: 2880px)': new Constraint('width', 2880),
    '(orientation: landscape) and (min-width: 2881px)':                         new Constraint('width', Infinity),
}

class ImageData {
    constructor(bloblUrl, mainColor) {
        this.bloblUrl = bloblUrl;
        this.mainColor = mainColor;
    }
}

class ImageMetadata {
    constructor(location, date) {
        this.location = location;
        this.date = date;
    }
}

const allImages = {
    'shenzhen_1.jpg': 		new ImageMetadata('Shenzhen, China', 'August 2014'),
    'geneva_1.jpg': 		new ImageMetadata('Geneva, Switzerland', 'June 2016'),
    'newyork_2.jpg': 		new ImageMetadata('New York, USA', 'August 2015'),
    'thun_1.jpg': 			new ImageMetadata('Thun, Switzerland', 'May 2016'),
    'montreux_1.jpg': 		new ImageMetadata('Montreux, Switzerland', 'October 2016',),
    'dubai_1.jpg': 			new ImageMetadata('Dubai, UAE', 'June 2017'),
    'new_delhi_1.jpg': 		new ImageMetadata('New Delhi, India', 'June 2017'),
    'kyoto_1.jpg': 			new ImageMetadata('Kyoto, Japan', 'October 2017'),
    'singapore_1.jpg': 		new ImageMetadata('Singapore', 'January 2019'),
    'nuggets_point_1.jpg': 	new ImageMetadata('Nuggets Point, New Zealand', 'January 2019'),
    'whaikiti_beach_1.jpg': new ImageMetadata('Whaikiti Beach, New Zealand', 'January 2019'),
    'malibu_1.jpg': 		new ImageMetadata('Malibu, USA', 'March 2019'),
    'lhc_1.jpg': 			new ImageMetadata('LHC, France / Switzerland', 'August 2019'),
    'dents_du_midi_1.jpg': 	new ImageMetadata('Dents du Midi, Switzerland', 'January 2020')
};

const cache = new Map();

let currentConstraint = null;

async function printRandomBackground(parent, imageName, constraint) {
    const metadata = allImages[imageName];

    let url = `images/${imageName}?${constraint.toQueryString()}`

    if (!webpPromiseRan) {
        console.log('Awaiting the webp promise');
        webpSupported = await p;
        webpPromiseRan = true;
    }

    let accept = ['image/jpeg'];

    console.log(`webp supported: ${webpSupported}`)

    if (webpSupported) {
        accept = ['image/webp'].concat(accept);
    }

    let imageData = cache.get(url);

    if (imageData != null) {
        console.log('Using cache')
    } else {
        console.log('Fetching ' + url);

        const response = await fetch(url, {
            headers: new Headers({'Accept': accept.join(',')})
        });

        const blob = await response.blob();

        const objUrl = URL.createObjectURL(blob);
        const mainColor = response.headers.get('X-Main-Color');

        imageData = new ImageData(objUrl, mainColor);

        cache.set(url, imageData);
    }

    const div = document.createElement('div');
    div.style.backgroundImage = `url("${imageData.bloblUrl}")`;
    div.style.backgroundPosition = 'center';
    parent.appendChild(div);

    div.classList.add('visible');
    div.classList.add('top');

    document.querySelector('#where').innerHTML = metadata.location;
    document.querySelector('#when').innerHTML = metadata.date;

    document.querySelector('meta[name=theme-color]').setAttribute('content', imageData.mainColor);

    currentConstraint = constraint;
}

(function() {

    "use strict";

    var	$body = document.querySelector('body');

    // Methods/polyfills.

        // classList | (c) @remy | github.com/remy/polyfills | rem.mit-license.org
            !function(){function t(t){this.el=t;for(var n=t.className.replace(/^\s+|\s+$/g,"").split(/\s+/),i=0;i<n.length;i++)e.call(this,n[i])}function n(t,n,i){Object.defineProperty?Object.defineProperty(t,n,{get:i}):t.__defineGetter__(n,i)}if(!("undefined"==typeof window.Element||"classList"in document.documentElement)){var i=Array.prototype,e=i.push,s=i.splice,o=i.join;t.prototype={add:function(t){this.contains(t)||(e.call(this,t),this.el.className=this.toString())},contains:function(t){return-1!=this.el.className.indexOf(t)},item:function(t){return this[t]||null},remove:function(t){if(this.contains(t)){for(var n=0;n<this.length&&this[n]!=t;n++);s.call(this,n,1),this.el.className=this.toString()}},toString:function(){return o.call(this," ")},toggle:function(t){return this.contains(t)?this.remove(t):this.add(t),this.contains(t)}},window.DOMTokenList=t,n(Element.prototype,"classList",function(){return new t(this)})}}();

        // canUse
            // window.canUse=function(p){if(!window._canUse)window._canUse=document.createElement("div");var e=window._canUse.style,up=p.charAt(0).toUpperCase()+p.slice(1);return p in e||"Moz"+up in e||"Webkit"+up in e||"O"+up in e||"ms"+up in e};

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

    let selected = new URLSearchParams(window.location.search).get('img')

    if (selected == null) {
        const keys = Object.keys(allImages) // .filter(e => e != currentFile);
        selected = keys[Math.floor(Math.random()*keys.length)];
    }

    // Register all media query listeners
    for (let [q, c] of Object.entries(queries)) {
        const m = window.matchMedia(q)

        if (m.matches) {
            printRandomBackground($wrapper, selected, c);
        }

        m.addListener(e => {
            if (e.matches && (currentConstraint == undefined || currentConstraint.requiresUpdate(c))) {
                printRandomBackground($wrapper, selected, c);
            }
        });
    }
})();
