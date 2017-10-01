var colorSet = ["#373854","#493267","#7bb3ff","#e86af0","#7bb3ff","#9e379f"];
var g; // graph
var tx_unconfirmed_timer;
var tx_unconfirmed_timer_last_block;
var txid;

function checkConfirmed() {
    $.ajax({url: "/api/confirmed/"+$("#input_txid").val(), success: function(result){
        if(result.statusCode=="200"){
            if(tx_unconfirmed_timer_last_block == null && result.data.tx_confirmed_in_block > 1){ // tx is confirmed
                new Audio('mp3/attention-seeker.mp3').play();
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
    g.setSelection(false, ''+index,false);
}

function setCursorText(date,key,value) {
    date = new Date(date*1000);
    var time = ("0" + date.getHours()).slice(-2) + ":" + ("0" + date.getMinutes()).slice(-2);
    if(value>0){
        $("#cursortext").css("visibility","visible");
        $("#cursortext_tally").text(value);
        $("#cursortext_SpB").text(key);
        $("#cursortext_date").text(time);
    }else{
        $("#cursortext").css("visibility","hidden");
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
                if(result.data.block!=null){
                    $("#card_info_tx_confirmed_block").text(result.data.block);
                    $("#card_info_tx_confirmed").css("display", "block");
                    $("#card_info_tx_confirmed_sound").css("display", "none");
                }else {
                    $("#card_info_tx_confirmed").css("display", "none");
                    $("#card_info_tx_confirmed_sound").css("display", "block");
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
                changeGraphColor(result.data.rate.toFixed(0));
            }else{
                alert(result.statusCode+" - "+result.message);
            }
        }});
    }else {
        alert("invalid input: '" + $("#input_txid").val() + "' - length: [" +  $("#input_txid").val().length + "]");
    }
}

var makeGraph = function(data, isStacked) {
    var chart = document.getElementById('chart');
    var div = document.createElement('div');
    div.className = "chartclass";
    div.style.display = 'inline-block';
    chart.appendChild(div);

    var labels = data[1];
    g = new Dygraph(
        div,
        data,
        {
            width: 1000,
            height: 650,
            colors: colorSet,
            fillAlpha: 1,
            strokeWidth: 4,
            strokeBorderWidth: 0,
            highlightCircleSize: 0,
            stackedGraph: isStacked,
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
                setCursorText(x,g.getHighlightSeries(),g.rolledSeries_[g.rolledSeries_.length-g.getHighlightSeries()-1][row][1]);
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
            },
            clickCallback: function(e, x, points){ // TODO FIX: this somehow dosn't work. The callback is never called.
                if (g.isSeriesLocked()) {
                    g.clearSelection();
                } else {
                    g.setSelection(g.getSelection(), g.getHighlightSeries(), true);
                }
            }
        });
    };


    window.onload = function () {

        // register onClickListener for the 'load tx id'-button
        $("#button_load_tx_info").click(function(){
            loadTXinfo();
        });

        // loads stacked graph with data from https://mempool.observer/data.csv
        makeGraph('data.csv', true);

        // if load_txid was defined by the ejs renderer load the tx
        // load_txid is the permalink txid
        if (typeof load_txid !== 'undefined') {
            $("#input_txid").val(load_txid)
            loadTXinfo();
        }

        // onChangeListener for the confirm sound checkbox
        $('#card_info_tx_confirmed_sound_checkbox').change(function(){
            if($(this).is(':checked')) {
                new Audio('mp3/attention-seeker.mp3').play(); // plays audio once
                tx_unconfirmed_timer = setInterval(checkConfirmed, 60000); // checks with the backend (over ajax) every 60 seconds if the transaction is confirmed
            } else {
                window.clearInterval(tx_unconfirmed_timer) // disables the timer
            }
        });
    }
