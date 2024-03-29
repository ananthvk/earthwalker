// Package getplaces serves the get_places template.
package getplaces

import (
	"html/template"
	"log"
	"net/http"

	"gitlab.com/glatteis/earthwalker/util"
)

var getPlaces = template.Must(template.ParseFiles(util.StaticPath()+"/templates/main_template.html.tmpl", util.StaticPath()+"/templates/get_places/get_places.html.tmpl"))

// ServeGetPlaces serves the get_places template.
func ServeGetPlaces(w http.ResponseWriter, r *http.Request) {
	err := getPlaces.Execute(w, struct{}{})
	if err != nil {
		log.Println(err)
		http.Error(w, "there was some kind of internal error, sorry!", http.StatusUnprocessableEntity)
		return
	}
}
