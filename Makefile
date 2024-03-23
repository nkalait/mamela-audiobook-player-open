BINARY_NAME=mamela_audiobook_player
LIB=app/lib/mac

build:
	mkdir -p ${LIB}
	cp /Users/nada/Dev/mamela/lib/mac/libbass.dylib /Users/nada/Dev/mamela/${LIB}
	cp /Users/nada/Dev/mamela/lib/mac/libbass_aac.dylib /Users/nada/Dev/mamela/${LIB}

	GOARCH=amd64 GOOS=darwin go build -o app/${BINARY_NAME}-darwin main.go
#	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux main.go
#	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows main.go

	install_name_tool -change @loader_path/libbass.dylib @loader_path/lib/mac/libbass.dylib app/mamela_audiobook_player-darwin
	app/${BINARY_NAME}-darwin

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm app/${BINARY_NAME}-darwin
#	rm ${BINARY_NAME}-linux
#	rm ${BINARY_NAME}-windows
