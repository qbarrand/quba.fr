all: quba-fr

quba-fr: $(wildcard cmd/**/*.go internal/**/*.go pkg/**/*.go go.mod go.sum)
	go build -o $@ ./cmd
