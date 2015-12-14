package main

import (
	"github.com/btcsuite/btcutil"
	"github.com/soapboxsys/ombudslib/ombutil"
	"github.com/soapboxsys/ombudslib/pubrecdb"
)

type pubRecManager struct {
	// New Bulleitn Channel
	bltnChan chan *ombutil.Bulletin
	// TODO New Endorsement Chan
	db     *pubrecdb.PublicRecord
	server *server
}

func (m *pubRecManager) AcceptBlock(blk *btcutil.Block) {

}

func (m *pubRecManager) AttachBlock(blk *btcutil.Block) {

}

func (m *pubRecManager) DetachBlock(blk *btcutil.Block) {

}

func (m *pubRecManager) ProcessTx(tx *btcutil.Tx) {

}

func newPubRecManager(server *server) (*pubRecManager, error) {

	var db *pubrecdb.PublicRecord
	var err error

	db, err = pubrecdb.LoadDB(cfg.PubRecFile)
	if err != nil {
		db, err = pubrecdb.InitDB(cfg.PubRecFile, activeNetParams.Params)
		if err != nil {
			precLog.Errorf("Failed to Init db! %s", err)
			return nil, err
		}
	}
	precLog.Info("Successfully loaded pubrecord.db!")

	m := &pubRecManager{
		server: server,
		db:     db,
	}
	return m, nil
}

// The PublicRecordManager must be started as a separate go routine.
func (m *pubRecManager) Start() {
	precLog.Info("Starting PubRecord db Manager.")

}

func (m *pubRecManager) Stop() {
	precLog.Info("Stopping PubRecord db Manager.")
}
