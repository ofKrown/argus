package util

import (
	"log"
)

// LogError : logs error ¯\_(ツ)_/¯
func LogError(err error) {
	if err != nil {
		log.Fatalln(err);
	}
}