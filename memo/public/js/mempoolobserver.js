var colorSet = ["#373854","#493267","#7bb3ff","#e86af0","#7bb3ff","#9e379f"];
const FEE_SPACING = 1.05;
var graph;
var tx_unconfirmed_timer;
var tx_unconfirmed_timer_last_block;
var txid;

function checkConfirmed() {
    $.ajax({url: "/api/confirmed/"+txid, success: function(result){
        if(result.statusCode=="200"){
            // tx is confirmed
            if(tx_unconfirmed_timer_last_block == null && result.data.tx_confirmed_in_block > 1){

                if (Notification.permission === "granted") {

                    const title = 'Transaction confirmed!';

                    const options = {
                      body: txid.substring(0, 8) + "... just confirmed!\nblock: " + result.data.tx_confirmed_in_block,
                      icon: '/img/og_preview.png',
                      sound: '/mp3/attention-seeker.mp3' // October 2017: no browser supports sounds in notifications yet
                    };

                    var notification = new Notification(title,options);
                }

                new Audio('/mp3/attention-seeker.mp3').play();

                tx_unconfirmed_timer_last_block = result.data.tx_confirmed_in_block;
                $('#card_info_tx_confirmed_sound_checkbox').prop("checked", false );
                $("#card_info_tx_confirmed_block").text(result.data.tx_confirmed_in_block);
                $("#card_info_tx_confirmed_sound").css("display", "none");
                $("#card_info_tx_confirmed").css("display", "block");
                window.clearInterval(tx_unconfirmed_timer) // disables the timer
            }
        }else{
            alert(result.statusCode+" - "+result.message);
        }
    }});
}

function changeGraphColor(index) {
    graph.setSelection(false, ''+index,false);
}

function setCursorTextFeelevel(date,key,value) {
    date = new Date(date*1000);
    var time = ("0" + date.getHours()).slice(-2) + ":" + ("0" + date.getMinutes()).slice(-2);
    if(value>0){
        $("#cursortext_detailed").css("visibility","visible");
        $("#cursortext_detailed_tally").text(value);
        $("#cursortext_detailed_SpB").text(key);
        $("#cursortext_detailed_date").text(time);
    }else{
        $("#cursortext_detailed").css("visibility","hidden");
    }
}
function setCursorTextBucketlevel(date,key,value) {
    date = new Date(date*1000);
    var time = ("0" + date.getHours()).slice(-2) + ":" + ("0" + date.getMinutes()).slice(-2);
    if(value>0){
        let fee = calcFeeForBucket(key);
        $("#cursortext_bucket").css("visibility","visible");
        $("#cursortext_buckets_tally").text(value);
        $("#cursortext_buckets_bucket").text(key + " (" + fee.toFixed(2) + "s/B - " + (fee * FEE_SPACING).toFixed(2) + "s/B)");
        $("#cursortext_buckets_date").text(time);
    }else{
        $("#cursortext_bucket").css("visibility","hidden");
    }
}
function setCursorTextValuelevel(date,key,value) {
    date = new Date(date*1000);
    var time = ("0" + date.getHours()).slice(-2) + ":" + ("0" + date.getMinutes()).slice(-2);
    if(value>0){
        $("#cursortext_detailed_value").css("visibility","visible");
        $("#cursortext_detailed_value_value").text(value);
        $("#cursortext_detailed_value_SpB").text(key);
        $("#cursortext_detailed_value_date").text(time);
    }else{
        $("#cursortext_detailed_value").css("visibility","hidden");
    }
}
function setCursorTextSizelevel(date,key,value) {
    date = new Date(date*1000);
    var time = ("0" + date.getHours()).slice(-2) + ":" + ("0" + date.getMinutes()).slice(-2);
    if(value>0){
        $("#cursortext_detailed_size").css("visibility","visible");
        $("#cursortext_detailed_size_size").text((value/1000000).toFixed(6));
        $("#cursortext_detailed_size_SpB").text(key);
        $("#cursortext_detailed_size_date").text(time);
    }else{
        $("#cursortext_detailed_size").css("visibility","hidden");
    }
}

