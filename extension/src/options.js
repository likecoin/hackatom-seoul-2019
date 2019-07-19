const app = new Vue({
  el: '#app',
  created() {
    chrome.runtime.sendMessage({
      action: 'fetchInfo',
    }, (response) => {
      console.log(response);
      const { tx, address } = response;
      this.address = address;
      this.tx = tx;
    });
  },
  computed: {
    civicTxList() {
      return this.tx.filter(t => t.tx.type === 'auth/StdTx');
    },
  },
  methods: {
    getTxHash(tx) {
      return tx.txhash;
    },
    getTxTo(tx) {
      return tx.tx.value.msg[0].value.to_address;
    },
    getTxAmount(tx) {
      return tx.tx.value.msg[0].value.amount[0].amount;
    },
  },
  data: {
    address: 'Hello Vue!',
    tx: [],
  },
});
