package common

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// QueryDelegatorDelegations returns delegator's list of delegations
func QueryDelegatorDelegations(cliCtx context.CLIContext, queryRoute string, delegatorAddr sdk.AccAddress) ([]byte, error) {
	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", stakingtypes.QuerierRoute, stakingtypes.QueryDelegatorDelegations),
		cliCtx.Codec.MustMarshalJSON(stakingtypes.NewQueryDelegatorParams(delegatorAddr)),
	)
	return res, err
}

// WithdrawAllDelegatorRewards builds a multi-message slice to be used
// to withdraw all delegations rewards for the given delegator.
func WithdrawAllDelegatorRewards(cliCtx context.CLIContext, queryRoute string, delegatorAddr sdk.AccAddress) ([]sdk.Msg, error) {
	// retrieve the comprehensive list of all delegations
	bz, err := QueryDelegatorDelegations(cliCtx, queryRoute, delegatorAddr)
	if err != nil {
		return nil, err
	}

	var delegations stakingtypes.DelegationResponses
	if err := cliCtx.Codec.UnmarshalJSON(bz, &delegations); err != nil {
		return nil, err
	}

	// build multi-message transaction
	var msgs []sdk.Msg
	for _, delegation := range delegations {
		if delegation.Balance.LTE(sdk.NewInt(1)) {
			continue
		}

		msg := types.NewMsgWithdrawDelegatorReward(delegatorAddr, delegation.GetValidatorAddr())
		if err := msg.ValidateBasic(); err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}

	return msgs, nil
}
