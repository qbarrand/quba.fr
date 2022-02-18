all: img-out server webapp

image-resizer: $(shell find . -name '*.go' -type f) go.mod go.sum
	go build -o $@ ./cmd/image-resizer

img-out: image-resizer $(wildcard img-src/*)
	mkdir -p $@
	./$< -img-out-dir $@ -img-in-dir img-src -processor vips

webapp:
	npx webpack --mode production

server: $(shell find . -name '*.go' -type f) go.mod go.sum
	go build -o $@ ./cmd/server

clean:
	rm -fr img-out
	rm -f aot-images
	rm -f server
