/* 
  monitor-filter.js includes functionallity to template the filters on
  mempool.observer/monitor from the gFilters global below. 
*/

const filterStates = {
  hide: "0",
  inactive: "1",
  show: "2",
  lessEqual: "0",
  greaterEqual: "1",
  equal: "2", 
};

let gFilters = {

  segwitSpending: {title:"SegWit spending", id:"segwitSpending", state: filterStates.inactive, type:"tri-state-switch", queryStringCode: "b", isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.show && tx.spendsSegWit) || (filter.state == filterStates.hide && !tx.spendsSegWit) || filter.state == filterStates.inactive){return true} return false
  }},

  multisigSpending: {title:"Multisig spending", id:"multisigSpending", state: filterStates.inactive, type:"tri-state-switch", queryStringCode: "c", isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.show && tx.spendsMultisig) || (filter.state == filterStates.hide && !tx.spendsMultisig) || filter.state == filterStates.inactive){return true}  return false
  }},

  rbf: {title:"Replace By Fee", id:"rbf", state: filterStates.inactive, type:"tri-state-switch", queryStringCode: "a", isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.show && tx.signalsRBF) || (filter.state == filterStates.hide && !tx.signalsRBF) || filter.state == filterStates.inactive ){return true} return false
  }},

  bip69:{title:"BIP-69 compliant", id:"bip69", state: filterStates.inactive, type:"tri-state-switch", queryStringCode: "d", isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.show && tx.isBIP69) || (filter.state == filterStates.hide && !tx.isBIP69) || filter.state == filterStates.inactive){return true} return false
  }},

  size:{title:"Size (vByte)", id:"size", state: filterStates.greaterEqual, type:"withinput", queryStringCode: "j", input: {value: "0", type: "number", step: 1, min: 0,  label: "transaction size", queryStringCode: "ji"}, isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.greaterEqual && tx.size >= filter.input.value) || (filter.state == filterStates.lessEqual && tx.size <= filter.input.value) || (filter.state == filterStates.equal && tx.size == filter.input.value)){ return true } return false
  }},
  
  fee:{title:"Fee (sat)", id:"fee", state: filterStates.greaterEqual, type:"withinput", queryStringCode: "i", input: {value: "0", type: "number", step: 1, min: 0, label: "transaction fees", queryStringCode: "ii"}, isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.greaterEqual && tx.fee >= filter.input.value) || (filter.state == filterStates.lessEqual && tx.fee <= filter.input.value) || (filter.state == filterStates.equal && tx.fee == filter.input.value)){ return true } return false
  }},
  
  locktime:{title:"Locktime", id:"locktime", state: filterStates.greaterEqual, type:"withinput", queryStringCode: "k", input: {value: "0", type: "number", step: 1, min: 0, label: "transaction locktime", queryStringCode: "ki"}, isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.greaterEqual && tx.locktime >= filter.input.value) || (filter.state == filterStates.lessEqual && tx.locktime <= filter.input.value) || (filter.state == filterStates.equal && tx.locktime == filter.input.value)){ return true } return false
  }},

  version:{title:"Version", id:"version", state: filterStates.greaterEqual, type:"withinput", queryStringCode: "v", input: {value: "1", type: "number", step: 1, min: 1, label: "version", queryStringCode: "vi"}, isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.greaterEqual && tx.version >= filter.input.value) || (filter.state == filterStates.lessEqual && tx.version <= filter.input.value) || (filter.state == filterStates.equal && tx.version == parseInt(filter.input.value))){ return true } return false
  }},

  seperatorInputs: {title: "Inputs", type:"separator"},

  inputcount:{title:"Input count", id:"inputcount", state: filterStates.greaterEqual, type:"withinput", queryStringCode: "f", input: {value: "0", type: "number", step: 1, min: 0, label: "input count", queryStringCode: "fi"}, isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.greaterEqual && tx.inputCount >= filter.input.value) || (filter.state == filterStates.lessEqual && tx.inputCount <= filter.input.value) || (filter.state == filterStates.equal && tx.inputCount == parseInt(filter.input.value))){ return true } return false
  }},

  multisigs:{title:"Multisig spending", id:"multisigs",  type:"multisigInput", selectFilled: false,  input: {value: "0-of-0", queryStringCode: "X"}, isVisibleFunc: function (filter, tx) {
    if (filter.input.value == "0-of-0" || filter.input.value == ""){
      return true // don't do anything on default
    }
    if (tx.multisigsSpend != null){
      if (tx.multisigsSpend[filter.input.value]!=undefined) {
        return true
      }
    }
    return false
  }},

  spendsP2PKH:{title:"<abbr title='Pay-to-Public-Key-Hash'>P2PKH</abbr>&nbsp;spending", id:"spendsP2PKH", state: filterStates.inactive, type:"tri-state-switch", queryStringCode: "n", isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.show && tx.spends["P2PKH"]!=undefined) || (filter.state == filterStates.hide && tx.spends["P2PKH"]==undefined) || filter.state == filterStates.inactive ){return true} return false
  }},

  spendsP2SH:{title:"<abbr title='Pay-to-Script-Hash'>P2SH</abbr>&nbsp;spending", id:"spendsP2SH", state: filterStates.inactive, type:"tri-state-switch", queryStringCode: "r", isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.show && tx.spends["P2SH"]!=undefined) || (filter.state == filterStates.hide && tx.spends["P2SH"]==undefined) || filter.state == filterStates.inactive ){return true} return false
  }},
  
  spendsP2SH_P2WPKH:{title:"<abbr title='Nested Pay-to-Witness-Public-Key-Hash or P2SH-P2WPKH'>Nested P2WPKH</abbr>&nbsp;spending", id:"spendsP2SH_P2WPKH", state: filterStates.inactive, type:"tri-state-switch", queryStringCode: "o", isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.show && tx.spends["P2SH_P2WPKH"]!=undefined) || (filter.state == filterStates.hide && tx.spends["P2SH_P2WPKH"]==undefined) || filter.state == filterStates.inactive ){return true} return false
  }},

  spendsP2SH_P2WSH:{title:"<abbr title='Nested Pay-to-Witness-Script-Hash or P2SH-P2WSH'>Nested P2WSH</abbr>&nbsp;spending", id:"spendsP2SH_P2WSH", state: filterStates.inactive, type:"tri-state-switch", queryStringCode: "s", isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.show && tx.spends["P2SH_P2WSH"]!=undefined) || (filter.state == filterStates.hide && tx.spends["P2SH_P2WSH"]==undefined) || filter.state == filterStates.inactive ){return true} return false
  }},

  spendsP2WPKH:{title:"<abbr title='Pay-to-Witness-Public-Key-Hash'>P2WPKH</abbr>&nbsp;spending", id:"spendsP2WPKH", state: filterStates.inactive, type:"tri-state-switch", queryStringCode: "p", isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.show && tx.spends["P2WPKH"]!=undefined) || (filter.state == filterStates.hide && tx.spends["P2WPKH"]==undefined) || filter.state == filterStates.inactive ){return true} return false
  }},

  spendsP2WSH:{title:"<abbr title='Pay-to-Witness-Script-Hash'>P2WSH</abbr>&nbsp;spending", id:"spendsP2WSH", state: filterStates.inactive, type:"tri-state-switch", queryStringCode: "t", isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.show && tx.spends["P2WSH"]!=undefined) || (filter.state == filterStates.hide && tx.spends["P2WSH"]==undefined) || filter.state == filterStates.inactive ){return true} return false
  }},
  
  spendsP2PK:{title:"<abbr title='Pay-to-Public-Key'>P2PK</abbr>&nbsp;spending", id:"spendsP2PK", state: filterStates.inactive, type:"tri-state-switch", queryStringCode: "m", isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.show && tx.spends["P2PK"]!=undefined) || (filter.state == filterStates.hide && tx.spends["P2PK"]==undefined) || filter.state == filterStates.inactive ){return true} return false
  }},

  spendsP2MS:{title:"<abbr title='Pay-to-Multisig'>P2MS</abbr>&nbsp;spending", id:"spendsP2MS", state: filterStates.inactive, type:"tri-state-switch", queryStringCode: "q", isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.show && tx.spends["P2MS"]!=undefined) || (filter.state == filterStates.hide && tx.spends["P2MS"]==undefined) || filter.state == filterStates.inactive ){return true} return false
  }},

  seperatorOutputs: {title: "Outputs", type:"separator"},

  outputcount:{title:"Output count", id:"outputcount", state: filterStates.greaterEqual, type:"withinput", queryStringCode: "g", input: {value: "0", type: "number", step: 1, min: 0, label: "output count", queryStringCode: "gi"}, isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.greaterEqual && tx.outputCount >= filter.input.value) || (filter.state == filterStates.lessEqual && tx.outputCount <= filter.input.value) || (filter.state == filterStates.equal && tx.outputCount == filter.input.value)){ return true } return false
  }},

  outputsum:{title:"Output sum (BTC)", id:"outputsum", state: filterStates.greaterEqual, type:"withinput", queryStringCode: "h", input: {value: "0", type: "number", step: 0.0000001, min: 0, label: "output sum", queryStringCode: "hi"}, isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.greaterEqual && tx.outputValue >= filter.input.value*100000000) || (filter.state == filterStates.lessEqual && tx.outputValue <= filter.input.value*100000000) || (filter.state == filterStates.equal && tx.outputValue == filter.input.value*100000000)){ return true } return false
  }},

  paystoP2PKH:{title:"paying to&nbsp;<abbr title='Pay-to-Public-Key-Hash'>P2PKH</abbr>", id:"paystoP2PKH", state: filterStates.inactive, type:"tri-state-switch", queryStringCode: "v", isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.show && tx.paysTo["P2PKH"]!=undefined) || (filter.state == filterStates.hide && tx.paysTo["P2PKH"]==undefined) || filter.state == filterStates.inactive ){return true} return false
  }},

  paystoP2SH:{title:"paying to&nbsp;<abbr title='Pay-to-Script-Hash'>P2SH</abbr>", id:"paystoP2SH", state: filterStates.inactive, type:"tri-state-switch", queryStringCode: "x", isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.show && tx.paysTo["P2SH"]!=undefined) || (filter.state == filterStates.hide && tx.paysTo["P2SH"]==undefined) || filter.state == filterStates.inactive ){return true} return false
  }},

  paystoP2WPKH:{title:"paying to&nbsp;<abbr title='Pay-to-Witness-Public-Key-Hash'>P2WPKH</abbr>", id:"paystoP2WPKH", state: filterStates.inactive, type:"tri-state-switch", queryStringCode: "w", isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.show && tx.paysTo["P2WPKH"]!=undefined) || (filter.state == filterStates.hide && tx.paysTo["P2WPKH"]==undefined) || filter.state == filterStates.inactive ){return true} return false
  }},

  paystoP2WSH:{title:"paying to&nbsp;<abbr title='Pay-to-Witness-Script-Hash'>P2WSH</abbr>", id:"paystoP2WSH", state: filterStates.inactive, type:"tri-state-switch", queryStringCode: "O", isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.show && tx.paysTo["P2WSH"]!=undefined) || (filter.state == filterStates.hide && tx.paysTo["P2WSH"]==undefined) || filter.state == filterStates.inactive ){return true} return false
  }},

  paystoP2PK:{title:"paying to&nbsp;<abbr title='Pay-to-Public-Key'>P2PK</abbr>", id:"paystoP2PK", state: filterStates.inactive, type:"tri-state-switch", queryStringCode: "u", isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.show && tx.paysTo["P2PK"]!=undefined) || (filter.state == filterStates.hide && tx.paysTo["P2PK"]==undefined) || filter.state == filterStates.inactive ){return true} return false
  }},

  paystoP2MS:{title:"paying to&nbsp;<abbr title='Pay-to-Multisig'>P2MS</abbr>", id:"paystoP2MS", state: filterStates.inactive, type:"tri-state-switch", queryStringCode: "I", isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.show && tx.paysTo["P2MS"]!=undefined) || (filter.state == filterStates.hide && tx.paysTo["P2MS"]==undefined) || filter.state == filterStates.inactive ){return true} return false
  }},

  paystoOPRETURN:{title:"paying to&nbsp;OP_RETURN", id:"paystoOPRETURN", state: filterStates.inactive, type:"tri-state-switch", queryStringCode: "z", isVisibleFunc: function (filter, tx) {
    if ((filter.state == filterStates.show && tx.paysTo["OPRETURN"]!=undefined) || (filter.state == filterStates.hide && tx.paysTo["OPRETURN"]==undefined) || filter.state == filterStates.inactive ){return true} return false
  }},

  opreturndata:{title:"OP_RETURN contains", id:"opreturndata", state: filterStates.greaterEqual, type:"onlyinput", queryStringCode: "l", input: {value: "", queryStringCode: "li", type: "text", label: "opreturn contains"}, isVisibleFunc: function (filter, tx) {
    if(filter.input.value.length > 0){ if (tx.opreturnData.length > 0){ if (tx.opreturnData.indexOf(filter.input.value) != -1){ return true }} return false } return true 
  }},

};

