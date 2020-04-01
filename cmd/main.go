//Author: Liam Wilcox
//I built this basic text-adventure game as a learning exercise for golang
//Goals
//	-Multi-Language Support
//	-Multi-Player Support
//	-Coding Best Practices
//  -Save/Load Game State
//
//Known Weaknesses
//	-Technical Error messages are always in English

//TODO
//Learn logging for golang

//Add Take Functionality
//Add Use Functionality
//Add Close Functionality?
//Player help

//Remove player name parm and have a player object in the yaml.
//Start with something in your inventory?

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/wilcox-liam/text-game/pkg"
	"os"
	"strings"
)

const langDefault = "default"
const playerNameDefault = "default"
const saveStateDefault = "no-state"
const confDir = "../conf/"
const saveDir = "../saves/"

// TODO(wilcox-liam): Log Mode
// commandLineOptions parses and returns the options provided.
func commandLineOptions() (string, string, string) {
	lang := flag.String("lang", "en", "Game Language")
	playerName := flag.String("name", "Jazminne", "Player Name")
	saveState := flag.String("state", saveStateDefault, "Save State Name")
	flag.Parse()
	validateLanguage(*lang)
	return *lang, *playerName, *saveState
}

// TODO(wilcox-liam): check for valid conf files to initialise list?
// validLanguages returns a slice of languages the game supports.
func validLanguages() []string {
	return []string{"en", "es"}
}

// validateLanguage checks if a provided language is valid. If not the game exits.
func validateLanguage(lang string) {
	if !contains(validLanguages(), lang) {
		fmt.Println("Unknown Language")
		os.Exit(1)
	}
}

// languages presents the games valid languages to a user and returns the users choice.
func language() string {
	validLanguages := validLanguages()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Language? ", validLanguages, ": ")
	lang, _ := reader.ReadString('\n')
	lang = strings.TrimSpace(lang)
	validateLanguage(lang)
	fmt.Println()
	return lang
}

// player asks the user to enter a player name
// Bug(wilcox-liam): Should you always play as Jazminne? (and Cassandra)
func player(game *textgame.Game) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(game.GameDictionary["stringAskName"])
	name, _ := reader.ReadString('\n')
	return name
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

func main() {
	lang, playerName, saveState := commandLineOptions()

	if lang == langDefault {
		lang = language()
	}

	var game *textgame.Game
	if saveState == saveStateDefault {
		game = textgame.LoadGameState(confDir + lang + ".yaml")
	} else {
		game = textgame.LoadGameState(saveDir + saveState + ".yaml")
	}

	if playerName == playerNameDefault {
		playerName = player(game)
		fmt.Println()
	}
	game.Player = &textgame.Player{playerName, nil}

	game.PlayGame()
}