function loadTXinfo(){
    if($("#input_txid").val().length==64){
        $.ajax({url: "/api/tx/"+$("#input_txid").val(), success: function(result){
            if(result.statusCode=="200"){
                txid = result.data.hash;
                tx_unconfirmed_timer_last_block = result.data.block;
                $("#card_info_tx_id").text(result.data.hash);
                $("#card_info_tx_permalink").attr("href", "https://mempool.observer/"+result.data.hash);
                if(result.data.block!=null){ // transaction confirmed
                    $("#card_info_tx_confirmed_block").text(result.data.block);
                    $("#card_info_tx_confirmed").css("display", "block");
                    $("#card_info_tx_confirmed_sound").css("display", "none");
                }else { // transaction unconfirmed (block = null -> not inclued in a block)
                    $("#card_info_tx_confirmed").css("display", "none");
                    $("#card_info_tx_confirmed_sound").css("display", "block");
                    changeGraphColor(result.data.rate.toFixed(0));
                }
                $("#card_info_tx_input").text("");
                $("#card_info_tx_output").text("");
                result.data.inputs.forEach(function(input) {
                    $("#card_info_tx_input").append("<p>"+input.address[0]+"<br><span class='red_text'>"+(input.amount*0.00000001).toFixed(8)+" BTC</span></p>" );
                });

                $("#card_info_tx_middle_SpB").text(result.data.rate.toFixed(0)+" sat/byte");
                $("#card_info_tx_middle_fee").text(result.data.fee +" sat");
                $("#card_info_tx_middle_size").text(result.data.size+" byte");
                $("#card_info_tx_middle_timestamp").text((new Date(result.data.timestamp * 1000)).toUTCString());

                result.data.outputs.forEach(function(output) {
                    $("#card_info_tx_output").append( "<p>"+output.address[0]+"<br><span class='red_text'>"+(output.amount*0.00000001).toFixed(8)+" BTC</span></p>" );
                });

                $("#card_info_tx").css("display", "block");

            }else{
                alert(result.statusCode+" - "+result.message);
            }
        }});
    }else {
        alert("invalid input: '" + $("#input_txid").val() + "' - length: [" +  $("#input_txid").val().length + "]");
    }
}

var buildGraph = function(data,options) {
    var chart = document.getElementById("graphdiv");
    var div = document.createElement('div');
    div.className = "chartclass";
    div.style.display = 'inline-block';
    chart.appendChild(div);
    var labels = data[1];
    graph = new Dygraph(div, data, options)
}

var options_detailed = {
    width: 1000,
    height: 650,
    colors: colorSet,
    fillAlpha: 1,
    strokeWidth: 4,
    strokeBorderWidth: 0,
    highlightCircleSize: 0,
    stackedGraph: true,
    stackedGraphNaNFill:"none",
    rightGap: 0,

    legend: "never",
    ylabel: "transactions in mempool",
    xlabel: "time",

    highlightSeriesOpts: {
        strokeWidth: 5,
        strokeBorderWidth: 2,
        strokeBorderColor: "#FDD",
        highlightCircleSize: 1
    },
    interactionModel:  {},
    highlightSeriesBackgroundAlpha: 0.5,
    highlightSeriesBackgroundColor: "#000",
    highlightCallback: function(e, x, pts, row) {
        setCursorTextFeelevel(x,graph.getHighlightSeries(),graph.rolledSeries_[graph.rolledSeries_.length-graph.getHighlightSeries()-1][row][1]);
    },
    unhighlightCallback: function(e) {
        $("#cursortext").css("visibility","hidden");
    },
    axes: {
        x: {
            axisLabelFormatter: function(d, gran, opts) {
                return Dygraph.dateAxisLabelFormatter(new Date((d/60).toFixed(0)*60000), gran, opts);
            }
        },
        y: {
            axisLabelFormatter: function(y) {
                if(y>999){
                    return + y/1000 + 'k';
                }else{
                    return y;
                }
            },
            axisLabelWidth: 50,
            includeZero:true
        }
    }
}

var options_detailed_size = {
    width: 1000,
    height: 650,
    colors: colorSet,
    fillAlpha: 1,
    strokeWidth: 4,
    strokeBorderWidth: 0,
    highlightCircleSize: 0,
    stackedGraph: true,
    stackedGraphNaNFill:"none",
    rightGap: 0,

    legend: "never",
    ylabel: "size in mempool [MB]",
    xlabel: "time",

    highlightSeriesOpts: {
        strokeWidth: 5,
        strokeBorderWidth: 2,
        strokeBorderColor: "#FDD",
        highlightCircleSize: 1
    },
    interactionModel:  {},
    highlightSeriesBackgroundAlpha: 0.5,
    highlightSeriesBackgroundColor: "#000",
    highlightCallback: function(e, x, pts, row) {
        setCursorTextSizelevel(x,graph.getHighlightSeries(),graph.rolledSeries_[graph.rolledSeries_.length-graph.getHighlightSeries()-1][row][1]);
    },
    unhighlightCallback: function(e) {
        $("#cursortext").css("visibility","hidden");
    },
    axes: {
        x: {
            axisLabelFormatter: function(d, gran, opts) {
                return Dygraph.dateAxisLabelFormatter(new Date((d/60).toFixed(0)*60000), gran, opts);
            }
        },
        y: {
            axisLabelFormatter: function(y) {
                return + y/1000000;
            },
            axisLabelWidth: 50,
            includeZero:true
        }
    }
}

