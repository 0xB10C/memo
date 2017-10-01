var express = require('express');
var router = express.Router();
var Backend = require('../lib/backend.js');

var backend = new Backend();

var TX_PATTERN = /[a-fA-F0-9]{64}/;

router.get('/tx/:tx_id', function(req, res, next) {
    var input = req.params.tx_id;
    var txid="";
    if(TX_PATTERN.test(input)){
        backend.getTxInfo(input,function(data){
            res.send(data);
        });
    }else {
        var info= {
            statusCode: 503,
            message: "invalid request"
        };
        res.send(info);
    }
});

router.get('/confirmed/:tx_id', function(req, res, next) {
    var input = req.params.tx_id;
    var txid="";
    if(TX_PATTERN.test(input)){
        backend.isTxConfirmed(input,function(data){
            res.send(data);
        });
    }else {
        var info= {
            statusCode: 503,
            message: "invalid request"
        };
        res.send(info);
    }
});

module.exports = router;
