package textgame

import (
	"bufio"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

// Bug(wilcox-liam): Error messages here are not multi-lingual.
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
	game.sanityCheck()
	game.initialiseGameState()
	return &game
}

// SaveGameState saves a game state to a file to be continued later.
func SaveGameState(g Game, stateName string) {
	//TODO
}

// sanityCheck validates the game data for any obvious inconsitencies or errors.
func (g *Game) sanityCheck() {
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
func (g *Game) initialiseGameState() {
	g.CurrentRoom = g.GetRoomByID(1)
	g.setExits()
}

// setExits converts all RoomID's provided by the yaml configuration into a pointer to that room.
func (g *Game) setExits() {
	for i, room := range g.Rooms {
		for j, exit := range room.Exits {
			g.Rooms[i].Exits[j].Room = g.GetRoomByID(exit.RoomID)
			if g.Rooms[i].Exits[j].Room == nil {
				fmt.Println("Invalid Room ID", exit.RoomID, "in exit", room.Name)
				os.Exit(1)
			}
			return
		}
	}
}

// expandCommand takes a user entered shortcut and expands it into the full game command
// using the Game Dictionary provided in the yaml configuration.
func (g *Game) expandShortcut(words []string) []string {
	for i, word := range words {
		word = strings.ToLower(word)
		lookup := strings.ToLower(g.GameDictionary["shortcuts"][word])
		if lookup != "" {
			words[i] = lookup
		}
	}
	return words
}

// Does marshal copy case or lower case it?
// updateState updates the game state with user provided input.
// returns true if the go command executed sucessfully.
// Bug(wilcox-liam): Is the bool necessary or should I check the error type instead?
// Bug(wilcox-liam): Only expand the first word. Expand the second word if first word = go
func (g *Game) UpdateGameState(input string) (bool, error) {
	words := strings.Split(input, " ")
	words = g.expandShortcut(words)

	var command string
	var object string
	if len(words) == 0 {
		return false, errors.New(fmt.Sprintf(g.GameDictionary["errors"]["invalidCommand"], input))
	}
	command = strings.ToLower(words[0])
	if len(words) > 1 {
		object = strings.ToLower(strings.Join(words[1:], " "))
	}

	if command == strings.ToLower(g.GameDictionary["commands"]["go"]) {
		return g.Go(object)
	} else if command == strings.ToLower(g.GameDictionary["commands"]["examine"]) {
		return false, g.Examine(object)
	} else if command == strings.ToLower(g.GameDictionary["commands"]["refresh"]) {
		fmt.Println(g.GameDictionary["strings"]["refreshing"])
		return true, nil
	} else if command == strings.ToLower(g.GameDictionary["commands"]["inventory"]) {
		fmt.Println(g.Player.GetItemOptions())
		return false, nil
	} else if command == strings.ToLower(g.GameDictionary["commands"]["help"]) {
		fmt.Println(g.Help())
		return false, nil
	} else if command == strings.ToLower(g.GameDictionary["commands"]["open"]) {
		return false, g.Open(object)
	} else if command == strings.ToLower(g.GameDictionary["commands"]["take"]) {
		return false, g.Take(object)
	} else if command == strings.ToLower(g.GameDictionary["commands"]["use"]) {
		return false, g.Use(object, "")
	} else {
		return false, errors.New(fmt.Sprintf(g.GameDictionary["errors"]["invalidCommand"], input))
	}
}

// PlayGame contains the game logic and game loop for playing the textgame.
func (g *Game) PlayGame() {
	fmt.Println(g.Name)
	fmt.Println()
	fmt.Println(g.Description)

	var roomChanged = true
	var err error
	reader := bufio.NewReader(os.Stdin)
	for {
		//Only display Room information if the room has changed
		if roomChanged {
			fmt.Println(g.CurrentRoom.Name)
			fmt.Println()
			fmt.Println(g.CurrentRoom.Description)
			fmt.Println(g.CurrentRoom.GetDirections())
			fmt.Println(g.CurrentRoom.GetObjectOptions())
			fmt.Println()
		}

		fmt.Print(g.GameDictionary["strings"]["command"])
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		fmt.Println()
		roomChanged, err = g.UpdateGameState(input)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println()
	}
}
