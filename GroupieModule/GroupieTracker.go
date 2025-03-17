package GroupieModule

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	Concerts     []ConcertData
}

type ConcertData struct {
	ConcertPlace string
	ConcertDate  string
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
		searchResult := Search(r.FormValue("userInput"), Artists, searchType)
		if r.FormValue("dateStart") != "" {
			data := SearchPageData{
				Groups: DateFilter(r.FormValue("dateStart"), r.FormValue("dateEnd"), searchResult),
			}
			tmpl, _ := template.ParseFiles("templates/search.html")
			tmpl.Execute(w, data)
		} else {
			data := SearchPageData{
				Groups: searchResult,
			}
			tmpl, _ := template.ParseFiles("templates/search.html")
			tmpl.Execute(w, data)
		}
	} else {
		data := SearchPageData{
			Groups: Artists,
		}
		tmpl, _ := template.ParseFiles("templates/search.html")
		tmpl.Execute(w, data)
	}
}

func DateFilter(dateMin string, dateMax string, infos []GroupInfos) []GroupInfos {
	layout := "2006-01-02"
	layoutAPI := "02-01-2006"
	startDate, _ := time.Parse(layout, dateMin)
	endDate, _ := time.Parse(layout, dateMax)
	var result []GroupInfos
	for i := 0; i < len(infos); i++ {
		parsedDate, _ := time.Parse(layoutAPI, infos[i].FirstAlbum)
		if parsedDate.After(startDate) && parsedDate.Before(endDate) {
			result = append(result, infos[i])
		}
	}
	return result
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
	var concertPlacesData, concertDatesData = ArtistConcertData(Artists[groupId].Locations, Artists[groupId].Relations)
	var AllConcerts []ConcertData
	for i := 0; i < len(concertPlacesData); i++ {
		AllConcerts = append(AllConcerts, ConcertData{concertPlacesData[i], concertDatesData[i]})
	}
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
		Concerts:     AllConcerts,
	}
	tmpl, _ := template.ParseFiles("templates/groupinfo.html")
	tmpl.Execute(w, data)
}

func RsrcInit() bool {
	if !FetchAllArtists("https://groupietrackers.herokuapp.com/api/artists", &Artists) {
		print("Erreur lors de la tentative de récupération des artistes")
		return false
	}
	return true
}
