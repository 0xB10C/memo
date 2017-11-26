#!/bin/bash

# mempool.observer base dir
BASEDIR="/path/to/memo"

cd $BASEDIR

if [ $1 -eq 3 ] # mempool stats only
then
  python ./script/mempool_stats_to_db.py
else
  python ./script/mempool_to_db.py # read mempool data and write to db
  python ./script/database_to_dygraph_csv.py $1 # generate csv files
fi
