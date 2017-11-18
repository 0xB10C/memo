#!/usr/bin/python
import sys
import sqlite3 as db
import json

CSV_SEPERATOR = ","
conn = None


def dbToCSV(cur,s_type,filepath,sql_key,sql_view):

    csv_buffer = "x" + CSV_SEPERATOR # csv file buffer with X for the x-axis labels and a CSV_SEPERATOR
    type_list = []
    view_data = {}

    # fetch reversed s_type from  the according view
    # and write it in the first line of the csv file buffer
    string_fetch_type = "SELECT " + s_type + " FROM " + sql_view + " GROUP BY " + s_type + " ORDER BY " + s_type + " DESC"
    cur.execute(string_fetch_type)
    type_rows = cur.fetchall()


    for d_type in type_rows:
        type_list.append(d_type[0])
        csv_buffer += str(d_type[0]) + CSV_SEPERATOR
    csv_buffer = csv_buffer[:-1]
    csv_buffer += "\n"

    string_fetch_view = "SELECT statetime, " + s_type + ", " + sql_key + " FROM " + sql_view + " ORDER BY statetime, " + s_type + " DESC"
    cur.execute(string_fetch_view)
    rows = cur.fetchall()


    for row in rows:
        statetime = row[0]
        d_type = row[1]
        tally = row[2]

        if statetime not in view_data:
            view_data[statetime] = {}

        view_data[statetime][d_type] = tally

    for key, kvpairs in view_data.iteritems():
        csv_buffer += str(key) + CSV_SEPERATOR
        for d_type in type_list:
            if d_type in kvpairs:
                csv_buffer += str(kvpairs[d_type]) + CSV_SEPERATOR
            else:
                csv_buffer += CSV_SEPERATOR
        csv_buffer = csv_buffer[:-1]
        csv_buffer += "\n"

    with open(filepath, 'w') as outfile:
            outfile.write(csv_buffer)
    pass


try:
    conn = db.connect('./db/memo.sqlite3')
    cur = conn.cursor()

    if int(sys.argv[1]) >= 4:
        dbToCSV(cur,"spb","./memo/public/dyn/feelevel_amount4h.csv","tally","v_4hData_feelevel")
        dbToCSV(cur,"spb","./memo/public/dyn/feelevel_size4h.csv","size","v_4hData_feelevel")
        dbToCSV(cur,"spb","./memo/public/dyn/feelevel_value4h.csv","value","v_4hData_feelevel")

        dbToCSV(cur,"bucket","./memo/public/dyn/bucket_amount4h.csv","tally","v_4hData_bucketlevel")
        dbToCSV(cur,"bucket","./memo/public/dyn/bucket_size4h.csv","size","v_4hData_bucketlevel")
        dbToCSV(cur,"bucket","./memo/public/dyn/bucket_value4h.csv","value","v_4hData_bucketlevel")

    if int(sys.argv[1]) >= 24:
        dbToCSV(cur,"spb","./memo/public/dyn/feelevel_size24h.csv","size","v_24hData_feelevel")
        dbToCSV(cur,"spb","./memo/public/dyn/feelevel_value24h.csv","value","v_24hData_feelevel")
        dbToCSV(cur,"spb","./memo/public/dyn/feelevel_amount24h.csv","tally","v_24hData_feelevel")

        dbToCSV(cur,"bucket","./memo/public/dyn/bucket_amount24h.csv","tally","v_24hData_bucketlevel")
        dbToCSV(cur,"bucket","./memo/public/dyn/bucket_size24h.csv","size","v_24hData_bucketlevel")
        dbToCSV(cur,"bucket","./memo/public/dyn/bucket_value24h.csv","value","v_24hData_bucketlevel")

    if int(sys.argv[1]) >= 168:
        dbToCSV(cur,"spb","./memo/public/dyn/feelevel_amount7d.csv","tally","v_7dData_feelevel")
        dbToCSV(cur,"spb","./memo/public/dyn/feelevel_size7d.csv","size","v_7dData_feelevel")
        dbToCSV(cur,"spb","./memo/public/dyn/feelevel_value7d.csv","value","v_7dData_feelevel")

        dbToCSV(cur,"bucket","./memo/public/dyn/bucket_amount7d.csv","tally","v_7dData_bucketlevel")
        dbToCSV(cur,"bucket","./memo/public/dyn/bucket_size7d.csv","size","v_7dData_bucketlevel")
        dbToCSV(cur,"bucket","./memo/public/dyn/bucket_value7d.csv","value","v_7dData_bucketlevel")


except db.Error, e:
    print "DBError in database_to_dygraph_csv.py: %s" % e
    if conn:
        conn.close()
    sys.exit(1)

finally:
    if conn:
        conn.close()
