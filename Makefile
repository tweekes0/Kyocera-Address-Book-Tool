test:
	go test -race db/*
	go test -race importer/*
	go test -race prompt/*

run:
	go run main.go

clean: 
	rm -rf Database
	rm kyocera-ab-tool*
	rm bin/*

compile:
	GOOS=linux; GOARCH=386; go build -o bin/kyocera-ab-tool-linux-x86
	GOOS=windows; GOARCH=386; go build -o bin/kyocera-ab-tool-windows-x86
	
	GOOS=linux; GOARCH=arm64; go build -o bin/kyocera-ab-tool-linux-arm64
	GOOS=windows; GOARCH=arm64; go build -o bin/kyocera-ab-tool-windows-arm64

	