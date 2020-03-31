package textgame

import (
	"errors"
	"fmt"
)

// Go handles user input commandGo and will set CurrentRoom to the new room.
func (g Game) Go(where string) (bool, error) {
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
func (g Game) Examine(name string) error {
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
func (g Game) Open(name string) error {
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

func (g Game) Take(name string) {

}

func (g Game) Use(name string, on string) {

}
