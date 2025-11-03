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
	cp lib/mac/libbassopus.dylib ${LIB_MAC}

#	GOARCH=amd64 GOOS=darwin go build -tags debug -o app/${BINARY_NAME_MAC}-darwin main.go
	GOOS=darwin CGO_LDFLAGS="-L lib/mac -rpath lib/mac" go build -tags=working -o app/${BINARY_NAME_MAC}-darwin main.go

#	install_name_tool -change @loader_path/libbass.dylib @loader_path/lib/mac/libbass.dylib app/mamela_audiobook_player-darwin
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
	
	GOOS=linux GOARCH=arm64 CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc go build -tags working -tags=prod_linux64 -o app/${BINARY_NAME_LINUX64} main.go
	cd app && ./${BINARY_NAME_LINUX64}

#########################################################################
#########################################################################
#########################################################################
#########################################################################
BINARY_NAME_LINUXARM64=mamela_audiobook_player_linuxarm64
LIB_LINUXARM64=app/lib
build_linuxarm64:
	mkdir -p ${LIB_LINUXARM64}
	cp lib/arm64/linux/libbass.so app
	cp lib/arm64/linux/libbass_aac.so ${LIB_LINUXARM64}
	cp lib/arm64/linux/libbassopus.so ${LIB_LINUXARM64}
	
	GOOS=linux go build -tags working -tags=prod_linux64 -o app/${BINARY_NAME_LINUXARM64} main.go
	cd app && ./${BINARY_NAME_LINUXARM64}

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
VERSION := v1.0.0
NUM_VERSION := 1.0.0
APP_NAME := Mamela
APP_ID := mamela.co.ls
ICON := Icon.png
LIB_DIR := lib/mac
DIST_DIR := dist
pack_mac:
	@echo "🚀 Building $(APP_NAME) $(VERSION) for macOS..."
	FYNE_LDFLAGS="-X main.Version=$(VERSION)" \
	CGO_LDFLAGS="-L $(LIB_DIR)" \
	fyne package -os darwin \
		-appID $(APP_ID) \
		--tags prod_mac \
		--release \
		-icon $(ICON) \
		-name "$(APP_NAME)" \
		-appVersion "$(NUM_VERSION)"

	@echo "🧱 Preparing bundle..."
	mv mamela.app $(APP_NAME).app
	mkdir -p $(APP_NAME).app/Contents/lib/mac
	cp $(LIB_DIR)/libbass.dylib $(APP_NAME).app/Contents/lib/mac/
	cp $(LIB_DIR)/libbassopus.dylib $(APP_NAME).app/Contents/lib/mac/
	mkdir -p $(APP_NAME).app/Contents/db
	install_name_tool -add_rpath "@loader_path/../lib/mac" $(APP_NAME).app/Contents/MacOS/mamela

	@echo "📦 Packaging..."
	rm -rf $(DIST_DIR)
	mkdir -p $(DIST_DIR)
	mv $(APP_NAME).app $(DIST_DIR)/$(APP_NAME).app
	cd $(DIST_DIR) && zip -r9 $(APP_NAME)_$(VERSION)_macos_intel.zip $(APP_NAME).app && rm -rf $(APP_NAME).app
	

	@echo "✅ Done → $(DIST_DIR)/$(APP_NAME)_$(VERSION)_macos_intel.zip"

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
