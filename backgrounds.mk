imgen: $(shell find . -name '*.go' -type f) go.mod go.sum
	go build -o $@ ./cmd/imgen

BACKGROUNDS_DIR := backgrounds

$(BACKGROUNDS_DIR)/backgrounds.json: imgen
	mkdir -p $(BACKGROUNDS_DIR)
	./$< -out-dir $(BACKGROUNDS_DIR) -in-dir img-src
