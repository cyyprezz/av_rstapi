package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"

	"github.com/gorilla/mux"
	_ "github.com/nakagami/firebirdsql"
)

//Initialisierung der Applikation
type App struct {
	Router *mux.Router
	DB     *sqlx.DB
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@%s", user, password, dbname)

	var err error
	a.DB, err = sqlx.Open("firebirdsql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	fmt.Println("API RUNS DATABASE: " + connectionString)
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

// Festlegung der Routen- neue können hinzugefügt werden
// Für die Routen müssen Handler und Models mit Datenbankzugriffen programmiert werden
// Die Routen kann man nach starten der API im Browser,Postman etc.. aufrufen
// z.B http://localhost/lager -> Alle Lager o. http://localhost/artikel/10 -> Artikel mit ID 10
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/lager", a.getLagers).Methods("GET")
	a.Router.HandleFunc("/lager/{id:[0-9]+}", a.getLager).Methods("GET")
	a.Router.HandleFunc("/artikel", a.getArtikels).Methods("GET")
	a.Router.HandleFunc("/artikel/{id:[0-9]+}", a.getArtikel).Methods("GET")
	a.Router.HandleFunc("/artikel/{id:[0-9]+}", a.updateArtikel2).Methods("PUT")
	a.Router.HandleFunc("/artikellager/{id:[0-9]+}", a.getEinzellagerbyArtikel).Methods("GET")
	a.Router.HandleFunc("/lagerumbuchungen", a.createLagerumbuchung).Methods("POST")

}

//Dieser Handler holt den Ist-Bestand von einem Einzellager
//ausgehend von der Artikel-ID
func (a *App) getEinzellagerbyArtikel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Artikel-ID")
		return
	}
	u := Artikel{ID: id}
	products, err := u.getEinzelLagerbyArtikelID(a.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Lager not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, products)

}

//Dieser Handler stellt alle Artikel berreit
func (a *App) getArtikels(w http.ResponseWriter, r *http.Request) {
	products, err := getArtikels(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, products)
}

//Dieser Handler stellt einen spezifischen Artikel berreit
func (a *App) getArtikel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Artikel-ID")
		return
	}
	u := Artikel{ID: id}
	if err := u.GetArtikel(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Artikel not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, u)
}

//Mit diesem Handler kann einem Artikel eine EAN-Nummer hinzugefügt werden
func (a *App) updateArtikel2(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadGateway, "Invalid Artikel ID")
		return
	}
	var u Artikel2
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	u.ID = id

	if err := u.updateArtikel2(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, u)
}

//Dieser Handler stellt alle Lager berreit
func (a *App) getLagers(w http.ResponseWriter, r *http.Request) {
	products, err := getLagers(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, products)
}

//Ein spezfisches Lager
func (a *App) getLager(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Lager-ID")
		return
	}
	u := lager{ID: id}
	if err := u.getLager(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Lager not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, u)
}

//Hier wird via Insert Into eine Lagerumbuchung im System erstellt.
func (a *App) createLagerumbuchung(w http.ResponseWriter, r *http.Request) {
	var u Lagerumbuchungen
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := u.createLagerumbuchung(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, u)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

//ABCDEFG
