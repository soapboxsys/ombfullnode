package main

import (
	"github.com/NSkelsey/ahimsadb"
	"github.com/NSkelsey/btcsubprotos"
	"github.com/NSkelsey/protocol/ahimsa"
	"github.com/conformal/btcnet"
	"github.com/conformal/btcutil"
	"github.com/conformal/btcwire"
)

type PubRecManager struct {
	txChan    chan *btcutil.Tx
	blkChan   chan *btcutil.Block
	netParams *btcnet.Params
	db        *ahimsadb.PublicRecord
}

func (m *PubRecManager) handleBlockPush(blk *btcutil.Block) {

	hash, _ := blk.Sha()

	// Store block in the database
	if err := m.db.StoreBlock(blk); err != nil {
		precLog.Errorf("Failed to store blk %s", err)
	}
}

// Hand
func (m *PubRecManager) handleTxPush(tx *btcutil.Tx, blk *btcutil.Block) {

	if btcsubprotos.IsBulletin(tx.MsgTx()) {
		var blkSha *btcwire.ShaHash
		if blk != nil {
			blkSha, _ = blk.Sha()
		}

		precLog.Info("Storing bltn ", tx.Sha())
		bltn, err := ahimsa.NewBulletin(tx.MsgTx(), blkSha, m.netParams)
		if err != nil {
			precLog.Errorf("Could create bulletin from tx: %s", err)
		}

		if err := m.db.StoreBulletin(bltn); err != nil {
			precLog.Errorf("Failed to store bulletin: %s", err)
		}
	}
}

func newPubRecManager(net *btcnet.Params) *PubRecManager {

	var db *ahimsadb.PublicRecord
	var err error
	db, err = ahimsadb.LoadDB("./pubrecord.db")
	if err != nil {
		db, err = ahimsadb.InitDB("./pubrecord.db")
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

// Must be run as a seperate goroutine
func (m *PubRecManager) Start() {

	precLog.Info("Starting pubrecord db manager")
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
