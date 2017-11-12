#!/bin/bash

# mempool.observer base dir
BASEDIR="/path/to/memo"

cd $BASEDIR
python ./script/mempool_to_db.py
python ./script/bucketlevel_to_dygraph_csv.py
python ./script/database_to_dygraph_csv.py
