// Constants
const NEXT_BLOCK_LABELS = ["1 vMB", "2 vMB", "3 vMB"]

var isTabActive = true;
var updateInterval = 30000

// use the https://mempool.observer endpoint only when on https://mempool.observer/ 
// otherwise use https://dev.mempool.observer (e.g. when on localhost)
var apiHost = window.location.hostname == "mempool.observer" ? "https://mempool.observer" : "https://dev.mempool.observer"

//var apiHost = "http://localhost:23485"

// State
var state = {
  currentMempool: {
    chart: null,
    elementId: "card-current-mempool",
    isScrolledIntoView: true,
    data: {
      processedMempool: null,
      currentTx: null,
      timeLastUpdated: null,
    },
  },
  historicalMempool: {
    chart: null,
    elementId: "card-historical-mempool",
    isScrolledIntoView: false,
    data: {
      processedMempool: null,
      timeLastUpdated: null,
      timeframe: 2,
      bySelector: "byCount",
    },
  },
  pastBlocks: {
    chart: null,
    elementId: "card-past-blocks",
    isScrolledIntoView: false,
    data: {
      processedBlocks: null,
      timer: null,
      timeLastUpdated: null,
    },
  },
  transactionStats: {
    chart: null,
    elementId: "card-transaction-stats",
    isScrolledIntoView: false,
    data: {
      type: "percentage",
      processedStats: null,
      timeLastUpdated: null,
    },
  },
}

var cards = [state.currentMempool, state.historicalMempool,  state.pastBlocks, state.transactionStats]

window.onload = function () {

  console.log()


  // Add event listeners to the search bar
  document.getElementById('button-lookup-txid').addEventListener('click', currentMempoolCard.handleTxSearch)
  document.getElementById('input-lookup-txid').addEventListener("keyup", function (event) {
    if (event.keyCode === 13) { // Handle enter
      event.preventDefault()
      document.getElementById('button-lookup-txid').click()
    }
  })
  document.getElementById('random-tx').addEventListener('click', currentMempoolCard.loadRandomTransactionFromApi)

  window.onfocus = function () {
    drawChart(state.currentMempool.elementId)
  };

  document.addEventListener("visibilitychange", visibilityChangeHandler, false);

  // add scroll listener 
  $(document).scroll(scrollEventHandler);


  if ((localStorage.getItem('mode') || "light") == 'light') {
    $('body').attr('data-theme', 'light');
  } else {
    $('body').attr('data-theme', 'dark')
  }


  // init data reload loop
  reloadData()
}

function visibilityChangeHandler() {
  if (!document.hidden) {
    // the page is now visible 
    // call scrollEventHandler() which already implements much of the functionality we need
    scrollEventHandler()
  }
}

function scrollEventHandler() {

  // handle scroll offset over 60px from top to animate the icon
  if ($(".navbar").offset().top > 60) {
    $(".navbar").addClass("scrolled");
  } else {
    $(".navbar").removeClass("scrolled");
  }

  cards.forEach(card => {
    let el = document.getElementById(card.elementId)

    let scrolledIntoView = isScrolledIntoView(el)

    if (card.isScrolledIntoView == false && scrolledIntoView == true) {
      // redraw chart if card just scrolled into view
      card.isScrolledIntoView = scrolledIntoView
      drawChart(card.elementId)
    } else if (card.isScrolledIntoView == true && scrolledIntoView == false) {
      // destroy chart if card just scrolled out of view
      card.isScrolledIntoView = scrolledIntoView
      if (card.chart != null) {
        card.chart = card.chart.destroy();
        //console.error("Destroying", card.elementId)
      }
    }

    // debug card is scrolled into View
    /*
    console.log(card.elementId, card.isScrolledIntoView)
    if(card.isScrolledIntoView){
      el.style.setProperty("background-color", "lime", "important");
    } else {
      el.style.setProperty("background-color", "orange", "important");
    }
    */

  });

}

