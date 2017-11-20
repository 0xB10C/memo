#!/usr/bin/env python3
import bitcoin
import bitcoin.rpc
import sqlite3 as db
import time
import sys

CONF_FILE_PATH = None
DATABASE_LOCATION_STRING = './db/memo.sqlite3'

rpc = bitcoin.rpc.RawProxy(None,None,CONF_FILE_PATH,180)

# bitcoin core tx output types
# https://github.com/bitcoin/bitcoin/blob/3c098a8aa0780009c11b66b1a5d488a928629ebf/src/script/standard.cpp#L21
dict_vout_type = {  'pubkeyhash': 0,
                    'pubkey': 0,
                    'scripthash': 0,
                    'witness_v0_keyhash': 0,
                    'witness_v0_scripthash': 0,
                    'nulldata': 0,
                    'multisig': 0,
                    'nonstandard': 0,
                    'witness_unknown': 0 }

dict_vin_segwit = { 'count_segwit': 0,
                    'count_segwit_mixed': 0,
                    'count_non_segwit': 0 }

# https://stackoverflow.com/a/312464
def chunks(l, n):
    """Yield successive n-sized chunks from l."""
    for i in range(0, len(l), n):
        yield l[i:i + n]

rawmempool = rpc.getrawmempool(False);

for chunk in chunks(rawmempool, 2000): # creates chunks of 2000 txids

    commands = [ {"method": "getrawtransaction", "params": [txid,True]} for txid in chunk]
    results = rpc._batch(commands)

    for tx in results:
        hasSegWitInput, hasNonSegWitInput = False, False
        for vout in tx['result']['vout']: # for each output
            vout_type = vout['scriptPubKey']['type']
            dict_vout_type[vout_type] = dict_vout_type[vout_type] + 1 # increment type counter
        for vin in tx['result']['vin']: # for each input
            if 'txinwitness' in vin: # check if segwit input
                hasSegWitInput = True
            else:
                hasNonSegWitInput = True
        if hasSegWitInput and hasNonSegWitInput:
            dict_vin_segwit['count_segwit_mixed'] += 1
        elif hasSegWitInput and not hasNonSegWitInput:
            dict_vin_segwit['count_segwit'] += 1
        elif not hasSegWitInput and hasNonSegWitInput:
            dict_vin_segwit['count_non_segwit'] += 1

# current timestamp for the new state
measurementtime_string = str(int(time.time()))

try:
    conn = db.connect(DATABASE_LOCATION_STRING)
    cur = conn.cursor()

    # Foreign key support is not enabled in SQLite by default. You need to enable it manually each time you connect to the database using the pragma
    # -- https://stackoverflow.com/questions/5890250/on-delete-cascade-in-sqlite3
    enable_foreign_key_support_string  = "PRAGMA foreign_keys = ON"
    cur.execute(enable_foreign_key_support_string)

    insert_mempool_stats_string = "INSERT INTO Stats (measurement_time, type_pubkey, type_witness_v0_scripthash, type_pubkeyhash, type_nulldata, type_scripthash, type_witness_v0_keyhash, type_multisig, type_nonstandard, type_witness_unknown, count_segwit, count_non_segwit, count_segwit_mixed) "
    insert_mempool_stats_string += "VALUES (" + measurementtime_string + "," + str(dict_vout_type["pubkey"]) + "," + str(dict_vout_type["witness_v0_scripthash"]) + "," + str(dict_vout_type["pubkeyhash"]) +"," + str(dict_vout_type["nulldata"]) + "," + str(dict_vout_type["scripthash"]) + "," + str(dict_vout_type["witness_v0_keyhash"]) + "," + str(dict_vout_type["multisig"]) + "," + str(dict_vout_type["nonstandard"]) + "," + str(dict_vout_type["witness_unknown"]) + "," + str( dict_vin_segwit['count_segwit']) + "," + str(dict_vin_segwit['count_non_segwit']) + "," + str(dict_vin_segwit['count_segwit_mixed']) + ");"
    cur.execute(insert_mempool_stats_string)

    conn.commit()


except db.Error, e:

    print "DBError in mempool_stats_to_db: %s" % e.args[0]
    sys.exit(1)

finally:
    if conn:
        conn.close()
