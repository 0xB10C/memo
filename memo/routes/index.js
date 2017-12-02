var express = require('express');
var router = express.Router();
const config = require('../config.js');



/* GET home page. */
router.get('/', function(req, res, next) {
        res.render('index.ejs', {load_txid:"", GAToken: config.GAnalyticsToken});
});

router.get('/:tx_id', function(req, res, next) {
    var invalid
    var tx_pattern = /[a-fA-F0-9]{64}/;
    var input = req.params.tx_id;
    var txid="";
    if(tx_pattern.test(input)){
        res.render('index.ejs', {load_txid:input, GAToken: config.GAnalyticsToken});
    }else{
        res.status(404);
        res.render('error.html')
    }
});

module.exports = router;