function isScrolledIntoView(el) {
  const offset = 100
  let elemTop = el.getBoundingClientRect().top;
  let elemBottom = el.getBoundingClientRect().bottom;
  let elemHeight = elemBottom - elemTop

  let isVisible = (elemTop + elemHeight + offset >= 0) && (elemBottom - elemHeight - offset <= window.innerHeight);
  return isVisible;
}

// this object inhibits all functions used by or for the current mempool card
const currentMempoolCard = {
  processDataForChart: function (response) {
    state.currentMempool.data.timeLastUpdated = response.timestamp

    const mempoolSize = +(response.mempoolSize / 1000000).toFixed(2)
    const patternAreas = {
      '0to10': [],
      '11to100': [],
      '101to1k': [],
      'from1001': []
    }

    const lines = []
    const blocks = []
    const rowData = [
      [],
      []
    ]

    for (var feerate in response.feerateMap) {
      rowData[0].push(feerate.toString())
      rowData[1].push(response.feerateMap[feerate])

      // TODO: (0xb10c) Why do we do this and what does it?
      log1pOfCount = response.feerateMap[feerate]
      if (feerate <= 10) {
        patternAreas['0to10'].push(log1pOfCount)
      } else if (feerate <= 100) {
        patternAreas['11to100'].push(log1pOfCount)
      } else if (feerate <= 1000) {
        patternAreas['101to1k'].push(log1pOfCount)
      } else {
        patternAreas['from1001'].push(log1pOfCount)
      }
    }


    for (var position in response.megabyteMarkers) {
      if (position < 3) {
        blocks.push(response.megabyteMarkers[position])
        // add lines to show estimated next blocks on the mempool graph
        lines.push({
          value: response.megabyteMarkers[position],
          class: 'block-grid',
        })
      }
    }

    let colorPattern = currentMempoolCard.generateColorPattern(patternAreas)

    // Sum all txs to get the total number of tx in the mempool
    const sum = Object.values(rowData[1]).reduce((a, b) => a + b, 0)
    return {
      "yTickCache": {},
      "feerateMap": response.feerateMap,
      "mempoolSize": mempoolSize,
      "blocks": blocks,
      "colorPattern": colorPattern,
      "lines": lines,
      "rowData": rowData,
      "sum": sum
    }
  },
  generateColorPattern: function (patternAreas) {
    const patternColors = [
      ["#57e0fb", "#55ff00"],
      ["#55ff00", "#febf00"],
      ["#febf00", "#ff339c"],
      ["#ff339c", "#7705ec"]
    ]

    var colorPattern = []
    var c_counter = 0
    for (area in patternAreas) {
      logLimits = chroma.limits(patternAreas[area], 'e', patternAreas[area].length);
      pattern = chroma.scale(patternColors[c_counter])
        .mode('lch').classes(logLimits)
        .colors(patternAreas[area].length);

      colorPattern = colorPattern.concat(pattern)
      c_counter++
    }

    return colorPattern
  },
  calcYTick: function (yValue) {
    processedMempool = state.currentMempool.data.processedMempool

    // look if we have already processed this yValue (c3js does this multiple times somehow)
    if (processedMempool.yTickCache[yValue] != null || yValue == 0) {
      return yValue
    }

    var txCounter = 0
    var feerateAtY = 0
    for (var feerate in processedMempool.feerateMap) {
      txCounter += processedMempool.feerateMap[feerate]
      if (txCounter > yValue) {
        feerateAtY = feerate
        break;
      }
    }

    // add a new, invisible line with the feerate as description
    processedMempool.lines.push({
      value: yValue,
      class: 'hidden-feerate-grid',
      text: feerateAtY + ' sat/vbyte',
      position: 'middle',
    })

    processedMempool.yTickCache[yValue] = feerateAtY;

    return yValue
  },
  draw: async function () {
    let processed = state.currentMempool.data.processedMempool

    chartSetting = {
      bindto: '#current-mempool-chart',
      data: {
        rows: processed.rowData,
        type: 'bar',
        groups: [processed.rowData[0]],
        order: null,
      },
      point: {
        show: false
      },
      legend: {
        show: false
      },
      tooltip: {
        grouped: false,
        format: {
          name: function (name) {
            return name + ' sat/vbyte'
          },
          value: function (value) {
            return value + ' transactions'
          },
        }
      },
      size: {
        height: 750 // css #current-mempool-chart min-height: 750px needs to be changed too
      },
      padding: {
        top: 20
      },
      color: {
        pattern: processed.colorPattern
      },
      grid: {
        y: {
          lines: processed.lines
        }
      },
      axis: {
        y: {
          padding: {
            top: 0
          },
          show: true,
          label: {
            text: 'unconfirmed tx count'
          },
          tick: {
            format: function (y) {
              return currentMempoolCard.calcYTick(y)
            },
          },
        },
        y2: {
          outer: false,
          padding: {
            top: 0,
            bottom: 0
          },
          default: [0, processed.sum],
          label: {
            text: ''
          },
          show: true,
          tick: {
            format: function (d) {
              return NEXT_BLOCK_LABELS[processed.blocks.indexOf(d)]
            },
            values: processed.blocks
          }
        }
      }
    }

    // properly destroy chart and generate new chart
    if (state.currentMempool.chart) {
      state.currentMempool.chart = state.currentMempool.chart.destroy();
    }
    state.currentMempool.chart = c3.generate(chartSetting)

    // draw the tx the chart if it's unconfirmed
    let tx = state.currentMempool.data.currentTx;
    if (tx != null) {
      if (tx.blockHeight) {
        $('#tx-eta-data').html(`Confirmed (block ${tx.blockHeight}, ${minutes_since_confirmation} minutes ago)`) // FIXME: tx-eta-data is not used anymore
      } else {
        const feeRate = Math.floor(tx.fee / tx.vSize)
        currentMempoolCard.drawUserTxByFeeRate(feeRate)
      }
    }

  },
  updateCard: function (processed) {
    const spanTxCount = document.getElementById('current-mempool-count')
    const spanMempoolSize = document.getElementById('current-mempool-size')

    const txCountInMempool = processed.sum
    const txSizeInMempool = processed.mempoolSize

    spanTxCount.innerHTML = txCountInMempool
    spanMempoolSize.innerHTML = txSizeInMempool

    timeSinceLastUpdate = 0
    currentMempoolCard.updateCardLastUpdated()
  },
  updateCardLastUpdated: function () {
    // calc minutes from milliseconds
    const minutes = Math.floor((Date.now() - (state.currentMempool.data.timeLastUpdated * 1000)) / 1000 / 60)
    document.getElementById('current-mempool-last-update').innerHTML = (minutes)
  },
  handleTxSearch: async function () {
    $('#input-lookup-txid').removeClass("is-invalid") // remove invalid tx message 
    $('#current-mempool-tx-data').hide() // hide the current tx data

    inputTxid = document.getElementById('input-lookup-txid').value

    // Check if the input has the format of a txid
    if (/^[a-fA-F0-9]{64}$/.test(inputTxid) == false) {

      $('#invalid-feedback').html('Invalid Bitcoin transaction id.') // Set invalidt tx message
      $('#input-lookup-txid').addClass("is-invalid") // Show invalid tx message

    } else {

      try {
        console.log(await getTxFromApi(inputTxid))
        state.currentMempool.data.currentTx = await getTxFromApi(inputTxid)
        let tx = state.currentMempool.data.currentTx

        if (tx.confirmations) {
          let minutes_since_confirmation = Math.floor((Date.now() - tx.blockTime * 1000) / 1000 / 60)
          let error = `The transaction is already confirmed and therefore not in the mempool. (block ${tx.blockHeight}, ${minutes_since_confirmation} minutes ago)`
          console.error(error)
          $('#invalid-feedback').html(error)
          $('#input-lookup-txid').addClass("is-invalid")
        } else { // tx is unconfirmed
          const vSize = tx.vSize
          const feeRate = Math.floor(tx.fee / vSize)
          currentMempoolCard.drawUserTxByFeeRate(feeRate)
          currentMempoolCard.displayTransactionData(tx.fee, vSize, feeRate)
        }

      } catch (error) {
        console.error(error)
        $('#invalid-feedback').html(error)
        $('#input-lookup-txid').addClass("is-invalid")
      }

    }
  },
  drawUserTxByFeeRate: function (feeRate) {
    const position = currentMempoolCard.getUserTxPositionInChartByFeeRate(feeRate)

    // Remove existing line(s)
    state.currentMempool.chart.ygrids.remove({
      class: 'red-line'
    });

    // Draw a new line, but wait 20ms for c3.js to not to remove it directly
    setTimeout(function () {
      state.currentMempool.chart.ygrids.add([{
        value: position,
        text: 'Your transaction',
        class: 'red-line',
        position: 'middle'
      }]);
    }, 500)

  },
  getUserTxPositionInChartByFeeRate: function (feeRate) {
    let position = 0
    let processedMempool = state.currentMempool.data.processedMempool

    for (index in processedMempool.rowData[0]) {
      if (processedMempool.rowData[0][index] < feeRate) {
        position += processedMempool.rowData[1][index]
      } else if (processedMempool.rowData[0][index] == feeRate) {
        // get ruffly the middle position of the bar
        position += Math.round(processedMempool.rowData[1][index] / 2)
      } else {
        break;
      }
    }

    return position
  },
  displayTransactionData: function (fee, vsize, feeRate) {
    $('#current-mempool-tx-data').css("display", "flex");
    $('#current-mempool-tx-data-fee').html(fee)
    $('#current-mempool-tx-data-size').html(vsize)
    $('#current-mempool-tx-data-feerate').html(feeRate)
  },
  loadRandomTransactionFromApi: async function () {
    try {
      const recentTxs = await axios.get('https://blockstream.info/api/mempool/recent')
      const feePayingTxs = recentTxs.data.filter(tx => tx.fee > 0).map(tx => tx.txid)
      var randomTx = feePayingTxs[Math.floor(Math.random() * feePayingTxs.length)];
      document.getElementById('input-lookup-txid').value = randomTx;
      currentMempoolCard.handleTxSearch()
    } catch (e) {
      console.error(e)
    }
  }
}

