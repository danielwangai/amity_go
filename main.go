package main

import (
	"fmt"

	"amity_go/app"
)

func main() {
	rm1, _ := app.CreateRoom("Valhalla", "office")
	fmt.Println(*rm1)
	fmt.Println("========")
	rm2, _ := app.CreateRoom("Ruby", "office")
	rm4, _ := app.CreateRoom("Python", "office")
	_ = rm2
	_ = rm4
	// rm2 := app.Room{Name: "Hogwarts"}
	// rm2.CreateRoom("office")
	// p1 := app.Person{FirstName: "one", LastName: "two"}
	p1, _ := app.AddPerson("Daniel", "Maina", "staff", "yes")
	fmt.Println("OCC -> ", rm1.Occupants)
	fmt.Println("OCC -> ", rm2.Occupants)
	app.ListPeople("all")
	app.ListRooms("all")
	fmt.Println("========")
	fmt.Println(app.Rooms)
	fmt.Println("========")
	app.ListRoomDetail(rm1.Id)
	fmt.Println("========")
	rm3, _ := app.GetLivingSpaceFromPersonId(p1.Id)
	fmt.Println("====>", rm3)
	fmt.Println(app.GetAllocatedPeople())
	fmt.Println(">>>>>>>>>><<<")
	fmt.Println(app.People)
	app.ReallocatePerson(p1.Id, rm2.Id, "office")
	fmt.Println(app.Rooms)

}
