// import 'babel-polyfill';
import * as cosmosKeys from '@lunie/cosmos-keys';
import Cosmos from '@lunie/cosmos-api';
import { Buffer } from 'buffer';

const CHAIN_ID = 'likechain-local-testnet';

const { getSeed, getNewWalletFromSeed, signWithPrivateKey } = cosmosKeys;
const api = new Cosmos('http://10.100.0.110:1317', CHAIN_ID);

let globalSigner;
let globalAddress;

function normalizeAddress(address) {
  return address.replace(/\s/g, '');
}

function getStorage(key) {
  return new Promise(((resolve) => {
    chrome.storage.sync.get(key, (result) => {
      resolve(result[key]);
    });
  }));
}
function setStorage(key, value) {
  return new Promise(((resolve) => {
    chrome.storage.sync.set({ [key]: value }, (result) => {
      resolve(result);
    });
  }));
}

function parseMnemonic(mnemonic) {
  const wallet = getNewWalletFromSeed(mnemonic.replace(/\n/g, ''));
  const publicKey = Buffer.from(wallet.publicKey, 'hex');
  const privateKey = Buffer.from(wallet.privateKey, 'hex');
  const signer = (signMessage) => {
    const signature = signWithPrivateKey(signMessage, privateKey);
    return { signature, publicKey };
  };
  const address = wallet.cosmosAddress;
  return {
    wallet, publicKey, privateKey, signer, address,
  };
}

async function sendTx(msgCallPromise) {
  const { send } = await msgCallPromise;
  const res = await send({ gas: '200000' }, globalSigner);
  const { included } = res;
  await included();
}

async function MsgLIKE(
  senderAddress,
  {
    toAddress,
    sourceURL,
    likeCount,
  },
) {
  const message = {
    type: 'civicliker/like',
    value: {
      liker: senderAddress,
      likee: toAddress,
      url: sourceURL,
      count: likeCount,
    },
  };
  return {
    message,
    simulate: ({ memo = undefined }) => api.simulate(senderAddress, { message, memo }),
    send: (
      { gas, gasPrices, memo = undefined },
      signer,
    ) => api.send(senderAddress, { gas, gasPrices, memo }, message, signer),
  };
}

async function like(transferTo, sourceURL, likeCount = '1') {
  const from = normalizeAddress(globalAddress);
  const toAddress = normalizeAddress(transferTo);
  const msgPromise = MsgLIKE(from, {
    toAddress,
    sourceURL,
    likeCount,
  });
  sendTx(msgPromise);
}

function notify(payload, sender, sendResponse) {
  if (payload.action === 'civicLike') {
    // TODO: sign and broadcast
    const { wallet, sourceURL } = payload;
    console.log(`liking ${wallet} ${sourceURL}`);
    like(wallet, sourceURL).then(sendResponse({}));
    return true;
  }
  if (payload.action === 'fetchInfo') {
    Promise.all([
      fetch(`http://10.100.0.110:1317/civicliker/like-history/${globalAddress}`)
        .then(resp => resp.json()),
      fetch(`http://10.100.0.110:1317/subscription/subscription/${globalAddress}/1`)
        .then(resp => resp.json()),
    ]).then((results) => {
      sendResponse({
        LikedList: results[0],
        LIKE: results[1],
      });
    });

    return true;
  }
  return false;
}
chrome.runtime.onMessage.addListener(notify);

async function main() {
  let mnemonic = await getStorage('wallet');
  if (!mnemonic) {
    mnemonic = getSeed();
    await setStorage('wallet', mnemonic);
  }
  console.log(mnemonic);
  const { address, signer } = parseMnemonic(mnemonic);
  globalAddress = address;
  globalSigner = signer;
  console.log(await api.get.txs(address));
}

main();
