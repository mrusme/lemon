.PHONY: build

build:
	go build

build-rpi:
	GOARCH=arm64 GOOS=linux go build

copy-rpi:
	scp ./lemon l3m0n:~/lemon

