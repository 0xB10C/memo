/* 
  monitor-draw.js includes functionallity to draw the mempool.observer/monitor
  scatterplot. 
*/

var gDivChart = document.getElementById("chart")
var gCanvasCtx = null
var gData = null
var gQuadtree = null
var gRecentFeerateAPIData
var gBlockEntriesData

var gFeerateAPILines
var gLoadingSpinner

const yFormat = d3.format(".2s");
const timeFormat = d3.timeFormat("%H:%M");
const margin = {top: 20, right: 20, bottom: 50, left: 50};
var xScale = d3.scaleTime()
var yScale = d3.scaleLog()

// return the chart height and width minus margin
var width = function() {return gDivChart.clientWidth - margin.left - margin.right};
var height = function() {return gDivChart.clientHeight - margin.top - margin.bottom};

// return the x- and y-Axis values for a given data element 
var xValue = function (d) {return d.entryTime * 1000;}
var yValue = function (d) {return d.fee / d.size;}

// return the scaled x- and y-Axis values for a given data element
var xScaledValue = function (d) {return xScale(xValue(d))}
var yScaledValue = function (d) {return yScale(yValue(d))}

let color1 = "#b10c00"
let color2 = "lightgray"

let blockInclusionColorMap = {
  0: '#e630b9',
  1: '#2d8da5',
  2: '#638f28',
  3: '#f62e64',
  4: '#3480ed',
  5: '#299366',
  6: '#b87228',
  7: '#c149f1',
  8: '#2b9089',
  9: '#918328'
}

// coloring functions for the dots 
var cPlain = function (d) {return "#b10c00";}
var cFeerate = function (d) {return feerateColorScale(yValue(d));}
var cSegWit = function (d) {return d.spendsSegWit ? color1 : color2;}
var cLocktime = function (d) {return d.locktime > 0 ?  (d.locktime >= 500000000? "green" : color1) : color2;}
var cRBF = function (d) {return d.signalsRBF ? color1 : color2;}
var cVersion1 = function (d) {return d.version == 1 ? color1 : color2}
var cVersion2 = function (d) {return d.version == 2 ? color1 : color2}
var cOpReturn = function (d) {return d.opreturnLength ? color1 : color2;}
var cBIP69 = function (d) {return d.isBIP69 ? color1 : color2;}
var cMultisig = function (d) {return d.spendsMultisig ? color1 : color2;}
var cLNUnilateralClose = function (d) {return (d.locktime >= 500000000 && d.locktime < 600000000) ? color1 : "transparent";}
var cBlockInclusion = function(d) {
  if (d.block != null) {
    return blockInclusionColorMap[d.block % 10]
  }

  const shortTXID = d.txid.substring(0, 16)
  for (const block of gBlockEntriesData) {

    // skip if the block was found before the transaction was broadcast
    if (block.timestamp < d.entryTime) {
      continue
    }

    if (binarySearch(block.shortTXIDs, shortTXID) != -1) {
      d.block = block.height
      return blockInclusionColorMap[block.height % 10]
    }

  }
  return color2
 }
 var cUnconfirmed = function(d) {
  if (d.block != null) {
    return color2
  }

  const shortTXID = d.txid.substring(0, 16)
  for (const block of gBlockEntriesData) {

    // skip if the block was found before the transaction was broadcast
    if (block.timestamp < d.entryTime) {
      continue
    }

    if (binarySearch(block.shortTXIDs, shortTXID) != -1) {
      d.block = block.height
      return color2
    }

  }
  return color1
 }

