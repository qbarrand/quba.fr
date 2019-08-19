/*
    Dimension by HTML5 UP
    html5up.net | @ajlkn
    Free for personal and commercial use under the CCA 3.0 license (html5up.net/license)
*/

const sizes = {
    XS: {
        index: 0,
        suffix: '_xs',
    },
    S: {
        index: 1,
        suffix: '_s',
    },
    M: {
        index: 2,
        suffix: '_m',
    },
    L: {
        index: 3,
        suffix: '_l',
    },
    XL: {
        index: 4,
        suffix: '_xl',
    },
    FULL: {
        index: 5,
        suffix: '_full'
    },
}

const queries = {
    // XS: 480px
    //  S: 736px
    //  M: 980px
    //  L: 1280px
    // XL: 1690px

    // Portrait
    '(orientation: portrait) and (max-height: 480px)': sizes.XS,
    '(orientation: portrait) and (min-height: 481px) and (max-height: 736px)': sizes.S,
    '(orientation: portrait) and (min-height: 737px) and (max-height: 980px)': sizes.M,
    '(orientation: portrait) and (min-height: 981px) and (max-height: 1280px)': sizes.L,
    '(orientation: portrait) and (min-height: 1281px) and (max-height: 1690px)': sizes.XL,
    '(orientation: portrait) and (min-height: 1691px)': sizes.FULL,

    // Landscape
    '(orientation: landscape) and (max-width: 480px)': sizes.XS,
    '(orientation: landscape) and (min-width: 481px) and (max-width: 736px)': sizes.S,
    '(orientation: landscape) and (min-width: 737px) and (max-width: 980px)': sizes.M,
    '(orientation: landscape) and (min-width: 981px) and (max-width: 1280px)': sizes.L,
    '(orientation: landscape) and (min-width: 1281px) and (max-width: 1690px)': sizes.XL,
    '(orientation: landscape) and (min-width: 1691px)': sizes.FULL
}

let currentIndex;

