BASE_DIR := .
SRC := $(BASE_DIR)/cmd/cli
TARGET_DIR := $(BASE_DIR)/target/bin/
BINARY_NAME := $(TARGET_DIR)/nes

build:
	@echo "Building for Linux"
	mkdir -p $(TARGET_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -v -o $(BINARY_NAME) $(SRC)

build-windows:
	@echo "Building for Windows 64-bit"
	mkdir -p $(TARGET_DIR)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 \
	CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ \
	go build -ldflags "-extldflags=-static" -v -o $(BINARY_NAME).exe $(SRC)

build-windows-32:
	@echo "Building for Windows 32-bit"
	mkdir -p $(TARGET_DIR)
	GOOS=windows GOARCH=386 CGO_ENABLED=1 \
	CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ \
	go build -v -o $(BINARY_NAME)_32.exe $(SRC)

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe $(BINARY_NAME)_32.exe