// radius functions for the dots
var rUniform = function (d) {return 2}
var rSize = function (d) {return Math.log2(d.size * 0.01);}
var rFee = function (d) {return Math.log10(d.fee);}
var rInputCount = function (d) {return Math.log2(d.inputCount*1.2);}
var rOutputCount = function (d) {return Math.log2(d.outputCount*1.2);}
var rInputOutputCount = function (d) {return Math.log2(d.outputCount + d.inputCount);}
var rInputperOutputCount = function (d) {return d.inputCount >= d.outputCount ? Math.log2(d.inputCount / d.outputCount) : 0;} // we don't want to return a negative radius -> log2(x) is negativ if x < 1 -> x is < 1 if d.outputCount > d.inputCount (greater-than check is faster than calc log(a/b)) 
var rOutputperInputCount = function (d) {return d.outputCount >= d.inputCount ? Math.log2(d.outputCount / d.inputCount) : 0;} // similar as in rInputperOutputCount()  
var rOutputSum = function (d) {return Math.log1p(d.outputValue/10000000);}
var rDustOutput = function(d) {let value = Math.log2(1/d.outputValue*100000); return (value < 0) ? 0 : value;}
var rOPReturnDataLenght = function (d) {return d.opreturnLength > 0 ? d.opreturnLength/10 : 0}

// sets default color and radius function 
var currentColorFunction = cPlain;
var currentRadiusFunction = rSize;

