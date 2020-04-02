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

const ConfDir = "../conf/"
const SaveDir = "../saves/"

// LoadGameState restores a game state from a file into memory.
func LoadGameState(fileName string) (*Game, error) {
	path := fileName + ".yaml"
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to find save state %s", path))
	}

	var game Game
	err = yaml.Unmarshal(yamlFile, &game)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
		os.Exit(1)
	}
	game.sanityCheck()
	game.initialiseGameState()
	return &game, nil
}

// SaveGameState saves a game state to a file to be continued later.
func SaveGameState(g *Game, stateName string) error {
	g.SavedGame = true
	d, err := yaml.Marshal(g)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}
	path := SaveDir+stateName+".yaml"
	err = ioutil.WriteFile(path, d, 0644)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to write file %s", path))
	}
	fmt.Println(g.GameDictionary["strings"]["saveSucessful"])
	return nil
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
	g.CurrentRoom = g.GetRoomByID(g.CurrentRoomID)
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
// Bug(wilcox-liam): Look through this function again.
func (g *Game) UpdateGameState(input string) (*Game, error) {
	words := strings.Split(input, " ")
	words = g.expandShortcut(words)

	var command string
	var object string
	if len(words) == 0 {
		return g, errors.New(fmt.Sprintf(g.GameDictionary["errors"]["invalidCommand"], input))
	}
	command = strings.ToLower(words[0])
	if len(words) > 1 {
		object = strings.ToLower(strings.Join(words[1:], " "))
	}

	if command == strings.ToLower(g.GameDictionary["commands"]["go"]) {
		return g, g.Go(object)
	} else if command == strings.ToLower(g.GameDictionary["commands"]["examine"]) {
		return g, g.Examine(object)
	} else if command == strings.ToLower(g.GameDictionary["commands"]["refresh"]) {
		fmt.Println(g.GameDictionary["strings"]["refreshing"])
		return g, nil
	} else if command == strings.ToLower(g.GameDictionary["commands"]["inventory"]) {
		fmt.Println(g.Player.GetItemOptions())
		return g, nil
	} else if command == strings.ToLower(g.GameDictionary["commands"]["help"]) {
		fmt.Println(g.Help())
		return g, nil
	} else if command == strings.ToLower(g.GameDictionary["commands"]["save"]) {
		return g, SaveGameState(g, object)
	} else if command == strings.ToLower(g.GameDictionary["commands"]["load"]) {
		g, err := LoadGameState(SaveDir + object)
		if err == nil {
			fmt.Println(g.GameDictionary["strings"]["loadSucessful"])
		}
		return g, err
	} else if command == strings.ToLower(g.GameDictionary["commands"]["open"]) {
		return g, g.Open(object)
	} else if command == strings.ToLower(g.GameDictionary["commands"]["take"]) {
		return g, g.Take(object)
	} else if command == strings.ToLower(g.GameDictionary["commands"]["use"]) {
		return g, g.Use(object, "")
	} else {
		return g, errors.New(fmt.Sprintf(g.GameDictionary["errors"]["invalidCommand"], input))
	}
}

// PlayGame contains the game logic and game loop for playing the textgame.
// Bug(wilcox-liam): Is replacing the game from within the game loop super weird?
// TODO: consider removing this parm - check using some game data.
func (g *Game) PlayGame() {
	//Do not display the welcome text if loading a saved game
	if g.SavedGame == false {
		fmt.Println(fmt.Sprintf(g.GameDictionary["strings"]["welcome"], g.Player.Name, g.Name))
		fmt.Println()
		fmt.Println(g.GameDictionary["strings"]["helpAdvice"])
		fmt.Println()
		fmt.Println(g.Name)
		fmt.Println()
		fmt.Println(g.Description)
	}

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
		g, err = g.UpdateGameState(input)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println()
	}
}
