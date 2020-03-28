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

//Cant remember if I need these
//////////////////////////////////////////////////////////////////
type ItemContainer interface {
	GetItems() *[]Item
}

func (r Room) GetItems() *[]Item {
	return r.Items
}

func (i Item) GetItems() *[]Item {
	return i.Items
}

func (r Room) GetItemByName(name string) *Item {
	return getItemByName(name, r)
}

func getItemByName(name string, ic ItemContainer) *Item {
	for _, Item := range *ic.GetItems() {
		if Item.Name == name {
			return &Item
		}
		subItem := getItemByName(name, Item)
		if subItem != nil {
			return subItem
		}
	}
	return nil
}

//////////////////////////////////////////////////////////////////

//Needed
func (g Game) GetRoomByID(id int) *Room {
	for _, Room := range *g.Rooms {
		if Room.ID == id {
			return &Room
		}
	}
	return nil
}

func (r Room) GetDirections() string {
	directions := "Directions: "
	for _, exit := range *r.Exits {
		directions += "[" + exit.Direction + "]"
	}
	return directions
}

func (r Room) GetItemOptions() string {
	return "Objects: " + getItemOptions(r)
}

func getItemOptions(ic ItemContainer) string {
	var options string
	for _, Item := range *ic.GetItems() {
		options += " [" + Item.Name
		if Item.Open {
			options += getItemOptions(Item)
		}
		options += "]"

	}
	return options
}
