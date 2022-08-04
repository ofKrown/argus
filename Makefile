# Requires GO 1.15.2 (wget -c https://golang.org/dl/go1.15.2.linux-amd64.tar.gz)

install-dependencies:
	go get github.com/eiannone/keyboard; go get github.com/gookit/color
run:
	@go run main.go
build:
	@go build -o argus main.go && echo "✔️ Successfully compiled argus"

# builds and copies argus to /usr/local/bin
install: build
	@sudo rm /usr/local/bin/argus && sudo cp argus /usr/local/bin/argus && echo "✔️ Argus installed"