import resolve from 'rollup-plugin-node-resolve';
import commonjs from 'rollup-plugin-commonjs';

export default {
  input: 'src/background.js',
  output: {
    file: 'src/background.dist.js',
    format: 'esm',
  },
  // always put chromeExtension() before other plugins
  plugins: [resolve(), commonjs({
    namedExports: {
      '@lunie/cosmos-keys': ['getSeed', 'getNewWalletFromSeed', 'signWithPrivateKey'],
    },
  })],
};
