all:
	go build -o build/rurik.exe src/*.go

wasm:
	GOOS=js GOARCH=wasm go build -o build/rurik.exe src/*.go

perf:
	go tool pprof --pdf build/cpu.pprof > build/shit.pdf

clean:
	rm -rf build/*
