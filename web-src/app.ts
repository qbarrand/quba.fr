/*
    Eventually by HTML5 UP
    html5up.net | @ajlkn
    Free for personal and commercial use under the CCA 3.0 license (html5up.net/license)
*/

import './css/style.css'
import '@fortawesome/fontawesome-free/css/all.min.css'

import breakpoints from '../config/breakpoints.json'

import {generateConstraints} from "./ts/constraint";
import {BackgroundManager} from "./ts/bgmgr";

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

    const imgName = new URLSearchParams(window.location.search).get('img')

    const bgmgr = new BackgroundManager(imgName)

    bgmgr.addEventListener('change', ice => {
        $wrapper.style.backgroundImage = `url(${ice.url})`
        $wrapper.style.backgroundPosition = 'center'
        $wrapper.style.backgroundSize = 'cover'

        document.querySelector('#where').innerHTML = ice.location;
        document.querySelector('#when').innerHTML = ice.date;
        document.querySelector('meta[name=theme-color]').setAttribute('content', ice.mainColor);
    })

    generateConstraints(breakpoints.widths, breakpoints.heights).forEach(
        c => c.addEventListener('active', bgmgr.updateConstraint.bind(bgmgr))
    )

    dispatchEvent(
        new Event('resize')
    )
})();
