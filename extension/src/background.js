import * as cosmosKeys from '@lunie/cosmos-keys';

const { getSeed, getNewWalletFromSeed, signWithPrivateKey } = cosmosKeys;

function hexToArrayBuffer(hex) {
  if (typeof hex !== 'string') {
    throw new TypeError('Expected input to be a string');
  }

  if ((hex.length % 2) !== 0) {
    throw new RangeError('Expected string to be an even number of characters');
  }
  const view = new Uint8Array(hex.length / 2);
  for (let i = 0; i < hex.length; i += 2) {
    view[i / 2] = parseInt(hex.substring(i, i + 2), 16);
  }
  return view.buffer;
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
  const publicKey = hexToArrayBuffer(wallet.publicKey);
  const privateKey = hexToArrayBuffer(wallet.privateKey);
  const signer = (signMessage) => {
    const signature = signWithPrivateKey(signMessage, privateKey);
    return { signature, publicKey };
  };
  const address = wallet.cosmosAddress;
  return {
    wallet, publicKey, privateKey, signer, address,
  };
}

function notify(payload) {
  if (payload.action === 'civicLike') {
    // TODO: sign and broadcast
    const { wallet, sourceURL } = payload;
    console.log(`liking ${wallet} ${sourceURL}`);
  }
}
chrome.runtime.onMessage.addListener(notify);

async function main() {
  let mnemonic = await getStorage('wallet');
  if (!mnemonic) {
    mnemonic = getSeed();
    await setStorage('wallet', mnemonic);
  }
  parseMnemonic(mnemonic);
}

main();
