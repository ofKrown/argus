# Argus

## Implemented features

- Fully configurable and nestable menue via `json` file
- APIs
  - Harvest
    - list todays or yesterdays time entries
    - start a predefined time entry, projectId & taskId need to be provided in config
    - stop currently running time entry
    - continue latest time entry which is not equal to a provied taskId (e.g. continue latest non-daily time entry)
  - Clubhouse
    - list active stories (your user is assigned and it is optionally in a specific state)

## Upcoming features

- More documentation
- dynamic menue command aliases & user defined aliases as an argument for argus
  - e.g.: you have a menu entry for Harvest with the key `h` and a sub menue with an entry for starting a time entry with the key 's', argus would accept an argument like `argus hs` and execute the specified action
  - aliases can be assigned to all menu entries, e.g: you configure the alias `stop` for stopping the currently running Harvest time entry `argus stop` would execute the specified action
- expand Harvest API functions
- expand Clubhouse API functions
- Git integration: hey argus, create a feature branch bases on this clubhouse story
- Add Bitbucket API
- cross working APIs: hey argus list my clubhouse stories, great, start a harvest time entry for that one
- fancy stuff with slack integration: hey argus, watch the CI build of the current branch and tell me in slack when it is finished or it fails

# Argus API

- Clubhouse - requires configured `ClubhouseConfiguration`
- `api:clubhouse:listCurrent`

  - List current

- Harvest
- `api:harvest:listToday`
- `api:harvest:listYesterday`
- `api:harvest:showCompany`
- `api:harvest:showMe`
- `api:harvest:stopActive`

- `api:harvest:startTask`
- `api:harvest:continueMostRecentNonDaily`

# Config documentation

`.argus.json`

```json
{
    "harvest": {...}, // HarvestConfiguration
    "clubhouse": {...}, // ClubhouseConfiguration
    "menue": {...}, // Menu
    "debug": false, // primarily print out web request & responses
}
```

`HarvestConfiguration`

> Get your personal Harvest API access token here: https://id.getharvest.com/developers

> You can get your Harvest Account ID here: https://id.getharvest.com/accounts, your accounts will be listed, just hover over the wanted entry, the accountId will be included in the shown link

```json
{
  "harvestAPIURL": "https://api.harvestapp.com/api/v2", // preconfigured
  "harvestToken": "", // your Harvest API Token
  "harvestAccountID": "", // your Harvest Account ID
  "showDetails": false // if set to true argus will display project and task ids, this makes it easier to find wanted ids for preconfiguring actions
}
```

`ClubhouseConfiguration`

> Argus currently supports only a single clubhouse workspace. You can create a clubhouse API token here, you will need to replace `WORKSPACE` with one of your workspaces: `https://app.clubhouse.io/WORKSPACE/settings/account/api-tokens`

```json
{
  "clubhouseAPIURL": "https://api.clubhouse.io/api/v3/search/stories", // preconfigured
  "clubhouseToken": "", // your clubhouse API Token
  "clubhouseUser": "", // your clubhouse username
  "clubhouseStorieState": "" // only show stories in the specified state, e.g.: "In Development"
}
```

`Menu`

```json
{
    "name": "Menu Name",
    "entries": [{...}],
}

```

- `name`, `string`, **required**, display name of your menu/submenu
- `entries`, `Array<MenuEntry>`, **required**

`MenuEntry`

```json
{
    "name": "Menu Entry Name", // required, display name of the menu entry
    "key": "a", // required, basically every single character a keyboard can create, commonly used are a-z and 0-9
    "alias": "doit", // optional
    "action": "api:harvest:listToday", // optional, an argus api function
    "actionArgument": "12345", // optional, if the defined action requires an argument it can be defined here
    "menu": {...}, // optional, define a submenu, type `Menu`
}

```
