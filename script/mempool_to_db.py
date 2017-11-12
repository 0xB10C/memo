#!/usr/bin/python
import sys
import time
import sqlite3 as db
import bitcoin
import bitcoin.rpc

# const
DATABASE_LOCATION_STRING = './db/memo.sqlite3'
CONF_FILE_PATH = None
BUCKETCOUNT = 200 # last bucket with more than 2^14 sat/byte
FEE_SPACING = 1.05 # exponential increase by 5% per bucket - https://github.com/bitcoin/bitcoin/commit/e5007bae35ce22036a816505038277d99c84e3f7#diff-8c0941572d1cdf184d1751f7b7f1db4eR109
BUCKETS = [] # stores a list of fee buckets used in core fee estimation

# var
conn = None
size = None
fee = None
bucket = None
rates = {} # stores a touple: (feerate, size, value)
bucketrates = {} # counting amount of tx in the specific bucket
rpc = bitcoin.rpc.RawProxy(None,None,CONF_FILE_PATH)


# fills the bucket list
bucketFillIndex = 0
BUCKETS.append(1) # initial bucket with 1 sat / byte
for bucketFillIndex in range(1,BUCKETCOUNT):
    BUCKETS.append(BUCKETS[bucketFillIndex-1] * FEE_SPACING) # append bucket with 5% more than the previous


# binary find bucket by feerate
def findBucketByFeerate(feerate):
    posFirst = 0
    posLast = len(BUCKETS)-1

    if feerate >= BUCKETS[posLast]:
        return posLast
    if feerate <= BUCKETS[posFirst]:
        return posFirst

    while posFirst<=posLast:
        posMid = (posFirst+posLast)//2
        if BUCKETS[posMid] <= feerate and BUCKETS[posMid+1] > feerate:
            return posMid
        if BUCKETS[posMid] > feerate:
            posLast = posMid - 1
        if BUCKETS[posMid] < feerate:
            posFirst = posMid + 1

rawmempool = rpc.getrawmempool(True);

for tx in rawmempool.items():
    fee = float(tx[1]["fee"])
    size = int(tx[1]["size"])
    rate = int((fee*100000000/size)+.5)
    bucket = findBucketByFeerate(fee*100000000/size)

    # tx per feerate
    if rate in rates:
        rates[rate] = rates[rate][0] + 1, rates[rate][1] + size, rates[rate][2] + fee
    else:
        rates[rate] = 1, size, fee

    # tx per bucket
    if bucket in bucketrates:
        bucketrates[bucket] = bucketrates[bucket] + 1
    else:
        bucketrates[bucket] = 1

# current timestamp for the new state
statetime_string = str(int(time.time()))

try:
    conn = db.connect(DATABASE_LOCATION_STRING)

    cur = conn.cursor()

    # Foreign key support is not enabled in SQLite by default. You need to enable it manually each time you connect to the database using the pragma
    # -- https://stackoverflow.com/questions/5890250/on-delete-cascade-in-sqlite3
    enable_foreign_key_support_string  = "PRAGMA foreign_keys = ON"
    cur.execute(enable_foreign_key_support_string)

    # insert new state into db with current timestamp
    insert_state_string = "INSERT INTO State (statetime) VALUES (" + statetime_string + ");"
    cur.execute(insert_state_string)

    # get the state_id from the just inserted state
    query_state_string = "SELECT state_id FROM State WHERE statetime = "+statetime_string+";"
    cur.execute(query_state_string)
    state_id = cur.fetchone()[0]

    # insert Feelevel data into db
    insert_feelevel_string = "INSERT INTO Feelevel (spb,state_id,tally,size,value) VALUES "
    for key, value in rates.iteritems():
        insert_feelevel_string += "(" + str(key) + "," + str(state_id) + "," + str(value[0]) + "," + str(value[1]) + "," + str(value[2]) + "),"
    insert_feelevel_string = insert_feelevel_string[:-1] # removes the last ,
    cur.execute(insert_feelevel_string)

    # insert Bucketlevel data into db
    insert_bucket_string = "INSERT INTO Bucketlevel (bucket,state_id,tally) VALUES "
    for key, value in bucketrates.iteritems():
        insert_bucket_string += "(" + str(key) + "," + str(state_id) + "," + str(value) +"),"
    insert_bucket_string = insert_bucket_string[:-1] # removes the last ,
    cur.execute(insert_bucket_string)

    conn.commit()

except db.Error, e:

    print "DBError in mempool_to_db: %s" % e.args[0]
    sys.exit(1)

finally:
    if conn:
        conn.close()
