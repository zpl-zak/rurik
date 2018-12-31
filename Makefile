all:
	go build -o build/rurik.exe src/demo/*.go

win:
	CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o build/rurik.exe src/demo/*.go

wasm:
	CGO_ENABLED=1 GOOS=js GOARCH=wasm go build -o build/rurik.exe src/demo/*.go

perf:
	go tool pprof --pdf build/cpu.pprof > build/shit.pdf

clean:
	rm -rf build/*

play:
	./play.sh

bt: all play