package main

import (
	"fmt"

	"amity_go/app"
)

func main() {
	fmt.Println("Test")
	rm := app.Room{Id: "21321", Name: "Valhalla"}
	fmt.Println(rm.CreateRoom("office"))
}
