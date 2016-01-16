package main

import (
	"path"

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
	if !ombutil.PastPegDate(blk) {
		return
	}

	// Look at block height, if it is above the introduction threshold, Parse
	// the block. Otherwise just store the headers.
	ombBlk := ombutil.CreateUBlock(blk, precLog)
	err, ok := m.db.InsertUBlock(ombBlk)
	if err != nil || !ok {
		precLog.Errorf("Connecting Blk[%s] failed with: %s and: %s",
			blk.Sha(), err, ok)
		return
	}

	if len(ombBlk.Bulletins) > 0 || len(ombBlk.Endorsements) > 0 {
		precLog.Infof("Inserted Blk[%s] with %d bltns & %d endos",
			blk.Sha(), len(ombBlk.Bulletins), len(ombBlk.Endorsements))
	}
}

func (m *pubRecManager) DisconnectBlock(blk *btcutil.Block) {
	if !ombutil.PastPegDate(blk) {
		return
	}

	err, ok := m.db.DeleteBlockTip(blk.Sha())
	if err != nil || !ok {
		precLog.Errorf("Disconnecting Blk[%s] failed with: %s and: %s",
			blk.Sha(), ok, err)
		return
	}

	precLog.Infof("Disconnected Blk[%s]", blk.Sha())
}

func (m *pubRecManager) ProcessTx(tx *btcutil.Tx) {

}

func newPubRecManager(server *server) (*pubRecManager, error) {

	var db *pubrecdb.PublicRecord
	var err error

	dbName := "pubrecord.db"
	dbPath := path.Join(cfg.DataDir, dbName)

	db, err = pubrecdb.LoadDB(dbPath)
	if err != nil {
		db, err = pubrecdb.InitDB(dbPath, activeNetParams.Params)
		if err != nil {
			precLog.Errorf("Failed to Init db: %s", err)
			return nil, err
		}
	}
	precLog.Infof("Successfully loaded: %s", dbPath)

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
