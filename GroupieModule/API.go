package GroupieModule

import (
	"encoding/json"
	"net/http"
)

type Concerts struct {
	ID    int                 `json:"id"`
	Infos map[string][]string `json:"datesLocations"`
}
type Locations struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

func FetchAllArtists(str string, grps *[]GroupInfos) bool {
	var resp, err = http.Get(str)
	if err != nil {
		return false
	}
	err = json.NewDecoder(resp.Body).Decode(&grps)
	if err != nil {
		return false
	}
	return true
}

func ArtistConcertData(locURL string, relURL string) ([]string, []string) {
	var resp, err = http.Get(locURL)
	if err != nil {
		return nil, nil
	}

	var loc Locations
	err = json.NewDecoder(resp.Body).Decode(&loc)
	if err != nil {
		return nil, nil
	}

	resp, err = http.Get(relURL)
	if err != nil {
		return nil, nil
	}

	var rel Concerts
	err = json.NewDecoder(resp.Body).Decode(&rel)
	if err != nil {
		return nil, nil
	}

	var sortedConcertPlacesList []string
	var sortedConcertDatesList []string
	for i := 0; i < len(loc.Locations); i++ {
		sortedConcertPlacesList = append(sortedConcertPlacesList, loc.Locations[i])
		sortedConcertDatesList = append(sortedConcertDatesList, rel.Infos[loc.Locations[i]][0])
	}
	return sortedConcertPlacesList, sortedConcertDatesList
}
