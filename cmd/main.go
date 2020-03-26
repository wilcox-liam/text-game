package main

import (
	"bufio"
	"fmt"
	"github.com/wilcox-liam/text-game/pkg"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

const gameName = "My Game"
const confDir = "../conf/"

func main() {
	var current_room *textgame.Room
	validLanguages := validLanguages()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Language?", validLanguages)
	lang, _ := reader.ReadString('\n')
	lang = strings.TrimSpace(lang)

	if !contains(validLanguages, lang) {
		fmt.Println("Unknown Language")
		os.Exit(1)
	}
	gameStrings := gameStrings(lang)
	current_room = setInitialState(gameStrings)
	fmt.Println()
	for {
		fmt.Println(current_room.Name)
		fmt.Println()
		fmt.Println(current_room.Description)
		fmt.Println(current_room.GetDirections(gameStrings))
		fmt.Println(current_room.GetItemOptions())
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		fmt.Println()
		fmt.Println()

		current_room = updateState(current_room, input, gameStrings)
	}
	fmt.Println(gameStrings)
}

//Valid Game Languages
func validLanguages() []string {
	return []string{"en", "es"}
}

//Reads the gameStrings yaml file into memory for a given language
//Error Messages are always in English.
func gameStrings(lang string) map[string]string {
	fileName := confDir + lang + ".yaml"
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		os.Exit(1)
	}

	gameStrings := make(map[string]string)
	err = yaml.Unmarshal(yamlFile, &gameStrings)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
		os.Exit(1)
	}
	return gameStrings
}

//Helper function
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//Updates the Game State
func updateState(current_room *textgame.Room, input string, gameStrings map[string]string) *textgame.Room {
	var next_room *textgame.Room

	if strings.ToLower(input) == gameStrings["commandGoNorth"] {
		next_room = current_room.GoNorth()
	} else if strings.ToLower(input) == gameStrings["commandGoEast"] {
		next_room = current_room.GoEast()
	} else if strings.ToLower(input) == gameStrings["commandGoSouth"] {
		next_room = current_room.GoSouth()
	} else if strings.ToLower(input) == gameStrings["commandGoWest"] {
		next_room = current_room.GoWest()
	} else {
		fmt.Println(gameStrings["errorInvalidDirection"])
		return current_room
	}

	if next_room == nil {
		fmt.Printf(gameStrings["errorNoExit"], input)
	} else {
		return next_room
	}
	return current_room
}

//Sets the initial the Game State
func setInitialState(gameStrings map[string]string) *textgame.Room {
	room1ItemBed := textgame.Item{"bed", "This is your bed", false, false, nil}
	room1ItemPen := textgame.Item{"pen", "This is your pen", false, false, nil}	
	room1ItemBook := textgame.Item{"book", "This is your book", false, false, nil}		
	room1ItemDeskItems := []textgame.Item{room1ItemPen, room1ItemBook}
	room1ItemDesk := textgame.Item{"desk", "This is your desk", false, true, room1ItemDeskItems}
	room1Items := []textgame.Item{room1ItemBed, room1ItemDesk}

	room1 := textgame.Room{gameStrings["room1Name"], gameStrings["room1Description"], nil, nil, nil, nil, room1Items}
	room2 := textgame.Room{gameStrings["room2Name"], gameStrings["room2Description"], nil, nil, nil, nil, nil}
	room3 := textgame.Room{gameStrings["room3Name"], gameStrings["room3Description"], nil, nil, nil, nil, nil}
	room4 := textgame.Room{gameStrings["room4Name"], gameStrings["room4Description"], nil, nil, nil, nil, nil}
	room5 := textgame.Room{gameStrings["room5Name"], gameStrings["room5Description"], nil, nil, nil, nil, nil}
	room6 := textgame.Room{gameStrings["room6Name"], gameStrings["room6Description"], nil, nil, nil, nil, nil}
	room7 := textgame.Room{gameStrings["room7Name"], gameStrings["room7Description"], nil, nil, nil, nil, nil}
	room8 := textgame.Room{gameStrings["room8Name"], gameStrings["room8Description"], nil, nil, nil, nil, nil}
	room9 := textgame.Room{gameStrings["room9Name"], gameStrings["room9Description"], nil, nil, nil, nil, nil}
	room10 := textgame.Room{gameStrings["room10Name"], gameStrings["room10Description"], nil, nil, nil, nil, nil}
	room11 := textgame.Room{gameStrings["room11Name"], gameStrings["room11Description"], nil, nil, nil, nil, nil}
	room12 := textgame.Room{gameStrings["room12Name"], gameStrings["room12Description"], nil, nil, nil, nil, nil}
	room13 := textgame.Room{gameStrings["room13Name"], gameStrings["room13Description"], nil, nil, nil, nil, nil}
	room14 := textgame.Room{gameStrings["room14Name"], gameStrings["room14Description"], nil, nil, nil, nil, nil}

	//Exits
	room1.West = &room2

	room2.West = &room3
	room2.North = &room4
	room2.East = &room1
	room2.South = &room5

	room3.East = &room2

	room4.South = &room2

	room5.South = &room2
	room5.East = &room6
	room5.West = &room8
	//room5.South == driveway

	room6.West = &room5
	room6.North = &room7

	room7.South = &room6

	room8.East = &room5
	room8.North = &room9

	room9.South = &room8
	room9.North = &room10

	room10.South = &room9
	room10.East = &room11

	room11.West = &room10
	room11.North = &room12
	room11.East = &room13
	room11.South = &room14

	room12.South = &room11

	room13.West = &room11

	room14.North = &room11

	return &room1
}
