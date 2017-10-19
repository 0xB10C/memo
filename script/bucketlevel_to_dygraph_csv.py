#!/usr/bin/python
import sys
import sqlite3 as db
import json

conn = None

CSV_SEPERATOR = ","

states = []
bucket_list = []
try:
    conn = db.connect('./db/memo.sqlite3')
    cur = conn.cursor()
    with open('./memo/public/dyn/bucketlevel.csv', 'w') as outfile:

        cur.execute('SELECT bucket FROM Bucketlevel GROUP BY bucket')
        bucket_rows = cur.fetchall()
        outfile.write("x"+CSV_SEPERATOR)
        for bucket in reversed(bucket_rows):
            bucket_list.append(bucket[0])
            outfile.write(str(bucket[0])+CSV_SEPERATOR)
        outfile.write("\n")

        cur.execute('SELECT statetime FROM state')
        state_rows = cur.fetchall()

        for state in state_rows:
            states.append(state[0])

        for state in states:

            outfile.write(str(state)+CSV_SEPERATOR)

            select_string = 'SELECT bucket, tally FROM Bucketlevel NATURAL JOIN state WHERE statetime = '+ str(state)+" ORDER BY bucket DESC"
            cur.execute(select_string)
            kvpair = cur.fetchall()

            kv_dict = {}
            for pair in kvpair:
                kv_dict[pair[0]] = pair[1]

            for bucket in bucket_list:
                if bucket in kv_dict:
                    outfile.write(str(kv_dict[bucket])+CSV_SEPERATOR)
                else:
                    outfile.write(""+CSV_SEPERATOR)

            outfile.write("\n")

except db.Error, e:

    print "DBError in bucketlevel_to_dygraph_csv.py: %s" % e.args[0]
    sys.exit(1)

finally:

    if conn:
        conn.close()
