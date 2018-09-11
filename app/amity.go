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

func (room *Room) CreateRoom(roomType string) (Room, error) {
	if room.Name == "" {
		fmt.Println("Cannot create a room without a name")
		return Room{}, errors.New("Cannot create a room without a name")
	}
	for _, r := range Rooms {
		if room.Name == r.Name {
			return Room{}, errors.New("A room with the same name exists.")
		}
	}
	if roomType == "office" {
		// living spaces have 4 slots
		room.Capacity = 4
	} else if roomType == "living_space" {
		// living spaces have 6 slots
		room.Capacity = 6
	} else {
		return Room{}, errors.New("A room can only be living space or office.")
	}
	room.Id = uuid.Must(uuid.NewV4()).String()
	room.Category = roomType
	Rooms = append(Rooms, *room)
	return *room, nil
}

func (person *Person) AddPerson(personType, wantsAccomodation string) {
	/*
		personType: can be either fellow or staff
		wantsAccomodation: accepts Yes or No

		on successful addition of a person, randomly allocate rooms:-
			- office
			- living space - only for fellows opting for it. STRICTLY not for staff members
	*/
	if person.FirstName == "" || person.LastName == "" {
		fmt.Println("Name of person required.")
		return
	}
	person.Id = uuid.Must(uuid.NewV4()).String()
	person.Category = personType
	People = append(People, *person)
	// allocate room
	allocatedOffice, err := allocateRoom("office")
	if err != nil {
		fmt.Println(err)
	} else {
		allocatedOffice.Occupants = make([]Person, 6)
		allocatedOffice.Occupants = append(allocatedOffice.Occupants, *person)
		fmt.Println(person.FirstName + " " + person.LastName + " has been allocated to office " + allocatedOffice.Name)
	}
	if wantsAccomodation == "yes" && person.Category == "fellow" {
		allocatedLivingSpace, err := allocateRoom("livingSpace")
		if err != nil {
			fmt.Println(err)
		} else {
			allocatedLivingSpace.Occupants = make([]Person, 4)
			allocatedLivingSpace.Occupants = append(allocatedLivingSpace.Occupants, *person)
			fmt.Println(person.FirstName + " " + person.LastName + " has been allocated to living space " + allocatedOffice.Name)
		}
	}
	if personType == "staff" && wantsAccomodation == "yes" {
		fmt.Println("Staff members are not entitled to living spaces.")
	}
}

func allocateRoom(roomType string) (Room, error) {
	// Randomly allocates a room
	availableRooms, err := allocatableRooms(roomType)
	if err != nil {
		fmt.Println(err)
		return Room{}, err
	}
	allocatedRoom := availableRooms[rand.Intn(len(availableRooms)-1)]
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
		if room.Category == roomType {
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
		fmt.Println("List of all people.\n")
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
