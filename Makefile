all: img-out server webapp

image-resizer: $(wildcard cmd/image-resizer/*.go) go.mod go.sum
	go build -o $@ ./cmd/image-resizer

img-out: image-resizer $(wildcard img-src/*)
	mkdir -p $@
	./$< \
		-img-out-dir $@ \
		-img-in-dir img-src \
		-processor vips

webapp:
	npx webpack --mode production

server: $(wildcard cmd/server/*.go) go.mod go.sum
	go build -o $@ ./cmd/server

run: img-out server webapp
	./server \
		-img-out-dir img-out \
		-webroot-dir dist

clean:
	rm -fr img-out
	rm -f aot-images
	rm -f server
