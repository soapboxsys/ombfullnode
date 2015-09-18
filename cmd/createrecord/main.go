package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/database"
	_ "github.com/btcsuite/btcd/database/ldb"
	"github.com/btcsuite/btcd/limits"
	"github.com/btcsuite/go-flags"
	"github.com/soapboxsys/ombudslib/ombutil"
	"github.com/soapboxsys/ombudslib/pubrecdb"
)

var (
	helpStr = "This utility takes a leveldb database constructed by btcd or a\n" +
		"fork of it and produces the public record by reading and parsing blocks\n" +
		"stored within it."

	ombudsCoreHome     = ombutil.AppDataDir("ombudscore", false)
	ombFullNodeHome    = filepath.Join(ombudsCoreHome, "node")
	defaultDataDir     = filepath.Join(ombFullNodeHome, "data")
	defaultStartHeight = int64(0)
	activeNetParams    = &chaincfg.MainNetParams
)

type config struct {
	// data-directory
	DataDir  string `short:"b" long:"datadir" description:"Location of the btcd data directory"`
	TestNet3 bool   `long:"testnet" description:"Use the test network"`
	// Start height which defaults to 0
	StartHeight int64 `short:"s" long:"start" description:"The height which the pubrecdb starts at"`
	// overwrite
	Overwrite bool `long:"overwrite" description:"Overwrite if existing pubrecdb is present"`
}

func loadBlockDB(cfg config) (database.Db, error) {
	dbPath := filepath.Join(cfg.DataDir, "blocks_leveldb")

	log.Println(dbPath)
	db, err := database.OpenDB("leveldb", dbPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func loadPubRecDB(cfg config) (db *pubrecdb.PublicRecord, err error) {
	dbPath := filepath.Join(cfg.DataDir, "pubrecord.db")

	if cfg.Overwrite {
		db, err = pubrecdb.InitDB(dbPath)
		if err != nil {
			return nil, err
		}
	} else {
		db, err = pubrecdb.LoadDB(dbPath)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// Required or leveldb will not function properly
	if err := limits.SetLimits(); err != nil {
		os.Exit(1)
	}

	cfg := config{
		DataDir:     defaultDataDir,
		StartHeight: defaultStartHeight,
		Overwrite:   true,
		TestNet3:    false,
	}

	// Parse command line options.
	parser := flags.NewParser(&cfg, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); !ok || e.Type != flags.ErrHelp {
			parser.WriteHelp(os.Stderr)
		}
		log.Fatal(err)
	}

	if cfg.TestNet3 {
		cfg.DataDir = filepath.Join(cfg.DataDir, "testnet")
		activeNetParams = &chaincfg.TestNet3Params
	} else {
		cfg.DataDir = filepath.Join(cfg.DataDir, "mainnet")
	}

	// Load the existing blockdb
	ldb, err := loadBlockDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize or perhaps empty the pubrecord
	pdb, err := loadPubRecDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// With the level db start from the beginning and iteratively build out the
	// public record db
	curHeight := cfg.StartHeight
	s, err := ldb.FetchBlockShaByHeight(curHeight)
	if err != nil {
		// Assume from this point on that the db will behave when we query it...
		log.Fatal(err)
	}
	blk, _ := ldb.FetchBlockBySha(s)

	// Establish end of the reachable block heights
	_, finalHeight, _ := ldb.NewestSha()

	// Every time through this loop we take at least 3 round trips to disk.
	for blk.Height() < finalHeight {
		h := blk.Height()
		if h%100 == 0 {
			log.Printf("Got to height: %d, Processing...\n", h)
		}
		// Only scan blocks after start height
		bltns, err := ombutil.ProcessBlock(blk, activeNetParams)
		if err != nil {
			processBlkErr(h, err)
		}
		// Add to prec block then bulletins to prec
		if err = pdb.StoreBlock(blk); err != nil {
			processBlkErr(h, err)
		}

		for _, bltn := range bltns {
			if err = pdb.StoreBulletin(bltn); err != nil {
				processBlkErr(h, err)
			}
		}

		// Move the pointer to the next block by adding one
		curHeight += 1
		s, err := ldb.FetchBlockShaByHeight(curHeight)
		if err != nil {
			processBlkErr(curHeight, err)
		}
		blk, err = ldb.FetchBlockBySha(s)
		if err != nil {
			processBlkErr(curHeight, err)
		}
	}

	// Create a summary of both DBs to compare.
	// Or you know... don't

	// Fin.
	log.Println("We did it....!")
}

func processBlkErr(h int64, err error) {
	f := fmt.Errorf("Parsing blk: %d threw: %s", h, err)
	log.Fatal(f)
}
