all: build
build:
	go build -o ~/bin/zipfolder -ldflags="-s -w" -trimpath .