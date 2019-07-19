// import 'babel-polyfill';
import * as cosmosKeys from '@lunie/cosmos-keys';
import Cosmos from '@lunie/cosmos-api';
import { Buffer } from 'buffer';

const CHAIN_ID = 'likechain-cosmos-testnet-2';
const DENOM = 'nanolike';

const { getSeed, getNewWalletFromSeed, signWithPrivateKey } = cosmosKeys;
const api = new Cosmos('http://35.226.174.222:1317', CHAIN_ID);

let globalSigner;
let globalAddress;

function likeToNanolike(value) {
  return `${Number.parseInt(value, 10).toString()}000000000`;
}

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

function likeToAmount(value) {
  return { denom: DENOM, amount: likeToNanolike(value) };
}
async function sendTx(msgCallPromise) {
  const { simulate, send } = await msgCallPromise;
  const gas = (await simulate({})).toString();
  const { included } = await send({ gas }, globalSigner);
  await included();
}

async function transfer(transferTo, transferValue = 1) {
  const from = normalizeAddress(globalAddress);
  const toAddress = normalizeAddress(transferTo);
  const value = transferValue;
  const amount = likeToAmount(value);
  const msgPromise = api.MsgSend(from, { toAddress, amounts: [amount] });
  sendTx(msgPromise);
}

async function notify(payload) {
  if (payload.action === 'civicLike') {
    // TODO: sign and broadcast
    const { wallet, sourceURL } = payload;
    console.log(`liking ${wallet} ${sourceURL}`);
    await transfer(wallet);
  }
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
