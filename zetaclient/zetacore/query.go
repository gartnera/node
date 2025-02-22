package zetacore

import (
	"context"
	"fmt"
	"sort"
	"time"

	sdkmath "cosmossdk.io/math"
	tmhttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"google.golang.org/grpc"

	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
)

// GetCrosschainFlags returns the crosschain flags
func (c *Client) GetCrosschainFlags() (observertypes.CrosschainFlags, error) {
	client := observertypes.NewQueryClient(c.grpcConn)
	resp, err := client.CrosschainFlags(context.Background(), &observertypes.QueryGetCrosschainFlagsRequest{})
	if err != nil {
		return observertypes.CrosschainFlags{}, err
	}
	return resp.CrosschainFlags, nil
}

// GetBlockHeaderEnabledChains returns the enabled chains for block headers
func (c *Client) GetBlockHeaderEnabledChains() ([]lightclienttypes.HeaderSupportedChain, error) {
	client := lightclienttypes.NewQueryClient(c.grpcConn)
	resp, err := client.HeaderEnabledChains(context.Background(), &lightclienttypes.QueryHeaderEnabledChainsRequest{})
	if err != nil {
		return []lightclienttypes.HeaderSupportedChain{}, err
	}
	return resp.HeaderEnabledChains, nil
}

// GetRateLimiterFlags returns the rate limiter flags
func (c *Client) GetRateLimiterFlags() (crosschaintypes.RateLimiterFlags, error) {
	client := crosschaintypes.NewQueryClient(c.grpcConn)
	resp, err := client.RateLimiterFlags(context.Background(), &crosschaintypes.QueryRateLimiterFlagsRequest{})
	if err != nil {
		return crosschaintypes.RateLimiterFlags{}, err
	}
	return resp.RateLimiterFlags, nil
}

// GetChainParamsForChainID returns the chain params for a given chain ID
func (c *Client) GetChainParamsForChainID(externalChainID int64) (*observertypes.ChainParams, error) {
	client := observertypes.NewQueryClient(c.grpcConn)
	resp, err := client.GetChainParamsForChain(
		context.Background(),
		&observertypes.QueryGetChainParamsForChainRequest{ChainId: externalChainID},
	)
	if err != nil {
		return &observertypes.ChainParams{}, err
	}
	return resp.ChainParams, nil
}

