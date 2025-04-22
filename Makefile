all: backgrounds webapp

imgen: $(shell find . -name '*.go' -type f) go.mod go.sum
	go build -o $@ ./cmd/imgen

backgrounds: imgen
	mkdir -p $@
	./$< -out-dir $@ -in-dir img-src

fontawesome-subsets:
	make -C fa-src
	mv fa-src/fa-brands.woff2 fa-src/fa-solid.woff2 web-src/webfonts/

webapp: fontawesome-subsets backgrounds
	npx webpack --mode production
	cp -r backgrounds dist/backgrounds

clean:
	rm -fr imgen img-out
