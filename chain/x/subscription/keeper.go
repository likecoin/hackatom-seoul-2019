package subscription

import (
	"bytes"
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/likecoin/hackatom-seoul-2019/chain/x/subscription/types"
)

type PaymentHook func(sub *types.Subscription, ch types.Channel)
type PaymentHooks []PaymentHook

type Keeper struct {
	coinKeeper   bank.Keeper
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
	paymentHooks PaymentHooks
}

func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

var (
	ChannelKey      = []byte{0x00}
	SubscriptionKey = []byte{0x10}
)

func GetChannelKey(channelID uint64) []byte {
	buf := bytes.Buffer{}
	buf.Write(ChannelKey)
	binary.Write(&buf, binary.BigEndian, channelID)
	return buf.Bytes()
}

func (k Keeper) HasChannel(ctx sdk.Context, channelID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	key := GetChannelKey(channelID)
	return store.Has(key)
}

func (k Keeper) GetChannel(ctx sdk.Context, channelID uint64) types.Channel {
	store := ctx.KVStore(k.storeKey)
	key := GetChannelKey(channelID)
	if !store.Has(key) {
		return types.Channel{}
	}
	bz := store.Get(key)
	ch := types.Channel{}
	k.cdc.MustUnmarshalBinaryBare(bz, &ch)
	return ch
}

func (k Keeper) SetChannel(ctx sdk.Context, channel types.Channel) {
	store := ctx.KVStore(k.storeKey)
	key := GetChannelKey(channel.ID)
	store.Set(key, k.cdc.MustMarshalBinaryBare(channel))
}

func GetSubscriptionKey(subscriber sdk.AccAddress, channel uint64) []byte {
	buf := bytes.Buffer{}
	buf.Write(SubscriptionKey)
	buf.Write(subscriber.Bytes())
	binary.Write(&buf, binary.BigEndian, channel)
	return buf.Bytes()
}

func (k Keeper) IterateChannels(ctx sdk.Context, fn func(index int64, ch types.Channel) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, ChannelKey)
	defer iterator.Close()
	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		bz := iterator.Value()
		ch := types.Channel{}
		k.cdc.MustUnmarshalBinaryBare(bz, &ch)
		stop := fn(i, ch)
		if stop {
			break
		}
		i++
	}
}

func (k Keeper) HasSubscription(ctx sdk.Context, subscriber sdk.AccAddress, channelID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	key := GetSubscriptionKey(subscriber, channelID)
	return store.Has(key)
}

func (k Keeper) GetSubscription(ctx sdk.Context, subscriber sdk.AccAddress, channelID uint64) types.Subscription {
	store := ctx.KVStore(k.storeKey)
	key := GetSubscriptionKey(subscriber, channelID)
	if !store.Has(key) {
		return types.Subscription{}
	}
	bz := store.Get(key)
	sub := types.Subscription{}
	k.cdc.MustUnmarshalBinaryBare(bz, &sub)
	return sub
}

func (k Keeper) SetSubscription(ctx sdk.Context, sub types.Subscription) {
	store := ctx.KVStore(k.storeKey)
	key := GetSubscriptionKey(sub.Subscriber, sub.ChannelID)
	store.Set(key, k.cdc.MustMarshalBinaryBare(sub))
}

func (k Keeper) IterateSubscriptions(ctx sdk.Context, fn func(index int64, sub types.Subscription) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, SubscriptionKey)
	defer iterator.Close()
	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		bz := iterator.Value()
		sub := types.Subscription{}
		k.cdc.MustUnmarshalBinaryBare(bz, &sub)
		stop := fn(i, sub)
		if stop {
			break
		}
		i++
	}
}

func (k *Keeper) SetPaymentHooks(hooks PaymentHooks) {
	k.paymentHooks = hooks
}
