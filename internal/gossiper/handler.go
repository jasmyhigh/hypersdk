// Copyright (C) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package gossiper

import (
	"context"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/network/p2p"
	"github.com/ava-labs/avalanchego/utils/logging"
	"go.uber.org/zap"
)

var _ p2p.Handler = (*TxGossipHandler)(nil)

type Ready interface {
	IsReady() bool
}

type TxGossipHandler struct {
	p2p.NoOpHandler
	ready    Ready
	log      logging.Logger
	gossiper Gossiper
}

func NewTxGossipHandler(
	ready Ready,
	log logging.Logger,
	gossiper Gossiper,
) *TxGossipHandler {
	return &TxGossipHandler{
		ready:    ready,
		log:      log,
		gossiper: gossiper,
	}
}

func (t *TxGossipHandler) AppGossip(ctx context.Context, nodeID ids.NodeID, msg []byte) {
	if !t.ready.IsReady() {
		t.log.Debug("ignoring app gossip failed because vm is not ready")
		return
	}

	if err := t.gossiper.HandleAppGossip(ctx, nodeID, msg); err != nil {
		t.log.Warn("handle app gossip failed", zap.Error(err))
	}
}
