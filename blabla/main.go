package main

import (
	"errors"
	"fmt"
	
	"BinGo/binder"
)

func main() {
	err := binder.ErrBindOsFile("file")
	if errors.As(err, &binder.ErrBind) {
		fmt.Println("Binding error occurred")
	}
	if errors.As(err, &binder.ErrFileTypeNotSupported) {
		fmt.Println("File type not supported")
	}
	fmt.Println(err)
}
