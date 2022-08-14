package repository

import "gorm.io/gorm"

type Block struct {
	gorm.Model
	Number          uint64
	Hash            string `gorm:"type:char(66);unique;not null"`
	ParentHash      string `gorm:"type:char(66);not null"`
	UncleHash       string `gorm:"type:char(66);not null"`
	TxRootHash      string `gorm:"type:char(66);not null"`
	ReceiptRootHash string `gorm:"type:char(66);not null"`
	Miner           string `gorm:"type:char(42);not null"`
	StatusRoot      string `gorm:"type:char(66);not null"`
	GasLimit        uint64 `gorm:"not nul"`
	GasUsed         uint64 `gorm:"not nul"`

	ExtraData string `gorm:"type:text"`
	TxCount   int64  `gorm:"default:0"`

	Timestamp uint64  `gorm:"not nul"`
	Size      float64 `gorm:"default:0"`
}
