package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Open(dsn string) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", dsn)
	return db, err
}

func CreateTables(db *sql.DB) (err error) {
	stmt := `CREATE TABLE IF NOT EXISTS "attributes" (
		"id"          INTEGER NOT NULL,
		"type"        TEXT NOT NULL,
		"short"       TEXT NOT NULL,
		"name_de"     TEXT NOT NULL,
		"name_en"     TEXT NOT NULL,
		"location_id" INTEGER NOT NULL,
		PRIMARY KEY("id","type")
	)`
	if _, err = db.Exec(stmt); err != nil {
		return err
	}

	stmt = `CREATE TABLE IF NOT EXISTS "categories" (
		"id"          INTEGER NOT NULL,
		"name_de"     TEXT NOT NULL,
		"name_en"     TEXT NOT NULL,
		"location_id" INTEGER NOT NULL,
		PRIMARY KEY("id")
	)`
	if _, err = db.Exec(stmt); err != nil {
		return err
	}

	stmt = `CREATE TABLE IF NOT EXISTS "locations" (
		"id"   INTEGER NOT NULL,
		"name" TEXT NOT NULL,
		PRIMARY KEY("id")
	)`
	if _, err = db.Exec(stmt); err != nil {
		return err
	}

	stmt = `CREATE TABLE IF NOT EXISTS "meal_attributes" (
		"meal_id"        INTEGER NOT NULL,
		"attribute_id"   INTEGER NOT NULL,
		"attribute_type" TEXT NOT NULL,
		PRIMARY KEY("meal_id","attribute_id","attribute_type")
	)`
	if _, err = db.Exec(stmt); err != nil {
		return err
	}

	stmt = `CREATE TABLE IF NOT EXISTS "meals" (
		"id"            INTEGER NOT NULL,
		"name_de"       TEXT NOT NULL,
		"name_en"       TEXT NOT NULL,
		"category_id"   INTEGER NOT NULL,
		"price_student" REAL NOT NULL,
		"price_guest"   REAL NOT NULL,
		"co2_grams"     INTEGER,
		"co2_rating"    TEXT,
		"date"          TEXT NOT NULL,
		"location_id"   INTEGER NOT NULL,
		FOREIGN KEY("location_id") REFERENCES "locations"("id"),
		PRIMARY KEY("id" AUTOINCREMENT),
		FOREIGN KEY("category_id") REFERENCES "categories"("id")
	)`
	if _, err = db.Exec(stmt); err != nil {
		return err
	}

	stmt = `CREATE TABLE IF NOT EXISTS "nutrition" (
		"meal_id"       INTEGER NOT NULL,
		"kj"            INTEGER NOT NULL,
		"kcal"          INTEGER NOT NULL,
		"fat"           REAL NOT NULL,
		"saturated_fat" REAL NOT NULL,
		"carbohydrates" REAL NOT NULL,
		"sugar"         REAL NOT NULL,
		"protein"       REAL NOT NULL,
		"salt"          REAL NOT NULL,
		FOREIGN KEY("meal_id") REFERENCES "meals"("id")
	)`
	if _, err = db.Exec(stmt); err != nil {
		return err
	}

	return err
}
