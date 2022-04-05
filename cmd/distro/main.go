package main

import (
	"fmt"

	"github.com/inahga/inahgo/distro"
)

func main() {
	fmt.Printf("distribution: %s\n", distro.Detect())

	packman, err := distro.GetPackageManager()
	if err != nil {
		fmt.Println("package manager: unknown")
	} else {
		fmt.Printf("package manager: %s\n", packman.Command())
	}
}
