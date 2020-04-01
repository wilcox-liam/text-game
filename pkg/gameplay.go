package textgame

import (
	"errors"
	"fmt"
	"sort"
)

// Go handles user input commandGo and will set CurrentRoom to the new room.
func (g *Game) Go(where string) (bool, error) {
	exit := g.CurrentRoom.GetExitByDirection(where)
	if exit == nil {
		return false, errors.New(fmt.Sprintf(g.GameDictionary["errors"]["noExit"], g.CurrentRoom.Name, where))
	} else {
		g.CurrentRoom = exit.Room
	}
	return true, nil
}

// Examine will return the description of an object matching the
// provided name or direction.
func (g *Game) Examine(name string) error {
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
	return errors.New(fmt.Sprintf(g.GameDictionary["errors"]["noObject"], name, g.CurrentRoom.Name))
}

// Open will set the Open attribute of an item to true in the Current Room
// or the players inventory.
func (g *Game) Open(name string) error {
	item := getItemByName(name, g.CurrentRoom)
	if item == nil {
		item = getItemByName(name, g.Player)
	}
	if item == nil {
		return errors.New(fmt.Sprintf(g.GameDictionary["errors"]["noObject"], name, g.CurrentRoom.Name))
	}
	//return if item is already open or cannot be opened.
	if item.Open {
		return errors.New(fmt.Sprintf(g.GameDictionary["errors"]["itemOpen"], item.Name))
	}
	if item.Openable == false {
		return errors.New(fmt.Sprintf(g.GameDictionary["errors"]["itemNotOpenable"], item.Name))
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
		return errors.New(fmt.Sprintf(g.GameDictionary["errors"]["noItem"], name, g.CurrentRoom.Name))
	}
	if item.Takeable {
		g.Player.Inventory = append(g.Player.Inventory, *item)
		fmt.Println(fmt.Sprintf(g.GameDictionary["strings"]["itemAdded"], item.Name))
		return nil
	} else {
		return errors.New(fmt.Sprintf(g.GameDictionary["errors"]["itemNotTakeable"]))
	}
}

// Use actions the use function of an item in a players inventory or the room.
// Bug(wilcox-liam): NYI Using an item on another item.
func (g *Game) Use(name string, on string) error {
	item := g.CurrentRoom.GetItemByName(name)
	if item == nil {
		return errors.New(fmt.Sprintf(g.GameDictionary["errors"]["noItem"], name, g.CurrentRoom.Name))
	}
	if item.Useable {
		fmt.Println(item.UseString)
		return nil
	} else {
		return errors.New(fmt.Sprintf(g.GameDictionary["errors"]["itemNotUseable"]))
	}
}

// Help returns a list of ingame commands and shortcuts based on the Game Dictionary
func (g *Game) Help() string {
	//So the help options come out in the same order every time.
	sortedKeys := sortedKeys(g.GameDictionary["shortcuts"])

	helptext := "List of commands:"
	for _, key := range sortedKeys {
		value := g.GameDictionary["shortcuts"][key]
		helpstring := g.GameDictionary["helptext"][value]
		helptext += "\n (" + key + ") " + value + ": " + helpstring
	}
	return helptext
}

// sortedKeys is a helper function to sorts the keys in a map
func sortedKeys(m map[string]string) ([]string) {
    keys := make([]string, len(m))
    i := 0
    for k := range m {
        keys[i] = k
        i++
    }
    sort.Strings(keys)
    return keys
}