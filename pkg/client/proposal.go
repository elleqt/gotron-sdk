package client

import (
	"context"
	"fmt"

	"github.com/elleqt/gotron-sdk/pkg/common"
	"github.com/elleqt/gotron-sdk/pkg/proto/api"
	"github.com/elleqt/gotron-sdk/pkg/proto/core"
	"google.golang.org/protobuf/proto"
)

// ProposalsList return all network proposals
func (g *Client) ProposalsList(ctx context.Context) (*api.ProposalList, error) {
	return g.Client.ListProposals(ctx, new(api.EmptyMessage))
}

// ProposalCreate create proposal based on parameter list
func (g *Client) ProposalCreate(ctx context.Context, from string, parameters map[int64]int64) (*api.TransactionExtention, error) {
	var err error

	contract := &core.ProposalCreateContract{
		Parameters: parameters,
	}
	if contract.OwnerAddress, err = common.DecodeCheck(from); err != nil {
		return nil, err
	}

	tx, err := g.Client.ProposalCreate(ctx, contract)
	if err != nil {
		return nil, err
	}
	if proto.Size(tx) == 0 {
		return nil, fmt.Errorf("bad transaction")
	}
	if tx.GetResult().GetCode() != 0 {
		return nil, fmt.Errorf("%s", tx.GetResult().GetMessage())
	}
	return tx, nil
}

// ProposalApprove change URL info
func (g *Client) ProposalApprove(ctx context.Context, from string, id int64, confirm bool) (*api.TransactionExtention, error) {
	var err error

	contract := &core.ProposalApproveContract{
		ProposalId:    id,
		IsAddApproval: confirm,
	}
	if contract.OwnerAddress, err = common.DecodeCheck(from); err != nil {
		return nil, err
	}

	tx, err := g.Client.ProposalApprove(ctx, contract)
	if err != nil {
		return nil, err
	}
	if proto.Size(tx) == 0 {
		return nil, fmt.Errorf("bad transaction")
	}
	if tx.GetResult().GetCode() != 0 {
		return nil, fmt.Errorf("%s", tx.GetResult().GetMessage())
	}
	return tx, nil
}

func (g *Client) ProposalWithdraw(ctx context.Context, from string, id int64) (*api.TransactionExtention, error) {
	var err error

	contract := &core.ProposalDeleteContract{
		ProposalId: id,
	}
	if contract.OwnerAddress, err = common.DecodeCheck(from); err != nil {
		return nil, err
	}

	tx, err := g.Client.ProposalDelete(ctx, contract)
	if err != nil {
		return nil, err
	}
	if proto.Size(tx) == 0 {
		return nil, fmt.Errorf("bad transaction")
	}
	if tx.GetResult().GetCode() != 0 {
		return nil, fmt.Errorf("%s", tx.GetResult().GetMessage())
	}
	return tx, nil
}
