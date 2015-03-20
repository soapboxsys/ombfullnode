package main

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/soapboxsys/ombudslib/protocol/ombproto"
	"github.com/soapboxsys/ombudslib/pubrecdb"
)

type PubRecManager struct {
	txChan    chan *btcutil.Tx
	blkChan   chan *btcutil.Block
	netParams *chaincfg.Params
	db        *pubrecdb.PublicRecord
}

func (m *PubRecManager) handleBlockPush(blk *btcutil.Block) {

	// Store block in the database
	if err := m.db.StoreBlock(blk); err != nil {
		precLog.Errorf("Failed to store blk %s", err)
	}
}

// handleTxPush takes a transaction and stores the bulletin contained within in
// the PubRecManager's database.
func (m *PubRecManager) handleTxPush(tx *btcutil.Tx, blk *btcutil.Block) {

	if ombproto.IsBulletin(tx.MsgTx()) {
		var blkSha *wire.ShaHash
		if blk != nil {
			blkSha, _ = blk.Sha()
		}

		precLog.Info("Storing bltn ", tx.Sha())
		bltn, err := ombproto.NewBulletin(tx.MsgTx(), blkSha, m.netParams)
		if err != nil {
			precLog.Errorf("Could create bulletin from tx: %s", err)
		}

		if err := m.db.StoreBulletin(bltn); err != nil {
			precLog.Errorf("Failed to store bulletin: %s", err)
		}
	}
}

func newPubRecManager(net *chaincfg.Params) *PubRecManager {

	var db *pubrecdb.PublicRecord
	var err error

	// TODO resolve path
	db, err = pubrecdb.LoadDB(defaultDBPath)
	if err != nil {
		db, err = pubrecdb.InitDB(defaultDBPath)
		if err != nil {
			precLog.Errorf("Failed to Init db! %s", err)
		}
	}
	precLog.Info("Successfully loaded pubrecord.db!")

	return &PubRecManager{
		txChan:    make(chan *btcutil.Tx, 10000),
		blkChan:   make(chan *btcutil.Block, 100),
		netParams: net,
		db:        db,
	}
}

// The PublicRecordManager must be started as a separate go routine.
func (m *PubRecManager) Start() {

	precLog.Info("Starting PubRecord db Manager.")
	for {
		select {
		case tx := <-m.txChan:
			m.handleTxPush(tx, nil)
		case blk := <-m.blkChan:
			m.handleBlockPush(blk)
			for _, tx := range blk.Transactions() {
				m.handleTxPush(tx, blk)
			}
		}
	}

}
