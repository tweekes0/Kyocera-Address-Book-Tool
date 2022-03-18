test:
	go test -race db/*
	go test -race importer/*
	go test -race prompt/*

run:
	go run main.go

clean: 
	@if [ -d Database ]; then \
		rm -rf Database;\
	fi; \

	@if [ -f kyocera-ab-tool ]; then \
		rm kyocera-ab-tool*;\
	fi; 

	rm bin/*

compile:
	GOOS=linux; GOARCH=386; go build -o bin/kyocera-ab-tool-linux-x86
	GOOS=windows; GOARCH=386; go build -o bin/kyocera-ab-tool-windows-x86.exe
	GOOS=linux; GOARCH=amd64; go build -o bin/kyocera-ab-tool-linux-amd64
	GOOS=windows; GOARCH=amd64; go build -o bin/kyocera-ab-tool-windows-amd64.exe

	