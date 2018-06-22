package model

import (
	"github.com/ethereum/go-ethereum/p2p/discover"
	"time"
)

// Peers contains the Image data and related user who uploaded the image
type PeerNode struct {
	ID string `json:"id" gorm:"not null"`
	Peer *discover.Node `json:"peer" gorm:"not null"`
	LastSeen time.Time `json:"last_seen" gorm:"not null"`
}
