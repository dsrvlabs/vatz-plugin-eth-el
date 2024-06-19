.PHONY: build clear

build:
	@go build -o ./build/node_is_alived ./plugins/node_is_alived
	@go build -o ./build/block_sync ./plugins/block_sync

clean:
	@rm -rf ./build
