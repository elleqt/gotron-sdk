package client

import (
	"context"

	"github.com/elleqt/gotron-sdk/pkg/proto/api"
	"github.com/elleqt/gotron-sdk/pkg/proto/core"
)

// GetTransactionSignWeight queries transaction sign weight
func (g *GrpcClient) GetTransactionSignWeight(ctx context.Context, tx *core.Transaction) (*api.TransactionSignWeight, error) {
	result, err := g.Client.GetTransactionSignWeight(ctx, tx)
	if err != nil {
		return nil, err
	}
	return result, nil
}
