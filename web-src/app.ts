/*
    Eventually by HTML5 UP
    html5up.net | @ajlkn
    Free for personal and commercial use under the CCA 3.0 license (html5up.net/license)
*/

import './css/style.css'
import './css/fontawesome.css'
import './css/fontawesome-brands.css'
import './css/fontawesome-solid.css'

import breakpoints from '../config/breakpoints.json'

import {generateConstraints} from "./ts/constraint";
import {BackgroundManager} from "./ts/bgmgr";

(async function() {
    // Play initial animations on page load.
    window.onload = () => {
        window.setTimeout(() => {
            document.body.classList.remove('is-preload');
        }, 100);
    }

    const imgName = new URLSearchParams(window.location.search).get('img')
    const bgmgr = new BackgroundManager(imgName)

    bgmgr.addEventListener('change', ice => {
        document.getElementById('bg').style.backgroundImage = `url(${ice.url})`
        document.getElementById('where').innerHTML = ice.location;
        document.getElementById('when').innerHTML = ice.date;
        document.querySelector('meta[name=theme-color]').setAttribute('content', ice.mainColor);
    })

    generateConstraints(breakpoints.widths, breakpoints.heights).forEach(
        c => c.addEventListener('active', bgmgr.updateConstraint.bind(bgmgr))
    )

    dispatchEvent(
        new Event('resize')
    )
})();
