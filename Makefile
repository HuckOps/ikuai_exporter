MAIN_PATH="main.go"

.PHNOY:
clear:
	@rm -rf bin
build:
	@go mod tidy
	@mkdir bin
	@go build -o bin/main ${MAIN_PATH}