const historicalMempoolCard = {
  updateCardLastUpdated: function () {
    // calc seconds from milliseconds
    const minutes = Math.floor((Date.now() - (state.historicalMempool.data.timeLastUpdated)) / 1000 / 60)
    document.getElementById('historical-mempool-last-update').innerHTML = (minutes)
  },
  switchTimeframe: function (newTimeframe) {
    if (newTimeframe != state.historicalMempool.data.timeframe) {
      state.historicalMempool.data.timeframe = newTimeframe
      reloadHistoricalMempool()
    }
  },
  switchBySelector: function (newBySelector) {
    if (newBySelector != state.historicalMempool.data.bySelector) {
      state.historicalMempool.data.bySelector = newBySelector
      reloadHistoricalMempool()
    }
  },
  processDataForChart: function (response) {
    state.historicalMempool.data.timeLastUpdated = new Date();

    let rows = [
      ["x", 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 15, 18, 22, 27, 33, 41, 50, 62, 76, 93, 114, 140, 172, 212, 261, 321, 395, 486, 598, 736, 905, 1113, 1369, 1684, 2071, 2547, 3133, 3854, "3854+"]
    ]

    for (blockIndex in response) {
      let timestamp = response[blockIndex].timestamp * 1000
      let dataInBuckets = response[blockIndex].dataInBuckets

      rows.push([new Date(timestamp)].concat(dataInBuckets))
    }

    return {
      rows: rows,
    }
  },
  draw: async function () {

    let processed = state.historicalMempool.data.processedMempool

    if (processed == null) {
      return
    }

    chartSetting = {
      bindto: '#historical-mempool-chart',
      data: {
        x: "x",
        xFormat: '%Y-%m-%dT%H:%M:%S.%LZ',
        rows: processed.rows,
        type: "area-spline",
        groups: [processed.rows[0]],
        order: null,
      },
      axis: {
        x: {
          type: 'timeseries',
          tick: {
            format: state.historicalMempool.data.timeframe > 3 ? '%H:%M %d.%m.' : '%H:%M'
          }
        },
        y: {
          tick: {
            format: d3.format(".2s")
          }
        }
      },
      size: {
        height: 450 // css #historical-mempool-chart min-height: 450px needs to be changed too
      },
      padding: {
        top: 20,
        bottom: 20,
      },
      color: {
        pattern: ["#57e0fb", "#00e7fb", "#00edf5", "#00f2e9", "#00f7d6", "#00fbbe", "#00fda1", "#00ff7f", "#00ff55", "#55ff00", "#55ff00", "#79f900", "#92f300", "#a7ed00", "#b9e700", "#c8e000", "#d6da00", "#e2d300", "#edcc00", "#f6c600", "#febf00", "#febf00", "#ffb011", "#ffa022", "#ff9032", "#ff8041", "#ff6f50", "#ff5f5f", "#ff506e", "#ff427e", "#ff388d", "#ff339c", "#ff339c", "#f719a5", "#ec00af", "#df00ba", "#ce00c6", "#b800d2", "#9d00df", "#7705ec"]
      },
      legend: {
        show: false
      },
      point: {
        show: false,
      },
      tooltip: {
        grouped: false,
        contents: function (d, defaultTitleFormat, defaultValueFormat, color) {
          date = d[0].x
          let day = "0" + date.getDate();
          let month = "0" + (date.getMonth() + 1)
          let year = date.getFullYear()
          let hours = "0" + date.getHours();
          let minutes = "0" + date.getMinutes();
          let formattedTime = hours.substr(-2) + ':' + minutes.substr(-2) + " " + day.substr(-2) + "/" + month.substr(-2) + "/" + year

          var value = ""
          switch (state.historicalMempool.data.bySelector) {
            case "byCount":
              value = d[0].value + " transactions"
              break;
            case "byFee":
              value = d[0].value.toFixed(4) + " BTC"
              break;
            case "bySize":
              value = (d[0].value / 1000).toFixed(0) + " vkB"
              break;
          }

          return `
            <div class="c3-tooltip"><table><tbody>
              <tr><td class="text-center" colspan="2">${formattedTime}</td></tr>
              <tr><td>${d[0].id} sat/vbyte</td><td>${value}</td></tr>
            </tbody></table></div>
            `
        }
      },
    }

    // properly destroy chart and generate new chart
    if (state.historicalMempool.chart) {
      state.historicalMempool.chart = state.historicalMempool.chart.destroy();
    }
    state.historicalMempool.chart = c3.generate(chartSetting)

  }
}

