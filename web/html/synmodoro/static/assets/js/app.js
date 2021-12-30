var timer = document.getElementById('timer');
var type = document.getElementById('type');
var type_button = document.getElementById('timerbutton');
const wsurl = 'ws://' + window.location.host + '/api/ws';
const stateuri = 'http://' + window.location.host + '/api/state';
const starturi = 'http://' + window.location.host + '/api/start';
var request = new XMLHttpRequest();

function parseState(statestr){
  const timerinfo = JSON.parse(statestr);
  return timerinfo;
}

function draw(work_status, endtime){
  type.innerHTML = work_status;
  const now_unixsec = Date.now() / 1000;
  const end_unixsec = endtime;
  const diff = Math.max(0, end_unixsec - now_unixsec);
  const diff_min = Math.floor(diff / 60);
  const diff_sec = Math.floor(diff) % 60;
  timer.innerHTML = diff_min + ':' + diff_sec;
  return diff;
}

function onOpen(event){
  console.log("Connected.");
}

function  onMessage(event) {
  const timerinfo = JSON.parse(event.data);
  const endtimeunixsec = timerinfo.EndTimeUnixSec;
  const state = timerinfo.Status;

  var interval = setInterval(function(){
    const diff = draw(state, endtimeunixsec);
    if(diff <= 0) {
      clearInterval(interval);
    }
  }, 1000);
}

function onClick() {
  fetch(starturi, {
    method: 'POST',
    mode: 'cors',
    cache: 'no-cache',
    credentials: 'same-origin',
    headers: {
      'Content-Type': 'text/plain',
    },
    redirect: 'follow',
    body: "",
  }).then(response => {
    const ok = response.status;
    if(ok != 200){
      type_button.innerHTML = 'ERROR!';
    }
  });
}

function onError(event){
  console.log("WS Connection Error.");
  console.log(event);
}

function onClose(event){
  console.log("WS Connection Closed.");
  console.log(event);
}

window.onload = function(){
  fetch(stateuri).then(response => response.text()).then(response=> {
    const statestr = response;
    const state = parseState(statestr);
    const interval = setInterval(function(){
      const diff = draw(state.Status, state.EndTimeUnixSec);
      if(diff <= 0) {
        clearInterval(interval);
      }
    }, 1000);
    connection = new WebSocket(wsurl);
    connection.onopen = onOpen;
    connection.onmessage = onMessage;
    connection.onerror = onError;
    connection.onclose = onClose;

    type_button.onclick = onClick;
  });
}


