#!/usr/bin/python
import sys
import sqlite3 as db
import json

conn = None

CSV_SEPERATOR = ","

states = []
spb_list = []
try:
    conn = db.connect('./db/memo.sqlite3')
    cur = conn.cursor()
    with open('./memo/public/dyn/feelevel.csv', 'w') as outfile:

        cur.execute('SELECT spb FROM feelevel GROUP BY spb')
        spb_rows = cur.fetchall()
        outfile.write("x"+CSV_SEPERATOR)
        for spb in reversed(spb_rows):
            spb_list.append(spb[0])
            outfile.write(str(spb[0])+CSV_SEPERATOR)
        outfile.write("\n")

        cur.execute('SELECT statetime FROM state')
        state_rows = cur.fetchall()

        for state in state_rows:
            states.append(state[0])

        for state in states:

            outfile.write(str(state)+CSV_SEPERATOR)

            select_string = 'SELECT spb, tally FROM feelevel NATURAL JOIN state WHERE statetime = '+ str(state)+" ORDER BY spb DESC"
            cur.execute(select_string)
            kvpair = cur.fetchall()

            kv_dict = {}
            for pair in kvpair:
                kv_dict[pair[0]] = pair[1]

            for spb in spb_list:
                if spb in kv_dict:
                    outfile.write(str(kv_dict[spb])+CSV_SEPERATOR)
                else:
                    outfile.write(""+CSV_SEPERATOR)

            outfile.write("\n")

except db.Error, e:

    print "DBError in feelevel_to_dygraph_csv.py: %s" % e.args[0]
    sys.exit(1)

finally:

    if conn:
        conn.close()