const transactionstatsCard = {
  switchType: function (switchTo) {
    if (state.transactionStats.data.type != switchTo && (switchTo == "count" || switchTo == "percentage")) {
      state.transactionStats.data.type = switchTo;
      drawChart(state.transactionStats.elementId)
    }
  },
  updateCardLastUpdated: function () {
    // calc seconds from milliseconds
    const minutes = Math.floor((Date.now() - (state.transactionStats.data.timeLastUpdated)) / 1000 / 60)
    document.getElementById('transaction-stats-last-update').innerHTML = (minutes)
  },
  processDataForChart: function (response) {
    state.transactionStats.data.timeLastUpdated = new Date();

    // We process data for the two display modi: byCount and byPercentage

    let columnsCount = [
      ['timestamps'],
      ['Replace-By-Fee count'],
      ['SegWit spending count'],
      ['Unconfirmed Transaction count']
    ]

    let columnsPercentage = [
      ['timestamps'],
      ['Replace-By-Fee percentage'],
      ['SegWit spending percentage'],
    ]

    response.forEach(function(element) {
      columnsCount[0].push(new Date(element.timestamp * 1000))
      columnsCount[1].push(element.rbfCount)
      columnsCount[2].push(element.segwitCount)
      columnsCount[3].push(element.txCount)

      columnsPercentage[0].push(new Date(element.timestamp * 1000))
      columnsPercentage[1].push(((element.rbfCount/ element.txCount)*100).toFixed(2))
      columnsPercentage[2].push(((element.segwitCount / element.txCount)*100).toFixed(2))
    });

    return {
      columnsCount: columnsCount,
      columnsPercentage: columnsPercentage,
    }
  },
  draw: async function () {

    let processed = state.transactionStats.data.processedStats
    let chartType = state.transactionStats.data.type

    if (processed == null) {
      return
    }

    chartSetting = {
      bindto: '#transaction-stats-chart',
      data: {
        x: "timestamps",
        xFormat: '%Y-%m-%dT%H:%M:%S.%LZ',
        columns: chartType == "count" ? processed.columnsCount : processed.columnsPercentage,
        type: "spline",
        order: null,
      },
      axis: {
        x: {
          type: 'timeseries',
          tick: {
            format: '%H:%M'
          }
        },
        y: {
          tick: chartType == "count" ? {
            format: d3.format(".2s")
          } : {},
          min: -0,
          padding: {
            top: 20,
            bottom: 0
          }
        }
      },
      size: {
        height: 450 // css #transaction-stats-chart min-height: 450px needs to be changed too
      },
      legend: {
        show: true
      },
      color: {
        pattern: ["#ffa600", "#e62470", "#223399"]
      },
      point: {
        show: false,
      },
      tooltip: {
        contents: function (d) {
          date = d[0].x
          let day = "0" + date.getDate();
          let month = "0" + (date.getMonth() + 1)
          let year = date.getFullYear()
          let hours = "0" + date.getHours();
          let minutes = "0" + date.getMinutes();
          let formattedTime = hours.substr(-2) + ':' + minutes.substr(-2) + " " + day.substr(-2) + "/" + month.substr(-2) + "/" + year

          if (state.transactionStats.data.type == "count") {
            let segWitPercentage = ((d[1].value / d[2].value) * 100).toFixed(2)
            let rbfPercentage = ((d[0].value / d[2].value) * 100).toFixed(2)

            return `
            <div class="c3-tooltip"><table><tbody>
              <tr><td class="text-center" colspan="2">${formattedTime}</td></tr>
              <tr><td>Total Transaction</td><td>${d[2].value} tx</td></tr>
              <tr><td>SegWit spending</td><td>${d[1].value} tx (${segWitPercentage}%)</td></tr>
              <tr><td>BIP125-RBF signaling</td><td>${d[0].value} tx (${rbfPercentage}%)</td></tr>
            </tbody></table></div>
            `
          } else if (state.transactionStats.data.type == "percentage") {
            return `
            <div class="c3-tooltip"><table><tbody>
              <tr><td class="text-center" colspan="2">${formattedTime}</td></tr>
              <tr><td>SegWit spending</td><td>${d[1].value} %</td></tr>
              <tr><td>BIP125-RBF signaling</td><td>${d[0].value} %</td></tr>
            </tbody></table></div>
            `
          }

        },
        horizontal: true
      }
    }

    // properly destroy chart and generate new chart
    if (state.transactionStats.chart) {
      state.transactionStats.chart = state.transactionStats.chart.destroy();
    }
    state.transactionStats.chart = c3.generate(chartSetting)

  }
}

