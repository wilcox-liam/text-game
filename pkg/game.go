package textgame

type Room struct {
	Name  string
	Description string
	North *Room
	East  *Room
	South *Room
	West  *Room
	Items []Item
}

//Items should not contain a Sub-Item of the same name
type Item struct {
	Name string
	Description string
	Openable bool
	Open bool
	Items []Item
}

type ItemContainer interface {
	GetItems() []Item
}

func (r Room) GetItems() []Item {
	return r.Items
}

func (i Item) GetItems() []Item {
	return i.Items
}

//If a struct implements these functions, it can use this interface
// type Income interface {
//	Pay // an interface can inherit an interface
// 	calculate() int
// 	source() string
// }

// type Dining_Room struct {
// 	Room // Inheritance
// 	number_of_chairs int
// }

//TODO
//Add set direction helper functions

func (r Room) GoNorth() *Room {
	return r.North
}

func (r Room) GoEast() *Room {
	return r.East
}

func (r Room) GoSouth() *Room {
	return r.South
}

func (r Room) GoWest() *Room {
	return r.West
}

func (r Room) GetDirections(gameStrings map[string]string) string {
	directions := "Directions: "
	if r.North != nil {
		directions += "[" + gameStrings["commandGoNorth"] + "] "
	}
	if r.East != nil {
		directions += "[" + gameStrings["commandGoEast"] + "] "
	}
	if r.South != nil {
		directions += "[" + gameStrings["commandGoSouth"] + "] "
	}
	if r.West != nil {
		directions += "[" + gameStrings["commandGoWest"] + "] "
	}
	return directions
}

func (r Room) GetItemByName(name string) *Item {
	return getItemByName(name, r)
}

func getItemByName(name string, ic ItemContainer) *Item {
	for _, Item := range ic.GetItems() {
		if Item.Name == name {
			return &Item
		}
		if Item.Openable {
			subItem := getItemByName(name, Item)
			if subItem != nil {
				return subItem
			}			
		}
	}
	return nil
}

func (r Room) GetItemOptions() string {
	return "Objects: " + getItemOptions(r)
}

func getItemOptions(ic ItemContainer) string {
	var options string
	for _, Item := range ic.GetItems() {
		options += " [" + Item.Name
		if Item.Open {
			options += getItemOptions(Item)
		}
		options += "]"

	}
	return options
}

