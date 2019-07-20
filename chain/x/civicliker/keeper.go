package civicliker

import (
	"bytes"
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/likecoin/hackatom-seoul-2019/chain/x/civicliker/types"
	subTypes "github.com/likecoin/hackatom-seoul-2019/chain/x/subscription/types"
)

type SubscriptionKeeper interface {
	GetChannel(ctx sdk.Context, channelID uint64) subTypes.Channel
	GetSubscription(ctx sdk.Context, subscriber sdk.AccAddress, channelID uint64) subTypes.Subscription
	SetSubscription(ctx sdk.Context, sub subTypes.Subscription)
}

type LikeRecordsCache map[string][]types.LikeRecord // key: likers

type Keeper struct {
	coinKeeper            bank.Keeper
	subscriptionKeeper    SubscriptionKeeper
	storeKey              sdk.StoreKey
	cdc                   *codec.Codec
	blockLikeRecordsCache LikeRecordsCache
}

func NewKeeper(
	coinKeeper bank.Keeper, subscriptionKeeper SubscriptionKeeper, storeKey sdk.StoreKey, cdc *codec.Codec,
) Keeper {
	return Keeper{
		coinKeeper:            coinKeeper,
		subscriptionKeeper:    subscriptionKeeper,
		storeKey:              storeKey,
		cdc:                   cdc,
		blockLikeRecordsCache: make(map[string][]types.LikeRecord),
	}
}

var (
	LikeGlobalCountKey = []byte{0x00}
	LikeRecordKey      = []byte{0x10}
	LikeRecordCountKey = []byte{0x11}
	LikerStateKey      = []byte{0x20}
)

func (k Keeper) GetLikeGlobalCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(LikeGlobalCountKey) {
		return 0
	}
	bz := store.Get(LikeGlobalCountKey)
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetLikeGlobalCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(LikeGlobalCountKey, bz)
}

func (k Keeper) SetLikeRecord(ctx sdk.Context, r types.LikeRecord) {
	buf := bytes.Buffer{}
	buf.Write(LikeRecordKey)
	buf.Write(r.Liker.Bytes())
	buf.Write(hashUrl(r.Url))
	count := k.GetLikeGlobalCount(ctx)
	binary.Write(&buf, binary.BigEndian, count)
	store := ctx.KVStore(k.storeKey)
	key := buf.Bytes()
	store.Set(key, k.cdc.MustMarshalBinaryBare(r))
}

func (k Keeper) IterateLikerLikeHistory(ctx sdk.Context, liker sdk.AccAddress,
	fn func(index int64, r types.LikeRecord) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	buf := bytes.Buffer{}
	buf.Write(LikeRecordKey)
	buf.Write(liker.Bytes())
	prefix := buf.Bytes()
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()
	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		bz := iterator.Value()
		r := types.LikeRecord{}
		k.cdc.MustUnmarshalBinaryBare(bz, &r)
		stop := fn(i, r)
		if stop {
			break
		}
		i++
	}
}

func GetLikeRecordCountKey(liker sdk.AccAddress, url string) []byte {
	buf := bytes.Buffer{}
	buf.Write(LikeRecordCountKey)
	buf.Write(liker.Bytes())
	buf.Write(hashUrl(url))
	return buf.Bytes()
}

func (k Keeper) GetLikeRecordCount(ctx sdk.Context, liker sdk.AccAddress, url string) uint64 {
	store := ctx.KVStore(k.storeKey)
	key := GetLikeRecordCountKey(liker, url)
	if !store.Has(key) {
		return 0
	}
	bz := store.Get(key)
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetLikeRecordCount(ctx sdk.Context, liker sdk.AccAddress, url string, count uint64) {
	store := ctx.KVStore(k.storeKey)
	key := GetLikeRecordCountKey(liker, url)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(key, bz)
}

func GetLikerStateKey(liker sdk.AccAddress) []byte {
	buf := bytes.Buffer{}
	buf.Write(LikerStateKey)
	buf.Write(liker.Bytes())
	return buf.Bytes()
}

func (k Keeper) HasLikerState(ctx sdk.Context, liker sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	key := GetLikerStateKey(liker)
	return store.Has(key)
}

func (k Keeper) GetLikerState(ctx sdk.Context, liker sdk.AccAddress) types.LikerState {
	store := ctx.KVStore(k.storeKey)
	key := GetLikerStateKey(liker)
	if !store.Has(key) {
		return types.LikerState{}
	}
	bz := store.Get(key)
	s := types.LikerState{}
	k.cdc.MustUnmarshalBinaryBare(bz, &s)
	return s
}

func (k Keeper) SetLikerState(ctx sdk.Context, s types.LikerState) {
	store := ctx.KVStore(k.storeKey)
	key := GetLikerStateKey(s.Liker)
	store.Set(key, k.cdc.MustMarshalBinaryBare(s))
}

func (k Keeper) GetBlockLikeRecordsCache() LikeRecordsCache {
	return k.blockLikeRecordsCache
}
