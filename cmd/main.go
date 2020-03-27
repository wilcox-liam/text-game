//Author: Liam Wilcox
//I built this basic text-adventure game as a learning exercise for golang
//Goals
//	-Multi-Language Support
//	-Multi-Player Support
//	-Coding Best Practices
//
//Known Weaknesses
//	-Technical Error messages are always in English

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
	lang := language()
	gameStrings := gameStrings(lang)
	game := game(lang)
	sanityCheck(game)	
	setInitialState(game, gameStrings)
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(game.CurrentRoom.Name)
		fmt.Println()
		fmt.Println(game.CurrentRoom.Description)
		fmt.Println(game.CurrentRoom.GetDirections(gameStrings))
		fmt.Println(game.CurrentRoom.GetItemOptions())
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		fmt.Println()
		fmt.Println()

		updateState(game, input, gameStrings)
	}
	fmt.Println(gameStrings)
}

//Valid Game Languages
func validLanguages() []string {
	return []string{"en", "es"}
}

//Returns the Game Language
func language() string {
	validLanguages := validLanguages()
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Language?", validLanguages)
	lang, _ := reader.ReadString('\n')
	lang = strings.TrimSpace(lang)

	if !contains(validLanguages, lang) {
		fmt.Println("Unknown Language")
		os.Exit(1)
	}
	return lang
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

//Reads the game yam=l file into memory for a given language
//Error Messages are always in English.
func game(lang string) *textgame.Game {
	fileName := confDir + lang + "-game.yaml"
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		os.Exit(1)
	}

	var game textgame.Game
	err = yaml.Unmarshal(yamlFile, &game)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
		os.Exit(1)
	}
	return &game
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
func updateState(game *textgame.Game, input string, gameStrings map[string]string) {
	var nextRoom *textgame.Room

	if strings.ToLower(input) == gameStrings["commandGoNorth"] {
		nextRoom = game.CurrentRoom.GoNorth()
	} else if strings.ToLower(input) == gameStrings["commandGoEast"] {
		nextRoom = game.CurrentRoom.GoEast()
	} else if strings.ToLower(input) == gameStrings["commandGoSouth"] {
		nextRoom = game.CurrentRoom.GoSouth()
	} else if strings.ToLower(input) == gameStrings["commandGoWest"] {
		nextRoom = game.CurrentRoom.GoWest()
	} else {
		fmt.Println(gameStrings["errorInvalidDirection"])
	}

	if nextRoom == nil {
		fmt.Printf(gameStrings["errorNoExit"], input)
	} else {
		game.CurrentRoom = nextRoom
	}
}

//Sets the initial the Game State
//The first room is treated as the starting room.
func setInitialState(game *textgame.Game, gameStrings map[string]string){
	currentRoom := &game.Rooms[0]
	game.CurrentRoom = currentRoom
}

//Checks the game state for any data quality issues.
func sanityCheck(game *textgame.Game) {
	//Check all the Exits align
	//Check Rooms do not contain items of the same name
}
