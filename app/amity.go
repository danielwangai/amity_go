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
}

var Rooms []Room

func (room *Room) CreateRoom(roomType string) (Room, error) {
	if room.Name == "" {
		fmt.Println("Cannot create a room without a name")
		return Room{}, errors.New("Cannot create a room without a name")
	}
	fmt.Println(room.Name)
	if roomType == "office" || roomType == "living_space" {
		room.Id = uuid.Must(uuid.NewV4()).String()
		room.Category = roomType
		Rooms = append(Rooms, *room)
		return *room, nil
	}
	return Room{}, errors.New("A room can only be living space or office.")
}
