// Package textgame provides data structures and functions to support
// development of Text Adventure Games.
package textgame

// It is idiomatic to use a pointer receiver for a method that modifies a slice

import (
	"strings"
)

// Game provides the data structures to play a text-game
type Game struct {
	Name            string
	Description     string
	Dictionary      map[string]map[string]string
	Player          *player
	Rooms           []room
	CurrentRoomID   int
	CurrentRoom     *room
	SavedGame       bool
	DisplayRoomInfo bool
}

type player struct {
	Name      string
	Inventory []item
}

type room struct {
	ID          int
	Name        string
	Description string
	Exits       []exit
	Items       []item
	Entered		bool
	StoryString	string
}

//Should I use inheritance or interfaces for items and exits?
type exit struct {
	Name         string
	Description  string
	Locked       bool
	LockedString string
	UnlockedWith string
	UnlockString string

	RoomID    int
	Direction string
	GoString  string
}

type item struct {
	Name              string
	Description       string
	Locked            bool
	LockedString      string
	UnlockedWith      string
	UnlockString      string
	TakeableWith      string
	TakeableString    string
	NotTakeableString string

	Open       bool
	Openable   bool
	OpenString string
	Takeable   bool
	Useable    bool
	UseString  string
	Items      []item
}

// itemContainer is an interface for Room, Player and Item
type itemContainer interface {
	getItems() []item
	setItems(items []item)
}

// getItems returns a slice of items in a room.
func (r *room) getItems() []item {
	return r.Items
}

// getItems returns a slice of items in an item.
func (i *item) getItems() []item {
	return i.Items
}

// getItems returns a slice of items in a player's inventory.
func (p *player) getItems() []item {
	return p.Inventory
}

// setItems returns a slice of items in a room.
func (r *room) setItems(items []item) {
	r.Items = items
}

// setItems returns a slice of items in an item.
func (i *item) setItems(items []item) {
	i.Items = items
}

// setItems returns a slice of items in a player's inventory.
func (p *player) setItems(items []item) {
	p.Inventory = items
}

// getRoomByID returns a room matching a provided id.
func (g *Game) getRoomByID(id int) *room {
	for index, room := range g.Rooms {
		if room.ID == id {
			return &g.Rooms[index]
		}
	}
	return nil
}

// setCurrentRoom sets the room the player is currently in.
func (g *Game) setCurrentRoom(room *room) {
	g.DisplayRoomInfo = true
	g.CurrentRoom = room
	g.CurrentRoomID = room.ID
	g.CurrentRoom.Entered = true
}

// getExitByName returns an exit matching a provided name in a room.
// Ignores case.
func (r *room) getExitByName(name string) *exit {
	for index, exit := range r.Exits {
		if strings.ToLower(exit.Name) == strings.ToLower(name) {
			return &r.Exits[index]
		}
	}
	return nil
}

// getExitByDirection returns an exit matching a provided direction in a Room.
// Ignores case.
func (r *room) getExitByDirection(direction string) *exit {
	for index, exit := range r.Exits {
		if strings.ToLower(exit.Direction) == strings.ToLower(direction) {
			return &r.Exits[index]
		}
	}
	return nil
}

// getItemByName will return an item given a name if it is visible
// to the player.
func (g *Game) getItemByName(name string) *item {
	item := g.CurrentRoom.getItemByName(name)
	if item == nil {
		item = g.Player.getItemByName(name)
	}
	return item
}

// getItemByName returns an item in a room if it is visible to the player.
func (r *room) getItemByName(name string) *item {
	return getItemByName(name, r)
}

// getItemByName returns an item in a a player's inventory if it is visible to the player.
func (p *player) getItemByName(name string) *item {
	return getItemByName(name, p)
}

// getItemByName returns an item in a item if it is visible to the player.
func (i *item) getItemByName(name string) *item {
	return getItemByName(name, i)
}

// getItemByName returns an Item matching a provided name in an ItemContainer.
// Only returns an item if it is visible to the player. i.e not inside an unopened container.
// Ignores Case.
func getItemByName(name string, ic itemContainer) *item {
	items := ic.getItems()
	for index, item := range items {
		if strings.ToLower(item.Name) == strings.ToLower(name) {
			return &items[index]
		}
		if item.Open {
			subItem := item.getItemByName(name)
			if subItem != nil {
				return subItem
			}
		}

	}
	return nil
}

// getItemOptions returns a formatted string of all items within a room.
func (r *room) getItemOptions() string {
	return getItemOptions(r)
}

// GetItemOptions returns a formatted string of all items in a player's inventory.
func (p *player) getItemOptions() string {
	return getItemOptions(p)
}

// getItemOptions returns a formatted string of all items within an item.
func (i *item) getItemOptions() string {
	return getItemOptions(i)
}

// getItemOptions returns a formatted string of all items in an itemContainer.
func getItemOptions(ic itemContainer) string {
	var options string
	for _, item := range ic.getItems() {
		options += " [" + item.Name
		if item.Open {
			options += item.getItemOptions()
		}
		options += "]"
	}
	return options
}

// getExitOptions returns a formatted string of all Exits in a Room name.
func (r *room) getExitOptions() string {
	var exitNames string
	for _, exit := range r.Exits {
		exitNames += " [" + exit.Name + "]"
	}
	return exitNames
}

// getDirections returns a formatted string of all exits in a room by direction.
func (r *room) getDirections() string {
	var directions string
	for _, exit := range r.Exits {
		directions += "[" + exit.Direction + "] "
	}
	return directions
}

/////////////////////

// Pop returns and removes an item from a room
func (r *room) pop(name string) *item {
	return pop(name, r)
}

// Remove an item from a slice - maintaining order
func remove(slice []item, s int) []item {
	return append(slice[:s], slice[s+1:]...)
}

// Pop returns and removes an item from an ItemContainer
// Bug(wilcox-liam): Better abstraction, this became specific to the take command
func pop(name string, ic itemContainer) *item {
	items := ic.getItems()
	for index, item := range items {
		if strings.ToLower(item.Name) == strings.ToLower(name) {
			if item.Takeable {
				items = remove(items, index)
				ic.setItems(items)
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
