all: server webapp

aot-images: $(wildcard cmd/build/*.go)
	go build -o $@ ./cmd/build

img-out: aot-images $(wildcard img-src/*)
	mkdir -p $@
	./$< \
		-img-out-dir $@ \
		-img-in-dir img-src \
		-height-breakpoints 480,736,980,1280,1690 \
		-processor vips

webapp:
	npx webpack

server: $(wildcard cmd/server) go.mod go.sum
	go build -o $@ ./cmd/server

clean:
	rm -fr img-out
	rm -f aot-images
	rm -f server