const pastBlocksCard = {
  updateCardLastUpdated: function () {
    // calc seconds from milliseconds
    const minutes = Math.floor((Date.now() - (state.pastBlocks.data.timeLastUpdated)) / 1000 / 60)
    document.getElementById('past-blocks-last-update').innerHTML = (minutes)
  },
  processDataForChart: function (response) {
    state.pastBlocks.data.timeLastUpdated = new Date();

    let rows = [
      ["date", "block"]
    ]
    let lines = []
    let regions = []

    for (blockIndex in response) {
      
      let timestamp = response[blockIndex].timestamp * 1000
      let height = response[blockIndex].height

      rows.push([
        new Date(timestamp), 1
      ])

      lines.push({
        value: new Date(timestamp),
        text: "Block " + height,
      })
    }


    // add 10 minute grid lines

    var truncatedMin = parseInt(new Date().getMinutes()/10) * 10 // set the last digit from the minute value to zero
    var gridTime = (new Date(new Date().setMinutes(truncatedMin))).setSeconds(0)

    while (gridTime > response[9].timestamp * 1000) {
      lines.push({
        value: gridTime,
        class: "grid-10-min"
      })
      gridTime -= 10*60*1000;
    }

    lines.push({
      value: new Date(),
      text: 'Now',
      position: 'start',
      class: "red-line"
    })
    rows.push([new Date(), 1])

    // region from the last block to the current time
    // last block is here element zero
    regions.push({
      axis: 'x',
      start: response[0].timestamp * 1000,
      end: new Date(),
      class: 'time-since-last-block'
    }, )

    return {
      blocks: response,
      rows: rows,
      lines: lines,
      regions: regions,
      minHeight: response[9].height
    }
  },
  setTimer: function () {
    // clear timer if set
    if (state.pastBlocks.data.timer != null) {
      clearInterval(state.pastBlocks.data.timer);
    }

    // set new timer
    var sec = (new Date() / 1000 - state.pastBlocks.data.processedBlocks.blocks[0].timestamp).toFixed(0)

    function pad(val) {
      return val > 9 ? val : "0" + val;
    }
    state.pastBlocks.data.timer = setInterval(function () {
      ++sec
      $("#past-blocks-timer").html(pad(parseInt(sec / 60, 10)) + ":" + pad(sec % 60));
    }, 1000);

  },
  draw: async function () {

    let processed = state.pastBlocks.data.processedBlocks

    if (processed == null) {
      return
    }

    chartSetting = {
      bindto: '#past-blocks-chart',
      data: {
        x: 'date',
        xFormat: '%Y-%m-%dT%H:%M:%S.%LZ',
        rows: processed.rows
      },
      size: {
        height: 300 // css #past-blocks-chart min-height: 300px needs to be changed too
      },
      point: {
        show: true,
        r: 15,
        focus: {
          expand: {
            enabled: true,
            r: 17,
          }
        }
      },
      regions: processed.regions,
      padding: {
        top: 20,
        bottom: 20,
      },
      color: {
        pattern: ['#323232']
      },
      legend: {
        show: false
      },
      axis: {
        rotated: false,
        x: {
          type: 'timeseries',
          padding: {
            left: 2 * 60 * 1000,
            right: 2 * 60 * 1000 // 2 min in ms 
          },
          tick: {
            format: '%H:%M',
            fit: false,
          },
        },
        y: {
          show: false,
          max: 1.4,
          min: 0.8
        },
      },
      grid: {
        x: {
          lines: processed.lines
        }
      },
      tooltip: {
        contents: function (d, defaultTitleFormat, defaultValueFormat, color) {
          let block = state.pastBlocks.data.processedBlocks.blocks[9 - d[0].index]
          let foundMinAgo = (Date.now() - block.timestamp * 1000) / 1000 / 60
          return `
            <div class="c3-tooltip"><table><tbody>
              <tr><td>Height</td><td>${block.height}</td></tr>
              <tr><td>Transactions</td><td>${block.txCount}</td></tr>
              <tr><td>Size</td><td>${(block.size / 1000 / 1000).toFixed(2) + " MB"}</td></tr>
              <tr><td>Weight units</td><td>${(block.weight / 1000).toFixed(0) + " kWU"}</td></tr>
              <tr><td>Found</td><td>${foundMinAgo.toFixed(0) + " min ago"}</td></tr>
            </tbody></table></div>
            `
        }
      }
    }

    // properly destroy chart and generate new chart
    if (state.pastBlocks.chart) {
      state.pastBlocks.chart = state.pastBlocks.chart.destroy();
    }
    state.pastBlocks.chart = c3.generate(chartSetting)

  }
}