(function($) {
    $(function() {
        var $window = $(window),
            $body = $('body'),
            $wrapper = $('#wrapper'),
            $header = $('#header'),
            $footer = $('#footer'),
            $main = $('#main'),
            $main_articles = $main.children('article');

        // Disable animations/transitions until the page has loaded.
        $body.addClass('is-loading');

        $window.on('load', function() {
            window.setTimeout(function() {
                $body.removeClass('is-loading');
            }, 100);
        });

        // Fix: Placeholder polyfill.
        $('form').placeholder();

        // Fix: Flexbox min-height bug on IE.
        const ua = window.navigator.userAgent;

        if (ua.includes('MSIE') || ua.includes('Trident')) {

            var flexboxFixTimeoutId;

            $window.on('resize.flexbox-fix', function() {

                clearTimeout(flexboxFixTimeoutId);

                flexboxFixTimeoutId = setTimeout(function() {

                    if ($wrapper.prop('scrollHeight') > $window.height())
                        $wrapper.css('height', 'auto');
                    else
                        $wrapper.css('height', '100vh');

                }, 250);

            }).triggerHandler('resize.flexbox-fix');

        }

        // Nav.
        var $nav = $header.children('nav'),
            $nav_li = $nav.find('li');

        // Add "middle" alignment classes if we're dealing with an even number of items.
        if ($nav_li.length % 2 == 0) {

            $nav.addClass('use-middle');
            $nav_li.eq( ($nav_li.length / 2) ).addClass('is-middle');

        }

        // Main.
        var delay = 325,
            locked = false;

        // Methods.
        $main._show = function(id, initial) {

            var $article = $main_articles.filter('#' + id);

            // No such article? Bail.
                if ($article.length == 0)
                    return;

            // Handle lock.

                // Already locked? Speed through "show" steps w/o delays.
                    if (locked || (typeof initial != 'undefined' && initial === true)) {

                        // Mark as switching.
                            $body.addClass('is-switching');

                        // Mark as visible.
                            $body.addClass('is-article-visible');

                        // Deactivate all articles (just in case one's already active).
                            $main_articles.removeClass('active');

                        // Hide header, footer.
                            $header.hide();
                            $footer.hide();

                        // Show main, article.
                            $main.show();
                            $article.show();

                        // Activate article.
                            $article.addClass('active');

                        // Unlock.
                            locked = false;

                        // Unmark as switching.
                            setTimeout(function() {
                                $body.removeClass('is-switching');
                            }, (initial ? 1000 : 0));

                        return;

                    }

                // Lock.
                    locked = true;

            // Article already visible? Just swap articles.
                if ($body.hasClass('is-article-visible')) {

                    // Deactivate current article.
                        var $currentArticle = $main_articles.filter('.active');

                        $currentArticle.removeClass('active');

                    // Show article.
                        setTimeout(function() {

                            // Hide current article.
                                $currentArticle.hide();

                            // Show article.
                                $article.show();

                            // Activate article.
                                setTimeout(function() {

                                    $article.addClass('active');

                                    // Window stuff.
                                        $window
                                            .scrollTop(0)
                                            .triggerHandler('resize.flexbox-fix');

                                    // Unlock.
                                        setTimeout(function() {
                                            locked = false;
                                        }, delay);

                                }, 25);

                        }, delay);

                }

            // Otherwise, handle as normal.
                else {

                    // Mark as visible.
                        $body
                            .addClass('is-article-visible');

                    // Show article.
                        setTimeout(function() {

                            // Hide header, footer.
                                $header.hide();
                                $footer.hide();

                            // Show main, article.
                                $main.show();
                                $article.show();

                            // Activate article.
                                setTimeout(function() {

                                    $article.addClass('active');

                                    // Window stuff.
                                        $window
                                            .scrollTop(0)
                                            .triggerHandler('resize.flexbox-fix');

                                    // Unlock.
                                        setTimeout(function() {
                                            locked = false;
                                        }, delay);

                                }, 25);

                        }, delay);

                }

        };

        $main._hide = function(addState) {

            var $article = $main_articles.filter('.active');

            // Article not visible? Bail.
                if (!$body.hasClass('is-article-visible'))
                    return;

            // Add state?
                if (typeof addState != 'undefined'
                &&  addState === true)
                    history.pushState(null, null, '#');

            // Handle lock.

                // Already locked? Speed through "hide" steps w/o delays.
                    if (locked) {

                        // Mark as switching.
                            $body.addClass('is-switching');

                        // Deactivate article.
                            $article.removeClass('active');

                        // Hide article, main.
                            $article.hide();
                            $main.hide();

                        // Show footer, header.
                            $footer.show();
                            $header.show();

                        // Unmark as visible.
                            $body.removeClass('is-article-visible');

                        // Unlock.
                            locked = false;

                        // Unmark as switching.
                            $body.removeClass('is-switching');

                        // Window stuff.
                            $window
                                .scrollTop(0)
                                .triggerHandler('resize.flexbox-fix');

                        return;

                    }

                // Lock.
                    locked = true;

            // Deactivate article.
                $article.removeClass('active');

            // Hide article.
                setTimeout(function() {

                    // Hide article, main.
                        $article.hide();
                        $main.hide();

                    // Show footer, header.
                        $footer.show();
                        $header.show();

                    // Unmark as visible.
                        setTimeout(function() {

                            $body.removeClass('is-article-visible');

                            // Window stuff.
                                $window
                                    .scrollTop(0)
                                    .triggerHandler('resize.flexbox-fix');

                            // Unlock.
                                setTimeout(function() {
                                    locked = false;
                                }, delay);

                        }, 25);

                }, delay);


        };

        // Articles.
        $main_articles.each(function() {

            var $this = $(this);

            // Close.
                $('<div class="close">Close</div>')
                    .appendTo($this)
                    .on('click', function() {
                        location.hash = '';
                    });

            // Prevent clicks from inside article from bubbling.
                $this.on('click', function(event) {
                    event.stopPropagation();
                });

        });

        // Events.
        $body.on('click', function(event) {

            // Article visible? Hide.
                if ($body.hasClass('is-article-visible'))
                    $main._hide(true);

        });

        $window.on('keyup', function(event) {

            switch (event.keyCode) {

                case 27:

                    // Article visible? Hide.
                        if ($body.hasClass('is-article-visible'))
                            $main._hide(true);

                    break;

                default:
                    break;

            }

        });

        $window.on('hashchange', function(event) {

            // Empty hash?
                if (location.hash == ''
                ||  location.hash == '#') {

                    // Prevent default.
                        event.preventDefault();
                        event.stopPropagation();

                    // Hide.
                        $main._hide();

                }

            // Otherwise, check for a matching article.
                else if ($main_articles.filter(location.hash).length > 0) {

                    // Prevent default.
                        event.preventDefault();
                        event.stopPropagation();

                    // Show article.
                        $main._show(location.hash.substr(1));

                }

        });

        // Scroll restoration.
        // This prevents the page from scrolling back to the top on a hashchange.
        if ('scrollRestoration' in history)
            history.scrollRestoration = 'manual';
        else {

            var oldScrollPos = 0,
                scrollPos = 0,
                $htmlbody = $('html,body');

            $window
                .on('scroll', function() {

                    oldScrollPos = scrollPos;
                    scrollPos = $htmlbody.scrollTop();

                })
                .on('hashchange', function() {
                    $window.scrollTop(oldScrollPos);
                });

        }

        // Initialize.

        // Hide main, articles.
        $main.hide();
        $main_articles.hide();

        // Initial article.
        if (location.hash != ''
        &&  location.hash != '#')
            $window.on('load', function() {
                $main._show(location.hash.substr(1), true);
            });

        // Add an image
        var backgrounds = {
            "shenzhen_1": {
                "location": "Shenzhen, China",
                "date": "August 2014",
                "hex_color": "#5D0C1C"
            },

            "geneva_1": {
                "location": "Geneva, Switzerland",
                "date": "June 2016",
                "hex_color": "#737D86"
            },

            "newyork_2": {
                "location": "New York, USA",
                "date": "August 2015",
                "hex_color": "#808B8F"
            },

            "thun_1": {
                "location": "Thun, Switzerland",
                "date": "May 2016",
                "hex_color": "#597FA5"
            },

            "montreux_1": {
                "location": "Montreux, Switzerland",
                "date": "October 2016",
                "hex_color": "#778693"
            },

            "dubai_1": {
                "location": "Dubai, UAE",
                "date": "June 2017",
                "hex_color": "#514C44"
            },

            "lhc_1": {
                "location": "LHC, France / Switzerland",
                "date": "August 2019",
                "hex_color": "#817365"
            }
        };

        let imageId;

        const qsImageName = new URLSearchParams(location.search).get("bgimg");

        if (qsImageName != null) {
            imageId = qsImageName;
        } else {
            const keys = Object.keys(backgrounds);

            // Pick a random image
            imageId = keys[Math.floor(Math.random() * keys.length)]
        }

        const image = backgrounds[imageId];

        $('#bg_location').text(image.location);
        $('#bg_date').text(image.date);
        $('meta[name=theme-color]').attr('content', image.hex_color);

        $('#bg').after().css({
            // 'background-image': `url(images/bg/${imageId}${size.suffix}.jpg)`,
            '-moz-transform': 'scale(1.125)',
            '-webkit-transform': 'scale(1.125)',
            '-ms-transform': 'scale(1.125)',
            'transform': 'scale(1.125)',
            '-moz-transition': '-moz-transform 0.325s ease-in-out, -moz-filter 0.325s ease-in-out',
            '-webkit-transition': '-webkit-transform 0.325s ease-in-out, -webkit-filter 0.325s ease-in-out',
            '-ms-transition': '-ms-transform 0.325s ease-in-out, -ms-filter 0.325s ease-in-out',
            'transition': 'transform 0.325s ease-in-out, filter 0.325s ease-in-out',
            'background-position': 'center',
            'background-size': 'cover',
            'background-repeat': 'no-repeat',
            'z-index': '1'
        });

        function setBackgroundImage(size) {
            console.log(`Fetching ${imageId}${size.suffix}.jpg`)

            // Show its properties on the home page
            $('#bg').after().css({
                'background-image': `url(images/bg/${imageId}${size.suffix}.jpg)`,
            });

            currentIndex = size.index;
        }

        // Register all media query listeners
        for (let [q, size] of Object.entries(queries)) {
            const m = window.matchMedia(q)

            if (m.matches) {
                setBackgroundImage(size);
            }

            m.addListener(e => {
                if (e.matches && (currentIndex == undefined || size.index > currentIndex)) {
                    setBackgroundImage(size);
                }
            });
        }

        // Skills
        $('.skill-expand').click(function(event) {
            var $this = $(this);

            // Hide the three dots and the +
            $this
                .find('.skill-dots, .skill-plus')
                .hide();

            // Reset the cursor
            $this.css('cursor', 'inherit');

            // Show the content
            $this
                .next('p')
                .slideDown();
        });

        function sendForm() {

            var resultDiv = $('#contact-actions');

            $.ajax({
                url: 'https://formspree.io/quentin@quba.fr',
                method: 'POST',
                data: {
                        date: new Date().toString(),
                        email: $('#contact-email').val(),
                        message: $('#contact-body').val(),
                        name: $('#contact-name').val(),
                        _subject: 'New message from quba.fr',
                        _format: 'plain'
                    },
                dataType: 'json'
            })
            .always(function() { resultDiv.css('text-align', 'center').empty(); })
            .done(function(response) {

                resultDiv
                    .append('<h2>Thanks !</h2>')
                    .append("<p>We'll be in touch soon.</p>");
            })
            .fail(function(response) {
                var this_link = '<a href="mailto:quentin@quba.fr?subject=Fallback mailing method - quba.fr';
                var body = encodeURIComponent($('#contact-body').val())
                body += encodeURIComponent(`\n\n--\nTechnical information: ${response.status} / ${response.responseJSON.error}`);

                this_link += `&body=${body}" target="_blank">this link</a>`;

                resultDiv
                    .append('<h3>Something went wrong.</h3>')
                    .append(`<p>Please use ${this_link}.</p>'`);
            });

        } // function sendForm()

        // Tooltipster initialization
        var tooltipFields = $('.tooltip');
        tooltipFields.tooltipster({ trigger: 'custom' });

        // Validator initialization
        $('#contact-form').validate({
            submitHandler: sendForm,
            errorPlacement: function (error, element) {
                var errorText = $(error).text();

                if (errorText) {
                    var element = $(element)
                    element.tooltipster('content', errorText);
                    element.tooltipster('open');
                }
            },
            success: function (_, element) {
                $(element).tooltipster('close');
            }
        });

        $('#form-clear').click(function() {
            tooltipFields.tooltipster('hide');
            $('#contact-form')[0].reset();
        });

        // ServiceWorker for offline caching
        if ('serviceWorker' in navigator) {
            navigator.serviceWorker.register('sw.js', { scope: '/' })
        }
    });

})(jQuery);
