package blockchain

import (
	"sync/atomic"

	"MyCode/copernicus/btcutil"

	"github.com/btcboost/copernicus/algorithm"
	"github.com/btcboost/copernicus/mempool"
	"github.com/btcboost/copernicus/utils"
	"github.com/btcboost/copernicus/utxo"
)

type BlockMap struct {
	Data map[utils.Hash]*BlockIndex
}

// ChainState store the blockchain global state
type ChainState struct {
	ChainAcTive       Chain
	MapBlockIndex     BlockMap
	PindexBestInvalid *BlockIndex

	//* The set of all CBlockIndex entries with BLOCK_VALID_TRANSACTIONS (for itself
	//* and all ancestors) and as good as our current tip or better. Entries may be
	//* failed, though, and pruning nodes may be missing the data for the block.
	setBlockIndexCandidates *algorithm.CustomSet

	// All pairs A->B, where A (or one of its ancestors) misses transactions, but B
	// has transactions. Pruned nodes may have entries where B is missing data.
	MapBlocksUnlinked map[*BlockIndex][]*BlockIndex
}

// Global status for blockchain
var (
	//GChainState Global unique variables
	GChainState    ChainState
	GfImporting    atomic.Value
	GfReindex      atomic.Value
	GMaxTipAge     int64
	Gmempool       *mempool.Mempool
	GpcoinsTip     *utxo.CoinsViewCache
	Gpblocktree    *BlockTreeDB
	GminRelayTxFee utils.FeeRate
)

var (
	// GfHavePruned Pruning-related variables and constants, True if any block files have ever been pruned.
	GfHavePruned = false
	GfPruneMode  = false
	GfTxIndex    = false

	//GindexBestHeader Best header we've seen so far (used for getheaders queries' starting points)
	GindexBestHeader *BlockIndex

	//GfCheckForPruning Global flag to indicate we should check to see if there are block/undo files
	//* that should be deleted. Set on startup or if we allocate more file space when
	//* we're in prune mode.
	GfCheckForPruning    = false
	GfCheckpointsEnabled = DEFAULT_CHECKPOINTS_ENABLED
	GfCheckBlockIndex    = false
	GfRequireStandard    = true
	GfIsBareMultisigStd  = DEFAULT_PERMIT_BAREMULTISIG
)

const (
	// MAX_BLOCKFILE_SIZE The maximum size of a blk?????.dat file (since 0.8)  // 128 MiB
	MAX_BLOCKFILE_SIZE = 0x8000000
	// BLOCKFILE_CHUNK_SIZE The pre-allocation chunk size for blk?????.dat files (since 0.8)  // 16 MiB
	BLOCKFILE_CHUNK_SIZE = 0x1000000
	// UNDOFILE_CHUNK_SIZE The pre-allocation chunk size for rev?????.dat files (since 0.8) // 1 MiB
	UNDOFILE_CHUNK_SIZE                     = 0x100000
	DEFAULT_MIN_RELAY_TX_FEE btcutil.Amount = 1000
)

func init() {
	GChainState.MapBlockIndex.Data = make(map[utils.Hash]*BlockIndex)
	GChainState.MapBlocksUnlinked = make(map[*BlockIndex][]*BlockIndex)
	GChainState.setBlockIndexCandidates = algorithm.NewCustomSet(BlockIndexWorkComparator)
	GfReindex.Store(false)
	GfImporting.Store(false)
	GMaxTipAge = DEFAULT_MAX_TIP_AGE
	GminRelayTxFee.SataoshisPerK = int64(DEFAULT_MIN_RELAY_TX_FEE)
	Gmempool = mempool.NewMemPool(GminRelayTxFee)
}