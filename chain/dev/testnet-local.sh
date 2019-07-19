#!/bin/bash

set -e

CHAIN_ID="likechain-local-testnet"
MONIKER="local-dev"
PASSWORD="password"
SEED_VALIDATOR="knife what dinosaur unknown payment gallery stamp unfair turtle neither student aspect harsh divide subject mystery mandate once polar inspire wing dignity million harbor" # Address: cosmos16s47cyy5w6ja07w42s3yxe7p37pdvcrr39sc8e
SEED_FAUCET="sad ordinary multiply purpose add comfort warrior split wrestle ugly dismiss march buddy axis glove coral earth post pen object caught salute green accuse" # Address: cosmos134ckwu586qzgfhyx584rahc6lmc9vj8e6l8gu9

LIKE_HOME=$(dirname "$0")

mkdir -p "$LIKE_HOME"
pushd "$LIKE_HOME" > /dev/null
LIKE_HOME="$(pwd)"
popd > /dev/null

if [ ! -d "$LIKE_HOME/.likecli" ]; then
    mkdir -p "$LIKE_HOME/.likecli"
    printf "$PASSWORD\n$SEED_VALIDATOR\n" | \
        docker run --rm -i --volume "$LIKE_HOME/.likecli:/root/.likecli" likechain/likechain likecli keys add --recover validator
    printf "$PASSWORD\n$SEED_FAUCET\n" | \
        docker run --rm -i --volume "$LIKE_HOME/.likecli:/root/.likecli" likechain/likechain likecli keys add --recover faucet
fi

if [ ! -d "$LIKE_HOME/.liked" ]; then
    mkdir -p "$LIKE_HOME/.liked"
    docker run --rm --volume "$LIKE_HOME/.liked:/root/.liked" likechain/likechain liked init --chain-id "$CHAIN_ID" "$MONIKER" > /dev/null 2>&1
    cp "$LIKE_HOME/genesis.json" "$LIKE_HOME/.liked/config/genesis.json"

    VALIDATOR_ADDRESS=`docker run --rm --volume "$LIKE_HOME/.likecli:/root/.likecli" likechain/likechain likecli keys show validator -a`
    FAUCET_ADDRESS=`docker run --rm --volume "$LIKE_HOME/.likecli:/root/.likecli" likechain/likechain likecli keys show faucet -a`
    TM_PUBKEY=`docker run --rm --volume "$LIKE_HOME/.liked:/root/.liked" likechain/likechain liked tendermint show-validator`

    docker run --rm --volume "$LIKE_HOME/.liked:/root/.liked" likechain/likechain \
        liked add-genesis-account "$VALIDATOR_ADDRESS" 1000000000000000nanolike

    docker run --rm --volume "$LIKE_HOME/.liked:/root/.liked" likechain/likechain \
        liked add-genesis-account "$FAUCET_ADDRESS" 849000000000000000nanolike

    printf "$PASSWORD\n" | \
    docker run --rm -i --volume "$LIKE_HOME/.liked:/root/.liked" --volume "$LIKE_HOME/.likecli:/root/.likecli" likechain/likechain \
        liked gentx \
            --name validator \
            --amount 1000000000000000nanolike \
            --details "Only for local development" \
            --pubkey "$TM_PUBKEY"
    docker run --rm --volume "$LIKE_HOME/.liked:/root/.liked" likechain/likechain \
        liked collect-gentxs
else
    docker run --rm --volume "$LIKE_HOME/.liked:/root/.liked" likechain/likechain \
        liked unsafe-reset-all
fi
