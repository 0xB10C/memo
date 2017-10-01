#!/bin/bash

# mempool.observer base dir
BASEDIR="/path/to/mempool.observer"

# bitcoin-cli location
BITCOINCLI=/usr/bin/bitcoin-cli

cd $BASEDIR
$BITCOINCLI getrawmempool true | python ./script/mempool_to_db.py
python ./script/db_to_dygraph_csv.py
