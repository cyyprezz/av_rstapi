package main

// Hier werden die Models die für die App nötig sind programmiert
// Datenbankzugriffe ausgehend von den Models
// Models deffinieren die Datenfelder die später via Json abgerufen werden können

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type lager struct {
	ID          int    `json:"id"`
	LagerNr     string `json:"maskenkey"`
	Bezeichnung string `json:"bez"`
}

func (u *lager) getLager(db *sqlx.DB) error {
	statement := fmt.Sprintf("SELECT maskenkey,bez FROM blager WHERE id =%d", u.ID)
	return db.QueryRow(statement).Scan(&u.LagerNr, &u.Bezeichnung)
}

func getLagers(db *sqlx.DB) ([]lager, error) {
	statement := fmt.Sprintf("SELECT id,maskenkey,bez FROM blager")
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lagers := []lager{}

	for rows.Next() {
		var u lager
		if err := rows.Scan(&u.ID, &u.LagerNr, &u.Bezeichnung); err != nil {
			return nil, err
		}
		lagers = append(lagers, u)
	}
	return lagers, nil
}

type Artikel2 struct {
	ID  int    `json:"ID"`
	EAN string `json:"eancode"`
}

func (u *Artikel2) updateArtikel2(db *sqlx.DB) error {
	statment := fmt.Sprintf("UPDATE BSA SET eancode='%s' WHERE id=%d", u.EAN, u.ID)
	_, err := db.Exec(statment)
	return err
}

type Artikel struct {
	ID         int            `json:"ID"`
	ArtikelNr  string         `json:"maskenkey"`
	ArtikelBez string         `json:"artbez"`
	EAN        sql.NullString `json:"eancode"`
}

func (u *Artikel) GetArtikel(db *sqlx.DB) error {
	statement := fmt.Sprintf("SELECT maskenkey,artbez,eancode FROM BSA WHERE id=%d", u.ID)
	return db.QueryRow(statement).Scan(&u.ArtikelNr, &u.ArtikelBez, &u.EAN)
}

func getArtikelsTest(db *sqlx.DB) ([]Artikel, error) {
	aritkels := []Artikel{}
	if err := db.Select(&aritkels, "SELECT ID,maskenkey,artbez,eancode FROM BSA"); err != nil {
		return nil, err
	}
	return aritkels, nil
}

func getArtikels(db *sqlx.DB) ([]Artikel, error) {
	rows, err := db.Queryx("SELECT ID,maskenkey,artbez,eancode FROM BSA")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	artikels := []Artikel{}

	for rows.Next() {
		var u Artikel
		if err := rows.Scan(&u.ID, &u.ArtikelNr, &u.ArtikelBez, &u.EAN); err != nil {
			return nil, err
		}
		artikels = append(artikels, u)
	}
	return artikels, nil
}

type EinzelLager struct {
	ID         int    `json:"ID"`
	BSAID      int    `json:"BSA_ID_LINKKEY"`
	BlagerID   int    `json:"BLAGER_ID_LAGERNR"`
	ISTBestand string `json:"list"`
}

func (u *Artikel) getEinzelLagerbyArtikelID(db *sqlx.DB) ([]EinzelLager, error) {
	rows, err := db.Queryx("SELECT ID,BLAGER_ID_LAGERNR,LIST FROM BARTLH WHERE BSA_ID_LINKKEY=?", u.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	einzellager := []EinzelLager{}

	for rows.Next() {
		var a EinzelLager
		if err := rows.Scan(&a.ID, &a.BlagerID, &a.ISTBestand); err != nil {
			return nil, err
		}
		einzellager = append(einzellager, a)
	}
	return einzellager, nil
}

type Lagerumbuchungen struct {
	ID         int     `json:"ID"`
	BSAID      int     `json:"BSA_ID_ARTNR"`
	Menge      float32 `json:"MENGE"`
	InLagerID  int     `json:"BARTLH_ID_INLAGER"`
	VonLagerID int     `json:"BARTLH_ID_VONLAGER"`
}

func (u *Lagerumbuchungen) createLagerumbuchung(db *sqlx.DB) error {
	_, err := db.NamedExec(`INSERT INTO BLAGVE(BSA_ID_ARTNR,MENGE,BARTLH_ID_INLAGER,BARTLH_ID_VONLAGER,LLDRUCKEN,BMAND_ID) VALUES (:bsaid,:menge,:inlager,:vonlager,'N',1)`,
		map[string]interface{}{
			"bsaid":    u.BSAID,
			"menge":    u.Menge,
			"inlager":  u.InLagerID,
			"vonlager": u.VonLagerID,
		})
	if err != nil {
		return err
	}
	err = db.QueryRow("SELECT max(id) FROM BARTLH").Scan(&u.ID)

	if err != nil {
		return err
	}
	return nil
}
func test(){
	
}