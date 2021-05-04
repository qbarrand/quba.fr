/*
    Eventually by HTML5 UP
    html5up.net | @ajlkn
    Free for personal and commercial use under the CCA 3.0 license (html5up.net/license)
*/

"use strict";

const p = new Promise((resolve, _) => {
    const img = new Image()
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
        return other.direction !== this.direction || other.n > this.n
    }

    toQueryString() {
        if (this.n === Infinity) {
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
    constructor(bloblUrl, mainColor, location, date) {
        this.bloblUrl = bloblUrl;
        this.mainColor = mainColor;
        this.location = location;
        this.date = date;
    }
}

const cache = new Map();

let currentConstraint = null;

async function printRandomBackground(parent, imageName, constraint) {
    const url = `images/${imageName}?${constraint.toQueryString()}`
    const accept = ['image/jpeg'];

    if (await p) {
        accept.unshift('image/webp');
    }

    let imageData = cache.get(url);

    if (imageData != null) {
        console.debug('Using cache')
    } else {
        console.debug('Fetching ' + url);

        const response = await fetch(url, {
            headers: new Headers({'Accept': accept.join(', ')})
        })

        const blob = await response.blob();

        const unixSecs = response.headers.get('X-Quba-Date')
        const date = new Date(unixSecs * 1000)

        imageData = new ImageData(
            URL.createObjectURL(blob),
            response.headers.get('X-Quba-Main-Color'),
            response.headers.get('X-Quba-Location'),
            date.toLocaleDateString('en-EN', {year: 'numeric', month: 'long'})
        );

        cache.set(url, imageData);
    }

    const div = document.createElement('div');
    div.style.backgroundImage = `url("${imageData.bloblUrl}")`;
    div.style.backgroundPosition = 'center';
    parent.appendChild(div);

    div.classList.add('visible');
    div.classList.add('top');

    document.querySelector('#where').innerHTML = imageData.location;
    document.querySelector('#when').innerHTML = imageData.date;

    document.querySelector('meta[name=theme-color]').setAttribute('content', imageData.mainColor);

    currentConstraint = constraint;
}

(async function() {
    const $body = document.querySelector('body');

    // Play initial animations on page load.
    window.onload = () => {
        window.setTimeout(() => {
            $body.classList.remove('is-preload');
        }, 100);
    }

    const $wrapper = document.createElement('div');
    $wrapper.id = 'bg';
    $body.appendChild($wrapper);

    const res = await fetch('/images');
    const allImages = await res.json();

    let selected = new URLSearchParams(window.location.search).get('img')

    if (selected == null) {
        selected = allImages[Math.floor(Math.random()*allImages.length)];
    }

    // Register all media query listeners
    for (let [q, c] of Object.entries(queries)) {
        const m = window.matchMedia(q)

        if (m.matches) {
            printRandomBackground($wrapper, selected, c);
        }

        m.addEventListener('change', e => {
            if (e.matches && (currentConstraint === undefined || currentConstraint.requiresUpdate(c))) {
                printRandomBackground($wrapper, selected, c);
            }
        });
    }
})();
