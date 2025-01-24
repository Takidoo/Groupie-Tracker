package GroupieModule

import (
	"encoding/json"
	"net/http"
)

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
