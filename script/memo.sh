#!/bin/bash

# mempool.observer base dir
BASEDIR="/path/to/memo"

cd $BASEDIR

# read mempool data and write to db
python ./script/mempool_to_db.py

if [ $1 -eq 3 ] # mempool stat identifyer
then
  python ./script/mempool_stats_to_db.py
fi

# generate csv files
python ./script/bucketlevel_to_dygraph_csv.py
python ./script/database_to_dygraph_csv.py $1
