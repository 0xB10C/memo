const colorSet = ["#373854","#493267","#7bb3ff","#e86af0","#7bb3ff","#9e379f"];
const FEE_SPACING = 1.05;
const weekdays = ["Sunday","Monday","Tuesday","Wednesday","Thursday","Friday","Saturday"];

var graph;
var tx_unconfirmed_timer;
var tx_unconfirmed_timer_last_block;
var txid;
var isFullscreen = false;
var setNextStepFullscreenON = false;




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


function setCursorText(date,key,value,chartType) {
    if(value>0){
        date = new Date(date*1000);

        var time;

        switch (chartType.timespan) {
            case 168:
                time = weekdays[date.getDay()]+" "+("0" + date.getHours()).slice(-2) + ":" + ("0" + date.getMinutes()).slice(-2);
                break;
            default:
                time = ("0" + date.getHours()).slice(-2) + ":" + ("0" + date.getMinutes()).slice(-2);
        }

        switch (chartType.name) {
            case "amount":
                $("#cursortext_detailed").css("visibility","visible");
                $("#cursortext_detailed_tally").text(value);
                $("#cursortext_detailed_SpB").text(key);
                $("#cursortext_detailed_date").text(time);
                break;
            case "size":
                $("#cursortext_detailed_size").css("visibility","visible");
                $("#cursortext_detailed_size_size").text((value/1000000).toFixed(6));
                $("#cursortext_detailed_size_SpB").text(key);
                $("#cursortext_detailed_size_date").text(time);
                break;
            case "fee":
                $("#cursortext_detailed_value").css("visibility","visible");
                $("#cursortext_detailed_value_value").text(value);
                $("#cursortext_detailed_value_SpB").text(key);
                $("#cursortext_detailed_value_date").text(time);
                break;
            case "bucket":
                let fee = calcFeeForBucket(key);
                $("#cursortext_bucket").css("visibility","visible");
                $("#cursortext_buckets_tally").text(value);
                $("#cursortext_buckets_bucket").text(key + " (" + fee.toFixed(2) + "s/B - " + (fee * FEE_SPACING).toFixed(2) + "s/B)");
                $("#cursortext_buckets_date").text(time);
                break;
            default:
                o.ylabel = "no ylabel set in optionBuilder"
        }
    }else{
        $("#cursortext_detailed").css("visibility","hidden");
        $("#cursortext_bucket").css("visibility","hidden");
        $("#cursortext_detailed_value").css("visibility","hidden");
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

// creates a Dygraph option from a chartType
// sample chartType = {name="size",timespan=24}
function optionBuilder(chartType) {
    var o = { // default option
        width: 1040,
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
        highlightSeriesOpts: {
            strokeWidth: 5,
            strokeBorderWidth: 2,
            strokeBorderColor: "#FDD",
            highlightCircleSize: 1
        },
        interactionModel:  {},
        highlightSeriesBackgroundAlpha: 0.5,
        highlightSeriesBackgroundColor: "#000",
        axes: {
            x: {axisLabelFormatter: function(d, gran, opts) {return Dygraph.dateAxisLabelFormatter(new Date((d/60).toFixed(0)*60000), gran, opts);}}
        }
    }

    // set options depending on name
    switch (chartType.name) {
        case "amount":
            o.ylabel = "tx count";
            o.highlightCallback = function(e, x, pts, row) {setCursorText(x,graph.getHighlightSeries(),graph.rolledSeries_[graph.rolledSeries_.length-graph.getHighlightSeries()-1][row][1],chartType);}
            o.axes.y = {axisLabelFormatter: function(y) {if(y>=1000){return + y/1000 + 'k';}else{return y;}},axisLabelWidth: 50,includeZero:true}
            break;
        case "size":
            o.ylabel =  "mempool size [MB]";
            o.highlightCallback = function(e, x, pts, row) {setCursorText(x,graph.getHighlightSeries(),graph.rolledSeries_[graph.rolledSeries_.length-graph.getHighlightSeries()-1][row][1],chartType);}
            o.axes.y = {axisLabelFormatter: function(y) {return + y/1000000;},axisLabelWidth: 50,includeZero:true}
            break;
        case "fee":
            o.ylabel = "amount in fees [BTC]";
            o.highlightCallback = function(e, x, pts, row) {setCursorText(x,graph.getHighlightSeries(),graph.rolledSeries_[graph.rolledSeries_.length-graph.getHighlightSeries()-1][row][1],chartType);}
            o.axes.y = {axisLabelFormatter: function(y) {return y.toFixed(2);},axisLabelWidth: 50,includeZero:true}
            break;
        case "bucket":
            o.ylabel = "tx count per bucket";
            o.highlightCallback = function(e, x, pts, row) {setCursorText(x, graph.getHighlightSeries(),graph.rolledSeries_[graph.rolledSeries_.length-graph.getHighlightSeries()-2][row][1],chartType);}
            o.axes.y = {axisLabelFormatter: function(y) {if(y>=1000){return + y/1000 + 'k';}else{return y;}},axisLabelWidth: 50,includeZero:true}
            break;
        default:
            o.ylabel = "no ylabel set in optionBuilder"
    }

    // set options depending on timespan
    switch (chartType.timespan) {
        case 168: // 7 days
            o.xlabel = "datetime";
            break;
        default: // no special values for any other timespans yet
            o.xlabel = "time"
    }
    return o;
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

    buildGraph('/dyn/amount4h.csv', optionBuilder({name:"amount",timespan:4}));

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

            case "nav-amount-4h":
                graph.updateOptions($.extend(optionBuilder({name:"amount",timespan:4}), {file: "/dyn/amount4h.csv"}),false);break;
            case "nav-amount-24h":
                graph.updateOptions($.extend(optionBuilder({name:"amount",timespan:24}), {file: "/dyn/amount24h.csv"}),false);break;
            case "nav-amount-7d":
                graph.updateOptions($.extend(optionBuilder({name:"amount",timespan:168}), {file: "/dyn/amount7d.csv"}),false);break;


            case "nav-size-4h":
                graph.updateOptions($.extend(optionBuilder({name:"size",timespan:4}), {file: "/dyn/size4h.csv"}),false);break;
            case "nav-size-24h":
                graph.updateOptions($.extend(optionBuilder({name:"size",timespan:24}), {file: "/dyn/size24h.csv"}),false);break;
            case "nav-size-7d":
                graph.updateOptions($.extend(optionBuilder({name:"size",timespan:168}), {file: "/dyn/size7d.csv"}),false);break;


            case "nav-value-4h":
                graph.updateOptions($.extend(optionBuilder({name:"fee",timespan:4}), {file: "/dyn/value4h.csv"}),false);break;
            case "nav-value-24h":
                graph.updateOptions($.extend(optionBuilder({name:"fee",timespan:24}), {file: "/dyn/value24h.csv"}),false);break;
            case "nav-value-7d":
                graph.updateOptions($.extend(optionBuilder({name:"fee",timespan:168}), {file: "/dyn/value7d.csv"}),false);break;


            case "nav-bucket-tab":
                graph.updateOptions($.extend(optionBuilder({name:"bucket",timespan:0}), {file: "/dyn/bucketlevel.csv"}),false);
            break;

        }
    });

    // fullscreen listener
    $(document).on('webkitfullscreenchange mozfullscreenchange fullscreenchange', function(e){
        if (isFullscreen) {
            fullscreenOff();
        }

        if(setNextStepFullscreenON){
            isFullscreen=true;
            setNextStepFullscreenON = false;
        }
    });
}

