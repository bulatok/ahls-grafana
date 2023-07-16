.PHONY: build

BIN_DIR=./bin


build:
	@go build -o $(BIN_DIR)/$(service) cmd/$(service)/*

run-local:
	@make build && $(BIN_DIR)/$(service)

say:
	@echo $(foo)