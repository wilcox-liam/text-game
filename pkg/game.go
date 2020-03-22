package text_game

type Room struct  {
	Name string
	Enter string
	North *Room
	East *Room
	South *Room
	West *Room
}

func (r Room) Go_north() *Room {
	return r.North
}

func (r Room) Go_east() *Room {
	return r.East
}

func (r Room) Go_south() *Room {
	return r.South
}

func (r Room) Go_west() *Room {
	return r.West
}

func (r Room) Get_options() string {
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