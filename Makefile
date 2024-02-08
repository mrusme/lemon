.PHONY: build

build:
	go build

build-rpi:
	GOARCH=arm64 GOOS=linux go build

copy-rpi:
	scp ./lemon mrus@l3m0n:~/lemon

