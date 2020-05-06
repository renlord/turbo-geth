package apis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/consensus/ethash"
	"github.com/ledgerwatch/turbo-geth/consensus/misc"
	"github.com/ledgerwatch/turbo-geth/core"
	"github.com/ledgerwatch/turbo-geth/core/state"
	"github.com/ledgerwatch/turbo-geth/core/types"
	"github.com/ledgerwatch/turbo-geth/core/vm"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"github.com/ledgerwatch/turbo-geth/ethdb/remote/remotechain"
	"github.com/ledgerwatch/turbo-geth/params"
)

func RegisterRetraceAPI(router *gin.RouterGroup, e *Env) error {
	router.GET(":chain/:number", e.GetWritesReads)
	return nil
}

func (e *Env) GetWritesReads(c *gin.Context) {
	results, err := Retrace(c.Param("number"), c.Param("chain"), e.DB)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err) //nolint:errcheck
		return
	}
	c.JSON(http.StatusOK, results)
}

type WritesReads struct {
	Reads  []string `json:"reads"`
	Writes []string `json:"writes"`
}
type RetraceResponse struct {
	Storage WritesReads `json:"storage"`
	Account WritesReads `json:"accounts"`
}

func Retrace(blockNumber, chain string, remoteDB ethdb.KV) (RetraceResponse, error) {
	chainConfig := ReadChainConfig(remoteDB, chain)
	noOpWriter := state.NewNoopWriter()
	bn, err := strconv.Atoi(blockNumber)
	if err != nil {
		return RetraceResponse{}, err
	}
	block, err := GetBlockByNumber(remoteDB, uint64(bn))
	chainCtx := NewRemoteContext(remoteDB)
	if err != nil {
		return RetraceResponse{}, err
	}
	writer := state.NewChangeSetWriter()
	reader := NewRemoteReader(remoteDB, uint64(bn))
	intraBlockState := state.New(reader)

	if err = runBlock(intraBlockState, noOpWriter, writer, chainConfig, chainCtx, block); err != nil {
		return RetraceResponse{}, err
	}

	var output RetraceResponse
	accountChanges, _ := writer.GetAccountChanges()
	if err != nil {
		return RetraceResponse{}, err
	}
	for _, ch := range accountChanges.Changes {
		output.Account.Writes = append(output.Account.Writes, common.Bytes2Hex(ch.Key))
	}
	for _, ch := range reader.GetAccountReads() {
		output.Account.Reads = append(output.Account.Reads, common.Bytes2Hex(ch))
	}

	storageChanges, _ := writer.GetStorageChanges()
	for _, ch := range storageChanges.Changes {
		output.Storage.Writes = append(output.Storage.Writes, common.Bytes2Hex(ch.Key))
	}
	for _, ch := range reader.GetStorageReads() {
		output.Storage.Reads = append(output.Storage.Reads, common.Bytes2Hex(ch))
	}
	return output, nil
}

func runBlock(ibs *state.IntraBlockState, txnWriter state.StateWriter, blockWriter state.StateWriter,
	chainConfig *params.ChainConfig, bcb core.ChainContext, block *types.Block,
) error {
	header := block.Header()
	vmConfig := vm.Config{}
	engine := ethash.NewFullFaker()
	gp := new(core.GasPool).AddGas(block.GasLimit())
	usedGas := new(uint64)
	var receipts types.Receipts
	if chainConfig.DAOForkSupport && chainConfig.DAOForkBlock != nil && chainConfig.DAOForkBlock.Cmp(block.Number()) == 0 {
		misc.ApplyDAOHardFork(ibs)
	}
	for _, tx := range block.Transactions() {
		receipt, err := core.ApplyTransaction(chainConfig, bcb, nil, gp, ibs, txnWriter, header, tx, usedGas, vmConfig)
		if err != nil {
			return fmt.Errorf("tx %x failed: %v", tx.Hash(), err)
		}
		receipts = append(receipts, receipt)
	}
	// Finalize the block, applying any consensus engine specific extras (e.g. block rewards)
	if _, err := engine.FinalizeAndAssemble(chainConfig, header, ibs, block.Transactions(), block.Uncles(), receipts); err != nil {
		return fmt.Errorf("finalize of block %d failed: %v", block.NumberU64(), err)
	}

	ctx := chainConfig.WithEIPsFlags(context.Background(), header.Number)
	if err := ibs.CommitBlock(ctx, blockWriter); err != nil {
		return fmt.Errorf("commiting block %d failed: %v", block.NumberU64(), err)
	}
	return nil
}

func GetBlockByNumber(db ethdb.KV, number uint64) (*types.Block, error) {
	var block *types.Block
	err := db.View(context.Background(), func(tx ethdb.Tx) error {
		b, err := remotechain.GetBlockByNumber(tx, number)
		block = b
		return err
	})
	if err != nil {
		return nil, err
	}
	return block, nil
}

// ReadChainConfig retrieves the consensus settings based on the given genesis hash.
func ReadChainConfig(db ethdb.KV, chain string) *params.ChainConfig {
	var k []byte
	var data []byte
	switch chain {
	case "mainnet":
		k = params.MainnetGenesisHash[:]
	case "testnet":
		k = params.RopstenGenesisHash[:]
	case "rinkeby":
		k = params.RinkebyGenesisHash[:]
	case "goerli":
		k = params.GoerliGenesisHash[:]
	}
	_ = db.View(context.Background(), func(tx ethdb.Tx) error {
		b := tx.Bucket(dbutils.ConfigPrefix)
		d, _ := b.Get(k)
		data = d
		return nil
	})
	var config params.ChainConfig
	_ = json.Unmarshal(data, &config)
	return &config
}