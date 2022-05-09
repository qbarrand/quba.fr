all: img-out server webapp

image-resizer: $(shell find . -name '*.go' -type f) go.mod go.sum
	go build -o $@ $(GOTAGS) ./cmd/image-resizer

img-out: image-resizer $(wildcard img-src/*)
	mkdir -p $@
	./$< -img-out-dir $@ -img-in-dir img-src -processor vips

fontawesome-subsets:
	make -C fa-src
	mv fa-src/fa-brands.woff2 fa-src/fa-solid.woff2 web-src/webfonts/

webapp: fontawesome-subsets
	npx webpack --mode production

server: $(shell find . -name '*.go' -type f) go.mod go.sum
	go build -o $@ $(GOTAGS) ./cmd/server

clean:
	rm -fr img-out
	rm -f aot-images
	rm -f server