function reloadData() {

  // reload mempool chart data
  axios.get(apiHost + '/api/mempool')
    .then(function (response) {
      state.currentMempool.data.processedMempool = currentMempoolCard.processDataForChart(response.data)
      currentMempoolCard.updateCard(state.currentMempool.data.processedMempool)
      drawChart(state.currentMempool.elementId)
    });

  // reload past blocks data
  axios.get(apiHost + '/api/recentBlocks')
    .then(function (response) {
      state.pastBlocks.data.processedBlocks = pastBlocksCard.processDataForChart(response.data)
      pastBlocksCard.setTimer()
      drawChart(state.pastBlocks.elementId)
    });

  // reload transaction stats data
  axios.get(apiHost + '/api/transactionStats')
    .then(function (response) {
      state.transactionStats.data.processedStats = transactionstatsCard.processDataForChart(response.data)
      drawChart(state.transactionStats.elementId)
    });


  reloadHistoricalMempool()

  // reload data again in 30 seconds
  setTimeout(function () {
    reloadData()
  }, updateInterval);
}

function reloadHistoricalMempool() {
  axios.get(apiHost + '/api/historicalMempool/' + state.historicalMempool.data.timeframe + "/" + state.historicalMempool.data.bySelector)
    .then(function (response) {
      state.historicalMempool.data.processedMempool = historicalMempoolCard.processDataForChart(response.data)
      drawChart(state.historicalMempool.elementId)
    });
}

