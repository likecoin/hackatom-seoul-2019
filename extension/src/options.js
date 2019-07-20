const app = new Vue({
  el: '#app',
  created() {
    this.updateInfo();
    this.updateTimer = setInterval(() => this.updateInfo(), 5000);
  },
  beforeDestroy() {
    if (this.updateTimer) clearInterval(this.updateTimer);
  },
  computed: {
    civicTxList() {
      return this.tx
        .sort((a, b) => a.time - b.time);
    },
  },
  methods: {
    updateInfo() {
      chrome.runtime.sendMessage({
        action: 'fetchInfo',
      }, (response) => {
        const { LIKE, LikedList } = response;
        console.log(LIKE);
        console.log(LikedList);
        const {
          subscriber,
          remaining,
        } = LIKE;
        this.address = subscriber;
        this.remainingLIKE = remaining;
        this.tx = LikedList;
      });
    },
    getTxHash(tx) {
      return tx.tx_hash;
    },
    getTxTo(tx) {
      return tx.likee;
    },
    getTxAmount(tx) {
      return tx.coin_distributed.amount;
    },
    getLikeCount(tx) {
      return tx.count;
    },
  },
  data: {
    pagination: {
      descending: true,
      rowsPerPage: -1,
    },
    address: 'Loading...',
    remainingLIKE: 'Loading...',
    tx: [],
    BIG_DIPPER_HOST: '10.100.0.110:3000',
  },
});
