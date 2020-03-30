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
//Learn logging for golang
//Learn how to use errors

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

const langDefault = "default"
const playerNameDefault = "default"

// TODO(wilcox-liam): Log Mode
// commandLineOptions parses and returns the options provided.
func commandLineOptions() (string, string) {
	lang := flag.String("lang", "en", "Game Language")
	playerName := flag.String("name", "Jazminne", "Player Name")
	flag.Parse()
	return *lang, *playerName
}

// TODO(wilcox-liam): check for valid conf files to initialise list?
// validLanguages returns a slice of languages the game supports.
func validLanguages() []string {
	return []string{"en", "es"}
}

// languages presents the games valid languages to a user and returns the users choice.
func language() string {
	validLanguages := validLanguages()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Language? ", validLanguages, ": ")
	lang, _ := reader.ReadString('\n')
	lang = strings.TrimSpace(lang)
	if !contains(validLanguages, lang) {
		fmt.Println("Unknown Language")
		os.Exit(1)
	}
	fmt.Println()
	return lang
}

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
	lang, playerName := commandLineOptions()
	if lang == langDefault {
		lang = language()
	}

	game := game(lang)
	setExits(game)
	sanityCheck(game)
	setInitialState(game)

	if playerName == playerNameDefault {
		playerName = player(game)
		fmt.Println()
	}
	game.Player = &textgame.Player{playerName, nil}

	fmt.Println(game.Name)
	fmt.Println()
	fmt.Println(game.Description)
	fmt.Println()

	roomChanged := true
	reader := bufio.NewReader(os.Stdin)
	for {
		//Only display Room information if the room has changed
		if roomChanged {
			fmt.Println(game.CurrentRoom.Name)
			fmt.Println()
			fmt.Println(game.CurrentRoom.Description)
		}
		fmt.Println(game.CurrentRoom.GetDirections())
		fmt.Println(game.CurrentRoom.GetItemOptions())

		fmt.Println()
		fmt.Print(game.GameDictionary["stringCommand"])
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		fmt.Println()
		roomChanged = updateState(game, input)
		fmt.Println()
	}
}
