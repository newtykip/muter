SOURCES = main.go service.go websocket.go
OUT = muter.exe

run:
	go run $(SOURCES)

clean:
	@rm -f $(OUT)

build: clean
	go build -ldflags="-s -w" -o $(OUT) $(SOURCES)
	upx --best --lzma $(OUT)

installer: build
	makensis installer.nsi