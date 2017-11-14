#!/usr/bin/python
import sys
import sqlite3 as db
import json

CSV_SEPERATOR = ","
conn = None




def dbToCSV(cur,filepath,sql_key,sql_view):

    csv_buffer = "x" + CSV_SEPERATOR # csv file buffer with X for the x-axis labels and a CSV_SEPERATOR
    spb_list = []
    view_data = {}

    # fetch reversed spb (satoshi per byte) from  the according view
    # and write it in the first line of the csv file buffer
    string_fetch_spb = "SELECT spb FROM " + sql_view + " GROUP BY spb ORDER BY spb DESC"
    cur.execute(string_fetch_spb)
    spb_rows = cur.fetchall()

    for spb in spb_rows:
        spb_list.append(spb[0])
        csv_buffer += str(spb[0]) + CSV_SEPERATOR
    csv_buffer += "\n"

    string_fetch_view = "SELECT statetime, spb, " + sql_key + " FROM " + sql_view + " ORDER BY statetime, spb DESC"
    cur.execute(string_fetch_view)
    rows = cur.fetchall()

    for row in rows:
        statetime = row[0]
        spb = row[1]
        tally = row[2]

        if statetime not in view_data:
            view_data[statetime] = {}

        view_data[statetime][spb] = tally

    for key, kvpairs in view_data.iteritems():
        csv_buffer += str(key) + CSV_SEPERATOR
        for spb in spb_list:
            if spb in kvpairs:
                csv_buffer += str(kvpairs[spb]) + CSV_SEPERATOR
            else:
                csv_buffer += CSV_SEPERATOR
        csv_buffer += "\n"

    with open(filepath, 'w') as outfile:
            outfile.write(csv_buffer)
    pass


try:
    conn = db.connect('./db/memo.sqlite3')
    cur = conn.cursor()

    if int(sys.argv[1]) >= 4:
        dbToCSV(cur,"./memo/public/dyn/amount4h.csv","tally","v_4hData")
        dbToCSV(cur,"./memo/public/dyn/size4h.csv","size","v_4hData")
        dbToCSV(cur,"./memo/public/dyn/value4h.csv","value","v_4hData")

    if int(sys.argv[1]) >= 24:
        dbToCSV(cur,"./memo/public/dyn/size24h.csv","size","v_24hData")
        dbToCSV(cur,"./memo/public/dyn/value24h.csv","value","v_24hData")
        dbToCSV(cur,"./memo/public/dyn/amount24h.csv","tally","v_24hData")

    if int(sys.argv[1]) >= 168:
        dbToCSV(cur,"./memo/public/dyn/amount7d.csv","tally","v_7dData")
        dbToCSV(cur,"./memo/public/dyn/size7d.csv","size","v_7dData")
        dbToCSV(cur,"./memo/public/dyn/value7d.csv","value","v_7dData")

except db.Error, e:
    print "DBError in database_to_dygraph_csv.py: %s" % e.args[0]
    if conn:
        conn.close()
    sys.exit(1)

finally:
    if conn:
        conn.close()
