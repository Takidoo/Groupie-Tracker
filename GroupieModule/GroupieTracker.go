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
type GroupInfos struct { // Structure qui stocke les informations d'un groupe
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

type ConcertData struct { // Structure qui stocke les données d'un concert
	ConcertPlace string
	ConcertDate  string
}

var Artists []GroupInfos // Liste de tous les Artistes

func SearchPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		searchResult := Search(r.FormValue("userInput"), Artists) // Appelle la fonction qui gère la recherche
		filterType := r.FormValue("filterType")
		if filterType != "" { // Si les tris sont activés, alors on trie
			filterTypeint, _ := strconv.Atoi(filterType)
			searchResult = ArtistsSort(filterTypeint, searchResult)
		}
		data := SearchPageData{
			Groups: searchResult,
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

func ArtistsSort(filterType int, groups []GroupInfos) []GroupInfos {
	switch filterType {
	case 1: // Plus anciens
		n := len(groups)
		for i := 0; i < n-1; i++ {
			for j := 0; j < n-i-1; j++ {
				if groups[j].CreationDate > groups[j+1].CreationDate {
					groups[j], groups[j+1] = groups[j+1], groups[j]
				}
			}
		}
		return groups
		break
	case 2: // Plus récents
		n := len(groups)
		for i := 0; i < n-1; i++ {
			for j := 0; j < n-i-1; j++ {
				if groups[j].CreationDate < groups[j+1].CreationDate {
					groups[j], groups[j+1] = groups[j+1], groups[j]
				}
			}
		}
		return groups
		break
	case 3: // Nombre de membre >
		n := len(groups)
		for i := 0; i < n-1; i++ {
			for j := 0; j < n-i-1; j++ {
				if len(groups[j].Members) > len(groups[j+1].Members) {
					groups[j], groups[j+1] = groups[j+1], groups[j]
				}
			}
		}
		return groups
		break
	case 4: // Nombre de membres <
		n := len(groups)
		for i := 0; i < n-1; i++ {
			for j := 0; j < n-i-1; j++ {
				if len(groups[j].Members) < len(groups[j+1].Members) {
					groups[j], groups[j+1] = groups[j+1], groups[j]
				}
			}
		}
		return groups
	}
	return groups

}

func Search(input string, infos []GroupInfos) []GroupInfos { // Fonction de recherche qui gère tous les filtres
	var output []GroupInfos
	for i := 0; i < len(infos); i++ {
		if strings.Contains(string(infos[i].CreationDate), strings.ToLower(input)) { // Date de création
			output = append(output, infos[i])
			continue
		} else if strings.Contains(strings.ToLower(infos[i].FirstAlbum), strings.ToLower(input)) { // Date premier album
			output = append(output, infos[i])
			continue
		} else if strings.Contains(strings.ToLower(infos[i].Name), strings.ToLower(input)) { // Nom du groupe
			output = append(output, infos[i])
			continue
		} else {
			for k := 0; k < len(infos[i].Members); k++ { // Membres du groupe
				if strings.Contains(strings.ToLower(infos[i].Members[k]), strings.ToLower(input)) {
					output = append(output, infos[i])
					break
				}
			}
		}
	}
	return output
}
func InfoPage(w http.ResponseWriter, r *http.Request) {
	var groupId, err = strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || groupId > len(Artists) || groupId < 0 {
		tmpl, _ := template.ParseFiles("templates/error.html")
		tmpl.Execute(w, ErrorData{ErrorMessage: "Invalid ID"})
		return
	}
	groupId += -1
	var concertPlacesData, concertDatesData = ArtistConcertData(Artists[groupId].Locations, Artists[groupId].Relations) // récupération des concerts du groupe
	var AllConcerts []ConcertData
	for i := 0; i < len(concertPlacesData); i++ { // Création de la liste des concerts
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

func RsrcInit() bool { // Vérification des prérequis
	if !FetchAllArtists("https://groupietrackers.herokuapp.com/api/artists", &Artists) {
		print("Erreur lors de la tentative de récupération des artistes")
		return false
	}
	return true
}
