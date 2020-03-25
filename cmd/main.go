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
	room1 := textgame.Room{gameStrings["room1Name"], gameStrings["room1Description"], nil, nil, nil, nil}
	room2 := textgame.Room{"Dining Room", "Welcome to the Dining Room", nil, nil, nil, nil}
	room3 := textgame.Room{"Kitchen", "Welcome to the Kitchen", nil, nil, nil, nil}
	room4 := textgame.Room{"Lounge Room", "Welcome to the Lounge Room", nil, nil, nil, nil}
	room5 := textgame.Room{"Master Bedroom", "Welcome to the Master Bedroom", nil, nil, nil, nil}

	//Exits
	room1.East = &room4
	room1.North = &room2

	room2.South = &room1
	room2.East = &room3

	room3.West = &room2
	room3.South = &room4

	room4.North = &room3
	room4.West = &room1
	room4.East = &room5

	room5.West = &room4

	return &room1
}

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

	for {
		fmt.Println(current_room.Name)
		fmt.Println(current_room.Description)
		fmt.Println(current_room.GetOptions(gameStrings))

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		fmt.Println("")

		current_room = updateState(current_room, input, gameStrings)
	}
	fmt.Println(gameStrings)
}