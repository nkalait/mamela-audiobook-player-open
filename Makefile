.DEFAULT_GOAL := run 

#########################################################################
#########################################################################
#########################################################################
#########################################################################
BINARY_NAME_MAC=mamela_audiobook_player
LIB_MAC=app/lib/mac
build_mac:
	mkdir -p ${LIB_MAC}
	cp lib/mac/libbass.dylib ${LIB_MAC}
	cp lib/mac/libbass_aac.dylib ${LIB_MAC}
	cp lib/mac/libbassopus.dylib ${LIB_MAC}

#	GOARCH=amd64 GOOS=darwin go build -tags debug -o app/${BINARY_NAME_MAC}-darwin main.go
	GOARCH=amd64 GOOS=darwin go build -tags working -o app/${BINARY_NAME_MAC}-darwin main.go
#	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux main.go
#	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows main.go

	install_name_tool -change @loader_path/libbass.dylib @loader_path/lib/mac/libbass.dylib app/mamela_audiobook_player-darwin
	cd app && ./${BINARY_NAME_MAC}-darwin

#########################################################################
#########################################################################
#########################################################################
#########################################################################
# when building for linux ubuntu(using ubuntu), copy bass.h to /usr/include
# and libbass.so to /usr/lib64 and run ldconfig
# also run apt install libxxf86vm-dev.
# If having problems like "error while loading shared libraries: libbass.so: 
# cannot open shared object file: No such file or directory" then try placing
# libbass.so in /lib, also ldd the executable to see where it is loading
# some libraries from
BINARY_NAME_LINUX64=mamela_audiobook_player_linux64
LIB_LINUX64=app/lib
build_linux64:
	mkdir -p ${LIB_LINUX64}
	cp lib/linux64/libbass.so app
	cp lib/linux64/libbass_aac.so ${LIB_LINUX64}
	cp lib/linux64/libbassopus.so ${LIB_LINUX64}
	
	GOARCH=amd64 GOOS=linux go build -tags working -o app/${BINARY_NAME_LINUX64} -tags=prod_linux64 main.go
	cd app && ./${BINARY_NAME_LINUX64}

#########################################################################
#########################################################################
#########################################################################
#########################################################################
BINARY_NAME_WIN86=mamela_audiobook_player
LIB_WIN86=lib
build_win86:
	mkdir -p ${LIB_WIN86}
	cp lib/win32/bass.dll ${LIB_WIN86}
	cp lib/win32/bass_aac.dll ${LIB_WIN86}
	cp lib/win32/bassopus.dll ${LIB_WIN86}

	go build -tags prod_win -o app/${BINARY_NAME_WIN86}-win86 main.go
#########################################################################
#########################################################################
#########################################################################
#########################################################################
PACK_APP_NAME_MAC=Mamela.app
PACK_LIB_MAC=${PACK_APP_NAME_MAC}/Contents/lib/mac
APP_DIR_MAC=${PACK_APP_NAME_MAC}/Contents/MacOS
PACK_DB_DIR_MAC=${PACK_APP_NAME_MAC}/Contents/db
pack_mac:
	fyne package -os darwin -appID mamela.co.ls --tags prod_mac --release 
	mv mamela.app ${PACK_APP_NAME_MAC}
	mkdir -p ${PACK_LIB_MAC}
	cp lib/mac/libbass.dylib ${PACK_LIB_MAC}
	cp lib/mac/libbass_aac.dylib ${PACK_LIB_MAC}
	cp lib/mac/libbassopus.dylib ${PACK_LIB_MAC}
	mkdir -p ${PACK_DB_DIR_MAC}
	install_name_tool -change @loader_path/libbass.dylib @loader_path/../lib/mac/libbass.dylib ${APP_DIR_MAC}/mamela

#########################################################################
#########################################################################
#########################################################################
#########################################################################

PACK_APP_NAME_LINUX64=Mamela_linux64
PACK_LIB_LINUX64=${PACK_APP_NAME_LINUX64}/lib
APP_DIR_LINUX64=${PACK_APP_NAME_LINUX64}
PACK_DB_DIR_LINUX64=${PACK_APP_NAME_LINUX64}
pack_linux64:
# https://github.com/fyne-io/fyne-cross
	fyne package -os linux -appID mamela.co.ls --tags prod_linux64 --release 
#	mv mamela.app ${PACK_APP_NAME_LINUX64}
#	mkdir -p ${PACK_LIB_LINUX64}
#	cp lib/linux64/libbass.so ${PACK_LIB_LINUX64}
#	cp lib/linux64/libbass_aac.so ${PACK_LIB_LINUX64}
#	cp lib/linux64/libbassopus.so ${PACK_LIB_LINUX64}
#	mkdir -p ${PACK_DB_DIR_LINUX64}
#	install_name_tool -change @loader_path/libbass.dylib @loader_path/../lib/mac/libbass.dylib ${APP_DIR_LINUX64}/mamela

#	fyne-cross linux -arch=* -app-id="nada.co"


# run: build_mac
#	./${BINARY_NAME_MAC}

clean_mac:
	go clean
	rm app/${BINARY_NAME_MAC}-darwin

clean_linux64:
	go clean
	rm app/${BINARY_NAME_LINUX64}-linux64


#	rm ${BINARY_NAME}-windows
