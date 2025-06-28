package main

import (
	"fmt"
	"github.com/takeru-a/self-jvm/class_file"
	"os"
)

func main() {
	f, err := os.Open("MakeJVM.class")
	if err != nil {
		panic(err)
	}

	classFile, err := class_file.ReadClassFile(f)
	if err != nil {
		panic(err)
	}

	fmt.Println("loaded methods:")
	for _, m := range classFile.Methods() {
		fmt.Printf("- %s\n", m)
	}
}