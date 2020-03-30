package textgame

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

// Bug(wilcox-liam): Error messages here are not multi-lingual.
// Bug(wilcox-liam): Call sanityCheck and Initailise from here?
func LoadGameState(fileName string) *Game {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		os.Exit(1)
	}

	var game Game
	err = yaml.Unmarshal(yamlFile, &game)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
		os.Exit(1)
	}
	return &game
}

// SaveGameState saves a game state to a file to be continued later.
func SaveGameState(g Game, stateName string) {
	//TODO
}

// sanityCheck validates the game data for any obvious inconsitencies or errors.
func (g Game) sanityCheck() {
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
func (g Game) InitialiseGameState() {
	g.CurrentRoom = g.GetRoomByID(1)
	g.setExits()
}

// setExits converts all RoomID's provided by the yaml configuration into a pointer to that room.
// Bug(wilcox-liam) The data is not persisting in the game state. I don't understand pointers.
func (g Game) setExits() {
	for _, room := range *g.Rooms {
		for _, exit := range *room.Exits {
			exit.Room = g.GetRoomByID(exit.RoomID)
			if exit.Room == nil {
				fmt.Println("Invalid Room ID", exit.RoomID, "in exit", room.Name)
				os.Exit(1)
			}
			return
		}
	}
}

// Does marshal copy case or lower case it?
// updateState updates the game state with user provided input.
func (g Game) UpdateGameState(input string) error {
	words := strings.Split(input, " ")
	words = g.expandCommand(words)

	var command string
	var object string
	if len(words) == 0 {
		return errors.New(fmt.Sprintf(g.GameDictionary["errorInvalidCommand"], input))
	}
	command = strings.ToLower(words[0])
	if len(words) > 1 {
		object = strings.ToLower(strings.Join(words[1:], " "))
	}

	if command == strings.ToLower(g.GameDictionary["commandGo"]) {
		return g.Go(object)
	} else if command == strings.ToLower(g.GameDictionary["commandExamine"]) {
		return g.Examine(object)
	} else if command == strings.ToLower(g.GameDictionary["commandInventory"]) {
		fmt.Println(g.Player.GetItemOptions())
		return nil
	} else if command == strings.ToLower(g.GameDictionary["commandOpen"]) {
		return g.Open(object)
	} else {
		return errors.New(fmt.Sprintf(g.GameDictionary["errorInvalidCommand"], input))
	}
}
