#!/usr/bin/python
import sys
import time
import sqlite3 as db
import re

# const
DATABASE_LOCATION_STRING = './db/memo.sqlite3'


# var
conn = None
size = None
fee = None
rates = {}

for line in sys.stdin:
    re_size = re.search('(?<=\"size\": )(.*)(?=,)', line)
    if re_size:
        size = int(re_size.group(0))
    re_fee = re.search('(?<=\"fee\": )(.*)(?=,)',line)
    if re_fee:
        fee = float(re_fee.group(0))
        rate = int((fee*100000000/size)+.5)
        if rate in rates:
            rates[rate]=rates[rate]+1
        else:
            rates[rate] = 1

statetime_string = str(int(time.time()))

try:
    conn = db.connect(DATABASE_LOCATION_STRING)

    cur = conn.cursor()

    # Foreign key support is not enabled in SQLite by default. You need to enable it manually each time you connect to the database using the pragma
    # -- https://stackoverflow.com/questions/5890250/on-delete-cascade-in-sqlite3
    enable_foreign_key_support_string  = "PRAGMA foreign_keys = ON"
    cur.execute(enable_foreign_key_support_string)

    insert_state_string = "INSERT INTO State (statetime) VALUES (" + statetime_string + ");"
    cur.execute(insert_state_string)


    query_state_string = "SELECT state_id FROM State WHERE statetime = "+statetime_string+";"
    cur.execute(query_state_string)
    state_id = cur.fetchone()[0]


    insert_feelevel_string = "INSERT INTO Feelevel (spb,state_id,tally) VALUES "
    for key, value in rates.iteritems():
        insert_feelevel_string += "(" + str(key) + "," + str(state_id) + "," + str(value) +"),"
    insert_feelevel_string = insert_feelevel_string[:-1] # removes the last ,
    cur.execute(insert_feelevel_string)

    conn.commit()

except db.Error, e:

    print "DBError in mempool_to_db: %s" % e.args[0]
    sys.exit(1)

finally:
    if conn:
        conn.close()
