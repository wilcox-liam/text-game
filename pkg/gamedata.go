// textgame provides data structures and functions to support development of
// Text Adventure Games.
package textgame

import (
	"strings"
)

type Game struct {
	Name           string
	Description    string
	Player         *Player
	GameDictionary map[string]string
	Rooms          *[]Room
	CurrentRoom    *Room
}

type Player struct {
	Name      string
	Inventory *[]Item
}

type Room struct {
	ID          int
	Name        string
	Description string
	Exits       *[]Exit
	Items       *[]Item
}

type Exit struct {
	RoomID      int
	Name        string
	Direction   string
	Description string
	Locked      bool
	Room        *Room
}

type Item struct {
	Name        string
	Description string
	Open        bool
	Openable    bool
	OpenString  string
	Items       *[]Item
}

type ItemContainer interface {
	GetItems() *[]Item
}

// GetItems returns a slice of Items in a Room.
func (r Room) GetItems() *[]Item {
	return r.Items
}

// GetItems returns a slice of Items in an Item.
func (i Item) GetItems() *[]Item {
	return i.Items
}

// GetItems returns a slice of Items in a Player's Inventory.
func (p Player) GetItems() *[]Item {
	return p.Inventory
}

// GetRoomByID returns a Room matching a provided name.
func (g Game) GetRoomByID(id int) *Room {
	for _, room := range *g.Rooms {
		if room.ID == id {
			return &room
		}
	}
	return nil
}

// GetExitByName returns an Exit matching a provided name in a Room.
// Ignores case.
func (r Room) GetExitByName(name string) *Exit {
	for _, exit := range *r.Exits {
		if strings.ToLower(exit.Name) == strings.ToLower(name) {
			return &exit
		}
	}
	return nil
}

// GetExitByDirection returns an Exit matching a provided name in a Room.
// Ignores case.
func (r Room) GetExitByDirection(direction string) *Exit {
	//Possibly not needed if you learn how to use pointers dummy
	if r.Exits == nil {
		return nil
	}
	for _, exit := range *r.Exits {
		if strings.ToLower(exit.Direction) == strings.ToLower(direction) {
			return &exit
		}
	}
	return nil
}

// getItemByName returns an Item matching a provided name in an ItemContainer.
// Ignores Case
func getItemByName(name string, ic ItemContainer) *Item {
	//Possibly not needed if you learn how to use pointers dummy
	if ic.GetItems() == nil {
		return nil
	}

	for _, item := range *ic.GetItems() {
		if strings.ToLower(item.Name) == strings.ToLower(name) {
			return &item
		}
		subItem := getItemByName(name, item)
		if subItem != nil {
			return subItem
		}
	}
	return nil
}

// GetDirections returns a formatted string of all Exits in a Room.
func (r Room) GetDirections() string {
	directions := "Directions: "
	for _, exit := range *r.Exits {
		directions += "[" + exit.Direction + "]"
	}
	return directions
}

// GetItemOptions returns a formatted string of all Items in a Room.
func (r Room) GetItemOptions() string {
	return "Objects: " + getItemOptions(r)
}

// GetItemOptions returns a formatted string of all Items in a Player's Inventory.
func (p Player) GetItemOptions() string {
	return "Inventory: " + getItemOptions(p)
}

// getItemOptions returns a formatted string of all Items in an ItemContainer.
func getItemOptions(ic ItemContainer) string {
	//Possibly not needed if you learn how to use pointers dummy
	if ic.GetItems() == nil {
		return "[]"
	}
	var options string
	for _, item := range *ic.GetItems() {
		options += " [" + item.Name
		if item.Open {
			options += getItemOptions(item)
		}
		options += "]"

	}
	return options
}

func LoadGameState() *game {

}

func InitialiseGameState() *game {

}

func UpdateGameState() {
	
}