// draws charts for visible cards
function drawChart(id) {

  if (!document.hidden) {

    if (state.currentMempool.elementId == id && state.currentMempool.isScrolledIntoView) {
      // console.log("Drawing currentMempool chart")
      currentMempoolCard.draw()
    }

    if (state.historicalMempool.elementId == id && state.historicalMempool.isScrolledIntoView) {
      historicalMempoolCard.draw()
      // console.log("Drawing historicalMempool chart")
    }

    if (state.pastBlocks.elementId == id && state.pastBlocks.isScrolledIntoView) {
      pastBlocksCard.draw()
      // console.log("Drawing pastBlocks chart")
    }

    if (state.transactionStats.elementId == id && state.transactionStats.isScrolledIntoView) {
      transactionstatsCard.draw()
      // console.log("Drawing transactionStats chart")
    }

  } else {
    // TODO: notify the user that the chart hasn't been redrawn
    console.warn("Tab is not shown. Skipping chart refresh.")
  }

}


setInterval(function () {
  // Update the 'Time since last update' text for the current mempool card
  currentMempoolCard.updateCardLastUpdated()

  // Update the 'Time since last block' 
  pastBlocksCard.updateCardLastUpdated()
}, 10000);


function getTxFromApi(txId) {
  return axios.get(`https://api.bitaps.com/btc/v1/blockchain/transaction/${txId}`)
    .then(res => res.data.data)
    .catch(e => {
      console.error('Error getting data from explorer:', e)
      throw new Error('Could not find your transaction in the bitcoin network. Wait a few minutes before trying again.')
    })
}