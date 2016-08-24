/*
	Strata by HTML5 UP
	html5up.net | @n33co
	Free for personal and commercial use under the CCA 3.0 license (html5up.net/license)
*/

(function($) {

	var settings = {
		// Parallax background effect?
		parallax: true,

		// Parallax factor (lower = more intense, higher = less intense).
		parallaxFactor: 20
	};

	skel.breakpoints({
		xlarge: '(max-width: 1800px)',
		large: '(max-width: 1280px)',
		medium: '(max-width: 980px)',
		small: '(max-width: 736px)',
		xsmall: '(max-width: 480px)'
	});

	skel.layout({
		// For the only-xxx and not-xxx classes
		conditionals: true,
		grid: true
	});

	$(function() {

		var $window = $(window),
			$body = $('body'),
			$header = $('#header');

		// Disable animations/transitions until the page has loaded.
		$body.addClass('is-loading');

		$window.on('load', function() {
			$body.removeClass('is-loading');
		});

		// Touch?
		if (skel.vars.mobile) {
			// Turn on touch mode.
			$body.addClass('is-touch');

			// Height fix (mostly for iOS).
			window.setTimeout(function() {
				$window.scrollTop($window.scrollTop() + 1);
			}, 0);
		}

		// Fix: Placeholder polyfill.
		$('form').placeholder();

		// Prioritize "important" elements on medium.
		skel.on('+medium -medium', function() {
			$.prioritize(
				'.important\\28 medium\\29',
				skel.breakpoint('medium').active
			);
		});

		// Header.

		// Parallax background.

		// Disable parallax on IE (smooth scrolling is jerky), and on mobile platforms (= better performance).
		if (skel.vars.browser == 'ie' ||	skel.vars.mobile)
			settings.parallax = false;

		if (settings.parallax) {
			skel.on('change', function() {
				if (skel.breakpoint('medium').active) {
					$window.off('scroll.strata_parallax');
					$header.css('background-position', 'top left, center center');
				}
				else {
					$header.css('background-position', 'left 0px');

					$window.on('scroll.strata_parallax', function() {
						$header.css('background-position', 'left ' + (-1 * (parseInt($window.scrollTop()) / settings.parallaxFactor)) + 'px');
					});
				}
			});
		}

		// -------------------------------------------------------------------------
		// Header
		// -------------------------------------------------------------------------
		skel.on('change', function setHeaderPosition() {
			var headerPos = 'fixed';

			if(skel.breakpoint('xsmall').active
			|| skel.breakpoint('small').active
			|| skel.breakpoint('medium').active) {
				headerPos = 'absolute';
			}

			$('#header-caption').css('position', headerPos);

			if(skel.vars.stateId) {
				var m = skel.breakpoint('small').active ? 10 : 50;

				$('#china-poptrox').poptrox({
					preload: false,
					usePopupCaption: true,
					usePopupDefaultStyling: false,
					usePopupEasyClose: false,
					usePopupNav: true,
					windowMargin: m
				});
			}
		});

		// -------------------------------------------------------------------------
		// Background
		// -------------------------------------------------------------------------
		var backgrounds = [
			{
				"filename": "newyork_1.jpg",
				"location": "New York, USA",
				"date": "February 2014"
			},

			// {
			// 	"filename": "hongkong_1.jpg",
			// 	"location": "Hong Kong, China",
			// 	"date": "August 2014"
			// },

			{
				"filename": "beijing_1.jpg",
				"location": "Beijing, China",
				"date": "August 2014"
			},

			{
				"filename": "geneva_1.jpg",
				"location": "Geneva, Switzerland",
				"date": "September 2014"
			}
		];

		// Pick a random image
		var image = backgrounds[Math.floor(Math.random() * backgrounds.length)];

		// Set it as the background
		$('#header').css('background-image', 'url(css/images/overlay.png), url(images/bg/' + image.filename + ')');
		$('#header-caption-location').html(image.location);
		$('#header-caption-location').attr('href', 'images/fulls/' + image.filename);
		$('#header-caption-date').html(image.date);

		var poptrox_caption = image.location + ' - ' + image.date + ' - <i class="fa fa-cc"></i> BY';

		$('#header-caption').poptrox({
			caption: function($a) {
				return '<strong>' + $a.text() + ' - ' + $a.next('span').text() + ' - <i class="fa fa-cc"></i> BY Quentin Barrand</strong>';
			},
			preload: false,
			usePopupCaption: true,
			usePopupDefaultStyling: false,
			usePopupEasyClose: false,
		});

		// -------------------------------------------------------------------------
		// Sliders up / down
		// -------------------------------------------------------------------------
		function skill_slider_open() {
	    $(this)
				.nextAll('.table-skills')
				.slideDown('slow');

      $(this)
        .removeClass('fa-plus-circle')
        .addClass('fa-minus-circle')
        .text('Collapse');

      $(this)
        .off('click')
        .click(skill_slider_close);
    }

    function skill_slider_close() {
      $(this)
				.nextAll('.table-skills')
				.slideUp('slow');

      $(this)
        .removeClass('fa-minus-circle')
        .addClass('fa-plus-circle')
        .text('See a list');

      $(this)
        .off('click')
        .click(skill_slider_open);
    }

    $('.skill-slider-triggerer').click(skill_slider_open);

		$.fn.slider = function(scroll_down, scroll_up, cb_down, cb_up) {
			var t = $(this);

			var this_id = '#' + this.id;
			var slider_id = '#' + t.attr('data-slider');

			t.open_slider = function() {
				$(slider_id).slideDown(function() {

					// Go to the recently opened slider
					if(scroll_down) {
						smooth_scroll(slider_id);
					}

					// Prepare for the next click
					$(this_id)
						.off('click')
						.click(t.close_slider);

					// Finally, trigger callback
					if(cb_down) {
						cb_down();
					}
				});
			};

			t.close_slider = function() {
				// Go back to the triggering button
				if(scroll_up) {
					smooth_scroll(slider_id, 'fast');
				}

				$(slider_id).slideUp();

				// Prepare for the next click
				$(this_id)
					.off('click')
					.click(t.open_slider);

				// Finally, trigger callback
				if(cb_up) {
					cb_up();
				}
			}

			t.click(t.open_slider);
		}

		var hg_triggerer = $('#hackgyver-slider-button');

		hg_triggerer.slider(true, false, function() {
			// if(skel.isActive('small') || skel.isActive('xsmall')) {
			// 	hg_triggerer.addClass('nodisplay');
			// }
		});
	});

	$('#show-interview').click(function() {
		$(this).replaceWith('\
			<iframe \
				class="image fit player-wrapper" \
				src="https://www.youtube.com/embed/sH7gx7I1Juw" \
				frameborder="0" \
				allowfullscreen> \
			</iframe>');
	});

	// -------------------------------------------------------------------------
	// Contact form
	// -------------------------------------------------------------------------
	// Ajax call to send the form's content
	function sendForm() {

		$.ajax({
		    url: 'https://formspree.io/quentin@quba.fr',
		    method: 'POST',
		    data: {
					date: new Date().toString(),
					email: $('#contact-email').val(),
					message: $('#contact-body').val(),
					name: $('#contact-name').val()
				},
		    dataType: 'json'
		})
		.done(function(response) {
			var resultDiv = $('#contact-actions');
			resultDiv.empty();

			if(response.success) {
				resultDiv
					.append('<h2>Thanks !</h2>')
					.append("<p>We'll be in touch soon.</p>");
			} else {
				var this_link = '<a href="mailto:quentin@quba.fr?subject=Fallback mailing method - quba.fr';
				this_link += '&body=' + $('#contact-body').val() + '" target="_blank">this link</a>';

				resultDiv
					.append('<h2>Something went wrong.</h2>')
					.append('<p>Please use ' + this_link + '.</p>');
				}

			resultDiv
				.css('text-align', 'center');
		});
	}

	// Tooltipster initialization
	$('input, textarea').each(function(elem) {
		$(this).tooltipster({
			trigger: 'custom',
			onlyOne: false
		});
	});

	// Validator initialization
	var validator = $('#contact-form').validate({
		submitHandler: sendForm,
		errorPlacement: function (error, element) {
				$(element).tooltipster('update', $(error).text());
				$(element).tooltipster('show');
		},
		success: function (label, element) {
				$(element).tooltipster('hide');
		}
	});

	$('#form-send').click(function() {
			$('#contact-form').submit();
	});

	$('#form-clear').click(function() {
		$('input, textarea').tooltipster('hide');
		$('#contact-form')[0].reset();
	});

	// -------------------------------------------------------------------------
	// Privacy poptrox
	// -------------------------------------------------------------------------
	$('#privacy-poptrox').poptrox();


	// -------------------------------------------------------------------------
	// Misc
	// -------------------------------------------------------------------------

	// Smooth scrolling for links
	function smooth_scroll(id, speed, callback) {
		var	bh = $('body, html');

		var pos;

		if(id == "#top") {
			pos = 0;
		} else {
			// pos = Math.max($(id).offset().top - $('#nav').height() + 30, 0);
			pos = Math.max($(id).offset().top, 0);
		}

		bh.animate({ scrollTop: pos }, speed ? speed : 'slow', 'swing', callback);
	}

	$('a').click(function(e) {
		var h = $(this).attr('href');

		if(h == undefined) {
			// Not a link
			return;
		}

		if (h.charAt(0) == '#' && h.length > 1) {
			if(e) {
				// Cancel default actions on links
				e.preventDefault();
			}

			smooth_scroll(h);
		}
	});
})(jQuery);
