package GroupieModule

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type SearchPageData struct {
	Groups []GroupInfos
}
type ErrorData struct {
	ErrorMessage string
}
type GroupInfos struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

var Artists []GroupInfos

func SearchPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var searchType, err = strconv.Atoi(r.FormValue("searchType"))
		if err != nil {
			tmpl, _ := template.ParseFiles("templates/error.html")
			tmpl.Execute(w, ErrorData{ErrorMessage: "Invalid Filter"})
			return
		}
		data := SearchPageData{
			Groups: Search(r.FormValue("userInput"), Artists, searchType),
		}
		tmpl, _ := template.ParseFiles("templates/search.html")
		tmpl.Execute(w, data)
	} else {
		data := SearchPageData{
			Groups: Artists,
		}
		tmpl, _ := template.ParseFiles("templates/search.html")
		tmpl.Execute(w, data)
	}
}

func SearchByName(input string, infos []GroupInfos) []GroupInfos {
	var output []GroupInfos
	for i := 0; i < len(infos); i++ {
		if strings.Contains(strings.ToLower(infos[i].Name), strings.ToLower(input)) {
			output = append(output, infos[i])
		}
	}
	return output
}

func SearchByAlbumName(input string, infos []GroupInfos) []GroupInfos {
	var output []GroupInfos
	for i := 0; i < len(infos); i++ {
		if strings.Contains(strings.ToLower(infos[i].FirstAlbum), strings.ToLower(input)) {
			output = append(output, infos[i])
		}
	}
	return output
}

func SearchByCreationDate(input string, infos []GroupInfos) []GroupInfos {
	var output []GroupInfos
	for i := 0; i < len(infos); i++ {
		if strings.Contains(string(infos[i].CreationDate), strings.ToLower(input)) {
			output = append(output, infos[i])
		}
	}
	return output
}

func searchByMembers(input string, infos []GroupInfos) []GroupInfos {
	var output []GroupInfos
	for i := 0; i < len(infos); i++ {
		for k := 0; k < len(infos[i].Members); k++ {
			if strings.Contains(strings.ToLower(infos[i].Members[k]), strings.ToLower(input)) {
				output = append(output, infos[i])
			}
		}
	}
	return output
}

func Search(input string, infos []GroupInfos, searchType int) []GroupInfos {
	switch searchType {
	case 1:
		return SearchByName(input, infos)
	case 2:
		return SearchByAlbumName(input, infos)
	case 3:
		return SearchByCreationDate(input, infos)
	case 4:
		return searchByMembers(input, infos)

	}
	return infos
}
func InfoPage(w http.ResponseWriter, r *http.Request) {
	var groupId, err = strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || groupId > len(Artists) || groupId < 0 {
		tmpl, _ := template.ParseFiles("templates/error.html")
		tmpl.Execute(w, ErrorData{ErrorMessage: "Invalid ID"})
		return
	}
	groupId += -1
	data := GroupInfos{
		ID:           Artists[groupId].ID,
		Image:        Artists[groupId].Image,
		Name:         Artists[groupId].Name,
		Members:      Artists[groupId].Members,
		CreationDate: Artists[groupId].CreationDate,
		FirstAlbum:   Artists[groupId].FirstAlbum,
		Locations:    Artists[groupId].Locations,
		ConcertDates: Artists[groupId].ConcertDates,
		Relations:    Artists[groupId].Relations,
	}
	tmpl, _ := template.ParseFiles("templates/groupinfo.html")
	tmpl.Execute(w, data)
}

func RsrcInit() bool {
	if !FetchAllArtists("https://groupietrackers.herokuapp.com/api/artists", &Artists) {
		print("Erreur lors de la tentative de récupération des artistes")
		return false
	}
	//petite folie
	Artists = append(Artists, GroupInfos{ID: 53, Image: "https://images.mubicdn.net/images/film/44857/cache-34811-1445896251/image-w1280.jpg?size=800x", Name: "Michael Jackson", Members: []string{"Michael Jackson"}, CreationDate: 1958, FirstAlbum: "Off The Wall (Les autres ne comptent pas :)", Locations: "", ConcertDates: "", Relations: ""})

	return true
}
