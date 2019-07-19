import resolve from 'rollup-plugin-node-resolve';
import commonjs from 'rollup-plugin-commonjs';

export default {
  input: 'src/background.js',
  output: {
    file: 'src/background.dist.js',
    format: 'esm',
  },
  plugins: [resolve({ preferBuiltins: false }), commonjs({
    namedExports: {
      '@lunie/cosmos-keys': ['getSeed', 'getNewWalletFromSeed', 'signWithPrivateKey'],
    },
  })],
};
