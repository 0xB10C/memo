# memo - mempool.observer

## motivation

Seemingly stuck and longtime-unconfirmed transactions can be quite annoying for Bitcoin users.

Memo (mempool.observer) periodically takes a snapshot of a nodes memory pool. Recent data is then visualized and users can lookup the position of their transaction. The idea is to inform about the memory pool and to provide tool to double-check estimated fee/size-ratios.

>"What's going to happen to Bitcoin?" is the wrong question. The right question is "What are you going to contribute?" &mdash; <cite>[Greg Maxwell](https://twitter.com/nullc_)</cite>


## todo
* install.md
* more detail in the transaction info
* live updating view of the mempool visualization
* expand explanation + resources section, maybe even a complete overhaul
* replace bitaps api and use own node for transaction data
* more and different visualizations:
    * numberOfTx + sizeOfTx in the last X blocks
    * median / avg confirming fee in the last X blocks
    * fee estimators in comparison + lowest fee that was included in the past block

For everything else please open an issue or contact [me](https://twitter.com/0xb10c).

## thanks to
* https://bitaps.com: for the api
* http://dygraphs.com: for the fantastic js-charting library