async function redraw() {
  console.debug("started redraw()")

  d3.select(gDivChart).selectAll("*").remove();

  gLoadingSpinner = d3.select(gDivChart).append("div").attr("class", "loading-spinner")

  xScale.range([0, width()]);
  yScale.range([height(), 0]);

  var xAxis = d3.axisBottom(xScale).tickFormat(timeFormat);
  var yAxis = d3.axisLeft(yScale).tickFormat(yFormat);

  // The `svg` is appended before the canvas is, so that the svg is behind the canvas. 
  var svg = d3.select(gDivChart).append("svg")
    .attr("width", width() + margin.left + margin.right)
    .attr("height", height() + margin.top + margin.bottom)
    .style("left", "0px")
    .style("position", "relative");

  var chartArea = d3.select(gDivChart).append("div")
    .style("position", "absolute")
    .style("left", margin.left + "px")
    .style("top", "0px");

  var canvas = chartArea.append("canvas")
    .attr("width", width())
    .attr("height", height());
    
  gCanvasCtx = canvas.node().getContext("2d");

  var tooltip = d3.select(gDivChart).append("div")
    .attr("class", "tooltip")
    .style("opacity", 0);
  
  data = await loadEntryData()
  dataBlocks = await loadBlockEntriesData()

  var xMin = d3.min(data, function (d) {return xValue(d)})
  var xMax = d3.max(data, function (d) {return xValue(d)})

  // currently no used in favor of yMinQuantile and yMaxQuantile
  var yMin = d3.min(data, function (d) {return yValue(d)})
  var yMax = d3.max(data, function (d) {return yValue(d)})

  data.sort(function(a, b) {
    return yValue(a) - yValue(b);
  });

  var yMinQuantile = d3.quantile(data, 0.008, function (d) {return yValue(d)}) - 0.09;
  var yMaxQuantile = d3.quantile(data, 0.998, function (d) {return yValue(d)});

  xScale.domain([xMin, xMax]);
  yScale.domain([yMinQuantile, yMaxQuantile]);

  gQuadtree = d3.quadtree()
  .extent([[xScale(xMin), yScale(yMaxQuantile)], [xScale(xMax), yScale(yMinQuantile)]])
  .x(xScaledValue)
  .y(yScaledValue)


  svg.append("g")
      .attr("class", "axis axis--x")
      .attr("transform", "translate("+ margin.left + " ," + height() + ")")
      .call(xAxis);

  svg.append("g")
      .attr("class", "axis axis--y")
      .attr("transform", "translate("+ margin.left + ",0)")
      .call(yAxis);
      
  svg.append("text")
      .attr("transform", "rotate(-90)")
      .attr("y", 0)
      .attr("x", 0 - (height() / 2))
      .attr("dy", "1em")
      .style("text-anchor", "middle")
      .text("feerate (sat/vbyte)");  

  svg.append("text")
      .attr("transform",
            "translate(" + ((width() + margin.right + margin.left)/2) + " ," + 
                          (height() + margin.top + margin.bottom/1.2) + ")")
      .style("text-anchor", "middle")
      .text("arrival time");
  
  var highlight = chartArea.append("svg")
    .style("position", "absolute")
    .style("left", "0px")
    .style("top", "0px")
    .attr("width", width())
    .attr("height", height())
      .append("circle")
      .style("opacity", 0);

  gRecentBlockLines = svg.append("g")

  drawTransactions(data);
  fillQuadtree(data);

  gFeerateAPILines = svg.append("g");
  drawFeerateAPILines(xMin)

  canvas.on("mousemove",function(){  
    let mouse = d3.mouse(this);
    let closest = gQuadtree.find(mouse[0], mouse[1]);

    if (closest != undefined && isTransactionVisible(closest)) {
      highlight.attr("cx", xScale(xValue(closest)))
        .attr("cy", yScale(yValue(closest)))
        .attr("r", currentRadiusFunction(closest))
        .attr("stroke", currentColorFunction(closest))
        .attr("filter", "blur(0.5px) brightness(96%)")
        .transition().duration(1).style("opacity", 1);
  
      tooltip.html(formatTooltip(closest))
      let yPos = mouse[1] + 20
      let xPos = mouse[0] + ((mouse[0] <= gDivChart.clientWidth/2) ? 70 : -260);  
      tooltip.style("left", xPos + "px")
      tooltip.style("top", yPos + "px");

      tooltip.transition().duration(1).style("opacity", 1);

    }else{
      tooltip.transition().style("opacity", 0);
      highlight.transition().style("opacity", 0);
    }
  });

  canvas.on("mouseout",function(){
    highlight.transition().duration(100).style("opacity", 0);
    tooltip.transition().duration(100).style("opacity", 0);
  });

  canvas.on("click", function (d) {
    let mouse = d3.mouse(this);
    let closest = gQuadtree.find(mouse[0], mouse[1], 10);
    if(closest != undefined && isTransactionVisible(closest)){
      window.open("https://blockstream.info/tx/" + closest.txid, "_blank");
    }
  });

  d3.select("#select-radius").on('change', function(){
    setTimeout(function(){
      selected = d3.select("#select-radius").property('value');
      if (selected == "1") {
        deleteQueryStringParameter("R")
      }else{
        setQueryStringParameter("R",  selected);
      }
      descriptionRadius = d3.select("#span-radius-description")
      switch (selected) {
        case "0": // uniform
          currentRadiusFunction = rUniform;
          descriptionRadius.html("Every transaction has the same radius.")
          break;
        case "1": // size
          currentRadiusFunction = rSize;
          descriptionRadius.html("The radius is calculated based on transaction vsize.")
          break;
        case "2": //outputs
          currentRadiusFunction = rOutputCount;
          descriptionRadius.html("The radius is calculated based on the output count.")
          break;
        case "3": // inputs
          currentRadiusFunction = rInputCount;
          descriptionRadius.html("The radius is calculated based on the input count.")
          break;
        case "4": // inputsoutputs
          currentRadiusFunction = rInputOutputCount;
          descriptionRadius.html("The radius is calculated based on the <code>input + output</code> count.")
          break;
        case "7": // outputsum
          currentRadiusFunction = rOutputSum;
          descriptionRadius.html("The radius is calculated based on the sum of the output values. A bigger radius means a higher value transacted.")
          break;
        case "8": // dustouputs
          currentRadiusFunction = rDustOutput;
          descriptionRadius.html("The radius is calculated based on the inverse of the output values. A bigger radius means a smaller value transacted.")
          break;
        case "6": //outputperinput
          currentRadiusFunction = rOutputperInputCount;
          descriptionRadius.html("The radius is calculated based on the output per input ratio. A bigger radius means a growing UTXO set.")
          break;
        case "5": // inputperoutput
          currentRadiusFunction = rInputperOutputCount;
          descriptionRadius.html("The radius is calculated based on the input per output ratio. A bigger radius means a shrinking UTXO set.")
          break;
        case "9": // opreturnlength
          currentRadiusFunction = rOPReturnDataLenght;
          descriptionRadius.html("The radius is calculated based length of the data the OP_RETURN pushes.")
      }
      drawTransactions(data)
    }, 0)
  })

  d3.select("#select-estimator").on('change', function(){
    setTimeout(function(){
      drawFeerateAPILines(xMin)
    }, 0)
  })

  d3.select("#filter-select-multisig").on('change', function(){
    setTimeout(function(){
      filter = gFilters.multisigs
      selected = d3.select("#filter-select-multisig").property('value');
      filter.input.value = selected
      if (selected == ""){
        deleteQueryStringParameter(filter.input.queryStringCode)
      } else{
        setQueryStringParameter(filter.input.queryStringCode,  selected);
      }
      drawTransactions(data)
    }, 0)
  })

  fillMultisigSelect(data)

  // handles tri-state-switch changes
  d3.selectAll("input[type='range']").on("change", function (e, index, list){
    setTimeout(function() {
      input = list[index]
      
      filterKey = input.getAttribute("data-filter-id")
      var filter = gFilters[filterKey]
      
      switch (filter.type) {
        case "tri-state-switch":
          filter.state = input.value
          break;
        case "withinput":
          filter.state = input.value
          break;
        case "multisigInput":
          filter.state = input.value
          break;
      }
      
      if (filter.state == filterStates.inactive || filter.state == filterStates.greaterEqual) { 
        deleteQueryStringParameter(filter.queryStringCode) // delete query string on default
      }else{
        setQueryStringParameter(filter.queryStringCode,  filter.state);
      }

      drawTransactions(data);
    }, 0);
  });

    // handles filter input changes
    d3.selectAll("input[name='filter-input-freetext']").on("change", function (e, index, list){
      setTimeout(function() {
        input = list[index]
        
        filterKey = input.getAttribute("data-filter-id")
        var filter = gFilters[filterKey]
        filter.input.value = input.value
        
        if (filter.input.value == 0 || filter.input.value == "") { 
          deleteQueryStringParameter(filter.input.queryStringCode) // delete query string on default
        }else{
          setQueryStringParameter(filter.input.queryStringCode,  filter.input.value);
        }
  
        drawTransactions(data);
      }, 0);
    });

  d3.select("#select-highlight").on('change', function(){
    setTimeout(function(){
      selected = d3.select("#select-highlight").property('value');
      if (selected == "0") {
        deleteQueryStringParameter("H")
      }else{
        setQueryStringParameter("H",  selected);
      }
      descriptionFilter = d3.select("#span-highlight-description")
      switch (selected) {
        case "0": // all
          currentColorFunction = cPlain;
          descriptionFilter.html("No transactions are highlighted.");
          break;
        case "1": // segwit
          currentColorFunction = cSegWit;
          descriptionFilter.html("Transactions spending SegWit are highlighted");
          break;
        case "3": // locktime
          currentColorFunction = cLocktime;
          descriptionFilter.html("Transactions having a locktime greater than zero are highlighted.");
          break;
        case "4": // rbf
          currentColorFunction = cRBF;
          descriptionFilter.html("Transactions signaling explicit Replace-By-Fee are highlighted.");
          break;
        case "5": // opreturn
          currentColorFunction = cOpReturn;
          descriptionFilter.html("Transactions having a OP_RETURN output are highlighted.");
          break;
        case "6": // bip69
          currentColorFunction = cBIP69;
          descriptionFilter.html("Transactions compliant to <a href=\"https://github.com/bitcoin/bips/blob/master/bip-0069.mediawiki\">BIP-69</a> input and output sorting are highlighted.");
          break;
        case "2": // multisig
          currentColorFunction = cMultisig;
          descriptionFilter.html("Transactions spending Multisig are highlighted.");
          break;
        case "7": // version1
          currentColorFunction = cVersion1;
          descriptionFilter.html("Version 1 transactions are highlighted.");
          break;
        case "8": // version2
          currentColorFunction = cVersion2;
          descriptionFilter.html("Version 2 transactions are highlighted.");
          break;
        case "9": // block inclusion
          currentColorFunction = cBlockInclusion;
          descriptionFilter.html("Transactions are highlighted according to the block they were include in.");
          break;
        case "10": // unconfirmed
          currentColorFunction = cUnconfirmed;
          descriptionFilter.html("Unconfirmed transactions are highlighted.");
          break;
      }
      drawTransactions(data)
    }, 10);
  });
  
  gLoadingSpinner.remove()
  console.debug("finished redraw()")
}

