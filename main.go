package main

import (
	"./menu"
	"./configuration"
)


func main() {
	configuration.GetConfig();
	menu.Run();
}


