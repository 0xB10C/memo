// Constants
const NEXT_BLOCK_LABELS = ["1 vMB", "2 vMB", "3 vMB"]

var processedMempool = null
var blockchainTip = null

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
    data: {},
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
}

var cards = [state.currentMempool, state.historicalMempool, state.pastBlocks]

window.onload = function () {
  // Add event listeners to the search bar
  document.getElementById('button-lookup-txid').addEventListener('click', currentMempoolCard.handleTxSearch)
  document.getElementById('input-lookup-txid').addEventListener("keyup", function (event) {
    if (event.keyCode === 13) { // Handle enter
      event.preventDefault()
      document.getElementById('button-lookup-txid').click()
    }
  })
  document.getElementById('random-tx').addEventListener('click', currentMempoolCard.loadRandomTransactionFromApi)

  // add scroll listener 
  $(document).scroll(scrollEventHandler);

  // init data reload loop
  reloadData() 
}

function scrollEventHandler(){

  // handle scroll offset over 60px from top to animate the icon
  if ($(".navbar").offset().top > 60) {
    $(".navbar").addClass("scrolled");
  } else {
    $(".navbar").removeClass("scrolled");
  }
 
  cards.forEach(card => {
    let el = document.getElementById(card.elementId)
    
    let scrolledIntoView = isScrolledIntoView(el)

    if(card.isScrolledIntoView == false && scrolledIntoView == true) {
      // redraw chart if card just scrolled into view
      card.isScrolledIntoView = scrolledIntoView
      drawChart(card.elementId)
    } 
    
    else if (card.isScrolledIntoView == true && scrolledIntoView == false) {
      // destroy chart if card just scrolled out of view
      card.isScrolledIntoView = scrolledIntoView
      if(card.chart != null){
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
  processDataForChart: function(response) {
    state.currentMempool.data.timeLastUpdated = response.timestamp
  
    const mempoolSize = +(response.mempoolSize / 1000000).toFixed(2)
    const patternAreas = {'0to10': [], '11to100': [], '101to1k': [],'from1001': []}
  
    const lines = []
    const blocks = []
    const rowData = [[],[]]
  
    for (var feerate in response.mempoolData) {
      rowData[0].push(feerate.toString())
      rowData[1].push(response.mempoolData[feerate])
  
      // TODO: (0xb10c) Why do we do this and what does it?
      log1pOfCount = response.mempoolData[feerate]
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
  
  
    for (var position in response.positionsInGreedyBlocks) {
      if (position < 3) {
        blocks.push(response.positionsInGreedyBlocks[position])
        // add lines to show estimated next blocks on the mempool graph
        lines.push({
          value: response.positionsInGreedyBlocks[position]
        })
      }
    }
  
    let colorPattern = currentMempoolCard.generateColorPattern(patternAreas)
  
    // Sum all txs to get the total number of tx in the mempool
    const sum = Object.values(rowData[1]).reduce((a, b) => a + b, 0)
    return {
      "mempoolSize": mempoolSize,
      "blocks": blocks,
      "colorPattern": colorPattern,
      "lines": lines,
      "rowData": rowData,
      "sum": sum
    }
  },
  generateColorPattern: function(patternAreas) {
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
  draw: async function() {
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
        height: 750
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
            text: 'unconfirmed transactions'
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
            text: 'estimated blocks'
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
    if(state.currentMempool.chart){
      state.currentMempool.chart = state.currentMempool.chart.destroy(); 
    }
    state.currentMempool.chart = c3.generate(chartSetting)
  
    // draw the tx the chart if it's unconfirmed
    let tx = state.currentMempool.data.currentTx;
    if (tx != null) {
      if (tx.status.confirmed) {
        $('#tx-eta-data').html(`Confirmed (block ${tx.status.block_height}, ${minutes_since_confirmation} minutes ago)`) // FIXME: tx-eta-data is not used anymore
      } else {
        const feeRate = Math.floor(tx.fee / (tx.weight / 4))
        currentMempoolCard.drawUserTxByFeeRate(feeRate)
      }
    }
  
  },
  updateCard: function(processed) {
    const spanTxCount = document.getElementById('current-mempool-count')
    const spanMempoolSize = document.getElementById('current-mempool-size')
  
    const txCountInMempool = processed.sum
    const txSizeInMempool = processed.mempoolSize
  
    spanTxCount.innerHTML = txCountInMempool
    spanMempoolSize.innerHTML = txSizeInMempool
  
    timeSinceLastUpdate = 0
    currentMempoolCard.updateCardLastUpdated()
  },
  updateCardLastUpdated: function() {
  // calc minutes from milliseconds
  const minutes = Math.floor((Date.now() - (state.currentMempool.data.timeLastUpdated * 1000)) / 1000 / 60)
  document.getElementById('current-mempool-last-update').innerHTML = (minutes)
  },
  handleTxSearch: async function() {
    $('#input-lookup-txid').removeClass("is-invalid") // remove invalid tx message 
    $('#current-mempool-tx-data').hide() // hide the current tx data
  
    inputTxid = document.getElementById('input-lookup-txid').value
  
    // Check if the input has the format of a txid
    if (/^[a-fA-F0-9]{64}$/.test(inputTxid) == false) {
    
      $('#invalid-feedback').html('Invalid Bitcoin transaction id.') // Set invalidt tx message
      $('#input-lookup-txid').addClass("is-invalid") // Show invalid tx message
    
    } else {
  
      try {
        state.currentMempool.data.currentTx = await getTxFromApi(inputTxid)
        let tx = state.currentMempool.data.currentTx
  
        if (tx.status.confirmed) {
          let minutes_since_confirmation = Math.floor((Date.now() - tx.status.block_time * 1000) / 1000 / 60)
          let error = `The transaction is already confirmed and therefore not in the mempool. (block ${tx.status.block_height}, ${minutes_since_confirmation} minutes ago)`
          console.error(error)
          $('#invalid-feedback').html(error)
          $('#input-lookup-txid').addClass("is-invalid")
        } else { // tx is unconfirmed
          const vSize = (tx.weight / 4) // see Issue #11  
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
      } else if(processedMempool.rowData[0][index] == feeRate) {
        // get ruffly the middle position of the bar
        position += Math.round(processedMempool.rowData[1][index] / 2) 
      } else {
        break;
      }
    }
  
    return position
  },
  displayTransactionData: function (fee, vsize, feeRate) {
    $('#current-mempool-tx-data').show()
    $('#current-mempool-tx-data-fee').html(fee)
    $('#current-mempool-tx-data-size').html(vsize)
    $('#current-mempool-tx-data-feerate').html(feeRate)
  },
  loadRandomTransactionFromApi: async function() {
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

const pastBlocksCard = {
  updateCardLastUpdated: function () { 
    // calc seconds from milliseconds
    const minutes = Math.floor((Date.now() - (state.pastBlocks.data.timeLastUpdated)) / 1000 / 60) 
    document.getElementById('past-blocks-last-update').innerHTML = (minutes)
  },
  processDataForChart: function(response) {
    state.pastBlocks.data.timeLastUpdated = new Date();
    
    console.log(response)

    let rows = [["date", "block"]]
    let lines = []
    let regions = []
    
    for(blockIndex in response.data){
      let timestamp = response.data[blockIndex].receivedBlockTime * 1000
      let height = response.data[blockIndex].height
      
      rows.push([
          new Date(timestamp),  1
      ])

      lines.push({
        value: new Date(timestamp),
        text: "Block " + height,
      })
    
    }
    
    lines.push({value: new Date(), text: 'Now', position: 'start' , class: "red-line"})
    rows.push([new Date(), 1])

    // region from the last block to the current time
    // last block is here element zero
    regions.push({axis: 'x', start: response.data[0].receivedBlockTime * 1000, end: new Date(), class: 'time-since-last-block'},)

    return {
      blocks: response.data,
      rows: rows,
      lines: lines,
      regions: regions,
      minHeight: response.data[9].height
    }
  },
  setTimer: function (){
    // clear timer if set
    if(state.pastBlocks.data.timer != null) {
      clearInterval ( state.pastBlocks.data.timer );
    }

    // set new timer
    var sec = (new Date() / 1000 - state.pastBlocks.data.processedBlocks.blocks[0].receivedBlockTime).toFixed(0)
    function pad ( val ) { return val > 9 ? val : "0" + val; }
    state.pastBlocks.data.timer = setInterval( function(){
      ++sec
      $("#past-blocks-timer").html(pad(parseInt(sec/60,10)) + ":" + pad(sec%60));
    }, 1000);

  },
  draw: async function () {

    let processed = state.pastBlocks.data.processedBlocks

    if(processed == null){return}

    chartSetting = {
      bindto: '#past-blocks-chart',
      data: {
          x: 'date',
          xFormat: '%Y-%m-%dT%H:%M:%S.%LZ',
          rows: processed.rows
      },
      size: {
        height: 300
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
          let block = state.pastBlocks.data.processedBlocks.blocks[9-d[0].index]
          let miner = block.miner == null ? "Unknown" : block.miner.name
          let foundMinAgo = (Date.now() - block.receivedBlockTime * 1000) / 1000 / 60 
          return `
            <div class="c3-tooltip"><table><tbody>
              <tr><td>Height</td><td>${block.height}</td></tr>
              <tr><td>Transactions</td><td>${block.transactionsCount}</td></tr>
              <tr><td>Size</td><td>${(block.size / 1000 / 1000).toFixed(2) + " MB"}</td></tr>
              <tr><td>Miner</td><td>${miner}</td></tr>
              <tr><td>Found</td><td>${foundMinAgo.toFixed(0) + " min ago"}</td></tr>
            </tbody></table></div>
            `
        }
      }
    }
    
    // properly destroy chart and generate new chart
    if(state.pastBlocks.chart){
      state.pastBlocks.chart = state.pastBlocks.chart.destroy(); 
    }
    state.pastBlocks.chart = c3.generate(chartSetting)

  }
}

function reloadData() {
  
  // reload mempool chart data
  axios.get('https://mempool.observer/api/mempool')
    .then(function (response) {
      state.currentMempool.data.processedMempool = currentMempoolCard.processDataForChart(response.data)
      currentMempoolCard.updateCard(state.currentMempool.data.processedMempool)
      drawChart(state.currentMempool.elementId)
    });

  // reload last blocks data
  // using bitaps.com here since it's API gives the time first seen
  // e.g. blockstream.info only gives the miner timestamp, which 
  // is to flexibile for us here 
  axios.get('https://api.bitaps.com/btc/v1/blockchain/blocks/last/10')
  .then(function (response) {
    state.pastBlocks.data.processedBlocks = pastBlocksCard.processDataForChart(response.data)
    pastBlocksCard.setTimer()
    drawChart(state.pastBlocks.elementId)
  });

  // reload data again in 30 seconds
  setTimeout(function () {reloadData()}, 30000); 
}

// draws charts for visible cards
function drawChart(id) {
  
  if (!document.hidden) { 

    if (state.currentMempool.elementId == id && state.currentMempool.isScrolledIntoView){
      // console.log("Drawing currentMempool chart")
      currentMempoolCard.draw()
    }

    if (state.historicalMempool.elementId == id && state.historicalMempool.isScrolledIntoView) {
      // console.log("Drawing historicalMempool chart")
    }

    if (state.pastBlocks.elementId == id && state.pastBlocks.isScrolledIntoView) {
      pastBlocksCard.draw()
      // console.log("Drawing pastBlocks chart")
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


function getLastTenBlocksFromApi() {
  return axios.get(`https://blockstream.info/api/blocks`)
    .then(res => res.data)
    .catch(e => {
      console.error('Error getting data from explorer:', e)
      throw new Error('Could not load data from https://blockstream.info/api/blocks.')
    })
}

function getTxFromApi(txId) {
  return axios.get(`https://blockstream.info/api/tx/${txId}`)
    .then(res => res.data)
    .catch(e => {
      console.error('Error getting data from explorer:', e)
      throw new Error('Could not find your transaction in the bitcoin network. Wait a few minutes before trying again.')
    })
}