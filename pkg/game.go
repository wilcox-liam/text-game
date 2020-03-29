// textgame provides data structures and functions to support development of
// Text Adventure Games.
package textgame

type Game struct {
	Name           string
	Description    string
	GameDictionary map[string]string
	Rooms          *[]Room
	CurrentRoom    *Room
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

// GetRoomByID returns a Room matching a provided name.
func (g Game) GetRoomByID(id int) *Room {
	for _, room := range *g.Rooms {
		if room.ID == id {
			return &room
		}
	}
	return nil
}

// GetItemByName returns an Exit matching a provided name in a Room.
func (r Room) GetExitByName(name string) *Exit {
	for _, exit := range *r.Exits {
		if exit.Name == name {
			return &exit
		}
	}
	return nil
}

// GetItemByName returns an Item matching a provided name in a Room.
func (r Room) GetItemByName(name string) *Item {
	return getItemByName(name, r)
}

// getItemByName returns an Item matching a provided name in an ItemContainer.
func getItemByName(name string, ic ItemContainer) *Item {
	//Possibly not needed if you learn how to use pointers dummy
	if ic.GetItems() == nil {
		return nil
	}

	for _, item := range *ic.GetItems() {
		if item.Name == name {
			return &item
		}
		subItem := getItemByName(name, item)
		if subItem != nil {
			return subItem
		}
	}
	return nil
}

// examine will return the description of an object matching the provided name.
// Bug(wilcox-liam): Does not yet examine the users inventory.
func (r Room) Examine(name string) string {
	if r.Name == name {
		return r.Description
	}
	item := r.GetItemByName(name)
	if item != nil {
		return item.Description
	}
	exit := r.GetExitByName(name)
	if exit != nil {
		return exit.Description
	}

	return ""
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

// getItemOptions returns a formatted string of all Items in an ItemContainer.
func getItemOptions(ic ItemContainer) string {
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
