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
	availableRooms, err := GetRoomsWithAvailableSlots(roomType)
	if err != nil {
		fmt.Println(err)
		return Room{}, err
	}
	allocatedRoom := availableRooms[rand.Intn(len(availableRooms))]
	return allocatedRoom, nil
}

func GetRoomsWithAvailableSlots(roomType string) ([]Room, error) {
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

func GetPersonDetails(personId string) {
	person, err := getPersonById(personId)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s %s - %s\n", person.FirstName, person.LastName, person.Category)
}

func getRoomById(roomId string) (*Room, error) {
	for _, room := range Rooms {
		if room.Id == roomId {
			return &room, nil
		}
	}
	return nil, errors.New("Room matching ID not found.")
}

func getPersonById(personId string) (*Person, error) {
	for _, person := range People {
		if person.Id == personId {
			return &person, nil
		}
	}
	return nil, errors.New("Person matching ID not found.")
}

func GetOfficeFromPersonId(personId string) (*Room, error) {
	person, err := getPersonById(personId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Printf("%s %s - %s\n", person.FirstName, person.LastName, person.Category)
	for _, room := range Rooms {
		if room.Category == "office" {
			for _, p := range room.Occupants {
				if p.Id == person.Id {
					return &room, nil
				}
			}
		}
	}
	return nil, errors.New("Person matching ID is not allocated an office space.")
}

func GetLivingSpaceFromPersonId(personId string) (*Room, error) {
	// for fellows only
	person, err := getPersonById(personId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if person.Category == "fellow" {
		fmt.Printf("%s %s - %s\n", person.FirstName, person.LastName, person.Category)
		for _, room := range Rooms {
			if room.Category == "living_space" {
				for _, p := range room.Occupants {
					if p.Id == person.Id {
						return &room, nil
					}
				}
			}
		}
		return nil, errors.New("Person matching ID is not allocated a living space.")
	}
	return nil, errors.New("Staff members are not allocated living spaces.")
}

func GetAllocatedPeople() ([]Person, error) {
	allocatedPeople := make([]Person, 1)
	for _, room := range Rooms {
		for _, person := range room.Occupants {
			if len(person.Id) > 0 {
				allocatedPeople = append(allocatedPeople, person)
			}
		}
	}
	if len(allocatedPeople) > 0 {
		return allocatedPeople, nil
	}
	return []Person{}, errors.New("There are no allocated people.")
}

func ReallocatePerson(personId, newRoomId, roomType string) {
	newRoom, rmErr := getRoomById(newRoomId)
	if rmErr != nil {
		fmt.Println(rmErr)
		return
	}
	_, pErr := getPersonById(personId)
	if pErr != nil {
		fmt.Println(pErr)
		return
	}
	oldRoom := Room{}
	if roomType == "office" {
		room, err := GetOfficeFromPersonId(personId)
		if err != nil {
			fmt.Println("Reallocation unsuccessful.")
			fmt.Println(err)
			return
		}
		oldRoom = *room
	} else {
		room, err := GetOfficeFromPersonId(personId)
		if err != nil {
			fmt.Println("Reallocation unsuccessful.")
			fmt.Println(err)
			return
		}
		oldRoom = *room
	}
	if oldRoom.Category != newRoom.Category {
		fmt.Println("Cannot reallocate to a different room type.")
		return
	}
	if oldRoom.Id == newRoom.Id {
		fmt.Println("Cannot reallocate to the same room.")
		return
	}
	for k, p := range oldRoom.Occupants {
		if p.Id == personId {
			oldRoom.Occupants[k] = Person{}
			// reallocate
			if newRoom.getOccupiedSlots() > 0 {
				newRoom.Occupants[newRoom.getOccupiedSlots()-1] = p
			}
			newRoom.Occupants[newRoom.getOccupiedSlots()] = p
			break
		}
	}
	fmt.Println("Reallocation done successfully.")
}

/*
func GetUnallocatedPeople() ([]Person, error) {
	unAllocatedPeople := make([]Person, 1)
	allocated, err := GetAllocatedPeople()
	if err != nil {
		return People, nil
	}
	for _, person := range People {
		for _, p := range allocated {
			if person.Id == p.Id {
				continue
			}
			unAllocatedPeople = append(unAllocatedPeople, person)
		}
	}
	if len(unAllocatedPeople) > 0 {
		return unAllocatedPeople, nil
	}
	return []Person{}, errors.New("There are no unallocated people.")
}
*/
