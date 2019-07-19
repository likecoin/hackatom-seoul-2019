package contentdb

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/likecoin/hackatom-seoul-2019/chain/x/subscription/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
}

func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

func (k Keeper) HasContent(ctx sdk.Context, url string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(url))
}

func (k Keeper) GetContent(ctx sdk.Context, url string) types.Content {
	store := ctx.KVStore(k.storeKey)
	key := []byte(url)
	if !store.Has(key) {
		return types.Content{}
	}
	bz := store.Get(key)
	var content types.Content
	k.cdc.MustUnmarshalBinaryBare(bz, &content)
	return content
}

func (k Keeper) SetContent(ctx sdk.Context, content types.Content) {
	store := ctx.KVStore(k.storeKey)
	key := []byte(content.Url)
	store.Set(key, k.cdc.MustMarshalBinaryBare(content))
}
