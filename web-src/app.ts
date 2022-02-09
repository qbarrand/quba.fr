/*
    Eventually by HTML5 UP
    html5up.net | @ajlkn
    Free for personal and commercial use under the CCA 3.0 license (html5up.net/license)
*/

"use strict";

import './css/style.css'
import '@fortawesome/fontawesome-free/css/all.min.css'

import * as allImages from '../img-out/metadata.json'

const p = new Promise((resolve, _) => {
    const img = new Image()
    img.onload = () => resolve(img.width > 0 && img.height > 0)
    img.onerror = () => resolve(false)
    img.src = 'data:image/webp;base64,UklGRiIAAABXRUJQVlA4IBYAAAAwAQCdASoBAAEADsD+JaQAA3AAAAAA'
});

class Constraint {
    readonly direction: string
    readonly n: number

    constructor(direction: string, n: number    ) {
        this.direction = direction
        this.n = n
    }

    requiresUpdate(other): boolean {
        return other.direction !== this.direction || other.n > this.n
    }

    toQueryString(): string {
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
    readonly bloblUrl: string
    readonly mainColor: string
    readonly location: string
    readonly date: string


    constructor(bloblUrl: string, mainColor: string, location: string, date: string) {
        this.bloblUrl = bloblUrl;
        this.mainColor = mainColor;
        this.location = location;
        this.date = date;
    }
}

let currentConstraint = null;

async function printRandomBackground(parent, imageName, constraint) {
    const url = `images/${imageName}?${constraint.toQueryString()}`
    const accept = ['image/jpeg'];

    if (await p) {
        accept.unshift('image/webp');
    }

    console.debug('Fetching ' + url);

    const response = await fetch(url, {
        headers: new Headers({'Accept': accept.join(', ')})
    })

    const blob = await response.blob();

    const unixSecs = parseInt(
        response.headers.get('X-Quba-Date')
    )

    const date = new Date(unixSecs * 1000)

    const imageData = new ImageData(
        URL.createObjectURL(blob),
        response.headers.get('X-Quba-Main-Color'),
        response.headers.get('X-Quba-Location'),
        date.toLocaleDateString('en-EN', {year: 'numeric', month: 'long'})
    );

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

    let selectedKey = new URLSearchParams(window.location.search).get('img')

    if (selectedKey == null) {
        const allKeys = Object.keys(allImages)
        selectedKey = allKeys[Math.floor(Math.random()*allKeys.length)];
    }

    // const selectedImage = allImages[selectedKey]
    const selectedImage = selectedKey

    // Register all media query listeners
    for (let [q, c] of Object.entries(queries)) {
        const m = window.matchMedia(q)

        if (m.matches) {
            printRandomBackground($wrapper, selectedImage, c);
        }

        m.addEventListener('change', e => {
            if (e.matches && (currentConstraint === undefined || currentConstraint.requiresUpdate(c))) {
                printRandomBackground($wrapper, selectedImage, c);
            }
        });
    }
})();
