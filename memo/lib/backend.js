'use strict';

var http = require('https');

class Backend {

    constructor(){
        console.log("backend created");
    }


    makeHTTPrequest(in_path, fn){
        var options = {
            host: 'bitaps.com',
            path: in_path
        };

        var req = http.get(options, function(res) {
            if(res.statusCode == 200){
                var bodyChunks = [];
                res.on('data', function(chunk) {
                    bodyChunks.push(chunk);
                }).on('end', function() {
                    fn(JSON.parse(Buffer.concat(bodyChunks)));
                });
            }else{
                fn("");
            }
        });

        req.on('error', function(e) {
            console.log('ERROR: ' + e.message);
        });
    }

    //TODO refactor this to enable less load on bitaps api
    isTxConfirmed(txid, fn){
        this.getTxInfo(txid, function(txinfo){
            var confirmed;
            if(txinfo.statusCode==200) {
                confirmed = txinfo.data.block;
            }else {
                confirmed = null;
            }
            var retval = {
                statusCode: 200,
                message: "invalid request",
                data: {tx_confirmed_in_block: confirmed}
            }
            fn(retval);
        });
    }

    getTxInfo(txid,fn){
        this.makeHTTPrequest('/api/transaction/'+txid, function(result){
            if(result!=""){
                var info = {
                    statusCode: 200,
                    data:
                        {
                        size    :   result.size,
                        hash    :   result.hash,
                        fee     :   result.fee,
                        rate    :   result.fee/result.size,
                        inputs  :   result.input,
                        outputs :   result.output,
                        block   :   result.block,
                        timestamp:  result.timestamp
                    }
                };
            }else{
                var info= {
                    statusCode: 503,
                    message: "invalid request"
                };
            }
            fn(info);
        });
    }

};

module.exports = Backend;
