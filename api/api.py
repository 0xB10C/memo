import json
import datetime
import configparser
import mysql.connector

from flask import Flask, Response, Blueprint, jsonify

app = Flask(__name__)

config = configparser.ConfigParser()
config.read('config.ini')

db = mysql.connector.connect(
    host=config['DATABASE']['DBHOST'],
    user=config['DATABASE']['DBUSER'],
    passwd=config['DATABASE']['DBPASSWORD'],
    database=config['DATABASE']['DBNAME']
)

CORS = "*"


@app.route("/api/mempool/<path:by>", methods=['GET'])
def get_current_mempool(by):

    if by == "byCount" or by == "bySize":
        cursor = db.cursor()
        sql = "SELECT timestamp, %s FROM current_mempool WHERE id = 1" % (by)

        cursor.execute(sql)
        timestamp, data = cursor.fetchone()

        a = {'timestamp': (timestamp - datetime.datetime(1970, 1, 1)
                           ).total_seconds(), 'data': json.loads(data)}

        resp = jsonify(a)
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
