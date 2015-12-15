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

func (m *pubRecManager) ConnectBlock(blk *btcutil.Block) {
	// Look at block height, if it is above the introduction threshold, Parse
	// the block. Otherwise just store the headers.
	// TODO
	ombBlk := ombutil.CreateUBlock(blk)
	err := m.db.InsertOmbBlk(ombBlk)
	if err != nil {
		precLog.Errorf("Connecting Blk[%s] failed with: %s", blk.Sha(), err)
	}
}

func (m *pubRecManager) DisconnectBlock(blk *btcutil.Block) {
	err := m.db.RemoveBlk(blk.Sha())
	if err != nil {
		precLog.Errorf("Disconnecting Blk[%s] failed with: %s", blk.Sha(), err)
	}
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

func (m *pubRecManager) Start() {
	precLog.Info("Starting PubRecord db Manager.")

}

func (m *pubRecManager) Stop() {
	precLog.Info("Stopping PubRecord db Manager.")
}
