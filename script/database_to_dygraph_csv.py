#!/usr/bin/python
import sys
import sqlite3 as db
import json

CSV_SEPERATOR = ","
conn = None


def dbToCSV(cur,filepath,sql_key,sql_view):
    print "starting: " + filepath + " | " + sql_key + " | " + sql_view

    spb_list = []
    states = []

    with open(filepath, 'w') as outfile:

        # fetch spb (satoshi per byte) from  the arrording view
        # and write it in the first line of the csv file
        string_fetch_spb = "SELECT spb FROM " + sql_view + " GROUP BY spb"
        cur.execute(string_fetch_spb)
        spb_rows = cur.fetchall()
        outfile.write("x"+CSV_SEPERATOR)
        for spb in reversed(spb_rows):
            spb_list.append(spb[0])
            outfile.write(str(spb[0])+CSV_SEPERATOR)
        outfile.write("\n")

        # fetch all states in the view
        string_fetch_statetimes = "SELECT statetime FROM " + sql_view + " GROUP BY statetime"
        cur.execute(string_fetch_statetimes)
        state_rows = cur.fetchall()

        # transform sql rows into list
        for state in state_rows:
            states.append(state[0])

        for state in states:
            # begin the line with the timestamp (statetime)
            outfile.write(str(state)+CSV_SEPERATOR)

            string_fetch_keyvalue = "SELECT spb, " + sql_key + " FROM " + sql_view + " WHERE statetime = " + str(state) + " ORDER BY spb DESC"
            print string_fetch_keyvalue
            cur.execute(string_fetch_keyvalue)
            kvpair = cur.fetchall()

            # create dictionary with (spb => sql_key)
            kv_dict = {}
            for pair in kvpair:
                kv_dict[pair[0]] = pair[1]

            # write values to file
            for spb in spb_list:
                if spb in kv_dict:
                    outfile.write(str(kv_dict[spb])+CSV_SEPERATOR)
                else:
                    outfile.write(""+CSV_SEPERATOR)

            outfile.write("\n")

    pass


try:
    conn = db.connect('./db/memo.sqlite3')
    cur = conn.cursor()

    dbToCSV(cur,"./memo/public/dyn/amount4h.csv","tally","v_4hData")
    dbToCSV(cur,"./memo/public/dyn/amount24h.csv","tally","v_24hData")
    dbToCSV(cur,"./memo/public/dyn/amount7d.csv","tally","v_7dData")

    dbToCSV(cur,"./memo/public/dyn/size4h.csv","size","v_4hData")
    dbToCSV(cur,"./memo/public/dyn/size24h.csv","size","v_24hData")
    dbToCSV(cur,"./memo/public/dyn/size7d.csv","size","v_7dData")

    dbToCSV(cur,"./memo/public/dyn/value4h.csv","value","v_4hData")
    dbToCSV(cur,"./memo/public/dyn/value24h.csv","value","v_24hData")
    dbToCSV(cur,"./memo/public/dyn/value7d.csv","value","v_7dData")

except db.Error, e:
    print "DBError in database_to_dygraph_csv.py: %s" % e.args[0]
    if conn:
        conn.close()
    sys.exit(1)

finally:
    if conn:
        conn.close()