// GetChainParams returns all the chain params
func (c *Client) GetChainParams() ([]*observertypes.ChainParams, error) {
	client := observertypes.NewQueryClient(c.grpcConn)
	var err error

	resp := &observertypes.QueryGetChainParamsResponse{}
	for i := 0; i <= DefaultRetryCount; i++ {
		resp, err = client.GetChainParams(context.Background(), &observertypes.QueryGetChainParamsRequest{})
		if err == nil {
			return resp.ChainParams.ChainParams, nil
		}
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return nil, fmt.Errorf("failed to get chain params | err %s", err.Error())
}

// GetUpgradePlan returns the current upgrade plan
func (c *Client) GetUpgradePlan() (*upgradetypes.Plan, error) {
	client := upgradetypes.NewQueryClient(c.grpcConn)

	resp, err := client.CurrentPlan(context.Background(), &upgradetypes.QueryCurrentPlanRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Plan, nil
}

// GetAllCctx returns all cross chain transactions
func (c *Client) GetAllCctx() ([]*crosschaintypes.CrossChainTx, error) {
	client := crosschaintypes.NewQueryClient(c.grpcConn)
	resp, err := client.CctxAll(context.Background(), &crosschaintypes.QueryAllCctxRequest{})
	if err != nil {
		return nil, err
	}
	return resp.CrossChainTx, nil
}

// GetCctxByHash returns a cross chain transaction by hash
func (c *Client) GetCctxByHash(sendHash string) (*crosschaintypes.CrossChainTx, error) {
	client := crosschaintypes.NewQueryClient(c.grpcConn)
	resp, err := client.Cctx(context.Background(), &crosschaintypes.QueryGetCctxRequest{Index: sendHash})
	if err != nil {
		return nil, err
	}
	return resp.CrossChainTx, nil
}

// GetCctxByNonce returns a cross chain transaction by nonce
func (c *Client) GetCctxByNonce(chainID int64, nonce uint64) (*crosschaintypes.CrossChainTx, error) {
	client := crosschaintypes.NewQueryClient(c.grpcConn)
	resp, err := client.CctxByNonce(context.Background(), &crosschaintypes.QueryGetCctxByNonceRequest{
		ChainID: chainID,
		Nonce:   nonce,
	})
	if err != nil {
		return nil, err
	}
	return resp.CrossChainTx, nil
}

// GetObserverList returns the list of observers
func (c *Client) GetObserverList() ([]string, error) {
	var err error
	client := observertypes.NewQueryClient(c.grpcConn)

	for i := 0; i <= DefaultRetryCount; i++ {
		resp, err := client.ObserverSet(context.Background(), &observertypes.QueryObserverSet{})
		if err == nil {
			return resp.Observers, nil
		}
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return nil, err
}

// GetRateLimiterInput returns input data for the rate limit checker
func (c *Client) GetRateLimiterInput(window int64) (crosschaintypes.QueryRateLimiterInputResponse, error) {
	client := crosschaintypes.NewQueryClient(c.grpcConn)
	maxSizeOption := grpc.MaxCallRecvMsgSize(32 * 1024 * 1024)
	resp, err := client.RateLimiterInput(
		context.Background(),
		&crosschaintypes.QueryRateLimiterInputRequest{
			Window: window,
		},
		maxSizeOption,
	)
	if err != nil {
		return crosschaintypes.QueryRateLimiterInputResponse{}, err
	}
	return *resp, nil
}

// ListPendingCctx returns a list of pending cctxs for a given chainID
//   - The max size of the list is crosschainkeeper.MaxPendingCctxs
func (c *Client) ListPendingCctx(chainID int64) ([]*crosschaintypes.CrossChainTx, uint64, error) {
	client := crosschaintypes.NewQueryClient(c.grpcConn)
	maxSizeOption := grpc.MaxCallRecvMsgSize(32 * 1024 * 1024)
	resp, err := client.ListPendingCctx(
		context.Background(),
		&crosschaintypes.QueryListPendingCctxRequest{
			ChainId: chainID,
		},
		maxSizeOption,
	)
	if err != nil {
		return nil, 0, err
	}
	return resp.CrossChainTx, resp.TotalPending, nil
}

// ListPendingCctxWithinRatelimit returns a list of pending cctxs that do not exceed the outbound rate limit
//   - The max size of the list is crosschainkeeper.MaxPendingCctxs
//   - The returned `rateLimitExceeded` flag indicates if the rate limit is exceeded or not
func (c *Client) ListPendingCctxWithinRatelimit() ([]*crosschaintypes.CrossChainTx, uint64, int64, string, bool, error) {
	client := crosschaintypes.NewQueryClient(c.grpcConn)
	maxSizeOption := grpc.MaxCallRecvMsgSize(32 * 1024 * 1024)
	resp, err := client.ListPendingCctxWithinRateLimit(
		context.Background(),
		&crosschaintypes.QueryListPendingCctxWithinRateLimitRequest{},
		maxSizeOption,
	)
	if err != nil {
		return nil, 0, 0, "", false, err
	}
	return resp.CrossChainTx, resp.TotalPending, resp.CurrentWithdrawWindow, resp.CurrentWithdrawRate, resp.RateLimitExceeded, nil
}

// GetAbortedZetaAmount returns the amount of zeta that has been aborted
func (c *Client) GetAbortedZetaAmount() (string, error) {
	client := crosschaintypes.NewQueryClient(c.grpcConn)
	resp, err := client.ZetaAccounting(context.Background(), &crosschaintypes.QueryZetaAccountingRequest{})
	if err != nil {
		return "", err
	}
	return resp.AbortedZetaAmount, nil
}

// GetGenesisSupply returns the genesis supply
func (c *Client) GetGenesisSupply() (sdkmath.Int, error) {
	tmURL := fmt.Sprintf("http://%s", c.cfg.ChainRPC)
	s, err := tmhttp.New(tmURL, "/websocket")
	if err != nil {
		return sdkmath.ZeroInt(), err
	}
	res, err := s.Genesis(context.Background())
	if err != nil {
		return sdkmath.ZeroInt(), err
	}
	appState, err := genutiltypes.GenesisStateFromGenDoc(*res.Genesis)
	if err != nil {
		return sdkmath.ZeroInt(), err
	}
	bankstate := banktypes.GetGenesisStateFromAppState(c.encodingCfg.Codec, appState)
	return bankstate.Supply.AmountOf(config.BaseDenom), nil
}

// GetZetaTokenSupplyOnNode returns the zeta token supply on the node
func (c *Client) GetZetaTokenSupplyOnNode() (sdkmath.Int, error) {
	client := banktypes.NewQueryClient(c.grpcConn)
	resp, err := client.SupplyOf(context.Background(), &banktypes.QuerySupplyOfRequest{Denom: config.BaseDenom})
	if err != nil {
		return sdkmath.ZeroInt(), err
	}
	return resp.GetAmount().Amount, nil
}

// GetLastBlockHeight returns the last block height
func (c *Client) GetLastBlockHeight() ([]*crosschaintypes.LastBlockHeight, error) {
	client := crosschaintypes.NewQueryClient(c.grpcConn)
	resp, err := client.LastBlockHeightAll(context.Background(), &crosschaintypes.QueryAllLastBlockHeightRequest{})
	if err != nil {
		c.logger.Error().Err(err).Msg("query GetBlockHeight error")
		return nil, err
	}
	return resp.LastBlockHeight, nil
}

// GetLatestZetaBlock returns the latest zeta block
func (c *Client) GetLatestZetaBlock() (*tmservice.Block, error) {
	client := tmservice.NewServiceClient(c.grpcConn)
	res, err := client.GetLatestBlock(context.Background(), &tmservice.GetLatestBlockRequest{})
	if err != nil {
		return nil, err
	}
	return res.SdkBlock, nil
}

// GetNodeInfo returns the node info
func (c *Client) GetNodeInfo() (*tmservice.GetNodeInfoResponse, error) {
	var err error

	client := tmservice.NewServiceClient(c.grpcConn)
	for i := 0; i <= DefaultRetryCount; i++ {
		res, err := client.GetNodeInfo(context.Background(), &tmservice.GetNodeInfoRequest{})
		if err == nil {
			return res, nil
		}
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return nil, err
}

// GetBlockHeight returns the zetachain block height
func (c *Client) GetBlockHeight() (int64, error) {
	client := crosschaintypes.NewQueryClient(c.grpcConn)
	resp, err := client.LastZetaHeight(context.Background(), &crosschaintypes.QueryLastZetaHeightRequest{})
	if err != nil {
		return 0, err
	}
	return resp.Height, nil
}

// GetBaseGasPrice returns the base gas price
func (c *Client) GetBaseGasPrice() (int64, error) {
	client := feemarkettypes.NewQueryClient(c.grpcConn)
	resp, err := client.Params(context.Background(), &feemarkettypes.QueryParamsRequest{})
	if err != nil {
		return 0, err
	}
	if resp.Params.BaseFee.IsNil() {
		return 0, fmt.Errorf("base fee is nil")
	}
	return resp.Params.BaseFee.Int64(), nil
}

// GetBallotByID returns a ballot by ID
func (c *Client) GetBallotByID(id string) (*observertypes.QueryBallotByIdentifierResponse, error) {
	client := observertypes.NewQueryClient(c.grpcConn)
	return client.BallotByIdentifier(context.Background(), &observertypes.QueryBallotByIdentifierRequest{
		BallotIdentifier: id,
	})
}

// GetNonceByChain returns the nonce by chain
func (c *Client) GetNonceByChain(chain chains.Chain) (observertypes.ChainNonces, error) {
	client := observertypes.NewQueryClient(c.grpcConn)
	resp, err := client.ChainNonces(
		context.Background(),
		&observertypes.QueryGetChainNoncesRequest{Index: chain.ChainName.String()},
	)
	if err != nil {
		return observertypes.ChainNonces{}, err
	}
	return resp.ChainNonces, nil
}

// GetAllNodeAccounts returns all node accounts
func (c *Client) GetAllNodeAccounts() ([]*observertypes.NodeAccount, error) {
	client := observertypes.NewQueryClient(c.grpcConn)
	resp, err := client.NodeAccountAll(context.Background(), &observertypes.QueryAllNodeAccountRequest{})
	if err != nil {
		return nil, err
	}
	c.logger.Debug().Msgf("GetAllNodeAccounts: %d", len(resp.NodeAccount))
	return resp.NodeAccount, nil
}

// GetKeyGen returns the keygen
func (c *Client) GetKeyGen() (*observertypes.Keygen, error) {
	var err error
	client := observertypes.NewQueryClient(c.grpcConn)

	for i := 0; i <= ExtendedRetryCount; i++ {
		resp, err := client.Keygen(context.Background(), &observertypes.QueryGetKeygenRequest{})
		if err == nil {
			return resp.Keygen, nil
		}
		time.Sleep(DefaultRetryInterval * time.Second)
	}
	return nil, fmt.Errorf("failed to get keygen | err %s", err.Error())
}

// GetBallot returns a ballot by ID
func (c *Client) GetBallot(ballotIdentifier string) (*observertypes.QueryBallotByIdentifierResponse, error) {
	client := observertypes.NewQueryClient(c.grpcConn)
	resp, err := client.BallotByIdentifier(context.Background(), &observertypes.QueryBallotByIdentifierRequest{
		BallotIdentifier: ballotIdentifier,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetInboundTrackersForChain returns the inbound trackers for a chain
func (c *Client) GetInboundTrackersForChain(chainID int64) ([]crosschaintypes.InboundTracker, error) {
	client := crosschaintypes.NewQueryClient(c.grpcConn)
	resp, err := client.InboundTrackerAllByChain(
		context.Background(),
		&crosschaintypes.QueryAllInboundTrackerByChainRequest{ChainId: chainID},
	)
	if err != nil {
		return nil, err
	}
	return resp.InboundTracker, nil
}

// GetCurrentTss returns the current TSS
func (c *Client) GetCurrentTss() (observertypes.TSS, error) {
	client := observertypes.NewQueryClient(c.grpcConn)
	resp, err := client.TSS(context.Background(), &observertypes.QueryGetTSSRequest{})
	if err != nil {
		return observertypes.TSS{}, err
	}
	return resp.TSS, nil
}

// GetEthTssAddress returns the ETH TSS address
// TODO(revamp): rename to EVM
func (c *Client) GetEthTssAddress() (string, error) {
	client := observertypes.NewQueryClient(c.grpcConn)
	resp, err := client.GetTssAddress(context.Background(), &observertypes.QueryGetTssAddressRequest{})
	if err != nil {
		return "", err
	}
	return resp.Eth, nil
}

// GetBtcTssAddress returns the BTC TSS address
func (c *Client) GetBtcTssAddress(chainID int64) (string, error) {
	client := observertypes.NewQueryClient(c.grpcConn)
	resp, err := client.GetTssAddress(context.Background(), &observertypes.QueryGetTssAddressRequest{
		BitcoinChainId: chainID,
	})
	if err != nil {
		return "", err
	}
	return resp.Btc, nil
}

// GetTssHistory returns the TSS history
func (c *Client) GetTssHistory() ([]observertypes.TSS, error) {
	client := observertypes.NewQueryClient(c.grpcConn)
	resp, err := client.TssHistory(context.Background(), &observertypes.QueryTssHistoryRequest{})
	if err != nil {
		return nil, err
	}
	return resp.TssList, nil
}

// GetOutboundTracker returns the outbound tracker for a chain and nonce
func (c *Client) GetOutboundTracker(chain chains.Chain, nonce uint64) (*crosschaintypes.OutboundTracker, error) {
	client := crosschaintypes.NewQueryClient(c.grpcConn)
	resp, err := client.OutboundTracker(context.Background(), &crosschaintypes.QueryGetOutboundTrackerRequest{
		ChainID: chain.ChainId,
		Nonce:   nonce,
	})
	if err != nil {
		return nil, err
	}
	return &resp.OutboundTracker, nil
}

// GetAllOutboundTrackerByChain returns all outbound trackers for a chain
func (c *Client) GetAllOutboundTrackerByChain(
	chainID int64,
	order interfaces.Order,
) ([]crosschaintypes.OutboundTracker, error) {
	client := crosschaintypes.NewQueryClient(c.grpcConn)
	resp, err := client.OutboundTrackerAllByChain(
		context.Background(),
		&crosschaintypes.QueryAllOutboundTrackerByChainRequest{
			Chain: chainID,
			Pagination: &query.PageRequest{
				Key:        nil,
				Offset:     0,
				Limit:      2000,
				CountTotal: false,
				Reverse:    false,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if order == interfaces.Ascending {
		sort.SliceStable(resp.OutboundTracker, func(i, j int) bool {
			return resp.OutboundTracker[i].Nonce < resp.OutboundTracker[j].Nonce
		})
	}
	if order == interfaces.Descending {
		sort.SliceStable(resp.OutboundTracker, func(i, j int) bool {
			return resp.OutboundTracker[i].Nonce > resp.OutboundTracker[j].Nonce
		})
	}
	return resp.OutboundTracker, nil
}

// GetPendingNoncesByChain returns the pending nonces for a chain and current tss address
func (c *Client) GetPendingNoncesByChain(chainID int64) (observertypes.PendingNonces, error) {
	client := observertypes.NewQueryClient(c.grpcConn)
	resp, err := client.PendingNoncesByChain(
		context.Background(),
		&observertypes.QueryPendingNoncesByChainRequest{ChainId: chainID},
	)
	if err != nil {
		return observertypes.PendingNonces{}, err
	}
	return resp.PendingNonces, nil
}

// GetBlockHeaderChainState returns the block header chain state
func (c *Client) GetBlockHeaderChainState(chainID int64) (lightclienttypes.QueryGetChainStateResponse, error) {
	client := lightclienttypes.NewQueryClient(c.grpcConn)
	resp, err := client.ChainState(context.Background(), &lightclienttypes.QueryGetChainStateRequest{ChainId: chainID})
	if err != nil {
		return lightclienttypes.QueryGetChainStateResponse{}, err
	}
	return *resp, nil
}

// GetSupportedChains returns the supported chains
func (c *Client) GetSupportedChains() ([]chains.Chain, error) {
	client := observertypes.NewQueryClient(c.grpcConn)
	resp, err := client.SupportedChains(context.Background(), &observertypes.QuerySupportedChains{})
	if err != nil {
		return nil, err
	}
	return resp.GetChains(), nil
}

// GetAdditionalChains returns the additional chains
func (c *Client) GetAdditionalChains() ([]chains.Chain, error) {
	client := authoritytypes.NewQueryClient(c.grpcConn)
	resp, err := client.ChainInfo(context.Background(), &authoritytypes.QueryGetChainInfoRequest{})
	if err != nil {
		return nil, err
	}
	return resp.GetChainInfo().Chains, nil
}

// GetPendingNonces returns the pending nonces
func (c *Client) GetPendingNonces() (*observertypes.QueryAllPendingNoncesResponse, error) {
	client := observertypes.NewQueryClient(c.grpcConn)
	resp, err := client.PendingNoncesAll(context.Background(), &observertypes.QueryAllPendingNoncesRequest{})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Prove returns whether a proof is valid
func (c *Client) Prove(
	blockHash string,
	txHash string,
	txIndex int64,
	proof *proofs.Proof,
	chainID int64,
) (bool, error) {
	client := lightclienttypes.NewQueryClient(c.grpcConn)
	resp, err := client.Prove(context.Background(), &lightclienttypes.QueryProveRequest{
		BlockHash: blockHash,
		TxIndex:   txIndex,
		Proof:     proof,
		ChainId:   chainID,
		TxHash:    txHash,
	})
	if err != nil {
		return false, err
	}
	return resp.Valid, nil
}

// HasVoted returns whether an observer has voted
func (c *Client) HasVoted(ballotIndex string, voterAddress string) (bool, error) {
	client := observertypes.NewQueryClient(c.grpcConn)
	resp, err := client.HasVoted(context.Background(), &observertypes.QueryHasVotedRequest{
		BallotIdentifier: ballotIndex,
		VoterAddress:     voterAddress,
	})
	if err != nil {
		return false, err
	}
	return resp.HasVoted, nil
}

// GetZetaHotKeyBalance returns the zeta hot key balance
func (c *Client) GetZetaHotKeyBalance() (sdkmath.Int, error) {
	client := banktypes.NewQueryClient(c.grpcConn)
	address, err := c.keys.GetAddress()
	if err != nil {
		return sdkmath.ZeroInt(), err
	}
	resp, err := client.Balance(context.Background(), &banktypes.QueryBalanceRequest{
		Address: address.String(),
		Denom:   config.BaseDenom,
	})
	if err != nil {
		return sdkmath.ZeroInt(), err
	}
	return resp.Balance.Amount, nil
}