// The multisig select with the id `filter-select-multisig` gets filled with the currently possible
// n-of-m combinations present in the current transaction set. 
async function fillMultisigSelect(data){
  if (gFilters.multisigs.selectFilled == false) {
  var multisigs = {}
  data.forEach(function(d){
    if (d.spendsMultisig){
      for (var key in d.multisigsSpend) {
        if (multisigs[key] == null) {
          multisigs[key] = 0
        }
        multisigs[key]++
      }
    }
  });

  var selectElement = document.getElementById("filter-select-multisig")
    Object.keys(multisigs).sort().forEach(function(key){
      let option = document.createElement("option")
      option.value = key
      option.text = key + " (" + multisigs[key] + "x)"
      if (gFilters.multisigs.input.value == key){
        option.selected = true
      }
      selectElement.appendChild(option)
    });
  }

  gFilters.multisigs.selectFilled = true
}

async function fillQuadtree(data){
  gQuadtree.removeAll(gQuadtree.data())
  data.forEach(function(d){
    if (isTransactionVisible(d)){
      gQuadtree.add(d);
    }
  });
}

async function drawTransactions(data){
  gCanvasCtx.clearRect(0, 0, width(), height());
  const maxFeerate = yScale.domain()[1]
  var count = 0
  var countDrawn = 0
  var countOutOfBounds = 0
  data.forEach(function(d){
    count++
    if (isTransactionVisible(d)) {
      if (yValue(d) <= maxFeerate){
        countDrawn++
        let color = currentColorFunction(d);
        let radius = currentRadiusFunction(d);
        gCanvasCtx.beginPath();
        gCanvasCtx.fillStyle = color;
        gCanvasCtx.globalAlpha=0.8;
        gCanvasCtx.arc(xScaledValue(d), yScaledValue(d), radius, 0, 2 * Math.PI);
        gCanvasCtx.fill();
      } else {
        countOutOfBounds++
      }
    }
  });
  document.getElementById("span-transaction-drawn").innerText = countDrawn
  document.getElementById("span-transaction-drawn-p").innerText = ((countDrawn / count)*100).toFixed(2);
  document.getElementById("span-transaction-outofbounds").innerText = countOutOfBounds
  drawBlocks(dataBlocks);
  fillQuadtree(data)
}

