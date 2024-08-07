package client

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strconv"

	"github.com/elleqt/gotron-sdk/pkg/abi"
	"github.com/elleqt/gotron-sdk/pkg/address"
	"github.com/elleqt/gotron-sdk/pkg/common"
	"github.com/elleqt/gotron-sdk/pkg/proto/api"
	"github.com/elleqt/gotron-sdk/pkg/proto/core"
	"google.golang.org/protobuf/proto"
)

// UpdateEnergyLimitContract update contract enery limit
func (g *Client) UpdateEnergyLimitContract(ctx context.Context, from, contractAddress string, value int64) (*api.TransactionExtention, error) {
	fromDesc, err := address.Base58ToAddress(from)
	if err != nil {
		return nil, err
	}

	contractDesc, err := address.Base58ToAddress(contractAddress)
	if err != nil {
		return nil, err
	}

	ct := &core.UpdateEnergyLimitContract{
		OwnerAddress:      fromDesc.Bytes(),
		ContractAddress:   contractDesc.Bytes(),
		OriginEnergyLimit: value,
	}

	tx, err := g.Client.UpdateEnergyLimit(ctx, ct)
	if err != nil {
		return nil, err
	}

	if tx.Result.Code > 0 {
		return nil, fmt.Errorf("%s", string(tx.Result.Message))
	}

	return tx, err
}

// UpdateSettingContract change contract owner consumption ratio
func (g *Client) UpdateSettingContract(ctx context.Context, from, contractAddress string, value int64) (*api.TransactionExtention, error) {
	fromDesc, err := address.Base58ToAddress(from)
	if err != nil {
		return nil, err
	}

	contractDesc, err := address.Base58ToAddress(contractAddress)
	if err != nil {
		return nil, err
	}

	ct := &core.UpdateSettingContract{
		OwnerAddress:               fromDesc.Bytes(),
		ContractAddress:            contractDesc.Bytes(),
		ConsumeUserResourcePercent: value,
	}

	tx, err := g.Client.UpdateSetting(ctx, ct)
	if err != nil {
		return nil, err
	}

	if tx.Result.Code > 0 {
		return nil, fmt.Errorf("%s", string(tx.Result.Message))
	}

	return tx, err
}

// TriggerConstantContract and return tx result
func (g *Client) TriggerConstantContract(ctx context.Context, from, contractAddress, method, jsonString string) (*api.TransactionExtention, error) {
	var err error
	fromDesc := address.HexToAddress("410000000000000000000000000000000000000000")
	if len(from) > 0 {
		fromDesc, err = address.Base58ToAddress(from)
		if err != nil {
			return nil, err
		}
	}
	contractDesc, err := address.Base58ToAddress(contractAddress)
	if err != nil {
		return nil, err
	}

	param, err := abi.LoadFromJSON(jsonString)
	if err != nil {
		return nil, err
	}

	dataBytes, err := abi.Pack(method, param)
	if err != nil {
		return nil, err
	}

	ct := &core.TriggerSmartContract{
		OwnerAddress:    fromDesc.Bytes(),
		ContractAddress: contractDesc.Bytes(),
		Data:            dataBytes,
	}

	return g.triggerConstantContract(ctx, ct)
}

// triggerConstantContract and return tx result
func (g *Client) triggerConstantContract(ctx context.Context, ct *core.TriggerSmartContract) (*api.TransactionExtention, error) {
	return g.Client.TriggerConstantContract(ctx, ct)
}

