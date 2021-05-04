all: quba-fr

quba-fr: $(shell find . -type f -name '*.go') go.mod go.sum
	go build -o $@

clean:
	rm quba-fr
