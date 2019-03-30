import sys
import json
import configparser
import mysql.connector


def connect_to_database(config):
    db = mysql.connector.connect(
        host=config['DATABASE']['DBHOST'],
        user=config['DATABASE']['DBUSER'],
        passwd=config['DATABASE']['DBPASSWORD'],
        database=config['DATABASE']['DBNAME']
    )

    return db


if __name__ == "__main__":

    config = configparser.ConfigParser()
    config.read('config.ini')
    db = connect_to_database(config)


    mempool = ""
    for line in sys.stdin:
        mempool += line.rstrip()
        
    mempool = json.loads(mempool)

    feerates_count = {}
    feerates_size = {}

    for txid in mempool:
        tx = mempool[txid]
        feerate = int(
            round(float(tx['fee']) * 100000000 / float(tx['size'])))

        if feerate not in feerates_count:
            feerates_count[feerate] = 0
        feerates_count[feerate] += 1

        if feerate not in feerates_size:
            feerates_size[feerate] = 0
        feerates_size[feerate] += tx['size']

    cursor = db.cursor()

    sql = "UPDATE current_mempool SET byCount = %s, bySize = %s, timestamp = CURRENT_TIMESTAMP WHERE id = 1"
    val = (json.dumps(feerates_count), json.dumps(feerates_size))
    cursor.execute(sql, val)

    db.commit()
