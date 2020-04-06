// Package textgame provides data structures and functions to support
// development of Text Adventure Games.
package textgame

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

// ConfDir is the directory where new game configurations are stored.
const ConfDir = "../conf/"

// SaveDir is the directory where save games are stored.
const SaveDir = "../saves/"

// ReadLanguages lists all languages provided by the Game configuration.
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

// LoadGameState restores a Game state from a file into memory.
func LoadGameState(fileName string) (*Game, error) {
	path := fileName + ".yaml"
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to find save state %s", path)
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

// saveGameState saves a Game state to a file to be continued later.
func saveGameState(g *Game, stateName string) error {
	g.SavedGame = true
	d, err := yaml.Marshal(g)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}
	path := SaveDir + stateName + ".yaml"
	err = ioutil.WriteFile(path, d, 0644)
	if err != nil {
		return fmt.Errorf("Unable to write file %s", path)
	}
	fmt.Println(g.Dictionary["strings"]["saveSuccessful"])
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

// setInitialState initialises the Game state with information that cannot
// be provided by the yaml configuration file.
func (g *Game) initialiseGameState() {
	g.setCurrentRoom(g.getRoomByID(g.CurrentRoomID))
}

// expandCommand takes a user entered shortcut and expands it into the full game command
// using the Game Dictionary provided in the yaml configuration.
func (g *Game) expandShortcut(word string) string {
	lookup := strings.ToLower(g.Dictionary["shortcuts"][strings.ToLower(word)])
	if lookup != "" {
		return lookup
	}
	return word
}

// expandDirection takes a user entered shortcut and expands it into the full game direction
// using the Game Dictionary provided in the yaml configuration.
func (g *Game) expandDirection(word string) string {
	lookup := strings.ToLower(g.Dictionary["directions"][strings.ToLower(word)])
	if lookup != "" {
		return lookup
	}
	return word
}

// parseInput takes a user input and returns the command, Item and object
//<Command> <Object> [on <Object>]. Object names may includes spaces.
func (g *Game) parseInput(input string) (string, string, string, error) {
	words := strings.Split(input, " ")
	if len(words) == 0 {
		return "", "", "", fmt.Errorf(g.Dictionary["errors"]["invalidCommand"], input)
	}
	command := g.expandShortcut(words[0])

	var object string
	var objectTarget string
	if len(words) > 1 {
		object = strings.ToLower(strings.Join(words[1:], " "))
		if len(words) >= 4 && command == strings.ToLower(g.Dictionary["commands"]["use"]) {
			objects := strings.Split(object, " on ")
			if len(objects) != 2 {
				return "", "", "", fmt.Errorf(g.Dictionary["errors"]["invalidCommand"], input)
			}
			object = objects[0]
			objectTarget = objects[1]
		}
		if command == strings.ToLower(g.Dictionary["commands"]["go"]) ||
			command == strings.ToLower(g.Dictionary["commands"]["examine"]) {
			object = g.expandDirection(object)
		}
	}
	return command, object, objectTarget, nil
}

// updateGameState updates the game state with user provided input.
func (g *Game) updateGameState(input string) (*Game, error) {
	command, object, objectTarget, err := g.parseInput(input)
	if err != nil {
		return g, err
	}

	g.DisplayRoomInfo = false
	switch command {
	case strings.ToLower(g.Dictionary["commands"]["go"]):
		return g, g.goDirection(object)
	case strings.ToLower(g.Dictionary["commands"]["examine"]):
		return g, g.examine(object)
	case strings.ToLower(g.Dictionary["commands"]["refresh"]):
		fmt.Println(g.Dictionary["strings"]["refreshing"])
		g.DisplayRoomInfo = true
		return g, nil
	case strings.ToLower(g.Dictionary["commands"]["inventory"]):
		fmt.Println(g.Dictionary["strings"]["inventory"] + g.Player.getItemOptions())
		return g, nil
	case strings.ToLower(g.Dictionary["commands"]["help"]):
		fmt.Println(g.help())
		return g, nil
	case strings.ToLower(g.Dictionary["commands"]["save"]):
		return g, saveGameState(g, object)
	case strings.ToLower(g.Dictionary["commands"]["load"]):
		g, err := LoadGameState(SaveDir + object)
		if err == nil {
			fmt.Println(g.Dictionary["strings"]["loadSuccessful"])
		}
		return g, err
	case strings.ToLower(g.Dictionary["commands"]["quit"]):
		os.Exit(1)
	case strings.ToLower(g.Dictionary["commands"]["open"]):
		return g, g.open(object)
	case strings.ToLower(g.Dictionary["commands"]["take"]):
		return g, g.take(object)
	case strings.ToLower(g.Dictionary["commands"]["use"]):
		return g, g.use(object, objectTarget)
	default:
		return g, fmt.Errorf(g.Dictionary["errors"]["invalidCommand"], input)
	}
	return g, nil
}

// Play contains the game logic and game loop for playing the textgame.
// Bug(wilcox-liam): Is replacing the game from within the game loop super weird?
func (g *Game) Play() {
	//Do not display the welcome text if loading a saved game
	if g.SavedGame == false {
		fmt.Println(fmt.Errorf(g.Dictionary["strings"]["welcome"], g.Player.Name, g.Name))
		fmt.Println()
		fmt.Println(g.Dictionary["strings"]["helpAdvice"])
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
			fmt.Println(g.Dictionary["strings"]["directions"] + g.CurrentRoom.getDirections())
			fmt.Println(g.Dictionary["strings"]["exits"] + g.CurrentRoom.getExitOptions())
			fmt.Println(g.Dictionary["strings"]["items"] + g.CurrentRoom.getItemOptions())
			fmt.Println()
		}

		fmt.Print(g.Dictionary["strings"]["command"])
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		fmt.Println()
		g, err = g.updateGameState(input)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println()
	}
}
