package app

import (
	"errors"
	"fmt"
	"math/rand"

	uuid "github.com/satori/go.uuid"
)

type Person struct {
	Id            string
	FirstName     string
	LastName      string
	Category      string // fellow or staff
	IsAccomodated bool   // false
}

// Room models
type Room struct {
	Id        string
	Name      string
	Category  string // office or living space
	Capacity  int
	Occupants []Person
}

var Rooms []Room // store all created rooms
var People []Person

func CreateRoom(name, roomType string) (*Room, error) {
	if name == "" {
		fmt.Println("Cannot create a room without a name")
		return nil, errors.New("Cannot create a room without a name")
	}
	for _, r := range Rooms {
		if name == r.Name {
			return nil, errors.New("A room with the same name exists.")
		}
	}
	room := Room{Id: uuid.Must(uuid.NewV4()).String(), Name: name, Category: roomType}
	if roomType == "office" {
		// living spaces have 4 slots
		room.Capacity = 6
		room.Occupants = make([]Person, 6)
	} else if roomType == "living_space" {
		// living spaces have 6 slots
		room.Capacity = 4
		room.Occupants = make([]Person, 4)
	} else {
		return nil, errors.New("A room can only be living space or office.")
	}
	room.Category = roomType
	Rooms = append(Rooms, room)
	return &room, nil
}

func AddPerson(firstName, lastName, personType, wantsAccomodation string) (*Person, error) {
	/*
		personType: can be either fellow or staff
		wantsAccomodation: accepts Yes or No

		on successful addition of a person, randomly allocate rooms:-
			- office
			- living space - only for fellows opting for it. STRICTLY not for staff members
	*/
	if firstName == "" || lastName == "" {
		fmt.Println("Name of person required.")
		return nil, errors.New("Name of person required.")
	}
	person := Person{Id: uuid.Must(uuid.NewV4()).String(), FirstName: firstName, LastName: lastName, Category: personType}
	People = append(People, person)
	// allocate office
	allocatedOffice, err := getRandomRoom("office")
	if err != nil {
		fmt.Println(err)
	} else {
		allocatedOffice.allocateRoom(person)
		fmt.Println(person.FirstName + " " + person.LastName + " has been allocated to office " + allocatedOffice.Name)
	}
	if wantsAccomodation == "yes" && person.Category == "fellow" {
		allocatedLivingSpace, err := getRandomRoom("living_space")
		if err != nil {
			fmt.Println(err)
		} else {
			allocatedLivingSpace.allocateRoom(person)
			fmt.Println(person.FirstName + " " + person.LastName + " has been allocated to living space " + allocatedOffice.Name)
		}
	}
	if personType == "staff" && wantsAccomodation == "yes" {
		fmt.Println("Staff members are not entitled to living spaces.")
	}
	return &person, nil
}

func (room *Room) allocateRoom(person Person) {
	occupants := room.getOccupiedSlots()
	if occupants > 0 {
		room.Occupants[occupants-1] = person
		fmt.Println("OCcupants - ", *room)
	} else {
		room.Occupants[0] = person
		fmt.Println("OCcupants -> ", *room)
	}
}

func (room Room) getOccupiedSlots() int {
	count := 0
	for _, person := range room.Occupants {
		if len(person.Id) > 0 {
			count++
		}
	}
	return count
}

func getRandomRoom(roomType string) (Room, error) {
	// Randomly allocates a room
	availableRooms, err := allocatableRooms(roomType)
	if err != nil {
		fmt.Println(err)
		return Room{}, err
	}
	allocatedRoom := availableRooms[rand.Intn(len(availableRooms))]
	return allocatedRoom, nil
}

func allocatableRooms(roomType string) ([]Room, error) {
	// return a slice of rooms of type roomType having available slots.
	availableRooms := make([]Room, 1)
	if len(Rooms) == 0 {
		return []Room{}, errors.New("There are no rooms available.")
	}
	for _, room := range Rooms {
		if room.Category != roomType {
			continue
		}
		if room.Category == roomType && (room.Capacity-room.getOccupiedSlots()) > 0 {
			availableRooms = append(availableRooms, room)
		}
	}
	return availableRooms, nil
}

func ListPeople(personType string) {
	if len(People) > 0 {
		if personType == "fellow" || personType == "staff" {
			fmt.Println("List of all " + personType + "s")
			fmt.Println("ID\tFirst Name\tLast Name\tCategory")
			for _, person := range People {
				if person.Category == personType {
					fmt.Println(person.Id + "\t" + person.FirstName + "\t" + person.LastName + "\t")
				}
			}
			return
		}
		fmt.Println("List of all people.")
		fmt.Println("ID\tFirst Name\tLast Name\tCategory")
		for _, person := range People {
			fmt.Println(person.Id + "\t" + person.FirstName + "\t" + person.LastName + "\t" + person.Category)
		}
		return
	}
	fmt.Println("There are no People in the system.")
}

func ListRooms(roomType string) {
	if len(Rooms) > 0 {
		if roomType == "office" || roomType == "living_space" {
			fmt.Println("List of all " + roomType + "s")
			fmt.Println("ID\tRoom Name\tMax Capacity")
			for _, room := range Rooms {
				if room.Category == roomType {
					fmt.Println(room.Id + "\t" + room.Name + "\t" + string(room.Capacity) + "\t")
				}
			}
			return
		}
		fmt.Println("List of all rooms")
		fmt.Println("ID\tRoom Name\tMax Capacity\tRoom Type")
		for _, room := range Rooms {
			fmt.Println(room.Id + "\t" + room.Name + "\t" + string(room.Capacity) + "\t" + room.Category)
		}
		return
	}
	fmt.Println("There are no rooms added yet.")
}

func ListRoomDetail(roomId string) {
	room, err := getRoomById(roomId)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s - %s\n\n", room.Name, room.Category)
	if len(room.Occupants) == 0 {
		fmt.Println("Room has no occupants")
		return
	}
	fmt.Println("Occupants\n")
	for _, person := range room.Occupants {
		fmt.Printf("%s\t%s\t%s\t%s", person.Id, person.FirstName, person.LastName, person.Category)
	}
}

func getRoomById(roomId string) (*Room, error) {
	for _, room := range Rooms {
		if room.Id == roomId {
			return &room, nil
		}
	}
	return nil, errors.New("Room matching ID not found.")
}
