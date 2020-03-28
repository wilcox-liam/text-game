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
	"flag"
	"fmt"
	"github.com/wilcox-liam/text-game/pkg"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

const confDir = "../conf/"
const langDefault = "default"

//TODO
//Log Mode
//Language
//Auto Complete Exits
func commandLineOptions() string {
	lang := flag.String("lang", "en", "Game Language")
	flag.Parse()
	return *lang
}

//Valid Game Languages
func validLanguages() []string {
	return []string{"en", "es"}
}

//Asks the user to choose a language
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

//Reads the game yaml file into memory for a given language
//Error Messages are always in English.
func game(lang string) *textgame.Game {
	fileName := confDir + lang + ".yaml"
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

//Checks the game state for any data quality issues.
func sanityCheck(game *textgame.Game) {
	//Check Rooms do not contain items of the same name. Should always be done
	//Check Exits point to valid rooms
	//Check all the Exits align. May not always be true. e.g Jumping down a cliff is a 1 way exit.
}

//Sets the initial the Game State
//The first room is treated as the starting room.
func setInitialState(game *textgame.Game) {
	currentRoom := game.GetRoomByID(1)
	game.CurrentRoom = currentRoom
}

//Loops over all exits in the game and sets the *Room based on the RoomID
//TODO
//There is a pointer problem here
func setExits(game *textgame.Game) {
	for _, room := range *game.Rooms {
		for _, exit := range *room.Exits {
			exit.Room = game.GetRoomByID(exit.RoomID)
			if exit.Room == nil {
				fmt.Println("Invalid Room ID", exit.RoomID, "in exit", room.Name)
				os.Exit(1)
			}
			return
		}
	}
}

//Expands any shortcuts into their full command
func expandCommand(game *textgame.Game, words []string) []string {
	for i, word := range words {
		lookup := strings.ToLower(game.GameDictionary[word])
		if lookup != "" {
			words[i] = lookup
		}
	}
	return words
}

//Updates the Game State
func updateState(game *textgame.Game, input string) {
	words := strings.Split(input, " ")
	words = expandCommand(game, words)
	if strings.ToLower(words[0]) == strings.ToLower(game.GameDictionary["commandGo"]) && len(words) == 2 {
		updateStateGo(game, words[1])
	}
}

//Process the go direction input
func updateStateGo(game *textgame.Game, where string) {
	for _, exit := range *game.CurrentRoom.Exits {
		if strings.ToLower(exit.Direction) == strings.ToLower(where) {
			game.CurrentRoom = exit.Room
			return
		}
	}
	fmt.Print(game.GameDictionary["errorNoExit"], input)
}

func main() {
	lang := commandLineOptions()
	if lang == langDefault {
		lang = language()
	}

	game := game(lang)
	setExits(game)
	sanityCheck(game)
	setInitialState(game)

	fmt.Println(game.Name)
	fmt.Println()
	fmt.Println(game.Description)
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(game.CurrentRoom.Name)
		fmt.Println()
		fmt.Println(game.CurrentRoom.Description)
		fmt.Println(game.CurrentRoom.GetDirections())
		fmt.Println(game.CurrentRoom.GetItemOptions())

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		fmt.Println()
		fmt.Println()
		updateState(game, input)
	}
}
