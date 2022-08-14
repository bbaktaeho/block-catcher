package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/bbaktaeho/block-catcher/config"
	"github.com/bbaktaeho/block-catcher/database"
	"github.com/bbaktaeho/block-catcher/repository"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

const (
	DATA_COUNT = 3000
	RPC_URL    = "ws://hq_ws.gnd.devnet.kstadium.io:8546"
)

func init() {
	config.LoadEnvironmentFile(".env")

	d := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=latin1&parseTime=True",
		config.Env.DB_USER,
		config.Env.DB_PASSWORD,
		config.Env.DB_HOST,
		config.Env.DB_PORT,
		config.Env.DB_SCHEMA,
	)

	gormDB, err := database.ConnectGORM(mysql.Open(d), &gorm.Config{
		CreateBatchSize:        DATA_COUNT,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}

	db = gormDB
}

func getBlock(startNumber int, client *ethclient.Client, blocks chan<- *types.Block) {
	for i := startNumber; i < startNumber+DATA_COUNT; i++ {
		go func(number int64) {
			block, err := client.BlockByNumber(context.Background(), big.NewInt(number))
			if err != nil {
				log.Fatal("RPC ERROR:", err)
			}

			blocks <- block
		}(int64(i))
	}
}

func setBlock(blocks <-chan *types.Block) []*repository.Block {
	dbBlocks := make([]*repository.Block, 0, DATA_COUNT)

	for block := range blocks {
		dbBlock := &repository.Block{
			Number:          block.Number().Uint64(),
			Hash:            block.Hash().Hex(),
			ParentHash:      block.ParentHash().Hex(),
			UncleHash:       block.UncleHash().Hex(),
			TxRootHash:      block.TxHash().Hex(),
			ReceiptRootHash: block.ReceiptHash().Hex(),
			Miner:           block.Coinbase().Hex(),
			StatusRoot:      block.Root().Hex(),
			GasLimit:        block.GasLimit(),
			GasUsed:         block.GasUsed(),
			ExtraData:       hexutil.Encode(block.Extra()),
			TxCount:         int64(block.Transactions().Len()),
			Timestamp:       block.Time(),
			Size:            float64(block.Size()),
		}

		dbBlocks = append(dbBlocks, dbBlock)
		if len(dbBlocks) == DATA_COUNT {
			break
		}
	}

	return dbBlocks
}

func saveBlocks(dbBlocks []*repository.Block, start time.Time) {
	if err := db.Create(dbBlocks).Error; err != nil {
		log.Fatal("DB_ERR:", err)
	}
}

func main() {
	client, _ := ethclient.Dial(RPC_URL)

	for i := 0; i < 200; i++ {
		start := time.Now()
		blocks := make(chan *types.Block, DATA_COUNT)
		go getBlock(i*DATA_COUNT, client, blocks)
		dbBlocks := setBlock(blocks)
		close(blocks)
		go saveBlocks(dbBlocks, start)
	}
}

// headers := make(chan *types.Header)
// sub, err := client.SubscribeNewHead(context.Background(), headers)
// if err != nil {
// 	log.Fatal(err)
// }

// for {
// 	select {
// 	case err := <-sub.Err():
// 		log.Fatal(err)
// 	case header := <-headers:
// 		go func() {
// 			block, err := client.BlockByNumber(context.Background(), header.Number)
// 			if err != nil {
// 				log.Println(err)
// 			}

// 			log.Println("구독:", block.Number())
// 		}()
// 	}
// }
