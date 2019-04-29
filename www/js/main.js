// Constants
const NEXT_BLOCK_LABELS = ["1 vMB", "2 vMB", "3 vMB"]

var processedMempool = null
var blockchainTip = null

// State
var state = {
  currentMempool: {
    chart: null,
    isScrolledIntoView: true, // TODO: Maybe have everything default to false at start?
    data: {
      processedMempool: null,
      currentTx: null,
      timeLastUpdated: null,
    },
  },
  historicalMempool: {
    chart: null,
    isScrolledIntoView: true, // TODO: Maybe have everything default to false at start?
    data: {},
  },
  pastBlocks: {
    chart: null,
    isScrolledIntoView: true, // TODO: Maybe have everything default to false at start?
    data: {},
  },
}


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
 
}

function isScrolledIntoView(el) {
  let elemTop = el.getBoundingClientRect().top;
  let elemBottom = el.getBoundingClientRect().bottom;
  let isVisible = (elemTop >= 0) && (elemBottom <= window.innerHeight);
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
  draw: async function(processed) {
    chartSetting = {
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
  
    // TODO: can be removed
    // const tx = await getTxFromApi(currentTx.txid) // in sat/vbyte
    // currentTx = tx
  
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
      class: 'user-tx'
    });
  
    // Draw a new line, but wait 20ms for c3.js to not bug out
    setTimeout(function () {
      state.currentMempool.chart.ygrids.add([{
        value: position,
        text: 'Your transaction',
        class: 'user-tx',
        position: 'middle'
      }]);
    }, 20)
  
    // TODO: Might be a bug here since the tooltip doesn't show
    state.currentMempool.chart.tooltip.show({
      x: feeRate,
      index: 0,
      id: '1'
    })
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
    $('#tx-fee-data').html(fee)
    $('#tx-size-data').html(vsize)
    $('#tx-feerate-data').html(feeRate)
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


function reloadData() {
  
  // reload mempool chart data
  axios.get('https://mempool.observer/api/mempool')
    .then(function (response) {
      state.currentMempool.data.processedMempool = currentMempoolCard.processDataForChart(response.data)

      let processedMempool = state.currentMempool.data.processedMempool
      currentMempoolCard.updateCard(processedMempool)

      // only redraw if the tab is focused
      if (!document.hidden) { 
        currentMempoolCard.draw(processedMempool)
      } else {
        console.warn("Tab is not shown. Skipping chart refresh.")
        // TODO: notify the user that the chart hasn't been redrawn
      }

    });

  // reload last blocks data
  axios.get('https://blockstream.info/api/blocks')
  .then(function (response) {
    console.log(response.data)
  });

  // reload data again in 30 seconds
  setTimeout(function () {reloadData()}, 30000); 
}

setInterval(function () {
  // Update the 'Time since last update' text for the current mempool card
  currentMempoolCard.updateCardLastUpdated()

  // Update the 'Time since last block' 
  updateTimeSinceLastBlock()
}, 10000);

function updateTimeSinceLastBlock(){ // TODO: move to pastBlocksCard function object
  
}

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