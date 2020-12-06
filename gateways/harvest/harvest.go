package harvest

import (
	"fmt"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
	"strings"
	"strconv"
	"github.com/gookit/color"
	"../../configuration"
	"../../util"
)

// User : 
type User struct {
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	ID int64 `json:"id"`
	EMail string `json:"email"`
	Timezone string `json:"timezone"`
}

func (me User) String() string {
	return fmt.Sprintf("%s %s (%s) UserId: %d", me.FirstName, me.LastName, me.EMail, me.ID)
}

// Project : 
type Project struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
}
// Task : 
type Task struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
}

// TimeEntry : 
type TimeEntry struct {
	ID int64 `json:"id"`
	Project Project `json:"project"`
	Task Task `json:"task"`
	Hours float32 `json:"hours"`
	RoundedHours float32 `json:"rounded_hours"`
	Notes string `json:"notes"`
	StartTime string `json:"started_time"`
	EndTime string `json:"ended_time"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	SpentDate string `json:"spent_date"`
	IsRunning bool `json:"is_running"`
}

func (me TimeEntry) String() string {
	showDetails := configuration.GetConfig().Harvest.ShowDetails
	timeEntryInfo := ""

	timeRange := fmt.Sprintf("%5s - %5s / %2.2f Hours",me.StartTime, me.EndTime, me.Hours)
	if (me.IsRunning) {
		timeRange = fmt.Sprintf("%5s %s / %2.2f Hours",me.StartTime, color.FgLightBlue.Render(fmt.Sprintf("%7s", "Active")), me.Hours)
	}
	task := color.FgYellow.Render(fmt.Sprintf("%s", me.Task.Name))
	project := color.FgGreen.Render(fmt.Sprintf("[%s]", me.Project.Name))
	if (showDetails) {
		timeEntryInfo = color.FgWhite.Render(fmt.Sprintf("(TimeEntryId: %d) ", me.ID))
		task = color.FgYellow.Render(fmt.Sprintf("%s (TaskID: %d)", me.Task.Name, me.Task.ID))
		project = color.FgGreen.Render(fmt.Sprintf("[%s (ProjectID: %d)]", me.Project.Name, me.Project.ID))
	}
	if me.Project.ID == 19610249 {
		project = color.FgWhite.Render(fmt.Sprintf("[%s]", me.Project.Name))
		if (showDetails) {
			project = color.FgWhite.Render(fmt.Sprintf("[%s (ProjectID: %d)]", me.Project.Name, me.Project.ID))
		}
	}
	return fmt.Sprintf("%s%s: %s %s", timeEntryInfo, timeRange, task, project)
}

// TimeEntriesResponse : 
type TimeEntriesResponse struct {
	TimeEntries []TimeEntry `json:"time_entries"`
}

// StartTaskDTO :
type StartTaskDTO struct {
	ProjectID int64 `json:"project_id"`
	TaskID int64 `json:"task_id"`
	SpentDate string `json:"spent_date"`
}

// CompanyResponse :
type CompanyResponse struct {
	BaseURI string `json:"base_uri"`
	FullDomain string `json:"full_domain"`
	Name string `json:"name"`
	UseTimestamps bool `json:"wants_timestamp_timers"`
}

func (me CompanyResponse) String() string {
	name := color.FgGreen.Render(me.Name)
	url := color.FgWhite.Render(me.FullDomain)
	timestamps := "No"
	if (me.UseTimestamps) {
		timestamps = "Yes"
	}
	timestampsColored := color.FgWhite.Render(timestamps)
	return fmt.Sprintf("[%s] %s, Timestamps: %s", name, url, timestampsColored)
}

// ******* //
// Private //
// ******* //

func getHarvestBaseQuery() url.Values {
	config := configuration.GetConfig().Harvest
	var query url.Values = url.Values{}
	query.Add("access_token", config.HarvestToken)
	query.Add("account_id", config.HarvestAccountID)

	return query
}

func getMe() (user User){
	completeConfig := configuration.GetConfig()
	config := completeConfig.Harvest

	query := getHarvestBaseQuery()
	requestURL := fmt.Sprintf("%s/users/me?%s", config.HarvestAPIURL, query.Encode())
	
	responseBody := util.DoHTTPRequest(http.MethodGet, requestURL, nil)
	
	json.Unmarshal([]byte(responseBody), &user)
	util.DebugResponseData(fmt.Sprintf("%s", user))

	return
}

func getMyEntries(userID int64, startDate string) (timeEntriesResponse TimeEntriesResponse){
	completeConfig := configuration.GetConfig()
	config := completeConfig.Harvest

	query := getHarvestBaseQuery()
	query.Add("user_id", fmt.Sprintf("%d", userID))
	query.Add("from", startDate)

	requestURL := fmt.Sprintf("%s/time_entries?%s", config.HarvestAPIURL, query.Encode())

	responseBody := util.DoHTTPRequest(http.MethodGet, requestURL, nil)

	json.Unmarshal([]byte(responseBody), &timeEntriesResponse)
	util.DebugResponseData(fmt.Sprintf("%s", timeEntriesResponse))

	return
}

func getCompany() (companyResponse CompanyResponse) {
	completeConfig := configuration.GetConfig()
	config := completeConfig.Harvest
	
	query := getHarvestBaseQuery()
	requestURL := fmt.Sprintf("%s/time_entries?%s", config.HarvestAPIURL, query.Encode())

	responseBody := util.DoHTTPRequest(http.MethodGet, requestURL, nil)
	
	json.Unmarshal([]byte(responseBody), &companyResponse)
	util.DebugResponseData(fmt.Sprintf("%s", companyResponse))

	return
}

func startTimeEntry(projectID int64, taskID int64) (timeEntry TimeEntry) {
	start := time.Now()
	startDate := getFormatedDate(start)
	
	completeConfig := configuration.GetConfig()
	config := completeConfig.Harvest

	query := getHarvestBaseQuery()
	requestURL := fmt.Sprintf("%s/time_entries?%s", config.HarvestAPIURL, query.Encode())

	startTaskBody := &StartTaskDTO {
		ProjectID: projectID,
		TaskID: taskID,
		SpentDate: startDate,
	}

	requestBody, err := json.Marshal(startTaskBody)
	util.LogError(err)

	responseBody := util.DoHTTPRequest(http.MethodPost, requestURL, requestBody)

	json.Unmarshal([]byte(responseBody), &timeEntry)
	util.DebugResponseData(fmt.Sprintf("%s", timeEntry))

	return
}

func doTimeEntryPatch(action string, timeEntryID int64) (timeEntry TimeEntry) {
	completeConfig := configuration.GetConfig()
	config := completeConfig.Harvest
	query := getHarvestBaseQuery()
	requestURL := fmt.Sprintf("%s/time_entries/%d/%s?%s", config.HarvestAPIURL, timeEntryID, action, query.Encode())

	responseBody := util.DoHTTPRequest(http.MethodPatch, requestURL, nil)
		
	json.Unmarshal([]byte(responseBody), &timeEntry)
	util.DebugResponseData(fmt.Sprintf("%s", timeEntry))

	return 
}

func getFormatedDate(date time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d", date.Year(), date.Month(), date.Day())
}

// ****** //
// Public //
// ****** //

// APIShowMe :
func APIShowMe() {
	me := getMe()
	fmt.Println(me)
}

// APIShowCompany :
func APIShowCompany() {
	company := getCompany()
	fmt.Println(company)
}

// APIListTimeEntriesToday :
func APIListTimeEntriesToday() {
	APIListTimeEntries(time.Now())
}

// APIListTimeEntriesYesterday :
func APIListTimeEntriesYesterday() {
	APIListTimeEntries(time.Now().AddDate(0, 0, -1))
}

// APIListTimeEntries :
func APIListTimeEntries(start time.Time) {
	fmt.Println()
	me := getMe()
	startDate := getFormatedDate(start)
	timeEntries := getMyEntries(me.ID, startDate)
	for i := len(timeEntries.TimeEntries) -1; i >= 0; i-- {
		fmt.Println(timeEntries.TimeEntries[i])
	}
	if len(timeEntries.TimeEntries) == 0 {
		fmt.Println(color.FgYellow.Render("No time entries for today, yet..."))
	}
}

// APIStartTask :
func APIStartTask(argument string) {
	args := strings.Split(argument, ":")
	projectID, err := strconv.ParseInt(args[0], 10, 64)
	util.LogError(err)
	taskID, err := strconv.ParseInt(args[1], 10, 64)
	util.LogError(err)
	
	timeEntry := startTimeEntry(projectID, taskID)
	fmt.Println(fmt.Sprintf("%s: %s", color.FgGreen.Render("Started"), timeEntry))
}

// APIContinueMostRecentNonDaily :
func APIContinueMostRecentNonDaily(argument string) {
	fmt.Println()

	dailyTaskID, err := strconv.ParseInt(argument, 10, 64)
	util.LogError(err)
	
	me := getMe()
	start := time.Now()
	startDate := getFormatedDate(start)

	timeEntries := getMyEntries(me.ID, startDate)

	if len(timeEntries.TimeEntries) == 0 {
		fmt.Println(color.FgYellow.Render("No time entries for today, yet..."))
	} else {
		for i := 0; i < len(timeEntries.TimeEntries);  i++ {
			timeEntry := timeEntries.TimeEntries[i]
			if (timeEntry.Task.ID != dailyTaskID) {
				newTimeEntry := doTimeEntryPatch("restart", timeEntry.ID)
				fmt.Println(fmt.Sprintf("%s: %s", color.FgGreen.Render("Restarted"), newTimeEntry))
				break
			}
		}
	}
}

// APIDOTimeEntryPatch : 
func APIDOTimeEntryPatch(arguments string) {
	args := strings.Split(arguments, ":")
	timeEntryID, err := strconv.ParseInt(args[1], 10, 64)
	util.LogError(err)
	doTimeEntryPatch(args[0], timeEntryID)
}

// APIStopActive :
func APIStopActive() {
	fmt.Println()

	me := getMe()
	start := time.Now()
	startDate := getFormatedDate(start)

	timeEntries := getMyEntries(me.ID, startDate)

	if len(timeEntries.TimeEntries) == 0 {
		fmt.Println(color.FgYellow.Render("No time entries for today, yet..."))
	} else {
		timeEntry := timeEntries.TimeEntries[0]
		if (timeEntry.IsRunning) {
			newTimeEntry := doTimeEntryPatch("stop", timeEntry.ID)
			fmt.Println(fmt.Sprintf("%s: %s", color.FgGreen.Render("Stopped"), newTimeEntry))
		} else {
			fmt.Println(fmt.Sprintf("%s: %s", color.FgGreen.Render("Not Active"), timeEntry))
		}
		

	}
}