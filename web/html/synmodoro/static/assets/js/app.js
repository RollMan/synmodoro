var timer = document.getElementById('timer');
var type = document.getElementById('type');
var type_button = document.getElementById('timerbutton');
var notify_elemenet = document.getElementById('notify');
const wsurl = 'wss://' + window.location.host + '/api/ws';
const stateuri = 'https://' + window.location.host + '/api/state';
const starturi = 'https://' + window.location.host + '/api/start';
var request = new XMLHttpRequest();
const sound = new Audio("/sound/chaim.mp3");
const sound_slider = document.getElementById('volume')
const volume_check_button = document.getElementById('volumecheck')

function parseState(statestr){
  const timerinfo = JSON.parse(statestr);
  return timerinfo;
}

function draw(work_status, endtime){
  type.innerHTML = work_status;
  const now_unixsec = Date.now() / 1000;
  const end_unixsec = endtime;
  const diff = Math.max(0, end_unixsec - now_unixsec);
  const diff_min = String(Math.floor(diff / 60));
  const diff_sec = String(Math.floor(diff) % 60);
  timer.innerHTML = diff_min.padStart(2, '0') + ':' + diff_sec.padStart(2, '0');
  return diff;
}

function onOpen(event){
  console.log("Connected.");
  notify_elemenet.classList.remove('active');
}

function  onMessage(event) {
  notify_elemenet.classList.remove('active');
  const timerinfo = JSON.parse(event.data);
  const endtimeunixsec = timerinfo.EndTimeUnixSec;
  const state = timerinfo.Status;

  var interval = setInterval(function(){
    const diff = draw(state, endtimeunixsec);
    if(diff <= 0) {
      clearInterval(interval);
      const notification = new Notification("Finish " + state);
      sound.play();
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
      if (ok == 400){
        type_button.innerHTML = 'Timer already running.';
        setTimeout(function(){
          type_button.innerHTML = 'start timer';
        }, 3000);
      }else{
        type_button.innerHTML = 'ERROR!';
      }
    }
  });
}

function showErrorMessage(message) {
  notify_elemenet.classList.toggle('active');
  notify_elemenet.innerHTML = "Websocket connection error: " + message;

}

function onError(event){
  console.log("WS Connection Error.");
  console.log(event);
  showErrorMessage(event);
}

function onClose(event){
  console.log("WS Connection Closed.");
  console.log(event);
  showErrorMessage(event);
}

function volumeSliderHandler(e){
  sound.volume = sound_slider.value;
}

function volumeCheckButtonHandler(e){
  sound.pause();
  sound.currentTime=0.0;
  sound.play();
}

window.onload = function(){
  sound.volume = sound_slider.value;
  sound_slider.addEventListener("input", volumeSliderHandler);
  volume_check_button.addEventListener("click", volumeCheckButtonHandler);
  fetch(stateuri).then(response => response.text()).then(response=> {
    const statestr = response;
    const state = parseState(statestr);
    const interval = setInterval(function(){
      const diff = draw(state.Status, state.EndTimeUnixSec);
      if(diff <= 0) {
        clearInterval(interval);
        const notification = new Notification("Finish " + state);
        sound.play();
      }
    }, 1000);
    connection = new WebSocket(wsurl);
    connection.onopen = onOpen;
    connection.onmessage = onMessage;
    connection.onerror = onError;
    connection.onclose = onClose;

    type_button.onclick = onClick;
  });

  // Notification configuration
  if ("Notification" in window){
    let permission = Notification.permission;

    if (permission == "denied" || permission == "granted") {
      return;
    }

    Notification.requestPermission()
    .then(function(){
      let notification = new Notification("Notification enabled!");
    });
  }
}


