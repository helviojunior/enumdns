TARGET=./build
ARCHS=amd64 386 arm64 
GOOS=windows linux darwin
PACKAGENAME="github.com/helviojunior/intelparser"

ifdef VER
	VER := $(VER)
else
	VER := dev
endif

GITHASH=`git rev-parse --short HEAD`
BUILDENV=`go version | cut -d' ' -f 3,4 | sed 's/ /_/g'`
BUILDTIME=`date -u +'%Y-%m-%dT%H:%M:%SZ'`


.PHONY: help local windows linux mac all clean

default: local

help:		   ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'


get-version:
	@mkdir -p ${TARGET} ; \
	if [ "$(VER)" = "dev" ]; then \
		if [ -f ${TARGET}/.version ]; then \
			echo "Using cached release tag from .version" ; \
		else \
			echo "Fetching latest release tag from GitHub..." ; \
			curl -s https://api.github.com/repos/helviojunior/intelparser/releases/latest \
				| grep '"tag_name":' | head -n 1 | grep -oE 'v[0-9\.]+' > ${TARGET}/.version ; \
			echo "Saved release tag to .version" ; \
		fi ; \
		VER1=`cat ${TARGET}/.version` ; \
		echo "$${VER1}-dev" > ${TARGET}/.verstr ; \
	else \
		echo "VER is set to '$(VER)', skipping release tag fetch."; \
		echo "$(VER)" > ${TARGET}/.verstr ; \
	fi ; \
	echo "Version: $$(cat ${TARGET}/.verstr)"; \
	echo "-s -w \
-X=${PACKAGENAME}/internal/version.Version=$$(cat ${TARGET}/.verstr) \
-X=${PACKAGENAME}/internal/version.GitHash=${GITHASH} \
-X=${PACKAGENAME}/internal/version.GoBuildEnv=${BUILDENV} \
-X=${PACKAGENAME}/internal/version.GoBuildTime=${BUILDTIME} \
" > ${TARGET}/.ldflags


windows: get-version	## Make Windows x86 and x64 Binaries
	@for ARCH in ${ARCHS}; do \
		echo "Building for windows $${ARCH}.." ;\
		GOOS=windows GOARCH=$${ARCH} go build -a -ldflags "$$(cat ${TARGET}/.ldflags)" -o ${TARGET}/intelparser_windows_$${ARCH}.exe || exit 1 ;\
	done; \
	echo "Done."

linux: get-version	## Make Linux x86 and x64 Binaries
	@for ARCH in ${ARCHS}; do \
		echo "Building for linux $${ARCH}..." ; \
		GOOS=linux GOARCH=$${ARCH} go build -a -ldflags "$$(cat ${TARGET}/.ldflags)" -o ${TARGET}/intelparser_linux_$${ARCH} || exit 1 ;\
	done; \
	echo "Done."

mac: get-version	## Make Darwin (Mac) x86 and x64 Binaries
	@for ARCH in ${ARCHS}; do \
		if [ "$${ARCH}" != "386" ]; then \
			echo "Building for mac $${ARCH}..." ; \
			GOOS=darwin GOARCH=$${ARCH} go build -a -ldflags "$$(cat ${TARGET}/.ldflags)" -o ${TARGET}/intelparser_darwin_$${ARCH} || exit 1 ;\
		fi;\
	done; \
	echo "Done."

clean: ## Delete any binaries
	@rm -rf ${TARGET}/ ; \
	go clean -i -n ${PACKAGENAME} ; \
	echo "Done."

local: get-version	## Compila para a plataforma e arquitetura local detectada
	@OS=$$(uname | tr '[:upper:]' '[:lower:]') ; \
	ARCH=$$(if [ "$$(uname -m)" = "x86_64" ]; then echo "amd64"; else echo "arm64"; fi) ; \
	echo "Building for $$OS ($$ARCH)..." ; \
	GOOS=$$OS GOARCH=$$ARCH go build -ldflags "$$(cat ${TARGET}/.ldflags)" -o ${TARGET}/intelparser_$$OS\_$$ARCH || exit 1 ; \
	echo "Done."

all: ## Make Windows, Linux and Mac x86/x64 Binaries
all: clean windows linux mac
