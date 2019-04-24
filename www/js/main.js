// Constants
const NEXT_BLOCK_LABELS = ["1 vMB", "2 vMB", "3 vMB"]

// State 
var chart = null
var lastMempoolDataUpdate = 0
var processedMempool = null
var currentTx = null

function generateColorPattern(patternAreas) {

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
}

function processApiMempoolDataForChart(response) {

  console.log('Processing mempool data from', response.timestamp)
  lastMempoolDataUpdate = response.timestamp

  const mempoolSize = +(response.mempoolSize / 1000000).toFixed(2)

  const patternAreas = {
    '0to10': [],
    '11to100': [],
    '101to1k': [],
    'from1001': []
  }

  const ticks = []
  const lines = []
  const blocks = []
  const rowData = [
    [],
    []
  ]

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

  let colorPattern = generateColorPattern(patternAreas)

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
}

window.onload = function () {
  // Add event listeners to the search bar
  document.getElementById('button-lookup-txid').addEventListener('click', handleTxSearch)
  document.getElementById('input-lookup-txid').addEventListener("keyup", function (event) {
    if (event.keyCode === 13) { // Handle enter
      event.preventDefault()
      document.getElementById('button-lookup-txid').click()
    }
  })

  // add scroll listener for icon
  $(window).scroll(function () {
    if ($(".navbar").offset().top > 60) {
      $(".navbar").addClass("scrolled");
    } else {
      $(".navbar").removeClass("scrolled");
    }
  });

  // Get the mempool data 
  axios.get('https://mempool.observer/api/mempool')
    .then(function (response) {
      processedMempool = processApiMempoolDataForChart(response.data)
      draw(processedMempool)
      updateCurrentMempoolCard(processedMempool)
      redraw() // init redraw loop
    })
}

async function draw(processed) {
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

  if(chart){
    chart = chart.destroy(); // properly destroy chart
  }

  chart = c3.generate(chartSetting)
  
  if (currentTx != null) {
    // Check if the transaction has been confirmed 
    const tx = await getTxFromApi(currentTx.txid) // in sat/vbyte
    currentTx = tx
    if (tx.status.confirmed) {
      $('#tx-eta-data').html(`Confirmed (block ${tx.status.block_height}, ${minutes_since_confirmation} minutes ago)`)
    } else {
      const feeRate = Math.floor(tx.fee / (tx.weight / 4))
      drawTxIdInChartByFeeRate(currentTx.txid, feeRate)
    }
  }
}

function redraw() {
  setTimeout(function () {
    axios.get('https://mempool.observer/api/mempool')
      .then(function (response) {
        processed = processApiMempoolDataForChart(response.data)
        updateCurrentMempoolCard(processed)
        if (!document.hidden) { // only get new data to redraw if the tab is focused
          draw(processed)
        } else {
          console.warn("Tab is not shown. Skipping chart refresh.")
        }
      });
    redraw()
  }, 30000);
}

function updateCurrentMempoolCard(processed) { //TODO: Change name of function

  const spanTxCount = document.getElementById('current-mempool-count')
  const spanMempoolSize = document.getElementById('current-mempool-size')

  const txCountInMempool = processed.sum
  const txSizeInMempool = processed.mempoolSize

  spanTxCount.innerHTML = txCountInMempool
  spanMempoolSize.innerHTML = txSizeInMempool

  timeSinceLastUpdate = 0
  updateCurrentMempoolCardLastUpdated()
}

function updateCurrentMempoolCardLastUpdated() {
  // format as milliseconds since 1.1.1970 UTC
  // * 1000 to convert from seconds to milliseconds
  const millislastMempoolDataUpdate = lastMempoolDataUpdate * 1000
  const minutes = Math.floor((Date.now() - millislastMempoolDataUpdate) / 1000 / 60)

  document.getElementById('current-mempool-last-update').innerHTML = (minutes)
}

// Update the 'Time since last update' text
setInterval(function () {
  updateCurrentMempoolCardLastUpdated()
}, 10000);

async function handleTxSearch() {
  $('#input-lookup-txid').removeClass("is-invalid") // remove invalid tx message 
  $('#current-mempool-tx-data').hide() // hide the current tx data

  txId = document.getElementById('input-lookup-txid').value

  // Check if txid is invalid
  if (/^[a-fA-F0-9]{64}$/.test(txId) == false) {
    currentTx = null
    $('#invalid-feedback').html('Invalid Bitcoin transaction id.') // Set invalidt tx message
    $('#input-lookup-txid').addClass("is-invalid") // Show invalid tx message
  } else {
    try {
      const tx = await getTxFromApi(txId)
      if (tx.status.confirmed) {
        currentTx = null
        let minutes_since_confirmation = Math.floor((Date.now() - tx.status.block_time * 1000) / 1000 / 60)
        let error = `The transaction is already confirmed and therefore not in the mempool. (block ${tx.status.block_height}, ${minutes_since_confirmation} minutes ago)`
        console.error(error)

        $('#invalid-feedback').html(error)
        $('#input-lookup-txid').addClass("is-invalid")
      } else {
        currentTx = tx
        const feeRate = Math.floor(tx.fee / (tx.weight / 4)) // see Issue #11  
        drawTxIdInChartByFeeRate(tx.txid, feeRate)
        renderDataTable(tx.fee, (tx.weight / 4), feeRate)
      }

    } catch (error) {
      console.error(error)
      $('#invalid-feedback').html(error)
      $('#input-lookup-txid').addClass("is-invalid")
    }
  }
}

function drawTxIdInChartByFeeRate(txid, feeRate) {
  const position = getTxPostionInChartByFeeRate(feeRate)
  // Remove existing line in case
  chart.ygrids.remove({
    class: 'user-tx'
  });
  // Draw a new line, but wait 500ms for c3.js to not bug out
  setTimeout(function () {
    chart.ygrids.add([{
      value: position,
      text: txid.substring(0, 8) + "...",
      class: 'user-tx',
      position: 'middle'
    }]);
  }, 500)

  // TODO: Might be a bug here since the tooltip doesn't show
  chart.tooltip.show({
    x: feeRate,
    index: 0,
    id: '1'
  })
}

function getTxPostionInChartByFeeRate(feeRate) {
  let position = 0
  for (index in processedMempool.rowData[0]) {
    if (processedMempool.rowData[0][index] <= feeRate) {
      position += processedMempool.rowData[1][index]
    } else {
      break;
    }
  }
  return position
}

function getTxFromApi(txId) {
  return axios.get(`https://blockstream.info/api/tx/${txId}`)
    .then(res => res.data)
    .catch(e => {
      console.error('Error getting data from explorer:', e)
      throw new Error('Could not find your transaction in the bitcoin network. Wait a few minutes before trying again.')
    })
}

function renderDataTable(fee, vsize, feeRate) {
  $('#current-mempool-tx-data').show()
  $('#tx-fee-data').html(fee)
  $('#tx-size-data').html(vsize)
  $('#tx-feerate-data').html(feeRate)
}