async function drawBlocks(data) {
  if (currentColorFunction == cBlockInclusion) {
    data.forEach(function(d){
      let x = xScale(d.timestamp*1000)
      gCanvasCtx.beginPath();
      gCanvasCtx.moveTo(x, height());
      gCanvasCtx.lineTo(x, 0);
      gCanvasCtx.globalAlpha = 0.2;
      gCanvasCtx.stroke();
      gCanvasCtx.globalAlpha = 0.7;
      gCanvasCtx.fillStyle = "black";
      gCanvasCtx.save();
      gCanvasCtx.translate(x-4, height()-3);
      gCanvasCtx.rotate(-90 * (Math.PI / 180));
      gCanvasCtx.textAlign = "left";
      gCanvasCtx.fillText(`Block ${d.height}`, 0, 0);
      gCanvasCtx.restore();    
    });
  }
}

async function drawFeerateAPILines(xMin){
  gRecentFeerateAPIData = await loadRecentFeerateAPIData()
  let selectedId = d3.select("#select-estimator").property('value');
  if (selectedId == "0") {
    deleteQueryStringParameter("E")
  }else{
    setQueryStringParameter("E",  selectedId);
  }

  const descriptions = {
    0: "No feerate estimator overlayed.", // none
    1: "Feerates for a confirmation in half an hour, one hour and two hours are shown.", // bitcoinerlive
    2: "Feerates for a confirmation in two, four and six blocks shown.", // bitgocom
    3: "Feerates for a confirmation in two, three and six blocks shown.", // bitpaycom
    4: "Feerates for a priority and regular confirmation shown.", // blockchaininfo
    5: "The recommended feerate is shown. Feerate estimation data retrieved from <a href=\"https://blockchair.com/\">blockchair.com</a>.", // blockchaircom
    6: "High, medium and low feerate recomendations are shown.", // blockcyphercom
    7: "Feerates for a confirmation in two, three and six blocks shown.", // blockstreaminfo
    8: "Feerate recomendation for the next block is shown.", // btccom
    9: "The fastest, a half an hour and an hour feerate are shown.", // earncom
    10: "Feerates for a confirmation in one, three and six blocks shown.", // ledgercom
    11: "Feerates for two, four and ten blocks are shown.", // myceliumio
    12: "Feerates for two, four and six blocks are shown.", // trezorio
    13: "Feerates for two, four and six blocks are shown.", // wasabiwalletioEcon
    14: "The fastest, a half an hour and an hour feerate are shown.", // mempoolspace
  }

  const idToName = {
    0: "none",
    1: "bitcoinerlive", 
    2: "bitgocom", 
    3: "bitpaycom", 
    4: "blockchaininfo", 
    5: "blockchaircom", 
    6: "blockcyphercom", 
    7: "blockstreaminfo", 
    8: "btccom", 
    9: "earncom", 
    10: "ledgercom", 
    11: "myceliumio", 
    12: "trezorio", 
    13: "wasabiwalletioEcon", 
    14: "mempoolspace",
  }

  selected = idToName[selectedId]
  d3.select("#span-estimator-description").html(descriptions[selectedId]);
  gFeerateAPILines.selectAll("path").remove();

  let drawHighLine = selected != "none" // don't draw any lines when the selected is none
  let drawMedLine = drawHighLine && selected != "btccom" && selected != "blockchaircom" // don't draw a medium line for btc.com and blockchair.com
  let drawLowLine = drawMedLine && selected != "blockchaininfo" // don't draw a low line for blockchain.info API

  if ( drawHighLine ){
    let highLine = d3.line()
      .defined(function(d) {return isFeerateDefined(d[selected].high) && new Date(d.timestamp*1000) >= xMin})
      .x(function(d) { return xScale(new Date(d.timestamp*1000)); })
      .y(function(d) { return yScale(d[selected].high); })
      .curve(d3.curveStep)

    gFeerateAPILines.append("path")
      .datum(gRecentFeerateAPIData)
      .attr("class", "feerate-line-high")
      .attr("d", highLine)
      .attr("transform", "translate(" + margin.left + ", 0)");
  } 
  
  if ( drawMedLine ){
    let medLine = d3.line()
      .defined(function(d) {return isFeerateDefined(d[selected].med) && new Date(d.timestamp*1000) >= xMin})
      .x(function(d) { return xScale(new Date(d.timestamp*1000)); })
      .y(function(d) { return yScale(d[selected].med); }) 
      .curve(d3.curveStep)

    gFeerateAPILines.append("path")
      .datum(gRecentFeerateAPIData)
      .attr("class", "feerate-line-med")
      .attr("d",medLine) 
      .attr("transform", "translate(" + margin.left + ", 0)");
  } 
  
  if ( drawLowLine ){
    let lowLine = d3.line()
      .defined(function(d) {return isFeerateDefined(d[selected].low) && new Date(d.timestamp*1000) >= xMin})
      .x(function(d) { return xScale(new Date(d.timestamp*1000)); })
      .y(function(d) { return yScale(d[selected].low); })
      .curve(d3.curveStep)

    gFeerateAPILines.append("path")
      .datum(gRecentFeerateAPIData)
      .attr("class", "feerate-line-low")
      .attr("d",lowLine)
      .attr("transform", "translate(" + margin.left + ", 0)");
  } 
}

