import json
import datetime
import configparser


from flask import Flask, Response, Blueprint, jsonify
from flaskext.mysql import MySQL
import pymysql

config = configparser.ConfigParser()
config.read('config.ini')

app = Flask(__name__)
mysql = MySQL()

# MySQL configurations
app.config['MYSQL_DATABASE_USER'] = config['DATABASE']['DBUSER']
app.config['MYSQL_DATABASE_PASSWORD'] = config['DATABASE']['DBPASSWORD']
app.config['MYSQL_DATABASE_DB'] = config['DATABASE']['DBNAME']
app.config['MYSQL_DATABASE_HOST'] = config['DATABASE']['DBHOST']
mysql.init_app(app)

conn = mysql.connect()

CORS = "*"


def executeSQL(sql, conn):
    """ This function wraps the cursor.execute function to reconnect 
        when the database connection gets closed """
    try:
        cursor = conn.cursor()
        cursor.execute(sql)
        return cursor
    except pymysql.err.OperationalError:
        conn = mysql.connect()  # reconnect
        cursor = conn.cursor()
        cursor.execute(sql)
        return cursor


@app.route("/api/mempool/<path:by>", methods=['GET'])
def get_current_mempool(by):

    if by == "byCount" or by == "bySize":

        sql = "SELECT timestamp, %s FROM current_mempool WHERE id = 1" % (by)

        cursor = executeSQL(sql, conn)
        timestamp, data = cursor.fetchone()
        cursor.close()

        resp = jsonify({
            'timestamp': (timestamp - datetime.datetime(1970, 1, 1)).total_seconds(),
            'mempoolData': json.loads(data)
            })
        
        resp.status_code = 200
        return resp
    else:
        return Response("500 Internal Server Error" + by, status=500)


@app.after_request
def after_request(response):
    header = response.headers
    header['Access-Control-Allow-Origin'] = CORS
    return response


if __name__ == "__main__":
    app.run(host='0.0.0.0')
