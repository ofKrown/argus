package clubhouse

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"log"
	"fmt"
	"github.com/gookit/color"
	"../../configuration"
)

type colorizeDelegate func (...interface {}) string

func logError(err error) {
	if err != nil {
		log.Fatalln(err);
	}
}

// APIRequest :
type APIRequest struct {
	PageSize int `json:"page_size"`
	Query string `json:"query"`
}

// Story :
type Story struct {
	Story string `json:"name"`
	ID int `json:"id"`
	Type string  `json:"story_type"`
}

// APIResponse :
type APIResponse struct {
	Stories []Story `json:"data"`
}

// GetStories : does stuff, i hate golang
func GetStories(sortBy string, filter string) []Story {
	completeConfig := configuration.GetConfig();
	configuration := completeConfig.Clubhouse;
	// var baseurl string = "https://api.clubhouse.io/api/v3/search/stories?token=5f057780-3565-4db8-a99d-996ee967904e&page_size=25&query=owsner%3Abernhardmatz"
	baseurl := configuration.ClubhouseAPIURL + "?token=" + configuration.ClubhouseToken
	clubhouseQueryParts := []string{}
	clubhouseQueryParts = append(clubhouseQueryParts, "owner:" + configuration.ClubhouseUser)
	
	if (configuration.ClubhouseStorieState != "") {
		clubhouseQueryParts = append(clubhouseQueryParts, "state:" + configuration.ClubhouseStorieState)
	}

	if (filter != "") {
		clubhouseQueryParts = append(clubhouseQueryParts, filter)
	}
	
	if (sortBy != "") {
		clubhouseQueryParts = append(clubhouseQueryParts, "sort:" + sortBy)
	}

	fullClubHouseQuery := strings.Join(clubhouseQueryParts, " ")
	
	var query string = url.QueryEscape(fullClubHouseQuery)
	urlParts := []string{baseurl, "page_size=25", "query=" + query}
	
	fullURL := strings.Join(urlParts, "&");
	
	response, err := http.Get(fullURL)
	logError(err)
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	logError(err)

	var responseData APIResponse
	json.Unmarshal([]byte(responseBody), &responseData)

	return responseData.Stories
}

// APIListStories :
func APIListCurrentStories() {
	fmt.Println()
	stories := GetStories("changed", "state:\"In Development\"")
	for _, issue := range stories {
		var colorize colorizeDelegate
		switch issueType := issue.Type; issueType {
			case "bug":
				colorize = color.FgRed.Render
			case "feature":
				colorize = color.FgGreen.Render
			case "chore":
				colorize = color.FgBlue.Render
			default: 
				colorize = color.FgGray.Render
		}
		

		fmt.Printf("%s - %s - [%s]\n\r", colorize(fmt.Sprintf("%7s", issue.Type)), color.FgWhite.Render(fmt.Sprintf("%4d", issue.ID)), color.FgCyan.Render(issue.Story));
	}
}
