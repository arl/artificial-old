package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/aurelien-rainone/evolve/framework"
	_ "github.com/mattn/go-sqlite3"
)

const dbName = "evolution.db"

type sqliteObserver struct {
	db   *sql.DB   // sqlite db
	conn *sql.Conn // keep connection here as nobody will use it
	freq int       // backup every N generations
}

func newSqliteObserver(freq int, outDir string) (o *sqliteObserver, err error) {
	if freq == 0 {
		return nil, fmt.Errorf("sqliteObserver frequency can't be 0")
	}

	dbPath := path.Join(outDir, dbName)
	os.Remove(dbPath)

	o = &sqliteObserver{freq: freq}
	o.db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("can't open sqlite db: %v", err)
	}
	o.conn, err = o.db.Conn(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("can't open sqlite connection: %v", err)
	}

	sqlStmt := `
	CREATE TABLE generations (
		id INTEGER NOT NULL PRIMARY KEY,
		best_fitness REAL NOT NULL,
		mean_fitness REAL NOT NULL,
		fitness_stddev REAL NOT NULL,
		natural_fitness INTEGER NOT NULL,
		pop_size INTEGER NOT NULL,
		elite_count INTEGER NOT NULL,
		gen_number INTEGER NOT NULL,
		elapsed INTEGER NOT NULL);
	`
	//_, err = o.db.Exec(sqlStmt)
	_, err = o.conn.ExecContext(context.TODO(), sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("can't exec query: %q: %s", err, sqlStmt)
	}

	return o, nil
}

const genInsertStr = `INSERT INTO generations (best_fitness, mean_fitness, fitness_stddev, natural_fitness, pop_size, elite_count, gen_number, elapsed) values(?, ?, ?, ?, ?, ?, ?, ?)`

func (o *sqliteObserver) PopulationUpdate(data *framework.PopulationData) {
	genNum := data.GenerationNumber()
	if genNum%o.freq == 0 {
		//tx, err := o.db.Begin()
		tx, err := o.conn.BeginTx(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare(genInsertStr)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(
			data.BestCandidateFitness(),
			data.MeanFitness(),
			data.FitnessStandardDeviation(),
			data.IsNaturalFitness(),
			data.PopulationSize(),
			data.EliteCount(),
			data.GenerationNumber(),
			data.ElapsedTime()/time.Second,
		)
		if err != nil {
			log.Fatal(err)
		}
		tx.Commit()

	}
}

func (o *sqliteObserver) close() {
	if err := o.conn.Close(); err != nil {
		log.Printf("error closing connection: %v\n", err)
	}
	if err := o.db.Close(); err != nil {
		log.Printf("error closing database: %v\n", err)
	}
}

type bestObserver struct {
	freq   int    // print statistics every N generations
	outDir string // output directory
}

func newBestObserver(freq int, outDir string) (o *bestObserver, err error) {
	if freq == 0 {
		return nil, fmt.Errorf("bessObserver frequency can't be 0")
	}
	return &bestObserver{freq: freq, outDir: outDir}, nil
}

func (o *bestObserver) PopulationUpdate(data *framework.PopulationData) {
	generation := data.GenerationNumber()
	if generation%o.freq == 0 {
		// update best candidate
		log.Printf("Generation %d: best: %.2f mean: %.2f stddev: %.2f\n",
			data.GenerationNumber(), data.BestCandidateFitness(), data.MeanFitness(), data.FitnessStandardDeviation())
		saveToPng(
			path.Join(o.outDir, fmt.Sprintf("%d.png", generation)),
			data.BestCandidate().(*imageDNA).render())
	}
}