function fullscreenOn() {
    var i = document.getElementById("fullscreen-wrapper");
    if (i.requestFullscreen) {
        i.requestFullscreen();
    } else if (i.msRequestFullscreen) {
        i.msRequestFullscreen();
    } else if (i.mozRequestFullScreen) {
        i.mozRequestFullScreen();
    } else if (i.webkitRequestFullscreen) {
        i.webkitRequestFullscreen(Element.ALLOW_KEYBOARD_INPUT);
    }
    graph.resize($(window).width(),$(window).height());
    $('#fullscreen-wrapper').addClass("fullscreen");
    setNextStepFullscreenON = true; //HACK: somewhat dirty hack, because otherwise the listener would fire and direclty close the fullscreen again
}

function fullscreenOff() {
    var didRun = false;
    if (document.exitFullscreen) {
        document.exitFullscreen();
        didRun = true;
    } else if (document.msExitFullscreen) {
        document.msExitFullscreen();
        didRun = true;
    } else if (document.mozCancelFullScreen) {
        document.mozCancelFullScreen();
        didRun = true;
    } else if (document.webkitExitFullscreen) {
        document.webkitExitFullscreen();
        didRun = true;
    }
    if(didRun){
        $('#fullscreen-wrapper').removeClass( "fullscreen" );
        graph.resize(1040,650);
        isFullscreen=false;
    }
}

function toggleFullscreen() {
    var i = document.getElementById("fullscreen-wrapper");
    if (!i.fullscreenElement && !i.mozFullScreenElement && !i.webkitFullscreenElement && !i.msFullscreenElement && !isFullscreen) {
        fullscreenOn();
    } else {
        fullscreenOff();
    }
}
