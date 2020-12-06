package clubhouse

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"log"
	"fmt"
	"github.com/gookit/color"
	"../../configuration"
	"../../util"
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

	responseBody := util.DoHTTPRequest(http.MethodGet, fullURL, nil)

	var responseData APIResponse
	json.Unmarshal([]byte(responseBody), &responseData)
	util.DebugResponseData(fmt.Sprintf("%s", responseData))

	return responseData.Stories
}

// APIListCurrentStories :
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
