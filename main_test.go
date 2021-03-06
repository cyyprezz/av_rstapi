package main

//Tests für alle Komponenten der API - wenn die API auf die jeweilige Datenbank konfiguierit ist
//Kann man mit go test -v den Status der einzelenen Komponenten prüfen
//Der Test muss im Projektverzeichnis gestartet werden

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize("SYSDBA", "masterkey", "C:\\privbackup.FDB")

	ensureBLAGERExists()
	code := m.Run()
	ensureBSAExists()
	os.Exit(code)
}

func ensureBLAGERExists() {
	if _, err := a.DB.Exec(tableSelectQuery); err != nil {
		log.Fatal(err)
	}

}

func ensureBSAExists() {
	if _, err := a.DB.Exec(tableSelectQuery2); err != nil {
		log.Fatal(err)
	}
}

const tableSelectQuery2 = "SELECT * FROM BSA"
const tableSelectQuery = "SELECT * FROM BLAGER"

func TestLager(t *testing.T) {

	req, _ := http.NewRequest("GET", "/lager", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body == "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentLager(t *testing.T) {

	req, _ := http.NewRequest("GET", "/lager/45", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Lager not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Lager not found'. Got '%s'", m["error"])
	}
}

func TestGetLager(t *testing.T) {

	req, _ := http.NewRequest("GET", "/lager/12", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetArtikel(t *testing.T) {

	req, _ := http.NewRequest("GET", "/artikel/12", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestArtikel(t *testing.T) {

	req, _ := http.NewRequest("GET", "/artikel", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body == "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestLagerbyArtikel(t *testing.T) {

	req, _ := http.NewRequest("GET", "/artikellager/10", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

}

func TestUpdateArtikel(t *testing.T) {

	req, _ := http.NewRequest("GET", "/artikel/10", nil)
	response := executeRequest(req)

	var originalArtikel map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalArtikel)

	payload := []byte(`{"eancode": "888888888"}`)
	req, _ = http.NewRequest("PUT", "/artikel/10", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["id"] != originalArtikel["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalArtikel["id"], m["id"])
	}

	if m["eancode"] == originalArtikel["eancode"] {
		t.Errorf("Expected the eancode to change from '%v' to '%v'. Got '%v'", originalArtikel["eancode"], m["eancode"], m["eancode"])
	}
}

func TestCreateLagerumbuchung(t *testing.T) {
	payload := []byte(`{"BSA_ID_ARTNR":10,"MENGE":10,"BARTLH_ID_INLAGER":6,"BARTLH_ID_VONLAGER":14}`)

	req, _ := http.NewRequest("POST", "/lagerumbuchungen", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	fmt.Println(m)

}
