package main

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
