package textgame

// examine will return the description of an object matching the
// provided name or direction.
func (g Game) Examine(name string) string {
	// Room
	if g.CurrentRoom.Name == name {
		return g.CurrentRoom.Description
	}
	// Items in the Room
	item := getItemByName(name, g.CurrentRoom)
	if item != nil {
		return item.Description
	}
	// Items in the player inventory
	item = getItemByName(name, g.Player)
	if item != nil {
		return item.Description
	}
	// Exits in the Room
	exit := g.CurrentRoom.GetExitByName(name)
	if exit != nil {
		return "(" + exit.Direction + "): " + exit.Description
	}
	exit = g.CurrentRoom.GetExitByDirection(name)
	if exit != nil {
		return "(" + exit.Name + "): " + exit.Description
	}
	return ""
}

// Open will set the Open attribute of an item to true in the Current Room
// or the players inventory.
// Returns the item and also returns true if an item was closed and has been opened.
// Bug(wilcox-liam) Should probably return an error object instead.
func (g Game) Open(name string) (*Item, bool) {
	item := getItemByName(name, g.CurrentRoom)
	if item == nil {
		item = getItemByName(name, g.Player)
	}
	if item == nil {
		return nil, false
	}
	//return if item is already open or cannot be opened.
	if item.Open || item.Openable == false {
		return item, false
	}
	item.Open = true
	return item, true
}

func (g game) Take(name string) {

}

func (g game) User(name string, on string) {
	
}