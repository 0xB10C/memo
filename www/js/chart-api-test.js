var chart


function processApiMempoolDataForChart(response) {
  console.log(response.timestamp)

  patternColors = [
    ['darkslateblue', '667eea'],
    ['ff009f', 'ff9200'],
    ['orange', 'hotpink'],
    ['hotpink', 'darkblue'],
  ]
  const patternAreas = {
    '0to10': [],
    '11to100': [],
    '101to1k': [],
    'from1001': []
  }
  const lines = []
  const rowData = [[], []]

  for (var feerate in response.mempoolData) {
    rowData[0].push(feerate.toString())
    rowData[1].push(response.mempoolData[feerate])

    if (feerate <= 10){
      patternAreas['0to10'].push(Math.log1p(response.mempoolData[feerate]))
    } else if (feerate <= 100){
      patternAreas['11to100'].push(Math.log1p(response.mempoolData[feerate]))
    } else if (feerate <= 1000){
      patternAreas['101to1k'].push(Math.log1p(response.mempoolData[feerate]))
    } else {
      patternAreas['from1001'].push(Math.log1p(response.mempoolData[feerate]))
    }
  }

  var colorPattern = []
  var c_counter = 0
  for (area in patternAreas){
    logLimits = chroma.limits(patternAreas[area], 'l', patternAreas[area].length);
    pattern = chroma//.bezier(patternColors[c_counter])
    .scale(patternColors[c_counter]).mode('lch').classes(logLimits)
    .colors(patternAreas[area].length);

    colorPattern = colorPattern.concat(pattern)
    c_counter++
  }


  for (var position in response.positionsInGreedyBlocks) {
    if (position < 3) {
      lines.push({ value: response.positionsInGreedyBlocks[position], text: Number(position) + 1, position: 'start' })
    }
  }

  return {'colorPattern': colorPattern, "lines": lines, "rowData": rowData} 
}

window.onload = function () {
  axios.get('https://mempool.observer/api/mempool')
    .then(function (response) {
      processed = processApiMempoolDataForChart(response.data)
      draw(processed)
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
    point: { show: false },
    legend: { show: false },
    tooltip: { grouped: false },
    size: {
      height: 750,
      width: 350
    },
    color: { pattern: processed.colorPattern },
    grid: {
      y: {
        lines: processed.lines
      }
    },
  })
}

function redraw(){
  setTimeout(function () {
    axios.get('https://mempool.observer/api/mempool')
      .then(function (response) {
        processed = processApiMempoolDataForChart(response.data)
        draw(processed)
        redraw()
      });
  }, 60000);
  
}