function drawFilters() {
  var target = document.getElementById("filters-row")
  for (var key in gFilters) {
    if (gFilters.hasOwnProperty(key)) {
      var filter = gFilters[key]    
      if (filter.type != "separator"){
        target.appendChild(templateFilter(filter))
      } else {
        target.appendChild(templateSeparator(filter))
      }
    }
  }
}

function templateSeparator(filter) {
  let parent = document.createElement("div");
  parent.classList.add("col-12")
  let h5 = document.createElement("h5");
  h5.classList.add("my-0")
  h5.appendChild(document.createTextNode(filter.title))
  parent.appendChild(h5)
  let hr = document.createElement("hr");
  hr.classList.add("mt-2")
  hr.classList.add("mb-3")
  parent.appendChild(hr)
  return parent
}

function templateFilter(filter) {
  let parent = document.createElement("div");
  parent.classList.add("col-lg-6")
  switch (filter.type) {
    case "tri-state-switch": parent.appendChild(templateTriStateSwitchFilter(filter)); break;
    case "withinput": parent.appendChild(templateWithinputFilter(filter)); break;
    case "onlyinput": parent.appendChild(templateInputOnlyFilter(filter)); break;
    case "multisigInput": parent.appendChild(templateMultisigInputFilter(filter)); break;
  }
  return parent
}

