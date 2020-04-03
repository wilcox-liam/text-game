// textgame provides data structures and functions to support development of
// Text Adventure Games.
package textgame

// It is idiomatic to use a pointer receiver for a method that modifies a slice

import (
	"strings"
)

type Game struct {
	Name            string
	Description     string
	Player          *Player
	GameDictionary  map[string]map[string]string
	Rooms           []Room
	CurrentRoomID   int
	CurrentRoom     *Room
	SavedGame       bool
	DisplayRoomInfo bool
}

type Player struct {
	Name      string
	Inventory []Item
}

type Room struct {
	ID          int
	Name        string
	Description string
	Exits       []Exit
	Items       []Item
}

//Examinable
//on Useable
type Exit struct {
	RoomID      int
	Name        string
	Direction   string
	Description string
	Locked      bool
	Room        *Room
}

//Examinable
//Openable
//Takeable
//Useable on
type Item struct {
	Name        string
	Description string
	Open        bool
	Openable    bool
	OpenString  string
	Takeable    bool
	Useable     bool
	UseString   string
	Items       []Item
}

//Bug(wilcox-liam): Should maybe be pointer receivers
type ItemContainer interface {
	GetItems() []Item
	SetItems(items []Item)
}

// GetItems returns a slice of Items in a Room.
func (r *Room) GetItems() []Item {
	return r.Items
}

// GetItems returns a slice of Items in an Item.
func (i *Item) GetItems() []Item {
	return i.Items
}

// GetItems returns a slice of Items in a Player's Inventory.
func (p *Player) GetItems() []Item {
	return p.Inventory
}

// GetItems returns a slice of Items in a Room.
func (r *Room) SetItems(items []Item) {
	r.Items = items
}

// GetItems returns a slice of Items in an Item.
func (i *Item) SetItems(items []Item) {
	i.Items = items
}

// GetItems returns a slice of Items in a Player's Inventory.
func (p *Player) SetItems(items []Item) {
	p.Inventory = items
}

// GetRoomByID returns a Room matching a provided name.
func (g *Game) GetRoomByID(id int) *Room {
	for index, room := range g.Rooms {
		if room.ID == id {
			return &g.Rooms[index]
		}
	}
	return nil
}

// GetExitByName returns an Exit matching a provided name in a Room.
// Ignores case.
func (r *Room) GetExitByName(name string) *Exit {
	for index, exit := range r.Exits {
		if strings.ToLower(exit.Name) == strings.ToLower(name) {
			return &r.Exits[index]
		}
	}
	return nil
}

// GetExitByDirection returns an Exit matching a provided name in a Room.
// Ignores case.
func (r *Room) GetExitByDirection(direction string) *Exit {
	for index, exit := range r.Exits {
		if strings.ToLower(exit.Direction) == strings.ToLower(direction) {
			return &r.Exits[index]
		}
	}
	return nil
}

// getItemByName returns an Item matching a provided name in an ItemContainer.
// Only returns an item if it is visible to the player. i.e not inside an unopened container.
// Ignores Case.
func getItemByName(name string, ic ItemContainer) *Item {
	items := ic.GetItems()
	for index, item := range items {
		if strings.ToLower(item.Name) == strings.ToLower(name) {
			return &items[index]
		}
		if item.Open {
			subItem := getItemByName(name, &item)
			if subItem != nil {
				return subItem
			}
		}

	}
	return nil
}

//GetItemByName returns an Item in a room
func (r *Room) GetItemByName(name string) *Item {
	return getItemByName(name, r)
}

// Pop returns and removes an item from a room
func (r *Room) Pop(name string) *Item {
	return pop(name, r)
}

// Remove an item from a slice - maintaining order
func remove(slice []Item, s int) []Item {
    return append(slice[:s], slice[s+1:]...)
}

// Pop returns and removes an item from an ItemContainer
// Bug(wilcox-liam): Better abstraction, became specific to the take command
func pop(name string, ic ItemContainer) *Item {
	items := ic.GetItems()
	for index, item := range items {
		if strings.ToLower(item.Name) == strings.ToLower(name) {
			if item.Takeable {
				items = remove(items, index)
				ic.SetItems(items)			
			}
			return &item
		}
		if item.Open {
			subItem := pop(name, &items[index])
			if subItem != nil {
				return subItem
			}
		}
	}
	return nil
}

// GetDirections returns a formatted string of all Exits in a Room.
func (r *Room) GetDirections() string {
	directions := "Directions: "
	for _, exit := range r.Exits {
		directions += "[" + exit.Direction + "]"
	}
	return directions
}

// GetDirections returns a formatted string of all Exits in a Room.
func (r *Room) GetExitOptions() string {
	var exitNames string
	for _, exit := range r.Exits {
		exitNames += "[" + exit.Name + "]"
	}
	return exitNames
}

// GetObjectOptions returns a formatted string of all Items in a Room.
func (r *Room) GetObjectOptions() string {
	return "Objects: " + r.GetExitOptions() + getItemOptions(r)
}

// GetItemOptions returns a formatted string of all Items in a Player's Inventory.
func (p *Player) GetItemOptions() string {
	options := getItemOptions(p)
	if options == "" {
		options = " []"
	}
	return "Inventory:" + options
}

// getItemOptions returns a formatted string of all Items in an ItemContainer.
func getItemOptions(ic ItemContainer) string {
	var options string
	for _, item := range ic.GetItems() {
		options += " [" + item.Name
		if item.Open {
			options += getItemOptions(&item)
		}
		options += "]"
	}
	return options
}
