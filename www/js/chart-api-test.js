window.onload = function () {

  axios.get('https://mempool.observer/api/mempool')
    .then(function (response) {
      console.log(response.data.timestamp)
      draw(response.data)
    })
}

function draw(response) {

  const sizes = []
  const lines = []
  const colData = []
  const grpData = []

  for (var feerate in response.mempoolData) {
    colData.push([feerate.toString(), response.mempoolData[feerate]]);
    sizes.push(Math.log1p(response.mempoolData[feerate]))
    grpData.push(feerate)
  }

  logLimits = chroma.limits(sizes, 'l', sizes.length);
  var colorPattern = chroma.bezier(['#0d2738', 'hotpink'])
    .scale().mode('lch').classes(logLimits)
    .correctLightness().colors(colData.length);

  for (var position in response.positionsInGreedyBlocks) {
    lines.push({ value: position, text: position, position: 'start' })
  }

  var chart = c3.generate({
    data: {
      columns: colData,
      type: 'bar',
      groups: [grpData],
      order: null
    },
    point: { show: false },
    legend: { show: false },
    tooltip: { grouped: false },
    size: {
      height: 750,
      width: 300
    },
    color: { pattern: colorPattern },
    grid: {
      y: {
        lines: lines
      }
    }
  })
}