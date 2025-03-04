p

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
)

// tableau qui contient tous les groupes
type PageData struct {
	Groups []GroupInfos
}

// structure du Json
type GroupInfos struct { //objet
	ID           int
	Image        string   `json:"id"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

func start(w http.ResponseWriter, r *http.Request) {
	resp, _ := http.Get("https://groupietrackers.herokuapp.com/api/artists") //requête Get pour recupérer tous les groupes grace à l'API
	var infos []GroupInfos                                                   //création d'un tableau pour stocker les élément de la requête Get
	//Conversion de la réponse de la requête Get (type string) en objet go de type []GroupInfos
	json.NewDecoder(resp.Body).Decode(&infos) //récupère la structure de l'objet en fonction du json et va créer une instance 'copie' de l'objet

	//on prend le tableau de la 'func searchFunc' dans la méthode POST pour utiliser la fonction de recherche 'searchFunc' et renvoyer la valeur dans le tableau 'info'
	if r.Method == http.MethodPost { //on vérifie si la requête est de type POST: on compare si la requête 'r.Method' correspond à la methode POST 'http.MethodPost'

		/*fonction de recherche 'searchFunc': permet de trier et renvoyé le tableau*/
		/*'r.FormValue' récupère la valeur comprit dans 'userInput'*/
		infos = searchFunc(r.FormValue("userInput"), infos)
	}

	/*'data' récupère les données '{{}}'*/
	/*'PageData' est la structure utilisé pour contenir tous les groupes dans la structure json 'GroupInfos'*/
	/*'Groups' élément de type tableau contenant la structure d'objet json*/
	/*'infos' variable qui contient le tableau*/

	//data = PageData(structure json)
	//PageData est composé de tableau 'Groups = []GroupInfos' de départ
	//'infos' c'est la valeur des éléments (:les groupes ou artistes)

	data := PageData{ //Valorisation de chaque variable de la structure
		Groups: infos,
	}

	tmpl, _ := template.ParseFiles("templates/search.html") //Définition de la template à utiliser
	tmpl.Execute(w, data)                                   //Envoi de la page html modifié au client
}

// 'func searchFunc' créer un tableau qui doit etre utilisé dans la méthode POST
func searchFunc(input string, infos []GroupInfos) []GroupInfos { //fonction recherche:  prend en paramêtre la recherche utilisateur 'input' et le tableau d'artiste 'infos []GroupInfos' et renverra un nouveau tableau d'artiste trié []GroupInfos
	var output []GroupInfos           //'output' copie du tableau 'infos'
	for i := 0; i < len(infos); i++ { //on parcours l'entrée utilisateur
		if strings.Contains(infos[i].Name, input) { //on compare si l'entrée utilisateur 'input' fait partie du tableau de structure: grâce à l'indice on parcours le tableau d'artiste 'infos []GroupInfos' en fonction de la proriété de l'objet nom '.name' et on compare l'entrée utilisateur
			output = append(output, infos[i]) //si l'entrée utilisateur 'input' est bien dans le tableau de structure, on rajoute l'entrée utilisateur 'infos[]' au nouveau tableau 'output'
		}
	}
	return output //on renvoi le nouveau tableau 'output' mis à jour
}

func main() {
	/*'/' est la racine*/
	http.HandleFunc("/", start)

	http.ListenAndServe(":80", nil)
}
