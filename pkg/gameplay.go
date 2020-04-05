package textgame

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
)

// goDirection handles user input command.go and will set CurrentRoom to the new room.
func (g *Game) goDirection(where string) error {
	exit := g.CurrentRoom.getExitByDirection(where)
	if exit == nil {
		return errors.New(fmt.Sprintf(g.Dictionary["errors"]["noExit"], g.CurrentRoom.Name, where))
	} else {
		if exit.Locked {
			return errors.New(exit.LockedString)
		}
		nextRoom := g.getRoomByID(exit.RoomID)
		g.setCurrentRoom(nextRoom)
	}
	return nil
}

// examine will return the description of an object matching the
// provided name or direction.
func (g *Game) examine(name string) error {
	item := g.getItemByName(name)
	if item != nil {
		fmt.Println(item.Description)
		return nil
	}
	// Exits in the Room
	exit := g.CurrentRoom.getExitByName(name)
	if exit != nil {
		fmt.Println("(" + exit.Direction + "): " + exit.Description)
		return nil
	}
	exit = g.CurrentRoom.getExitByDirection(name)
	if exit != nil {
		fmt.Println("(" + exit.Name + "): " + exit.Description)
		return nil
	}
	return errors.New(fmt.Sprintf(g.Dictionary["errors"]["noObject"], name, g.CurrentRoom.Name))
}

// open will set the Open attribute of a visible item to true.
func (g *Game) open(name string) error {
	item := g.getItemByName(name)
	if item == nil {
		return errors.New(fmt.Sprintf(g.Dictionary["errors"]["noObject"], name, g.CurrentRoom.Name))
	}
	//return if item is already open or cannot be opened.
	if item.Open {
		return errors.New(fmt.Sprintf(g.Dictionary["errors"]["itemOpen"], item.Name))
	}
	if item.Openable == false {
		return errors.New(fmt.Sprintf(g.Dictionary["errors"]["itemNotOpenable"], item.Name))
	}
	if item.Locked == true {
		return errors.New(item.LockedString)
	}
	item.Open = true
	fmt.Println(item.OpenString)
	return nil
}

// take will remove an item from the room and add it to a players inventory.
// The item must be flagged as takeable.
func (g *Game) take(name string) error {
	item := g.CurrentRoom.pop(name)
	if item == nil {
		return errors.New(fmt.Sprintf(g.Dictionary["errors"]["noItem"], name, g.CurrentRoom.Name))
	}
	if item.Takeable {
		g.Player.Inventory = append(g.Player.Inventory, *item)
		fmt.Println(fmt.Sprintf(g.Dictionary["strings"]["itemAdded"], item.Name))
		return nil
	} else {
		return errors.New(fmt.Sprintf(g.Dictionary["errors"]["itemNotTakeable"]))
	}
}

// use actions the use function of an item in a players inventory or the room.
// An item can be used on an item or an exit.
func (g *Game) use(name string, on string) error {
	item := g.getItemByName(name)
	if item == nil {
		return errors.New(fmt.Sprintf(g.Dictionary["errors"]["noItem"], name, g.CurrentRoom.Name))
	}
	if on == "" {
		if item.Useable {
			fmt.Println(item.UseString)
			return nil
		}
		return errors.New(fmt.Sprintf(g.Dictionary["errors"]["itemNotUseable"]))
	}
	var unlockable unlockable
	unlockable = g.getItemByName(on)
	if isNil(unlockable) {
		unlockable = g.CurrentRoom.getExitByName(on)
		if isNil(unlockable) {
			return errors.New(fmt.Sprintf(g.Dictionary["errors"]["noItem"], on, g.CurrentRoom.Name))
		}
	}
	return g.useOn(item, unlockable)
}

// useOn actions the use function of an item in a players inventory or room or another item
// or exit in the players inventory or room.
func (g *Game) useOn(item *item, unlockable unlockable) error {
	if unlockable.locked() && unlockable.unlockedWith() == item.Name {
		unlockable.setLocked(false)
		fmt.Println(unlockable.unlockString())
		return nil
	}
	return errors.New(fmt.Sprintf(g.Dictionary["errors"]["cannotUseItem"], item.Name, unlockable.name()))
}

// isNil is a helper function to determine if an interface is nil
func isNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}

// help returns a list of ingame commands and shortcuts based on the Game Dictionary
func (g *Game) help() string {
	//So the help options come out in the same order every time.
	sortedKeys := sortedKeys(g.Dictionary["shortcuts"])

	helptext := "List of commands:"
	for _, key := range sortedKeys {
		value := g.Dictionary["shortcuts"][key]
		helpstring := g.Dictionary["helptext"][value]
		helptext += "\n (" + key + ") " + value + ": " + helpstring
	}
	return helptext
}

// sortedKeys is a helper function to sort the keys in a map
func sortedKeys(m map[string]string) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}
