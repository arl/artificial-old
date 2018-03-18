package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"path"
	"time"

	"github.com/aurelien-rainone/evolve/framework"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

const (
	dbName         = "evolution.db"
	createTableStr = `CREATE TABLE generations(
		id INTEGER NOT NULL PRIMARY KEY,
		best_fitness REAL NOT NULL,
		mean_fitness REAL NOT NULL,
		fitness_stddev REAL NOT NULL,
		natural_fitness INTEGER NOT NULL,
		pop_size INTEGER NOT NULL,
		elite_count INTEGER NOT NULL,
		gen_number INTEGER NOT NULL,
		elapsed INTEGER NOT NULL);`

	insertGenerationStr = `INSERT INTO generations(
		best_fitness,
		mean_fitness,
		fitness_stddev,
		natural_fitness,
		pop_size,
		elite_count,
		gen_number,
		elapsed)
		values(?, ?, ?, ?, ?, ?, ?, ?)`
)

type sqliteObserver struct {
	freq     int          // backup every N generations
	outDir   string       // output directory
	db       *sql.DB      // sqlite db
	sqlConn  *sql.Conn    // keep connection here, nobody else will use it
	unixLn   net.Listener // unix socket listener (use to signal viewer of new generations)
	unixConn net.Conn     // unix socket connection
}

func newSqliteObserver(freq int, outDir string) (o *sqliteObserver, err error) {
	if freq == 0 {
		return nil, fmt.Errorf("sqliteObserver frequency can't be 0")
	}

	o = &sqliteObserver{freq: freq, outDir: outDir}

	if err = o.setupSQL(); err != nil {
		return nil, err
	}
	if err = o.setupSignaling(); err != nil {
		return nil, err
	}
	return o, nil
}

func (o *sqliteObserver) setupSQL() error {
	dbPath := path.Join(o.outDir, dbName)
	os.Remove(dbPath)

	var err error
	o.db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("can't open sqlite db: %v", err)
	}
	o.sqlConn, err = o.db.Conn(context.TODO())
	if err != nil {
		return fmt.Errorf("can't open sqlite connection: %v", err)
	}
	_, err = o.sqlConn.ExecContext(context.TODO(), createTableStr)
	if err != nil {
		return fmt.Errorf("can't exec query: %q: %s", err, createTableStr)
	}
	return err
}

func (o *sqliteObserver) setupSignaling() error {
	var err error
	if o.unixLn == nil {
		o.unixLn, err = net.Listen("unix", path.Join(o.outDir, "generation.sock"))
		if err != nil {
			return err
		}
	}

	// listens for accepted connections in a go routine
	go func() {
		conn, err := o.unixLn.Accept()
		if err != nil {
			// handle error

			log.Error().Err(err).Msg("error accepting socket connection")
			o.unixConn = nil
			return
		}
		log.Info().Msg("accepted socket connection")
		o.unixConn = conn
	}()
	return nil
}

func (o *sqliteObserver) PopulationUpdate(data *framework.PopulationData) {
	genNum := data.GenerationNumber()
	if genNum%o.freq == 0 {

		// fill sql table with generation data
		tx, err := o.sqlConn.BeginTx(context.TODO(), nil)
		if err != nil {
			log.Fatal().Err(err).Msg("sql error")
		}
		stmt, err := tx.Prepare(insertGenerationStr)
		if err != nil {
			log.Fatal().Err(err).Msg("sql error")
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
			log.Fatal().Err(err).Msg("sql error")
		}
		tx.Commit()

		// signal external processes there is new data
		if o.unixConn != nil {
			if _, err = o.unixConn.Write([]byte("newdata")); err != nil {
				log.Error().Err(err).Msg("couldn't write to socket")
				o.unixConn.Close()
				o.unixConn = nil
				o.setupSignaling()
			}
		}
	}
}

func (o *sqliteObserver) close() {
	if err := o.sqlConn.Close(); err != nil {
		log.Error().Err(err).Msg("error closing database connection")
	}
	if err := o.db.Close(); err != nil {
		log.Error().Err(err).Msg("error closing database")
	}
	if o.unixLn != nil {
		if err := o.unixLn.Close(); err != nil {
			log.Error().Err(err).Msg("error closing unix socket listener")
		}
	}
	if o.unixConn != nil {
		if err := o.unixConn.Close(); err != nil {
			log.Printf("error closing unix socket connection: %v\n", err)
		}
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
		log.Info().Msgf("Generation %d: best: %.2f mean: %.2f stddev: %.2f",
			data.GenerationNumber(), data.BestCandidateFitness(), data.MeanFitness(), data.FitnessStandardDeviation())

		// save this generation's best candidate
		imgPath := path.Join(o.outDir, fmt.Sprintf("%d.png", generation))
		renderAndDiff(data.BestCandidate().(*imageDNA), &imgPath)
	}
}
