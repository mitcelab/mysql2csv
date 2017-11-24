MAJOR_VERSION := 0
MINOR_VERSION := 2

all: dist

dist: clean darwin linux windows

clean:
	rm -rf build
	mkdir build

darwin: clean
	GOOS=darwin go build -o build/darwin/mysql2csv main.go
	zip -jr build/mysql2csv-darwin-v$(MAJOR_VERSION).$(MINOR_VERSION).zip build/darwin/*
	rm -rf build/darwin

linux: clean
	GOOS=linux go build -o build/linux/mysql2csv main.go
	zip -jr build/mysql2csv-linux-v$(MAJOR_VERSION).$(MINOR_VERSION).zip build/linux/*
	rm -rf build/linux

windows: clean
	GOOS=windows go build -o build/windows/mysql2csv main.go
	zip -jr build/mysql2csv-windows-v$(MAJOR_VERSION).$(MINOR_VERSION).zip build/windows/*
	rm -rf build/windows
