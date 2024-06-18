.PHONY: build clear

build:
	@go build -o ./build/node_is_alived ./plugins/node_is_alived

clean:
	@rm -rf ./build
