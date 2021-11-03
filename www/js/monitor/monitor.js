/* 
  monitor.js is the main file for mempool.observer/monitor.
  It handles data loading, querystring parsing and start drawing
  the chart and the filters. 
*/

const gApiHost = window.location.origin

async function loadEntryData(){
  if(gData == null){
    await d3.json(gApiHost + "/api/getMempoolEntries").then(function (data) { 
      gData = data
      var xMin = d3.min(data, function (d) {return xValue(d)})
      var xMax = d3.max(data, function (d) {return xValue(d)})
      document.getElementById("span-transaction-loaded").innerText = data.length
      document.getElementById("span-minute-range").innerText = Math.trunc((xMax-xMin)/60/1000)
      return data
    });
  } 
  return gData
}

async function loadRecentFeerateAPIData(){  
  if ( gRecentFeerateAPIData == null ) {
    await d3.json(gApiHost + "/api/getRecentFeerateAPIData").then(function (data) { 
      gRecentFeerateAPIData = data
    });
  }
  return gRecentFeerateAPIData
}

async function loadBlockEntriesData(){  
  if ( gBlockEntriesData == null ) {
    await d3.json(gApiHost + "/api/getBlockEntries").then(function (data) { 
      for (const block of data) {
        block.shortTXIDs.sort();        
      }
      gBlockEntriesData = data
    });
  }
  return gBlockEntriesData
}

function scrollEventHandler() {
  // handle scroll offset over 60px from top to animate the mempool observer icon
  let navbar = document.getElementsByClassName("navbar")
  if (window.pageYOffset > 60) {
    setTimeout(function() {
      navbar[0].classList.add("scrolled")
    },0)
  } else {
    navbar[0].classList.remove("scrolled")
  }
}

function isMobile(){
  return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)
}

// getDecodedQueryString decodes the base64 encoded querystring which is the value of `filter`
function getDecodedQueryString() {
  let windowParams = new URLSearchParams(window.location.search);
  if (!windowParams.has("filter")){
    return ""
  }
  
  let filtersEncoded = windowParams.get("filter")
  filters = atob(filtersEncoded)
  return filters.toString()
}

// setDecodedQueryString encodes and then sets the passed querystring as the value of `filter`
function setDecodedQueryString(qstring){
  var windowParams = new URLSearchParams(window.location.search);
  filtersEncoded = btoa(qstring)
  if (filtersEncoded == ""){
    windowParams.delete("filter")
    history.pushState(null, '',  window.location.pathname);
  }else{
    windowParams.set("filter", filtersEncoded)
    history.replaceState(null, '',  window.location.pathname + '?' + windowParams.toString());
  }
}

/* The value of the filter querystring is the base64 encoded normal querystring. */


function deleteQueryStringParameter(parameter){
  if ('URLSearchParams' in window) { // some browsers don't support the URLSearchParams API (yet)
    let qstring = getDecodedQueryString()

    let filterParams = new URLSearchParams(qstring);
    filterParams.delete(parameter)

    setDecodedQueryString(filterParams.toString())
  }
}

function setQueryStringParameter(parameter, value){
  if ('URLSearchParams' in window) { // some browsers don't support the URLSearchParams API (yet)
    let qstring = getDecodedQueryString()
    
    let filterParams = new URLSearchParams(qstring);
    filterParams.set(parameter, value)

    setDecodedQueryString(filterParams.toString())
  }
}

function readInitalQueryString(){
  if ('URLSearchParams' in window) { // some browsers don't support the URLSearchParams API (yet)
  qstring = getDecodedQueryString()
  let filterParams = new URLSearchParams(qstring);
    for (var p of filterParams) {
      for (var key in gFilters) {
        var filter = gFilters[key]
        var value = sanitizeHTML(p[1])
        if (p[0] == filter.queryStringCode){
          filter.state = sanitizeHTML(value)
          break;
        } else if (filter.type == "withinput" || filter.type == "multisigInput" || filter.type == "onlyinput"){
          if (p[0] == filter.input.queryStringCode){
            filter.input.value = sanitizeHTML(value)
            break;
          }
        }
      }

    switch (p[0]) {
      case "H":
        setSelectValue("select-highlight", p[1]); break;
      case "R":
        setSelectValue("select-radius", p[1]); break;
      case "E":
        setSelectValue("select-estimator", p[1]); break;
      }
    }
  }
}

function setSelectValue(id, val) {
  let e = document.getElementById(id)
  e.value = val;
}

function disableEveryInput(b) {
  let inputs = document.querySelectorAll("input, select");
  for(var input in inputs){
    inputs[input].disabled = b
  }
}

window.onload = function () {
  disableEveryInput(true) // disable every input while the site is loading. gets enabled when redraw() finishes
  if (isMobile()) {
    document.getElementById("alert-mobile").style.display = "block";
  }
  readInitalQueryString()
  drawFilters()
  loadRecentFeerateAPIData() // preload current feerate API data
  loadBlockEntriesData() // preload current block entries data

  redraw().then(function(){
    // Draw for the first time to initialize.
    disableEveryInput(false);
    document.getElementById("select-highlight").dispatchEvent(new Event('change'));
    document.getElementById("select-radius").dispatchEvent(new Event('change'));
    document.getElementById("select-estimator").dispatchEvent(new Event('change'));
    // adding the scroll event handling for the mempoolobserver logo after the redraw helps
    // to smooth out the animation on mobile and other low powered devices
    document.addEventListener("scroll", scrollEventHandler)
  }) 
  window.addEventListener("resize", redraw); // Redraw based on the new size whenever the browser window is resized.
}

/*!
 * Sanitize and encode all HTML in a user-submitted string
 * (c) 2018 Chris Ferdinandi, MIT License, https://gomakethings.com
 * @param  {String} str  The user-submitted string
 * @return {String} str  The sanitized string
 */
var sanitizeHTML = function (str) {
	var temp = document.createElement('div');
	temp.textContent = str;
	return temp.innerHTML;
};
