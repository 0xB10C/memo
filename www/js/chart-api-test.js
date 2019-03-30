window.onload = function () {


  // Make a request for a user with a given ID
  axios.get('https://mempool.observer/api/mempool/byCount')
    .then(function (response) {
      console.log(response.data.timestamp)
      draw(response.data.data)
    })
}

function draw(data) {


  const sizes = []
  const colData = []
  const grpData = []

  for (var feerate in data) {
    colData.push([feerate.toString(), data[feerate]]);
    sizes.push(Math.log1p(data[feerate]))
    grpData.push(feerate)
  }

  logLimits = chroma.limits(sizes, 'l', sizes.length);
  var colorPattern = chroma.bezier(['#0d2738', 'hotpink'])
    .scale().mode('lch').classes(logLimits)
    .correctLightness().colors(colData.length);

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
        lines: [
          { value: 15000, text: 'next Block', position: 'start' },
          { value: 13000, text: '2nd next Block', position: 'start' },
        ]
      }
    }
  })
}