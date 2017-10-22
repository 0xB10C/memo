#!/bin/bash

# mempool.observer base dir
BASEDIR="/path/to/memo"

# bitcoin-cli location
BITCOINCLI=/usr/bin/bitcoin-cli

cd $BASEDIR
$BITCOINCLI getrawmempool true | python ./script/mempool_to_db.py
python ./script/bucketlevel_to_dygraph_csv.py
python ./script/feelevel_to_dygraph_csv.py
python ./script/sizelevel_to_dygraph_csv.py
python ./script/valuelevel_to_dygraph_csv.py