var options_detailed_value = {
    width: 1000,
    height: 650,
    colors: colorSet,
    fillAlpha: 1,
    strokeWidth: 4,
    strokeBorderWidth: 0,
    highlightCircleSize: 0,
    stackedGraph: true,
    stackedGraphNaNFill:"none",
    rightGap: 0,

    legend: "never",
    ylabel: "fees in mempool [BTC]",
    xlabel: "time",

    highlightSeriesOpts: {
        strokeWidth: 5,
        strokeBorderWidth: 2,
        strokeBorderColor: "#FDD",
        highlightCircleSize: 1
    },
    interactionModel:  {},
    highlightSeriesBackgroundAlpha: 0.5,
    highlightSeriesBackgroundColor: "#000",
    highlightCallback: function(e, x, pts, row) {
        setCursorTextValuelevel(x,graph.getHighlightSeries(),graph.rolledSeries_[graph.rolledSeries_.length-graph.getHighlightSeries()-1][row][1]);
    },
    unhighlightCallback: function(e) {
        $("#cursortext").css("visibility","hidden");
    },
    axes: {
        x: {
            axisLabelFormatter: function(d, gran, opts) {
                return Dygraph.dateAxisLabelFormatter(new Date((d/60).toFixed(0)*60000), gran, opts);
            }
        },
        y: {
            axisLabelFormatter: function(y) {
                return y.toFixed(2);
            },
            axisLabelWidth: 50,
            includeZero:true
        }
    }
}

var options_bucket = {
    width: 1000,
    height: 650,
    colors: colorSet,
    fillAlpha: 1,
    strokeWidth: 4,
    strokeBorderWidth: 0,
    highlightCircleSize: 0,
    stackedGraph: true,
    stackedGraphNaNFill:"none",
    rightGap: 0,

    legend: "never",
    ylabel: "transactions in mempool",
    xlabel: "time",

    highlightSeriesOpts: {
        strokeWidth: 5,
        strokeBorderWidth: 2,
        strokeBorderColor: "#FDD",
        highlightCircleSize: 1
    },
    interactionModel:  {},
    highlightSeriesBackgroundAlpha: 0.5,
    highlightSeriesBackgroundColor: "#000",
    highlightCallback: function(e, x, pts, row) {
        setCursorTextBucketlevel(
            x,
            graph.getHighlightSeries(), // bucket
            graph.rolledSeries_[graph.rolledSeries_.length-graph.getHighlightSeries()-2][row][1]
        );
    },
    unhighlightCallback: function(e) {
        $("#cursortext").css("visibility","hidden");
    },
    axes: {
        x: {
            axisLabelFormatter: function(d, gran, opts) {
                return Dygraph.dateAxisLabelFormatter(new Date((d/60).toFixed(0)*60000), gran, opts);
            }
        },
        y: {
            axisLabelFormatter: function(y) {
                if(y>999){
                    return + y/1000 + 'k';
                }else{
                    return y;
                }
            },
            axisLabelWidth: 50,
            includeZero:true
        }
    }
}

function calcFeeForBucket(bucket) {
    var feeForBucket = 1;
    for(i = 0; i<bucket; i++){
        feeForBucket = feeForBucket * 1.05;
    }
    return feeForBucket;
}

window.onload = function () {

    // register onClickListener for the 'load tx id'-button
    $("#button_load_tx_info").click(function(){
        loadTXinfo();
    });

    buildGraph('/dyn/feelevel.csv', options_detailed);

    // if 'load_txid' was defined by the ejs renderer load the tx
    // 'load_txid' is the permalink txid
    if (typeof load_txid !== 'undefined') {
        $("#input_txid").val(load_txid)
        loadTXinfo();
    }

    // onChangeListener for the confirm sound checkbox
    $('#card_info_tx_confirmed_sound_checkbox').change(function(){
        if($(this).is(':checked')) {
            new Audio('mp3/attention-seeker.mp3').play(); // plays audio once

            if (!("Notification" in window)) {
                alert("This browser does not support desktop notification");
            }

            if (Notification.permission !== 'denied') {
                Notification.requestPermission();
            }

            tx_unconfirmed_timer = setInterval(checkConfirmed, 60000); // checks with the backend (over ajax) every 60 seconds if the transaction is confirmed
        } else {
            window.clearInterval(tx_unconfirmed_timer) // disables the timer
        }
    });

    $('a[data-toggle="tab"]').on('shown.bs.tab', function (e) {
        switch (e.target.id) {
            case "nav-detailed-tab":
                graph.updateOptions($.extend(options_detailed, {file: "/dyn/feelevel.csv"}),false);
            break;
            case "nav-bucket-tab":
                graph.updateOptions($.extend(options_bucket, {file: "/dyn/bucketlevel.csv"}),false);
            break;
            case "nav-detailed-size-tab":
                graph.updateOptions($.extend(options_detailed_size, {file: "/dyn/sizelevel.csv"}),false);
            break;
            case "nav-detailed-value-tab":
                graph.updateOptions($.extend(options_detailed_value, {file: "/dyn/valuelevel.csv"}),false);
            break;
        }
    })

}
