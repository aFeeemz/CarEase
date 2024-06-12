package main

import (
	"FinalProject/initializers"
	"fmt"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDB()
}

func main() {
	fmt.Println("hehe")
}