// TriggerContract and return tx result
func (g *Client) TriggerContract(ctx context.Context, from, contractAddress, method, jsonString string,
	feeLimit, tAmount int64, tTokenID string, tTokenAmount int64) (*api.TransactionExtention, error) {
	fromDesc, err := address.Base58ToAddress(from)
	if err != nil {
		return nil, err
	}

	contractDesc, err := address.Base58ToAddress(contractAddress)
	if err != nil {
		return nil, err
	}

	param, err := abi.LoadFromJSON(jsonString)
	if err != nil {
		return nil, err
	}

	dataBytes, err := abi.Pack(method, param)
	if err != nil {
		return nil, err
	}

	ct := &core.TriggerSmartContract{
		OwnerAddress:    fromDesc.Bytes(),
		ContractAddress: contractDesc.Bytes(),
		Data:            dataBytes,
	}
	if tAmount > 0 {
		ct.CallValue = tAmount
	}
	if len(tTokenID) > 0 && tTokenAmount > 0 {
		ct.CallTokenValue = tTokenAmount
		ct.TokenId, err = strconv.ParseInt(tTokenID, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	return g.triggerContract(ctx, ct, feeLimit)
}

// triggerContract and return tx result
func (g *Client) triggerContract(ctx context.Context, ct *core.TriggerSmartContract, feeLimit int64) (*api.TransactionExtention, error) {
	tx, err := g.Client.TriggerContract(ctx, ct)
	if err != nil {
		return nil, err
	}

	if tx.Result.Code > 0 {
		return nil, fmt.Errorf("%s", string(tx.Result.Message))
	}
	if feeLimit > 0 {
		tx.Transaction.RawData.FeeLimit = feeLimit
		// update hash
		g.UpdateHash(tx)
	}
	return tx, err
}

// EstimateEnergy returns enery required
func (g *Client) EstimateEnergy(ctx context.Context, from, contractAddress, method, jsonString string,
	tAmount int64, tTokenID string, tTokenAmount int64) (*api.EstimateEnergyMessage, error) {
	fromDesc, err := address.Base58ToAddress(from)
	if err != nil {
		return nil, err
	}

	contractDesc, err := address.Base58ToAddress(contractAddress)
	if err != nil {
		return nil, err
	}

	param, err := abi.LoadFromJSON(jsonString)
	if err != nil {
		return nil, err
	}

	dataBytes, err := abi.Pack(method, param)
	if err != nil {
		return nil, err
	}

	ct := &core.TriggerSmartContract{
		OwnerAddress:    fromDesc.Bytes(),
		ContractAddress: contractDesc.Bytes(),
		Data:            dataBytes,
	}
	if tAmount > 0 {
		ct.CallValue = tAmount
	}
	if len(tTokenID) > 0 && tTokenAmount > 0 {
		ct.CallTokenValue = tTokenAmount
		ct.TokenId, err = strconv.ParseInt(tTokenID, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	return g.estimateEnergy(ctx, ct)
}

// triggerContract and return tx result
func (g *Client) estimateEnergy(ctx context.Context, ct *core.TriggerSmartContract) (*api.EstimateEnergyMessage, error) {
	tx, err := g.Client.EstimateEnergy(ctx, ct)
	if err != nil {
		return nil, err
	}

	if tx.Result.Code > 0 {
		return nil, fmt.Errorf("%s", string(tx.Result.Message))
	}

	return tx, err
}

// GetBandwidthPrices retrieves bandwidth prices
func (g *Client) GetBandwidthPrices(ctx context.Context) (*api.PricesResponseMessage, error) {
	result, err := g.Client.GetBandwidthPrices(ctx, new(api.EmptyMessage))
	if err != nil {
		return nil, fmt.Errorf("get bandwidth prices: %v", err)
	}

	return result, nil
}

// GetEnergyPrices retrieves energy prices
func (g *Client) GetEnergyPrices(ctx context.Context) (*api.PricesResponseMessage, error) {
	result, err := g.Client.GetEnergyPrices(ctx, new(api.EmptyMessage))
	if err != nil {
		return nil, fmt.Errorf("get energy prices: %v", err)
	}

	return result, nil
}

// DeployContract and return tx result
func (g *Client) DeployContract(ctx context.Context, from, contractName string,
	abi *core.SmartContract_ABI, codeStr string,
	feeLimit, curPercent, oeLimit int64,
) (*api.TransactionExtention, error) {

	var err error

	fromDesc, err := address.Base58ToAddress(from)
	if err != nil {
		return nil, err
	}

	if curPercent > 100 || curPercent < 0 {
		return nil, fmt.Errorf("consume_user_resource_percent should be >= 0 and <= 100")
	}
	if oeLimit <= 0 {
		return nil, fmt.Errorf("origin_energy_limit must > 0")
	}

	bc, err := common.FromHex(codeStr)
	if err != nil {
		return nil, err
	}

	ct := &core.CreateSmartContract{
		OwnerAddress: fromDesc.Bytes(),
		NewContract: &core.SmartContract{
			OriginAddress:              fromDesc.Bytes(),
			Abi:                        abi,
			Name:                       contractName,
			ConsumeUserResourcePercent: curPercent,
			OriginEnergyLimit:          oeLimit,
			Bytecode:                   bc,
		},
	}

	tx, err := g.Client.DeployContract(ctx, ct)
	if err != nil {
		return nil, err
	}
	if feeLimit > 0 {
		tx.Transaction.RawData.FeeLimit = feeLimit
		// update hash
		g.UpdateHash(tx)
	}
	return tx, err
}

// UpdateHash after local changes
func (g *Client) UpdateHash(tx *api.TransactionExtention) error {
	rawData, err := proto.Marshal(tx.Transaction.GetRawData())
	if err != nil {
		return err
	}

	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)
	tx.Txid = hash
	return nil
}

// GetContractABI return smartContract
func (g *Client) GetContractABI(ctx context.Context, contractAddress string) (*core.SmartContract_ABI, error) {
	var err error
	contractDesc, err := address.Base58ToAddress(contractAddress)
	if err != nil {
		return nil, err
	}

	sm, err := g.Client.GetContract(ctx, GetMessageBytes(contractDesc))
	if err != nil {
		return nil, err
	}
	if sm == nil {
		return nil, fmt.Errorf("invalid contract abi")
	}

	return sm.Abi, nil
}
