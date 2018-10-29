all:
	go build -o build/rurik.exe game/*.go

clean:
	rm -rf build/*
