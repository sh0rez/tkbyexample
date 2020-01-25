package main

import "fmt"

func main() {
	defer fmt.Println("hi");
	go func() { for {} }()
	select {}
}
