.PHONY: build
build: 
	CGO_ENABLED=1 GOARCH=arm GOOS=linux CC=${CROSS_TC}-gcc CXX=${CROSS_TC}-g++ go build -o ./build/kobowriter