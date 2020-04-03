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

//Take doesnt remove the item from the room
//Use doesnt support target items

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
const saveStateDefault = "no-state"

// TODO(wilcox-liam): Log Mode
// commandLineOptions parses and returns the options provided.
func commandLineOptions() (string, string) {
	lang := flag.String("lang", "en", "Game Language")
	saveState := flag.String("state", saveStateDefault, "Save State Name")
	flag.Parse()
	if *lang != langDefault {
		validateLanguage(*lang)
	}
	return *lang, *saveState
}

// validateLanguage checks if a provided language is valid. If not the game exits.
func validateLanguage(lang string) {
	if !contains(textgame.ReadLanguages(), lang) {
		fmt.Println("Unknown Language")
		os.Exit(1)
	}
}

// languages presents the games valid languages to a user and returns the users choice.
func language() string {
	validLanguages := textgame.ReadLanguages()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Language? ", validLanguages, ": ")
	lang, _ := reader.ReadString('\n')
	lang = strings.TrimSpace(lang)
	validateLanguage(lang)
	fmt.Println()
	return lang
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
	lang, saveState := commandLineOptions()
	if lang == langDefault {
		lang = language()
	}

	var game *textgame.Game
	var err error
	if saveState == saveStateDefault {
		game, err = textgame.LoadGameState(textgame.ConfDir + lang)
	} else {
		game, err = textgame.LoadGameState(textgame.SaveDir + saveState)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	game.Play()
}
