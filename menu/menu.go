package menu

import (
	"fmt"
	"github.com/gookit/color"
	"os"
	"github.com/eiannone/keyboard"
	"unicode/utf8"
	"../gateways/clubhouse"
	"../gateways/harvest"
	"../configuration"
)

func isExit(char rune, key keyboard.Key) bool {
	return key == keyboard.KeyEsc || char == '\x00'
}

func waitForKeypress() {
	keyboard.GetSingleKey()
}

func getNextKey() (char rune, key keyboard.Key, err error){
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		char, key, err = keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		if key == keyboard.KeyCtrlC || char == 'q' {
			os.Exit(0)
		}
		return
	}
}

func renderDivider(length int) {
	for i := 0; i < length; i++ {
		fmt.Print(color.FgCyan.Render("-"))
	}
	fmt.Println()
}

func renderMenu(menu *configuration.Menu, label string, isTop bool) {
	fmt.Println("")
	lineLength := 0	
	if (len(label) == 0) {
		fmt.Printf("\n%s\n", color.FgCyan.Render(menu.Name))
		lineLength = len(menu.Name)
	} else {
		fmt.Printf("\n%s\n", color.FgCyan.Render(label))
		lineLength = len(label)
	}

	renderDivider(lineLength)
	

	for _, entry := range menu.Entries {
		fmt.Printf("%s - %s \n", color.FgWhite.Render(entry.Key), entry.Name)
	}

	renderDivider(lineLength)
	if isTop {
		fmt.Println(color.FgYellow.Render("ESC/q - Quit"))
	} else {
		fmt.Println(color.FgYellow.Render("ESC - Back"))
		fmt.Println(color.FgYellow.Render("q - Quit"))
	}
	
}

func getActionMap() map[string]configuration.ActionDelegate {
	apiFunctions := make(map[string]configuration.ActionDelegate)

	apiFunctions["api:clubhouse:listCurrent"] = clubhouse.APIListCurrentStories
	apiFunctions["api:harvest:listToday"] = harvest.APIListTimeEntriesToday
	apiFunctions["api:harvest:listYesterday"] = harvest.APIListTimeEntriesYesterday
	apiFunctions["api:harvest:showCompany"] = harvest.APIShowCompany // Admin Permission required
	apiFunctions["api:harvest:showMe"] = harvest.APIShowMe
	apiFunctions["api:harvest:stopActive"] = harvest.APIStopActive
	
	return apiFunctions
}

func getActionWithArgumentMap() map[string]configuration.ActionWithArgumentDelegate {
	apiFunctions := make(map[string]configuration.ActionWithArgumentDelegate)
	
	apiFunctions["api:harvest:startTask"] = harvest.APIStartTask
	apiFunctions["api:harvest:continueMostRecentNonDaily"] = harvest.APIContinueMostRecentNonDaily
	
	return apiFunctions
}


func handleMenu(menu *configuration.Menu, parent string) {
	label := menu.Name
	if (len(parent) == 0) {
		label = menu.Name
	} else {
		label = fmt.Sprintf("%s -> %s", parent, menu.Name)
	}
	renderMenu(menu, label, len(parent) == 0)
	fmt.Print("? ")
	char, key, _ := getNextKey()
	for !isExit(char, key) {
		fmt.Printf("%c", char)
		
		selectedEntryIndex := -1
		for idx, entry := range menu.Entries {
			entryKey, _ := utf8.DecodeRune([]byte(entry.Key))
			if (entryKey == char) {
				selectedEntryIndex = idx
				break
			}
		}
		if (selectedEntryIndex >= 0) {
			selectedEntry := menu.Entries[selectedEntryIndex]
			if len(selectedEntry.Action) > 0 {
				fmt.Println()
				apiFunctions := getActionMap()
				apiFunction := apiFunctions[selectedEntry.Action]
				if apiFunction != nil {
					apiFunction()
				} else {
					apiWithArgumentFunctions := getActionWithArgumentMap()
					apiWithArgumentFunction := apiWithArgumentFunctions[selectedEntry.Action]
					if (apiWithArgumentFunction != nil) {
						apiWithArgumentFunction(selectedEntry.ActionArgument)
					} else {
						fmt.Print(color.FgRed.Render(fmt.Sprintf("\nAPI Function [%s] does not exist\n\n", selectedEntry.Action)))
						fmt.Println(color.FgWhite.Render("Available API Functions:"))
						for key := range apiFunctions {
							fmt.Printf(" - %s\n", key)
						}
						for key := range apiWithArgumentFunctions {
							fmt.Printf(" - %s\n", key)
						}

					}
				}
			}


			if(len(selectedEntry.Menu.Name) > 0) {
				handleMenu(&selectedEntry.Menu, label)
			}
		}

		renderMenu(menu, label, len(parent) == 0);
		fmt.Print("? ")
		char, _, _ = getNextKey()
	}
}

// Run : Main entry point for argus
func Run() {
	config := configuration.GetConfig()
	handleMenu(&config.Menu, "")
}