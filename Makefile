all:
	go build -o build/rurik.exe src/demo/*.go

archive:
	go build -o build/archives.exe src/archives/*.go

rel:
	go build -ldflags "-s -w" -o build/rurik.exe src/demo/*.go

win:
	CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o build/rurik.exe src/demo/*.go

wasm:
	CC=~/emsdk-portable/emscripten/1.38.24/emcc CGO_ENABLED=1 GOOS=js GOARCH=wasm go build -o build/rurik.exe src/demo/*.go

perf:
	go tool pprof --pdf build/cpu.pprof > build/shit.pdf

clean:
	rm -rf build/*

ship: rel
	./ship.sh

play:
	./play.sh

bt: all play