package steamcat

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"encoding/json"

	"github.com/andygrunwald/vdf"
)

// getAppsRes defines the response structure from the GET request for the global Steam Library
type getAppsRes struct {
	Applist struct {
		Apps []struct {
			ID   int    `json:"appid"`
			Name string `json:"name"`
		} `json:"apps"`
	} `json:"appList"`
}

type nameMap map[int]string

type vdfMap map[string]interface{}

// Game defines the necessary(ish) attributes needed for this to work
type Game struct {
	ID   int
	Name string
	Tags []string
}

// Library defines a collection of pointers to Games
type Library []*Game

// AppMap provides the mapping of Ids to Names for games in the global Steam Library.
var AppMap nameMap

var buff bytes.Buffer
var logger = log.New(&buff, "steamcat: ", log.Lshortfile)

// LoadVdfFrom takes a filepath and returns a Map allowing the key:pair values to
// be recovered and used.
func loadVDFFile(path string) (map[string]interface{}, error) {

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	parser := vdf.NewParser(f)
	vdfMap, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	return vdfMap, err
}

// Update refreshes the map of Ids to Names for the global Steam library
func (m nameMap) Update() error {
	jsonData := new(getAppsRes)
	endPoint := "https://api.steampowered.com/ISteamApps/GetAppList/v2/"
	res, err := http.Get(endPoint)
	if err != nil {
		return err
	}

	jsonByteData, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		logger.Println(err)
		return err
	}

	err = json.Unmarshal(jsonByteData, &jsonData)
	if err != nil {
		logger.Println(err)
		return err
	}

	for _, game := range jsonData.Applist.Apps {
		m[game.ID] = game.Name
	}

	return nil
}

func (g Game) hasTag(s string) bool {
	for _, tag := range g.Tags {
		if tag == s {
			return true
		}
	}
	return false
}

// TaggedWith returns a subset of the library's games which match a specified Tag
func (lib Library) TaggedWith(s string) Library {
	var filteredGames Library
	for i, thisGame := range lib {
		if thisGame.hasTag(s) {
			filteredGames = append(filteredGames, lib[i])
		}
	}
	return filteredGames
}

// RandomGameFrom returns a random game from the given library
func RandomGameFrom(lib Library) *Game {
	return lib[rand.Intn(len(lib))]
}

// GenerateLibraryFrom builds a library using the given file location.
func GenerateLibraryFrom(f string) (Library, error) {
	AppMap.Update()

	var lib Library

	vdfContents, err := loadVDFFile(f)
	if err != nil {
		return nil, err
	}

	vdfConfig, ok := vdfContents["UserRoamingConfigStore"].(vdfMap)
	vdfSoftware, ok := vdfConfig["Software"].(vdfMap)
	vdfValve, ok := vdfSoftware["Valve"].(vdfMap)
	vdfSteam, ok := vdfValve["Steam"].(vdfMap)
	vdfGames, ok := vdfSteam["apps"].(vdfMap)
	if !ok {
		vdfFormatError := fmt.Errorf("Assertion Error: Could not find \"apps\" at UserRoamingConfigStore.Software.Valve.Steam.apps")
		return nil, vdfFormatError
	}

	for k, games := range vdfGames {
		i, err := strconv.Atoi(k)
		if err != nil {
			log.Fatal(err)
		}
		info, ok := games.(vdfMap)
		if !ok {
			log.Printf("Game info format is incorrect:\nID: %d\nValue: %v", i, games)
		}
		var tags []string
		for _, attrValue := range info {
			if attrValue == "tags" {
				tagsValue, ok := attrValue.(vdfMap)
				if !ok {
					log.Printf("Tag format is incorrect: %d, %v", i, attrValue)
				} else {
					for _, tagValue := range tagsValue {
						tag, ok := tagValue.(string)
						if !ok {
							log.Printf("Tag value is incorrect: %d, %v", i, tagValue)
						}
						tags = append(tags, tag)
					}
				}
			}
		}
		lib = append(lib, &Game{ID: i, Name: AppMap[i], Tags: tags})
	}
	return lib, nil
}
