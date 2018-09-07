package app

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type Person struct {
	Id            int
	FirstName     string
	LastName      string
	Category      string // fellow or staff
	IsAccomodated bool   // false
}

// Room models
type Room struct {
	Id       string
	Name     string
	Category string // office or living space
	Capacity int
}

var Rooms []Room

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
