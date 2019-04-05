var chart
var timeSinceLastUpdate = 0

function generateColorPattern(patternAreas) {

  const patternColors = [
    ['LightGreen', 'YellowGreen'],
    ['YellowGreen', 'Bisque'],
    ['Bisque', 'Salmon'],
    ['Salmon', 'HotPink'],
  ]

  var colorPattern = []
  var c_counter = 0
  for (area in patternAreas) {
    logLimits = chroma.limits(patternAreas[area], 'l', patternAreas[area].length);
    pattern = chroma.scale(patternColors[c_counter])
      .mode('lch').classes(logLimits)
      .colors(patternAreas[area].length);

    colorPattern = colorPattern.concat(pattern)
    c_counter++
  }

  return colorPattern
}


function processApiMempoolDataForChart(response) {
  console.log('Processed data at timestamp: ', response.timestamp)

  const patternAreas = {
    '0to10': [],
    '11to100': [],
    '101to1k': [],
    'from1001': []
  }

  const lines = []
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

  // Draw lines to show estimated block sizes on the mempool graph
  for (var position in response.positionsInGreedyBlocks) {
    if (position == 0) {
      lines.push({
        value: response.positionsInGreedyBlocks[position],
        text: 'Next block (~1MB)',
        position: 'start'
      })
    }
    if (position > 0 && position < 3) {
      lines.push({
        value: response.positionsInGreedyBlocks[position],
        text: `${Number(position) + 1} blocks from now`,
        position: 'start'
      })
    }
  }

  let colorPattern = generateColorPattern(patternAreas)

  return {
    'colorPattern': colorPattern,
    "lines": lines,
    "rowData": rowData
  }
}

window.onload = function () {
  axios.get('https://mempool.observer/api/mempool')
    .then(function (response) {
      console.log(response.data)
      processed = processApiMempoolDataForChart(response.data)
      draw(processed)
      updateMainPage(processed)
      redraw()
    })
}

function draw(processed) {
  chart = c3.generate({
    data: {
      rows: processed.rowData,
      type: 'bar',
      groups: [processed.rowData[0]],
      order: null
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
        title: function (x, index) { return 'Tx: 129...231'; },
        name: function (name, ratio, id, index) { return name + ' sat/byte'; },
        value: function (value, ratio, id, index) {
          if (id <= 10) {
            return 'low'
          } else if (id <= 100) {
            return 'medium'
          } else if (id <= 1000) {
            return 'high'
          } else {
            return 'super high'
          }
        }
      }
    },
    size: {
      height: 750,
      width: 450
    },
    bar: {
      width: 140
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
      x : {
        padding: {
          left: 400,
          right: 400
        }
      },
      y: {
        label: {
          text: 'Unconfirmed transactions'
        },
        tick: {
          count: 1,
          values: [Object.values(processed.rowData[1]).reduce((a, b) => a + b, 0)] // Sum all txs
        }
      }
    }
  })
}

function redraw() {
  setTimeout(function () {
    axios.get('https://mempool.observer/api/mempool')
      .then(function (response) {
        processed = processApiMempoolDataForChart(response.data)
        draw(processed)
        updateMainPage(processed)
        redraw()
      });
  }, 60000);
}

function updateMainPage(processed) {
  const total = [Object.values(processed.rowData[1]).reduce((a, b) => a + b, 0)]
  document.getElementById('total-transactions').innerHTML = total + ' unconfirmed transasctions (23 MB).' // TODO: get size of mempool from API

  timeSinceLastUpdate = 0
  document.getElementById('last-update').innerHTML = 'last updated ' + timeSinceLastUpdate + ' minutes ago'
}

setInterval(function() {
  timeSinceLastUpdate += 1
  document.getElementById('last-update').innerHTML = 'last updated ' + timeSinceLastUpdate + ' minutes ago'
}, 60000);