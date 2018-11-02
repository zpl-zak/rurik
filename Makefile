all:
	go build -o build/rurik.exe src/*.go

wasm:
	GOOS=js GOARCH=wasm go build -o build/rurik.exe src/*.go

clean:
	rm -rf build/*
