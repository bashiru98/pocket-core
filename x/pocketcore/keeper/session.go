package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
)

func (k Keeper) IsSessionBlock(ctx sdk.Context) bool {
	frequency := k.posKeeper.SessionBlockFrequency(ctx)
	return ctx.BlockHeight()%frequency == 1
}

func (k Keeper) GetLatestSessionBlock(ctx sdk.Context) sdk.Context {
	sessionBlockHeight := (ctx.BlockHeight() % k.posKeeper.SessionBlockFrequency(ctx)) + 1
	return ctx.WithBlockHeight(sessionBlockHeight)
}

func (k Keeper) IsPocketSupportedBlockchain(ctx sdk.Context, chain string) bool {
	for _, c := range k.SupportedBlockchains(ctx) {
		if c == chain {
			return true
		}
	}
	return false
}
