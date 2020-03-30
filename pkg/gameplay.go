package textgame

import (
	"strings"
)

// Go handles user input commandGo and will set CurrentRoom to the new room.
// BUG(wilcox-liam) Return a string and an error instead of a bool
func (g Game) Go(where string) (string, err) {
	for _, exit := range *g.CurrentRoom.Exits {
		if strings.ToLower(exit.Direction) == strings.ToLower(where) {
			g.CurrentRoom = exit.Room
			return "", nil
		}
	}
	return g.GameDictionary["errorNoExit"], where, nil
}

// Examine will return the description of an object matching the
// provided name or direction.
func (g Game) Examine(name string) (string, err) {
	// Room
	if g.CurrentRoom.Name == name {
		return g.CurrentRoom.Description
	}
	// Items in the Room
	item := getItemByName(name, g.CurrentRoom)
	if item != nil {
		return item.Description, nil
	}
	// Items in the player inventory
	item = getItemByName(name, g.Player)
	if item != nil {
		return item.Description, nil
	}
	// Exits in the Room
	exit := g.CurrentRoom.GetExitByName(name)
	if exit != nil {
		return "(" + exit.Direction + "): " + exit.Description, nil
	}
	exit = g.CurrentRoom.GetExitByDirection(name)
	if exit != nil {
		return "(" + exit.Name + "): " + exit.Description, nil
	}
	return "", err(g.GameDictionary["errorNoObject"], name, g.CurrentRoom.Name)
}

// Open will set the Open attribute of an item to true in the Current Room
// or the players inventory.
// Bug(wilcox-liam) Should probably return an error object instead.
func (g Game) Open(name string) (string, err) {
	item := getItemByName(name, g.CurrentRoom)
	if item == nil {
		item = getItemByName(name, g.Player)
	}
	if item == nil {
		return g.GameDictionary["errorNoObject"], name, g.CurrentRoom.Name
	}
	//return if item is already open or cannot be opened.
	if item.Open {
		return g.GameDictionary["errorItemOpen"], item.Name
	}
	if item.Openable == false {
		return g.GameDictionary["errorItemNotOpenable"], item.Name
	}

	item.Open = true
	return item.OpenString
}

func (g Game) Take(name string) {

}

func (g Game) Use(name string, on string) {

}

// expandCommand takes a user entered shortcut and expands it into the full game command
// using the Game Dictionary provided in the yaml configuration.
// Bug(wilcox-liam) consider toLower
func (g Game) expandCommand(words []string) []string {
	for i, word := range words {
		lookup := strings.ToLower(g.GameDictionary[word])
		if lookup != "" {
			words[i] = lookup
		}
	}
	return words
}
