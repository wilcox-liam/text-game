package text-game

type Room struct  {
	Name string
	Enter string
	North *Room
	East *Room
	South *Room
	West *Room
	Items []*Item
}

type Item struct  {
	Name string
}

func (r Room) go_north() *Room {
	return r.North
}

func (r Room) go_east() *Room {
	return r.East
}

func (r Room) go_south() *Room {
	return r.South
}

func (r Room) go_west() *Room {
	return r.West
}

func (r Room) get_options() string {
	options := ""
	if r.North != nil {

		options += " [North] "
	}
	if r.East != nil {
		options += " [East] "
	}
	if r.South != nil {
		options += " [South] "
	}
	if r.West != nil {
		options += " [West] "
	}
	return options
}