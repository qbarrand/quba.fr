cachedFiles = [
    // html
    '/',
    '/index.html',
    //images
    '/images/cern_transparent_bg.png',
    '/images/ensiie_tsp_transparent_bg.png',
    '/images/iutbm_transparent_bg.png',
    '/images/overlay.png',
    '/images/quentinbarrand.jpg',
    '/images/sncf_transparent_bg.png',
    // images/bg
    '/images/bg/dubai_1.jpg',
    '/images/bg/geneva_1.jpg',
    '/images/bg/montreux_1.jpg',
    '/images/bg/newyork_2.jpg',
    '/images/bg/shenzhen_1.jpg',
    '/images/bg/thun_1.jpg',
    // images/icons
    '/images/icons/android-chrome-144x144.png',
    '/images/icons/android-chrome-192x192.png',
    '/images/icons/android-chrome-256x256.png',
    '/images/icons/android-chrome-36x36.png',
    '/images/icons/android-chrome-384x384.png',
    '/images/icons/android-chrome-48x48.png',
    '/images/icons/android-chrome-512x512.png',
    '/images/icons/android-chrome-72x72.png',
    '/images/icons/android-chrome-96x96.png',
    '/images/icons/apple-touch-icon-114x114.png',
    '/images/icons/apple-touch-icon-114x114-precomposed.png',
    '/images/icons/apple-touch-icon-120x120.png',
    '/images/icons/apple-touch-icon-120x120-precomposed.png',
    '/images/icons/apple-touch-icon-144x144.png',
    '/images/icons/apple-touch-icon-144x144-precomposed.png',
    '/images/icons/apple-touch-icon-152x152.png',
    '/images/icons/apple-touch-icon-152x152-precomposed.png',
    '/images/icons/apple-touch-icon-180x180.png',
    '/images/icons/apple-touch-icon-180x180-precomposed.png',
    '/images/icons/apple-touch-icon-57x57.png',
    '/images/icons/apple-touch-icon-57x57-precomposed.png',
    '/images/icons/apple-touch-icon-60x60.png',
    '/images/icons/apple-touch-icon-60x60-precomposed.png',
    '/images/icons/apple-touch-icon-72x72.png',
    '/images/icons/apple-touch-icon-72x72-precomposed.png',
    '/images/icons/apple-touch-icon-76x76.png',
    '/images/icons/apple-touch-icon-76x76-precomposed.png',
    '/images/icons/apple-touch-icon.png',
    '/images/icons/apple-touch-icon-precomposed.png',
    '/images/icons/browserconfig.xml',
    '/images/icons/favicon-16x16.png',
    '/images/icons/favicon-194x194.png',
    '/images/icons/favicon-32x32.png',
    '/images/icons/favicon.ico',
    '/images/icons/manifest.json',
    '/images/icons/mstile-144x144.png',
    '/images/icons/mstile-150x150.png',
    '/images/icons/mstile-310x150.png',
    '/images/icons/mstile-310x310.png',
    '/images/icons/mstile-70x70.png',
    '/images/icons/safari-pinned-tab.svg',
    // css
    '/assets/css/custom.css',
    '/assets/css/font-awesome.min.css',
    '/assets/css/ie9.css',
    '/assets/css/main.css',
    '/assets/css/tooltipster.bundle.min.css',
    // js
    '/assets/js/jquery.min.js',
    '/assets/js/jquery.validate.min.js',
    '/assets/js/main.js',
    '/assets/js/skel.min.js',
    '/assets/js/sw.js',
    '/assets/js/tooltipster.bundle.min.js',
    '/assets/js/util.js',
    // fonts
    '/assets/fonts/custom.css',
    '/assets/fonts/font-awesome.min.css',
    '/assets/fonts/ie9.css',
    '/assets/fonts/main.css',
    '/assets/fonts/tooltipster.bundle.min.css'
]

self.addEventListener('install', e => e.waitUntil(
    // open a new cache
    caches
        .open('my-pwa-cache')
        .then(cache => cache.addAll(cachedFiles))
    )
);
