//Author: Liam Wilcox
//I built this basic text-adventure game as a learning exercise for golang
//Goals
//	-Multi-Language Support
//	-Multi-Player Support
//	-Coding Best Practices
//
//Known Weaknesses
//	-Technical Error messages are always in English

//TODO
//Learn how to use pointers and slices you dummy
//Add a player struct
//Add Inventory
//Update Examine for inventory
//Learn logging for golang

//Add Open Functionality (Probably Openable bool)
//Add Take Functionality
//Add Use Functionality
//Player help

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

//TODO: Log Mode
// commandLineOptions parses and returns the options provided.
func commandLineOptions() string {
	lang := flag.String("lang", "en", "Game Language")
	flag.Parse()
	return *lang
}

// TODO check for valid conf files to initialise list?
// validLanguages returns a slice of languages the game supports.
func validLanguages() []string {
	return []string{"en", "es"}
}

// languages presents the games valid languages to a user and returns the users choice.
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

// game reads in initialises game data from a yaml file.
// Bug(wilcox-liam): Error messages here are not multi-lingual.
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

// contains is a helper function to return if a string appears in a slice.
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// sanityCheck validates the game data for any obvious inconsitencies or errors.
func sanityCheck(game *textgame.Game) {
	//Exits
		//Must have a name
		//Must have a unique name per Room
		//Must contain a direction
		//Must contain a description
		//Must contain a valid Room ID
	//Rooms
		//Must have a unique id
		//Must have a name
		//Must have a description
		//Must not have multiple exits with the same 'direction'
		//1 Room must have RoomID = 1
	//Items
		//Must be unique through the entire game
		//Must have a name
		//Must have a description
}

// setInitialState initalises the game state with information that cannot be provided by the yaml configuration file.
// The game room with ID 1 is set as the Room the Player begins the game in.
func setInitialState(game *textgame.Game) {
	currentRoom := game.GetRoomByID(1)
	game.CurrentRoom = currentRoom
}

// setExits converts all RoomID's provided by the yaml configuration into a pointer to that room.
// Bug(wilcox-liam) The data is not persisting in the game state. I don't understand pointers.
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

// expandCommand takes a user entered shortcut and expands it into the full game command
// using the Game Dictionary provided in the yaml configuration.
func expandCommand(game *textgame.Game, words []string) []string {
	for i, word := range words {
		lookup := strings.ToLower(game.GameDictionary[word])
		if lookup != "" {
			words[i] = lookup
		}
	}
	return words
}

// updateState updates the game state after user provided input.
func updateState(game *textgame.Game, input string) bool {
	words := strings.Split(input, " ")
	words = expandCommand(game, words)
	if strings.ToLower(words[0]) == strings.ToLower(game.GameDictionary["commandGo"]) && len(words) == 2 {
		return updateStateGo(game, words[1])
	} else if strings.ToLower(words[0]) == strings.ToLower(game.GameDictionary["commandExamine"]) && len(words) == 2 {
		examine(game, words[1])
	} else {
		fmt.Printf(game.GameDictionary["errorInvalidCommand"], input)
	}
	return false
}

// updateStateGo handles user input commandGo.
func updateStateGo(game *textgame.Game, where string) bool {
	for _, exit := range *game.CurrentRoom.Exits {
		if strings.ToLower(exit.Direction) == strings.ToLower(where) {
			game.CurrentRoom = exit.Room
			return true
		}
	}
	fmt.Printf(game.GameDictionary["errorNoExit"], where)
	return false
}

// examine handles the user input commandExamine.
func examine(game *textgame.Game, name string) {
	desc := game.CurrentRoom.Examine(name)
	if desc == "" {
		fmt.Printf(game.GameDictionary["errorNoObject"], name)
	} else {
		fmt.Println(desc)
	}
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

	var roomChanged = true
	reader := bufio.NewReader(os.Stdin)
	for {
		//Only display Room information if the room has changed
		if roomChanged {
			fmt.Println(game.CurrentRoom.Name)
			fmt.Println()
			fmt.Println(game.CurrentRoom.Description)
			fmt.Println(game.CurrentRoom.GetDirections())
			fmt.Println(game.CurrentRoom.GetItemOptions())
		}

		fmt.Println()
		fmt.Print(game.GameDictionary["stringCommand"])
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		fmt.Println()
		roomChanged = updateState(game, input)
	}
}
