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

// ReadLanguages lists all languages provided by the game configuration.
func ReadLanguages() []string {
	files, err := ioutil.ReadDir(ConfDir)
	if err != nil {
		fmt.Printf("No language files found.", err)
		os.Exit(1)
	}
	var langs []string
	for _, f := range files {
		langs = append(langs, strings.Replace(f.Name(), ".yaml", "", -1))
	}
	return langs
}

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
	path := SaveDir + stateName + ".yaml"
	err = ioutil.WriteFile(path, d, 0644)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to write file %s", path))
	}
	fmt.Println(g.GameDictionary["strings"]["saveSuccessful"])
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
	//Items
	//Must be unique through the entire game
	//Must have a name
	//Must have a description
	//Cannot have the same name as a direction
	//General
}

// setInitialState initialises the game state with information that cannot be provided by the yaml configuration file.
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
		}
	}
}

// expandCommand takes a user entered shortcut and expands it into the full game command
// using the Game Dictionary provided in the yaml configuration.
func (g *Game) ExpandShortcut(word string) string {
	lookup := strings.ToLower(g.GameDictionary["shortcuts"][strings.ToLower(word)])
	if lookup != "" {
		return lookup
	}
	return word
}

// ExpandDirection takes a user entered shortcut and expands it into the full game direction
// using the Game Dictionary provided in the yaml configuration.
func (g *Game) ExpandDirection(word string) string {
	lookup := strings.ToLower(g.GameDictionary["directions"][strings.ToLower(word)])
	if lookup != "" {
		return lookup
	}
	return word
}

// ParseInput takes a user input and returns the command, object and objectTarget test
//<Command> <Object> [on <Object>]. Object names may includes spaces.
func (g *Game) ParseInput(input string) (string, string, string, error) {
	words := strings.Split(input, " ")
	if len(words) == 0 {
		return "", "", "", errors.New(fmt.Sprintf(g.GameDictionary["errors"]["invalidCommand"], input))
	}
	command := g.ExpandShortcut(words[0])

	var object string
	var objectTarget string
	if len(words) > 1 {
		object = strings.ToLower(strings.Join(words[1:], " "))
		if len(words) >= 4 && command == strings.ToLower(g.GameDictionary["commands"]["use"]) {
			objects := strings.Split(object, " on ")
			if len(objects) != 2 {
				return "", "", "", errors.New(fmt.Sprintf(g.GameDictionary["errors"]["invalidCommand"], input))
			}
			object = objects[0]
			objectTarget = objects[1]
		}
		if command == strings.ToLower(g.GameDictionary["commands"]["go"]) ||
			command == strings.ToLower(g.GameDictionary["commands"]["examine"]) {
			object = g.ExpandDirection(object)
		}
	}
	return command, object, objectTarget, nil
}

// UpdateGameState updates the game state with user provided input.
func (g *Game) UpdateGameState(input string) (*Game, error) {
	command, object, objectTarget, err := g.ParseInput(input)
	if err != nil {
		return g, err
	}

	g.DisplayRoomInfo = false
	switch command {
	case strings.ToLower(g.GameDictionary["commands"]["go"]):
		g.DisplayRoomInfo = true
		return g, g.Go(object)
	case strings.ToLower(g.GameDictionary["commands"]["examine"]):
		return g, g.Examine(object)
	case strings.ToLower(g.GameDictionary["commands"]["refresh"]):
		fmt.Println(g.GameDictionary["strings"]["refreshing"])
		g.DisplayRoomInfo = true
		return g, nil
	case strings.ToLower(g.GameDictionary["commands"]["inventory"]):
		fmt.Println(g.Player.GetItemOptions())
		return g, nil
	case strings.ToLower(g.GameDictionary["commands"]["help"]):
		fmt.Println(g.Help())
		return g, nil
	case strings.ToLower(g.GameDictionary["commands"]["save"]):
		return g, SaveGameState(g, object)
	case strings.ToLower(g.GameDictionary["commands"]["load"]):
		g, err := LoadGameState(SaveDir + object)
		if err == nil {
			fmt.Println(g.GameDictionary["strings"]["loadSuccessful"])
		}
		g.DisplayRoomInfo = true
		return g, err
	case strings.ToLower(g.GameDictionary["commands"]["quit"]):
		os.Exit(1)
	case strings.ToLower(g.GameDictionary["commands"]["open"]):
		return g, g.Open(object)
	case strings.ToLower(g.GameDictionary["commands"]["take"]):
		return g, g.Take(object)
	case strings.ToLower(g.GameDictionary["commands"]["use"]):
		return g, g.Use(object, objectTarget)
	default:
		return g, errors.New(fmt.Sprintf(g.GameDictionary["errors"]["invalidCommand"], input))
	}
	return g, nil
}

// PlayGame contains the game logic and game loop for playing the textgame.
// Bug(wilcox-liam): Is replacing the game from within the game loop super weird?
func (g *Game) Play() {
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

	var err error
	reader := bufio.NewReader(os.Stdin)
	for {
		if g.DisplayRoomInfo {
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
