package keeper

import (
	"fmt"

	"github.com/andromedaprotocol/andromedad/x/distribution/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DistributeFromFeePool distributes funds from the distribution module account to
// a receiver address while updating the community pool
func (k Keeper) DistributeFromFeePool(ctx sdk.Context, amount sdk.Coins, receiveAddr sdk.AccAddress) error {
	feePool := k.GetFeePool(ctx)

	// NOTE the community pool isn't a module account, however its coins
	// are held in the distribution module account. Thus the community pool
	// must be reduced separately from the SendCoinsFromModuleToAccount call
	newPool, negative := feePool.CommunityPool.SafeSub(sdk.NewDecCoinsFromCoins(amount...))
	if negative {
		return fmt.Errorf("insufficient funds to distribute from the community pool")
	}

	feePool.CommunityPool = newPool

	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiveAddr, amount)
	if err != nil {
		return err
	}

	k.SetFeePool(ctx, feePool)
	return nil
}
