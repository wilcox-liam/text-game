package main

import (
	"fmt"
 	"bufio"
 	"os"
 	"strings"
 	"github.com/wilcox-liam/text-game"
)

const game_name = "My Game"

func main() {
	var current_room *Room

	current_room = set_initial_state()

  	reader := bufio.NewReader(os.Stdin)	
	fmt.Println("Hello and welcome to", game_name)

	for {
		fmt.Println(current_room.Name)
		fmt.Println("----")
		fmt.Println(current_room.Enter)
		fmt.Println(current_room.get_options())

		text, _ := reader.ReadString('\n')
  		text = strings.TrimSpace(text)
  		fmt.Println()

  		current_room = update_state(current_room, text)
	}
}

//Updates the Game State
func update_state(current_room *Room, input string) *Room {
	var next_room *Room

	if input == "n" {
		next_room = current_room.go_north()
	} else if input == "e" {
		next_room = current_room.go_east()
	} else if input == "s" {
		next_room = current_room.go_south()
	} else if input == "w" {
		next_room = current_room.go_west()
	} else {
		fmt.Println("You are directionally challenged!")
	}	

	if next_room == nil {
		fmt.Println("There is no direction to the", input)
	} else {
		current_room = next_room
	}
	return current_room
}

//Sets the initial game state
func set_initial_state() *Room {
	room1:= Room{"Room 1", "Welcome to room x", nil, nil, nil, nil, nil}
	room2 := Room{"Room 2", "Welcome to room 2", nil, nil, nil, nil, nil}
	room3 := Room{"Room 3", "Welcome to room 3", nil, nil, nil, nil, nil}
	room4 := Room{"Room 4", "Welcome to room 4", nil, nil, nil, nil, nil}

	room1.North = &room2
	room2.East = &room3
	room3.South = &room4
	room4.West = &room1

	return &room1
}