function isFeerateDefined(feerate){
  return !isNaN(feerate) && feerate > 1
}


function formatTooltipTableRow(name,value) {
  return `<tr><td>${name}</td><td style="max-width: 20ch; word-wrap:break-all;">${value}</td></tr>`
}

function formatTooltip(d){
  let shortTxID = d.txid.substring(0, 20) + "<small>...</small>"
  return `<div><table class="table table-sm table-bordered text-center"><tbody>
    <tr><th colspan=2>${shortTxID}</th></tr>
      ${formatTooltipTableRow("Entry time", new Date(xValue(d)).toLocaleTimeString())}
      ${formatTooltipTableRow("Feerate", yValue(d).toFixed(2) + " sat/vbyte")}
      ${formatTooltipTableRow("Size", d.size + " vbyte")}
      ${formatTooltipTableRow("Fee", d.fee + " sat")}
      ${formatTooltipTableRow("Version", d.version)}
      ${formatTooltipTableRow("Inputs", formatTooltipDicts(d.spends))}
      ${formatTooltipTableRow("Outputs", formatTooltipDicts(d.paysTo))}
      ${formatTooltipTableRow("Output amount", d.outputValue / 100000000 + " BTC")}
      ${formatTooltipTableRow("SegWit spending", d.spendsSegWit)}
      ${formatTooltipTableRow("Explicit RBF", d.signalsRBF)}
      ${formatTooltipTableRow("BIP69 compliant", d.isBIP69)}
      ${formatTooltipTableRow("Multisig spending", d.spendsMultisig)}
      ${d.spendsMultisig ? formatTooltipTableRow("Multisig inputs", formatTooltipDicts(d.multisigsSpend)) : ""}
      ${formatTooltipTableRow("OP_RETURN ", d.opreturnLength > 0 ? sanitizeHTML(d.opreturnData) : false)}
      ${d.opreturnLength > 0 ? formatTooltipTableRow("OP_RETURN length", d.opreturnLength + " bytes"): ""}
      ${formatTooltipTableRow("Locktime", d.locktime)}
    </tbody></table></div>
  `
}

function formatTooltipDicts(dict) {
  var formattedString = ""
  for (var key in dict) {
    formattedString = formattedString + `<span class="badge badge-secondary">${dict[key] + "x&nbsp;" + key}</span><br>` // produces e.g. `2x <key>` 
  }
  return formattedString
}

function isTransactionVisible(tx){
  if(currentRadiusFunction(tx) <= 0){
    return false;
  }

  for (var key in gFilters) {
    var filter = gFilters[key]
    if (filter.type != "separator"){
      if (!filter.isVisibleFunc(filter, tx)) {
        return false
      }
    }
  }

  return true
}

function binarySearch(items, value){
  var firstIndex  = 0,
    lastIndex   = items.length - 1,
    middleIndex = Math.floor((lastIndex + firstIndex)/2);

  while(items[middleIndex] != value && firstIndex < lastIndex){
    if (value < items[middleIndex]){
      lastIndex = middleIndex - 1;
    } else if (value > items[middleIndex]){
      firstIndex = middleIndex + 1;
    }
    middleIndex = Math.floor((lastIndex + firstIndex)/2);
  }
  return (items[middleIndex] != value) ? -1 : middleIndex;
}