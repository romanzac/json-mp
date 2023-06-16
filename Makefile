BINARY_NAME=json-mp

build:
	go build -o ${BINARY_NAME} main.go

test:
	go test -v ./mp/

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}



