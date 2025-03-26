const supportsWebp: Promise<boolean> = new Promise(resolve =>{
    const img = new Image()
    img.onload = () => resolve(img.width > 0 && img.height > 0)
    img.onerror = () => resolve(false)
    img.src = "data:image/webp;base64,UklGRiIAAABXRUJQVlA4IBYAAAAwAQCdASoBAAEADsD+JaQAA3AAAAAA"
})

const supportsAvif: Promise<boolean> = new Promise(resolve => {
    const img = new Image()
    img.onload = () => resolve(img.width > 0 && img.height > 0)
    img.onerror = () => resolve(false)
    img.src = "data:image/avif;base64,AAAAIGZ0eXBhdmlmAAAAAGF2aWZtaWYxbWlhZk1BMUIAAADybWV0YQAAAAAAAAAoaGRscgAAAAAAAAAAcGljdAAAAAAAAAAAAAAAAGxpYmF2aWYAAAAADnBpdG0AAAAAAAEAAAAeaWxvYwAAAABEAAABAAEAAAABAAABGgAAAB0AAAAoaWluZgAAAAAAAQAAABppbmZlAgAAAAABAABhdjAxQ29sb3IAAAAAamlwcnAAAABLaXBjbwAAABRpc3BlAAAAAAAAAAIAAAACAAAAEHBpeGkAAAAAAwgICAAAAAxhdjFDgQ0MAAAAABNjb2xybmNseAACAAIAAYAAAAAXaXBtYQAAAAAAAAABAAEEAQKDBAAAACVtZGF0EgAKCBgANogQEAwgMg8f8D///8WfhwB8+ErK42A="
})

import './css/style.css'
import './css/fontawesome.css'
import './css/fontawesome-brands.css'
import './css/fontawesome-solid.css'

import bg from '../img-out/backgrounds.json'
const backgrounds: Picture[] = bg

import cfg from '../config.json'
const config: Config = cfg

class Config {
    overflowPercent: number
}

class Source {
    filename: string
    length: number
}

class SourceMap {
    portrait: Source[]
    landscape: Source[]
}

class Picture {
    date: string
    location: string
    main_color: string
    tree: {
        avif: SourceMap
        webp: SourceMap
        jpeg: SourceMap
    }
}

const dateAttributeName = 'data-date';
const locationAttributeName = 'data-location';
const mainColorAttributeName = 'data-main-color';

function updateLegend(pic: HTMLPictureElement) {
    document.
    querySelector('meta[name=theme-color]').
    setAttribute(
        'content',
        pic.getAttribute(mainColorAttributeName),
    );
    
    document.getElementById('when').innerHTML = pic.getAttribute(dateAttributeName);
    document.getElementById('where').innerHTML = pic.getAttribute(locationAttributeName);
}

async function main() {
    const $body = document.querySelector('body');

    // Play initial animations on page load.
    window.addEventListener('load', () => {
        console.log('load')
        window.setTimeout(() => {
            $body.classList.remove('is-preload');
        }, 100);
    });
    
    const delay = 6000
    
    // Vars.
    let pos = 0, lastPos = 0, $bgs: Array<HTMLElement> = []

    const wrapper = document.getElementById('bg');

    backgrounds.sort(() => Math.random() - 0.5);

    const avif = await supportsAvif;
    const webp = await supportsWebp;

    const promises: Promise<void>[] = []
    
    for (const b of backgrounds) {
        // Create BG div.
        const bg = document.createElement('div')
        bg.setAttribute(dateAttributeName, b.date)
        bg.setAttribute(locationAttributeName, b.location)
        bg.setAttribute(mainColorAttributeName, b.main_color)

        let sm: SourceMap

        if (avif) {
            sm = b.tree.avif
        } else if (webp) {
            sm = b.tree.webp
        } else {
            sm = b.tree.jpeg
        }

        let sources: Source[]
        let screenLength: number

        if (window.screen.orientation.type.startsWith('portrait')) {
            screenLength = window.screen.height
            sources = sm.portrait
        } else {
            screenLength = window.screen.width
            sources = sm.landscape
        }

        let src: string

        for (const s of sources) {
            src = s.filename
            
            if (screenLength <= s.length) {
                break
            }
        }

        const p = fetch('/img-out/' + src)
            .then((response) => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.blob();
            })
            .then((blob) => {
                bg.style.backgroundImage = `url(${URL.createObjectURL(blob)})`
                $bgs.push(bg)
                wrapper.appendChild(bg);
            })
            .catch((error) => {
                console.error('Could not load image:', error);
            });

        promises.push(p)
    }

    Promise.any(promises).then(() => {
        // Main loop.
        $bgs[pos].classList.add('visible');
        $bgs[pos].classList.add('top');

        updateLegend($bgs[pos]);

        window.setInterval(() => {
            console.log($bgs)

            lastPos = pos;
            pos++;

            // Wrap to beginning if necessary.
            if (pos >= $bgs.length)
                pos = 0;

            // Swap top images.
            $bgs[lastPos].classList.remove('top');

            const current = $bgs[pos] as HTMLDivElement;
            current.classList.add('visible');
            current.classList.add('top');

            updateLegend(current);

            // Hide last image after a short delay.
            window.setTimeout((div: HTMLDivElement) => {
                div.classList.remove('visible');
            }, delay / 2, $bgs[lastPos]);
        }, delay);
    })
}

main()
