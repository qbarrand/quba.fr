const supportsWebp = new Promise(resolve =>{
    const img = new Image()
    img.onload = () => resolve(img.width > 0 && img.height > 0)
    img.onerror = () => resolve(false)
    img.src = "data:image/webp;base64,UklGRiIAAABXRUJQVlA4IBYAAAAwAQCdASoBAAEADsD+JaQAA3AAAAAA"
})

const supportsAvif = new Promise(resolve => {
    const img = new Image()
    img.onload = () => resolve(img.width > 0 && img.height > 0)
    img.onerror = () => resolve(false)
    img.src = "data:image/avif;base64,AAAAIGZ0eXBhdmlmAAAAAGF2aWZtaWYxbWlhZk1BMUIAAADybWV0YQAAAAAAAAAoaGRscgAAAAAAAAAAcGljdAAAAAAAAAAAAAAAAGxpYmF2aWYAAAAADnBpdG0AAAAAAAEAAAAeaWxvYwAAAABEAAABAAEAAAABAAABGgAAAB0AAAAoaWluZgAAAAAAAQAAABppbmZlAgAAAAABAABhdjAxQ29sb3IAAAAAamlwcnAAAABLaXBjbwAAABRpc3BlAAAAAAAAAAIAAAACAAAAEHBpeGkAAAAAAwgICAAAAAxhdjFDgQ0MAAAAABNjb2xybmNseAACAAIAAYAAAAAXaXBtYQAAAAAAAAABAAEEAQKDBAAAACVtZGF0EgAKCBgANogQEAwgMg8f8D///8WfhwB8+ErK42A="
})

import './css/style.css'
import './css/fontawesome.css'
import './css/fontawesome-brands.css'
import './css/fontawesome-solid.css'

import backgrounds from '../backgrounds/backgrounds.json'

const dateAttributeName = 'data-date';
const locationAttributeName = 'data-location';
const mainColorAttributeName = 'data-main-color';
const firstImageLoaded = new Event('firstImageLoaded')

function updateLegend(elem) {
    document.
    querySelector('meta[name=theme-color]').
    setAttribute(
        'content',
        elem.getAttribute(mainColorAttributeName),
    );
    
    document.getElementById('when').innerHTML = elem.getAttribute(dateAttributeName);
    document.getElementById('where').innerHTML = elem.getAttribute(locationAttributeName);
}

const $body = document.querySelector('body');

// Play initial animations on page load.
window.addEventListener('load', () => {
    window.setTimeout(() => {
        $body.classList.remove('is-preload');
    }, 100);
});

const delay = 6000

// Vars.
let pos = 0
let lastPos = 0
let bgDivs = []

const wrapper = document.getElementById('bg');

backgrounds.sort(() => Math.random() - 0.5);

if (!('connection' in navigator) || navigator.connection.saveData) {
    backgrounds = backgrounds.slice(0, 3)
}

const avif = await supportsAvif;
const webp = await supportsWebp;

for (const b of backgrounds) {
    // Create BG div.
    const bg = document.createElement('div')
    bg.setAttribute(dateAttributeName, b.date)
    bg.setAttribute(locationAttributeName, b.location)
    bg.setAttribute(mainColorAttributeName, b.main_color)

    const img = document.createElement('img')
    bg.appendChild(img)

    let sm

    if (avif) {
        sm = b.tree.avif
    } else if (webp) {
        sm = b.tree.webp
    } else {
        sm = b.tree.jpeg
    }

    let sources
    let screenLength

    if (window.screen.orientation.type.startsWith('portrait')) {
        screenLength = window.screen.height
        sources = sm.portrait
    } else {
        screenLength = window.screen.width
        sources = sm.landscape
    }

    let src

    for (const s of sources) {
        src = s.filename
        
        if (screenLength <= s.length) {
            break
        }
    }

    img.onload = () => {
        bgDivs.push(bg)
        wrapper.appendChild(bg);
        wrapper.dispatchEvent(firstImageLoaded)
    }

    img.src = '/backgrounds/' + src
}

function mainLoop() {
    // Main loop.
    bgDivs[pos].classList.add('visible');
    bgDivs[pos].classList.add('top');

    updateLegend(bgDivs[pos]);

    window.setInterval(() => {
        lastPos = pos;
        pos++;

        // Wrap to beginning if necessary.
        if (pos >= bgDivs.length)
            pos = 0;

        // Swap top images.
        bgDivs[lastPos].classList.remove('top');

        const current = bgDivs[pos];
        current.classList.add('visible');
        current.classList.add('top');

        updateLegend(current);

        // Hide last image after a short delay.
        window.setTimeout(div => {
            div.classList.remove('visible');
        }, delay / 2, bgDivs[lastPos]);
    }, delay);
}

wrapper.addEventListener('firstImageLoaded', mainLoop, { once: true })
