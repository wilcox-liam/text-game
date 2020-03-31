package textgame

import (
	"errors"
	"fmt"
)

// Go handles user input commandGo and will set CurrentRoom to the new room.
func (g *Game) Go(where string) (bool, error) {
	exit := g.CurrentRoom.GetExitByDirection(where)
	if exit == nil {
		return false, errors.New(fmt.Sprintf(g.GameDictionary["errorNoExit"], g.CurrentRoom.Name, where))
	} else {
		g.CurrentRoom = exit.Room
	}
	return true, nil
}

// Examine will return the description of an object matching the
// provided name or direction.
func (g *Game) Examine(name string) error {
	// Room
	if g.CurrentRoom.Name == name {
		fmt.Println(g.CurrentRoom.Description)
		return nil
	}
	// Items in the Room
	item := getItemByName(name, g.CurrentRoom)
	if item != nil {
		fmt.Println(item.Description)
		return nil
	}
	// Items in the player inventory
	item = getItemByName(name, g.Player)
	if item != nil {
		fmt.Println(item.Description)
		return nil
	}
	// Exits in the Room
	exit := g.CurrentRoom.GetExitByName(name)
	if exit != nil {
		fmt.Println("(" + exit.Direction + "): " + exit.Description)
		return nil
	}
	exit = g.CurrentRoom.GetExitByDirection(name)
	if exit != nil {
		fmt.Println("(" + exit.Name + "): " + exit.Description)
		return nil
	}
	return errors.New(fmt.Sprintf(g.GameDictionary["errorNoObject"], name, g.CurrentRoom.Name))
}

// Open will set the Open attribute of an item to true in the Current Room
// or the players inventory.
func (g *Game) Open(name string) error {
	item := getItemByName(name, g.CurrentRoom)
	if item == nil {
		item = getItemByName(name, g.Player)
	}
	if item == nil {
		return errors.New(fmt.Sprintf(g.GameDictionary["errorNoObject"], name, g.CurrentRoom.Name))
	}
	//return if item is already open or cannot be opened.
	if item.Open {
		return errors.New(fmt.Sprintf(g.GameDictionary["errorItemOpen"], item.Name))
	}
	if item.Openable == false {
		return errors.New(fmt.Sprintf(g.GameDictionary["errorItemNotOpenable"], item.Name))
	}

	item.Open = true
	fmt.Println(item.OpenString)
	return nil
}

// Take will remove an item from the room and add it to a players inventory.
// The item must be flagged as takeable.
func (g *Game) Take(name string) error {
	item := g.CurrentRoom.GetItemByName(name)
	if item == nil {
		return errors.New(fmt.Sprintf(g.GameDictionary["errorNoItem"], name, g.CurrentRoom.Name))
	}
	if item.Takeable {
		g.Player.Inventory = append(g.Player.Inventory, *item)
		fmt.Println(fmt.Sprintf(g.GameDictionary["stringItemAdded"], item.Name))	
		return nil
	} else {
		return errors.New(fmt.Sprintf(g.GameDictionary["errorItemNotTakeable"]))
	}
}

func (g *Game) Use(name string, on string) {

}

// Help returns a list of ingame commands and shortcuts
// Bug(wilcox-liam): Is not reading shortcuts from a config file. Loop over a map?
func (g *Game) Help() string {
	// Loop over the map, store the values where the key contains 'command'
	// Loop over the map again, and build the string with keys and values,
	// where value is in command slice.

	//Change gameDictionary to a map of maps?

	//Just need to identify the shortcut text
	var help string
	help += "List of commands:\n"
	help += "  " + g.GameDictionary["commandGo"] + "(g) <Direction>\n"
	help += "  " + g.GameDictionary["commandExamine"] + "(x) <Object> | <Direction>\n"
	help += "  " + g.GameDictionary["commandOpen"] + "(o) <Object>\n"
	help += "  " + g.GameDictionary["commandTake"] + "(t) <Object>\n"
	help += "  " + g.GameDictionary["commandUse"] + "(u) <Object> on <Object>\n"
	help += "  " + g.GameDictionary["commandInventory"] + "(i)\n"
	help += "  " + g.GameDictionary["commandHelp"] + "(h)\n"
	help += "  " + g.GameDictionary["commandRefresh"] + "(r)"	
	return help
}
