package utils

import (
	"SyncEthData/db"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	"math/big"
	"strconv"
	"time"
)

func TransformData(block *types.Block) {

	if block == nil {
		return
	}
	b := transferBlock(block)
	header := transferHeader(block.Header())

	transactions := []db.TRANSACTION{}
	if len(block.Transactions()) > 0 {
		for _, transaction := range block.Transactions() {
			trx := transferTrx(transaction, block.Header().Number)
			if trx != nil {
				transactions = append(transactions, *trx)
			}
		}
	}

	db.SaveData(b, header, &transactions)
}

func transferHeader(header *types.Header) *db.HEADER {
	var baseFee = big.NewInt(0)
	if header.BaseFee != nil {
		baseFee = header.BaseFee
	}
	result := db.HEADER{
		PARENTHASH:  header.ParentHash.Hex(),
		UNCLEHASH:   header.UncleHash.Hex(),
		COINBASE:    header.Coinbase.Hex(),
		ROOT:        header.Root.Hex(),
		TXHASH:      header.TxHash.Hex(),
		RECEIPTHASH: header.ReceiptHash.Hex(),
		BLOOM:       header.Bloom.Big().Int64(),
		DIFFICULTY:  header.Difficulty.Int64(),
		BLOCKNUMBER: header.Number.Int64(),
		GASLIMIT:    header.GasLimit,
		GASUSED:     header.GasUsed,
		TIME:        header.Time,
		EXTRA:       string(header.Extra),
		NONCE:       strconv.Itoa(int(header.Nonce.Uint64())),
		BASEFEE:     baseFee.Int64(),
		CREATETIME:  time.Now(),
		UPDATETIME:  time.Now(),
	}

	return &result
}

func transferTrx(trx *types.Transaction, num *big.Int) *db.TRANSACTION {
	msg, err := trx.AsMessage(types.NewLondonSigner(trx.ChainId()), nil)
	if err != nil {
		log.Error(err)
		return nil
	}
	result := db.TRANSACTION{
		TXDATA:      string(trx.Data()),
		HASH:        trx.Hash().Hex(),
		SIZE:        trx.Size().String(),
		FROMACCOUNT: msg.From().Hex(),
		BLOCKNUMBER: num.Int64(),
		CREATETIME:  time.Now(),
		UPDATETIME:  time.Now(),
	}
	return &result
}

func transferBlock(block *types.Block) *db.BLOCK {
	result := db.BLOCK{
		BLOCKNUM:   block.Number().Int64(),
		BLOCKHASH:  block.Hash().Hex(),
		BLOCKSIZE:  block.Size().String(),
		CREATETIME: time.Now(),
		UPDATETIME: time.Now(),
	}
	return &result
}
