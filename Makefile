install:
	go get github.com/eiannone/keyboard; go get github.com/gookit/color
run:
	@go run main.go
build:
	@go build -o argus main.go && echo "✔️  Done"

