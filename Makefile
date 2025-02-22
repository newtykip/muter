SOURCES = main.go service.go
OUT = muter.exe

run:
	go run $(SOURCES)

clean:
	@rm -f $(OUT)

build: clean
	go build -ldflags="-s -w" -o $(OUT) $(SOURCES)
	upx --best --lzma $(OUT)