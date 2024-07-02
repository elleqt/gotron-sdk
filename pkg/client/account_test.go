package client_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/elleqt/gotron-sdk/pkg/client"
	"github.com/elleqt/gotron-sdk/pkg/proto/core"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

var (
	conn                              *client.Client
	apiKey                            = "622ec85e-7406-431d-9caf-0a19501469a4"
	tronAddress                       = "grpc.nile.trongrid.io:50051"
	accountAddress                    = "TPpw7soPWEDQWXPCGUMagYPryaWrYR5b3b"
	accountAddressWitness             = "TGj1Ej1qRzL9feLTLhjwgxXF4Ct6GTWg2U"
	testnetNileAddressExample         = "TUoHaVjx7n5xz8LwPRDckgFrDWhMhuSuJM"
	testnetNileAddressDelegateExample = "TZ4UXDV5ZhNW7fb2AMSbgfAEZ7hWsnYS2g"
)

func TestMain(m *testing.M) {
	opts := make([]grpc.DialOption, 0)
	opts = append(opts, grpc.WithInsecure())

	conn = client.New(tronAddress)

	if err := conn.Start(opts...); err != nil {
		_ = fmt.Errorf("Error connecting GRPC Client: %v", err)
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestGetAccountDetailed(t *testing.T) {
	acc, err := conn.GetAccountDetailed(context.Background(), accountAddress)
	require.Nil(t, err)
	require.NotNil(t, acc.Allowance)
	require.NotNil(t, acc.Rewards)

	t.Skip()
	acc2, err := conn.GetAccountDetailed(context.Background(), accountAddressWitness)
	require.Nil(t, err)
	require.NotNil(t, acc2.Allowance)
	require.NotNil(t, acc2.Rewards)

}

func TestGetAccountDetailedV2(t *testing.T) {
	acc, err := conn.GetAccountDetailed(context.Background(), testnetNileAddressExample)

	require.Nil(t, err)
	require.NotNil(t, acc.Allowance)
	require.NotNil(t, acc.Rewards)

	require.NotNil(t, acc.MaxCanDelegateBandwidth)
	require.NotNil(t, acc.MaxCanDelegateEnergy)

}

func TestFreezeV2(t *testing.T) {
	t.Skip() // Only in testnet nile
	freezeTx, err := conn.FreezeBalanceV2(context.Background(), testnetNileAddressExample, core.ResourceCode_BANDWIDTH, 1000000)

	require.Nil(t, err)
	require.NotNil(t, freezeTx.GetTxid())

}

func TestUnfreezeV2(t *testing.T) {
	t.Skip() // Only in testnet nile
	unfreezeTx, err := conn.UnfreezeBalanceV2(context.Background(), testnetNileAddressExample, core.ResourceCode_BANDWIDTH, 1000000)

	require.Nil(t, err)
	require.NotNil(t, unfreezeTx.GetTxid())

}

func TestDelegate(t *testing.T) {
	t.Skip() // Only in testnet nile
	tx, err := conn.DelegateResource(context.Background(), testnetNileAddressExample, testnetNileAddressDelegateExample, core.ResourceCode_BANDWIDTH, 1000000, false, 10000)

	require.Nil(t, err)
	require.NotNil(t, tx.GetTxid())
}

func TestUndelegate(t *testing.T) {
	t.Skip() // Only in testnet nile
	tx, err := conn.UnDelegateResource(context.Background(), testnetNileAddressExample, testnetNileAddressDelegateExample, core.ResourceCode_BANDWIDTH, 1000000, false)

	require.Nil(t, err)
	require.NotNil(t, tx.GetTxid())
}

func TestDelegateMaxSize(t *testing.T) {
	t.Skip() // Only in testnet nile
	tx, err := conn.GetCanDelegatedMaxSize(context.Background(), testnetNileAddressExample, int32(core.ResourceCode_BANDWIDTH.Number()))

	require.Nil(t, err)
	require.GreaterOrEqual(t, tx.GetMaxSize(), int64(0))
}

func TestUnfreezeLeftCount(t *testing.T) {
	t.Skip() // Only in testnet nile
	tx, err := conn.GetAvailableUnfreezeCount(context.Background(), testnetNileAddressExample)

	require.Nil(t, err)
	require.GreaterOrEqual(t, tx.GetCount(), int64(0))
}
