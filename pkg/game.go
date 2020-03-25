package textgame

type Room struct {
	Name  string
	Description string
	North *Room
	East  *Room
	South *Room
	West  *Room
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

func (r Room) GetOptions(gameStrings map[string]string) string {
	options := ""
	if r.North != nil {
		options += "[" + gameStrings["commandGoNorth"] + "] "
	}
	if r.East != nil {
		options += "[" + gameStrings["commandGoEast"] + "] "
	}
	if r.South != nil {
		options += "[" + gameStrings["commandGoSouth"] + "] "
	}
	if r.West != nil {
		options += "[" + gameStrings["commandGoWest"] + "] "
	}
	return options
}