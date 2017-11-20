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

def statsToCSV():

    # transaction output type stats

    csv_buffer = "x" + CSV_SEPERATOR + "multisig" + CSV_SEPERATOR + "nonstandard" + CSV_SEPERATOR + "nulldata" + CSV_SEPERATOR + "pubkey" + CSV_SEPERATOR + "pubkeyhash" + CSV_SEPERATOR + "scripthash" + CSV_SEPERATOR + "witness_unknown" + CSV_SEPERATOR + "witness_v0_keyhash" + CSV_SEPERATOR + "witness_v0_scripthash" + "\n"

    string_fetch_types = "SELECT measurement_time, type_multisig, type_nonstandard, type_nulldata, type_pubkey, type_pubkeyhash, type_scripthash, type_witness_unknown, type_witness_v0_keyhash, type_witness_v0_scripthash FROM Stats"
    cur.execute(string_fetch_types)
    vout_type_rows = cur.fetchall()

    for row in vout_type_rows:

        csv_buffer += str(row[0]) + CSV_SEPERATOR + str(row[1]) + CSV_SEPERATOR + str(row[2]) + CSV_SEPERATOR + str(row[3]) + CSV_SEPERATOR + str(row[4]) + CSV_SEPERATOR + str(row[5]) + CSV_SEPERATOR + str(row[6]) + CSV_SEPERATOR + str(row[7]) + CSV_SEPERATOR + str(row[8]) + CSV_SEPERATOR + str(row[9]) + "\n"
    with open("./memo/public/dyn/stats_output_type.csv", 'w') as outfile:
        outfile.write(csv_buffer)

    # transaction segwit stats

    csv_buffer = "x" + CSV_SEPERATOR + "non-segwit" + CSV_SEPERATOR + "segwit-mixed" + CSV_SEPERATOR + "segwit" + "\n"

    string_fetch_segwit = "SELECT measurement_time, count_non_segwit, count_segwit_mixed, count_segwit FROM Stats"
    cur.execute(string_fetch_segwit)
    segwit_rows = cur.fetchall()

    for row in segwit_rows:
        csv_buffer += str(row[0]) + CSV_SEPERATOR + str(row[1]) + CSV_SEPERATOR + str(row[2]) + CSV_SEPERATOR + str(row[3]) + "\n"
    with open("./memo/public/dyn/stats_segwit.csv", 'w') as outfile:
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

        statsToCSV()

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
