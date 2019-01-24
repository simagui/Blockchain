package main

import "os"

func isFileExits(fileName string) bool {
	_, err := os.Stat(fileName)

	if os.IsNotExist(err) {
		return false
	}
	return true
}