function templateTriStateSwitchFilter(filter) {
  let template = document.createElement("div")
  template.innerHTML = `
  <div class="input-group input-group-sm mb-3">
    <span class="input-group-text form-control bg-light text-dark">${filter.title}</span>
    <input type="range" step="1" value="${filter.state}" min="0" max="2" name="filter-tristateswitch-${filter.id}" data-filter-id="${filter.id}" autocomplete="off" class="form-control tri-state-switch tri-state-switch-eyes">
  </div>
  `
  return template
}

function templateWithinputFilter(filter) {
  let template = document.createElement("div")
  template.innerHTML = `
  <div class="input-group input-group-sm mb-3">
  <div class="input-group-prepend">
    <span class="input-group-text bg-light text-dark">${filter.title}</span>
  </div>
  <input type="${filter.input.type}" value="${filter.input.value}" class="form-control" step="${filter.input.step}" min="${filter.input.min}" data-filter-id="${filter.id}" autocomplete="off" name="filter-input-freetext" aria-label="${filter.input.label}">
  <input type="range" step="${filter.step}" value="${filter.state}" min="0" max="2" autocomplete="off" data-filter-id="${filter.id}" class="form-control tri-state-switch tri-state-switch-equal">
  </div>
  `
  return template
}

function templateInputOnlyFilter(filter) {
  let template = document.createElement("div")
  template.innerHTML = `
  <div class="input-group input-group-sm mb-3">
  <div class="input-group-prepend">
    <span class="input-group-text bg-light text-dark">${filter.title}</span>
  </div>
  <input type="${filter.input.type}" value="${filter.input.value}" class="form-control" step="${filter.input.step}" min="${filter.input.min}" data-filter-id="${filter.id}" autocomplete="off" name="filter-input-freetext" aria-label="${filter.input.label}">
  </div>
  `
  return template
}

function templateMultisigInputFilter(filter) {
  let template = document.createElement("div")
  template.innerHTML = `
  <div class="input-group input-group-sm mb-3">
  <div class="input-group-prepend">
    <span class="input-group-text bg-light text-dark">${filter.title}</span>
  </div>
  <select class="custom-select form-control" id="filter-select-multisig">
    <option value="">No selection</option>
  </select>
  </div>
  `
  return template
}

