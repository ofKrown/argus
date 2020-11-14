package configuration

import (
	"os"
	"path"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/gookit/color"
)

// ActionDelegate : 
type ActionDelegate func()
// ActionWithArgumentDelegate :
type ActionWithArgumentDelegate func(argument string)

// Configuration : Configuration
type Configuration struct {
	Harvest HarvestConfiguration `json:"harvest"`
	Clubhouse ClubhouseConfiguration `json:"clubhouse"`
	Menu Menu `json:"menu"`
}

// HarvestConfiguration :
type HarvestConfiguration struct {
	HarvestAPIURL string `json:"harvestAPIURL"`
	HarvestToken string `json:"harvestToken"`
	HarvestAccountID string `json:"harvestAccountID"`
	ShowDetails bool `json:"showDetails"`
}

// ClubhouseConfiguration :
type ClubhouseConfiguration struct {
	ClubhouseAPIURL string `json:"clubhouseAPIURL"`
	ClubhouseToken string `json:"clubhouseToken"`
	ClubhouseUser string `json:"clubhouseUser"`
	ClubhouseStorieState string `json:"clubhouseStorieState"`
}

// MenuEntry :
type MenuEntry struct {
	Name string `json:"name"`
	Key string `json:"key"`
	Alias string `json:"alias"`
	Action string `json:"action"`
	ActionArgument string `json:"actionArgument"`
	Menu Menu `json:"menu"`
}

// Menu :
type Menu struct {
	Name string `json:"name"`
	Entries []MenuEntry `json:"entries"`
}

const configFileName =  ".argus.json";

func getConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return path.Join(homeDir, configFileName);
}
func getDefaultConfig() *Configuration {
	return &Configuration{
		Clubhouse: ClubhouseConfiguration{
			ClubhouseAPIURL: "https://api.clubhouse.io/api/v3/search/stories",
			ClubhouseStorieState: "In Development",
		},
		Harvest: HarvestConfiguration {
			HarvestAPIURL: "https://api.harvestapp.com/api/v2",
			HarvestToken: "Get your Personal Access Token here: https://id.getharvest.com/developers",
		},
		Menu: Menu {
			Name: "Argus",
			Entries: []MenuEntry{
				MenuEntry{
					Name: "Harvest",
					Key: "h",
					Menu: Menu{
						Name: "Harvest",
						Entries: []MenuEntry{
							MenuEntry{
								Name: "Show me",
								Key: "m",
								Action: "api:harvest:showMe",
							},
							MenuEntry{
								Name: "List Today",
								Key: "l",
								Action: "api:harvest:listToday",
							},
							MenuEntry{
								Name: "List Yesterday",
								Key: "y",
								Action: "api:harvest:listYesterday",
							},
							MenuEntry{
								Name: "Start Daily",
								Key: "d",
								Alias: "daily",
								Action: "api:harvest:startTask",
								ActionArgument: "HARVEST_PROJECT_ID_HERE:HARVEST_TASK_ID_HERE",
							},
							MenuEntry{
								Name: "Continue most recent non-daily",
								Key: "c",
								Alias: "continue,cnt",
								Action: "api:harvest:continueMostRecentNonDaily",
								ActionArgument: "HARVEST_TASK_ID_HERE",
							},
							MenuEntry{
								Name: "Stop currently running",
								Key: "s",
								Alias: "stop",
								Action: "api:harvest:stopActive",
							},
						},
					},
				},
			},
		},
	}
}

func printNewHere() {
	argusName := color.FgCyan.Render("Argus")
	argusLocation := color.FgCyan.Render("'~/.argus.json'")
	argusDoc := color.FgYellow.Render("http://github.com/bmatz/argus")

	fmt.Printf("\n%s Seems like you ran %s for the first time\n", color.FgGreen.Render("\nYeah!!!..."), argusName)
	fmt.Printf("%s created a configuration file for you.\n", argusName)
	fmt.Printf("You can find your new configuration file here -> %s\n", argusLocation)
	fmt.Printf("\nYou will need to set up some API tokens for each service you want %s to access.\n", argusName)
	fmt.Println("The configuration is seperated in sections based on the supported APIs.")
	fmt.Printf("E.g.: If you want %s to be able to list and handle your clubhouse stories you have to at least provide a cloubhouse API token.\n", argusName)
	fmt.Printf("\nPlease consider having a look in to %s config documentation which will explain everything in detail.\n", argusName)
	fmt.Printf("The documentation can be found here: %s\n", argusDoc);
}

func readConfig() *Configuration {
	configPath := getConfigPath();
	buffer, err := ioutil.ReadFile(configPath);
	if (err != nil) {
		fmt.Println("No configuration found, creating one...");
		newConfig := getDefaultConfig()
		jsonText, err := json.MarshalIndent(&newConfig, "", "\t");
		if (err != nil) {
			panic(err)
		}
		err = ioutil.WriteFile(configPath, jsonText, 0644)
		if (err != nil) {
			panic(err)
		}
		printNewHere()	
		os.Exit(0);
		
	}
	var configFromFile Configuration
	json.Unmarshal([]byte(buffer), &configFromFile)
	return &configFromFile
}

// GetConfig : 
func GetConfig() *Configuration {
	return readConfig()
}