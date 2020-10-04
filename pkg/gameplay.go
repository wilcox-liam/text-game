// Package textgame provides data structures and functions to support
// development of Text Adventure Games.
package textgame

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// goDirection handles user input command.go and will set CurrentRoom to the new room.
func (g *Game) goDirection(where string) error {
	exit := g.CurrentRoom.getExitByDirection(where)
	if exit == nil {
		exit = g.CurrentRoom.getExitByName(where)
		if exit == nil {
			return fmt.Errorf(g.Dictionary["errors"]["noExit"], g.CurrentRoom.Name, where)
		}
	}
	if exit.Locked {
		return errors.New(exit.LockedString)
	}
	nextRoom := g.getRoomByID(exit.RoomID)
	entered := nextRoom.Entered
	g.setCurrentRoom(nextRoom)
	if entered == false && nextRoom.StoryString != "" {
		fmt.Print(nextRoom.StoryString)
	}
	fmt.Printf(exit.GoString)
	//fmt.Println()
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
	return fmt.Errorf(g.Dictionary["errors"]["noObject"], name, g.CurrentRoom.Name)
}

// open will set the Open attribute of a visible item to true.
func (g *Game) open(name string) error {
	item := g.getItemByName(name)
	if item == nil {
		return fmt.Errorf(g.Dictionary["errors"]["noItem"], name, g.CurrentRoom.Name)
	}
	//return if item is already open or cannot be opened.
	if item.Open {
		return fmt.Errorf(g.Dictionary["errors"]["itemOpen"], item.Name)
	}
	if item.Openable == false {
		return fmt.Errorf(g.Dictionary["errors"]["itemNotOpenable"], item.Name)
	}
	if item.Locked == true {
		return errors.New(item.LockedString)
	}
	item.Open = true
	g.DisplayItemInfo = true
	fmt.Println(item.OpenString)
	return nil
}

// take will remove an item from the room and add it to a players inventory.
// The item must be flagged as takeable.
func (g *Game) take(name string) error {
	item := g.CurrentRoom.pop(name)
	if item == nil {
		return fmt.Errorf(g.Dictionary["errors"]["noItem"], name, g.CurrentRoom.Name)
	}
	if item.Takeable {
		g.DisplayItemInfo = true
		g.Player.Inventory = append(g.Player.Inventory, *item)
		fmt.Printf(g.Dictionary["strings"]["itemAdded"], item.Name)
		fmt.Println()
		return nil
	}
	if item.NotTakeableString != "" {
		return fmt.Errorf(item.NotTakeableString)
	}
	return fmt.Errorf(g.Dictionary["errors"]["itemNotTakeable"])
}

// use actions the use function of an item in a players inventory or the room.
// An item can be used on an item or an exit.
func (g *Game) use(name string, on string) error {
	item := g.getItemByName(name)
	if item == nil {
		return fmt.Errorf(g.Dictionary["errors"]["noItem"], name, g.CurrentRoom.Name)
	}
	if on == "" {
		if item.Useable {
			fmt.Println(item.UseString)
			return nil
		}
		return fmt.Errorf(g.Dictionary["errors"]["itemNotUseable"])
	}

	itemOn := g.getItemByName(on)
	if itemOn == nil {
		exit := g.CurrentRoom.getExitByName(on)
		if exit == nil {
			return fmt.Errorf(g.Dictionary["errors"]["noItem"], on, g.CurrentRoom.Name)
		}
		return g.useOnExit(item, exit)
	}
	return g.useOnItem(item, itemOn)
}

// useOnItem actions the use function of an item on another item.
func (g *Game) useOnItem(item *item, itemOn *item) error {
	if !itemOn.Takeable && strings.ToLower(itemOn.TakeableWith) == strings.ToLower(item.Name) {
		itemOn.Takeable = true
		fmt.Print(itemOn.TakeableString)
		fmt.Println()
		g.take(itemOn.Name)
		return nil
	}	
	if itemOn.Locked && strings.ToLower(itemOn.UnlockedWith) == strings.ToLower(item.Name) {
		itemOn.Locked = false
		if itemOn.UnlockName != "" {
			itemOn.Name = itemOn.UnlockName
			itemOn.Description = itemOn.UnlockDescription
		}
		fmt.Print(itemOn.UnlockString)
		fmt.Println()
		g.open(itemOn.Name)
		return nil
	}
	if itemOn.Takeable == false {
		return fmt.Errorf(itemOn.NotTakeableString)
	}

	return fmt.Errorf(g.Dictionary["errors"]["cannotUseItem"], item.Name, itemOn.Name)
}

// useOnExit actions the use function of an item on an exit.
func (g *Game) useOnExit(item *item, exit *exit) error {
	if exit.Locked && strings.ToLower(exit.UnlockedWith) == strings.ToLower(item.Name) {
		g.unlockExit(exit)
		return nil
	}
	return fmt.Errorf(g.Dictionary["errors"]["cannotUseItem"], item.Name, exit.Name)
}

// unlockExit unlocks a matching exit.
// Exits are not bi-directional. There is a separate exit object in the other room.
// When an exit is unlocked from one room, it should unlock the exit in the other room too.
func (g *Game) unlockExit(exit *exit) {
	exit.Locked = false
	room := g.getRoomByID(exit.RoomID)
	for index, e := range room.Exits {
		fmt.Println(e.Name, exit.Name)
		if e.Name == exit.Name {
			room.Exits[index].Locked = false
		}
	}
	if exit.UnlockName != "" {
		exit.Name = exit.UnlockName
		exit.Description = exit.UnlockDescription
	}
	fmt.Println(exit.UnlockString)
	fmt.Println()
	g.goDirection(exit.Direction)
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
		helptext += "\n" + key + ": " + value + ": " + helpstring
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
