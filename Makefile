FONTS := $(addprefix web-src/webfonts/,fa-brands.woff2 fa-solid.woff2)

all: $(FONTS) backgrounds/backgrounds.json
	mkdir -p dist
	npx webpack --mode production
	cp -r backgrounds dist/backgrounds

$(FONTS):
	make -C fa-src
	mv fa-src/fa-brands.woff2 fa-src/fa-solid.woff2 web-src/webfonts/

.PHONY: clean
clean:
	rm -fr imgen img-out

include backgrounds.mk
