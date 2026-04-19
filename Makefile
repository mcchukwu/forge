APP_NAME := forge
CMD_PATH := cmd/$(APP_NAME)/main.go
BIN_DIR := bin

.PHONY: run build clean

run:
	go run $(CMD_PATH) $(ARGS)

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_PATH)

clean:
	rm -rf $(BIN_DIR)
