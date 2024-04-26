BINARY_NAME=mamela_audiobook_player
LIB=app/lib/mac

PACK_APP_NAME_MAC=Mamela.app
PACK_LIB_MAC=${PACK_APP_NAME_MAC}/Contents/lib/mac
APP_DIR_MAC=${PACK_APP_NAME_MAC}/Contents/MacOS
PACK_DB_DIR_MAC=${PACK_APP_NAME_MAC}/Contents/db

.DEFAULT_GOAL := run 

build:
	mkdir -p ${LIB}
	cp lib/mac/libbass.dylib ${LIB}
	cp lib/mac/libbass_aac.dylib ${LIB}
	cp lib/mac/libbassopus.dylib ${LIB}

#	GOARCH=amd64 GOOS=darwin go build -tags debug -o app/${BINARY_NAME}-darwin main.go
	GOARCH=amd64 GOOS=darwin go build -o app/${BINARY_NAME}-darwin main.go
#	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux main.go
#	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows main.go

	install_name_tool -change @loader_path/libbass.dylib @loader_path/lib/mac/libbass.dylib app/mamela_audiobook_player-darwin
	cd app && ./${BINARY_NAME}-darwin

pack_mac:
	fyne package -os darwin -appID mamela.co.ls --tags prod_mac --release 
	mv mamela.app ${PACK_APP_NAME_MAC}
	mkdir -p ${PACK_LIB_MAC}
	cp lib/mac/libbass.dylib ${PACK_LIB_MAC}
	cp lib/mac/libbass_aac.dylib ${PACK_LIB_MAC}
	cp lib/mac/libbassopus.dylib ${PACK_LIB_MAC}
	mkdir -p ${PACK_DB_DIR_MAC}
	touch ${PACK_DB_DIR_MAC}/data.json
	chmod 777 ${PACK_DB_DIR_MAC}/data.json 
	install_name_tool -change @loader_path/libbass.dylib @loader_path/../lib/mac/libbass.dylib ${APP_DIR_MAC}/mamela

#	fyne-cross linux -arch=* -app-id="nada.co"


run: build
	./${BINARY_NAME}

clean:
	go clean
	rm app/${BINARY_NAME}-darwin
#	rm ${BINARY_NAME}-linux
#	rm ${BINARY_NAME}-windows
