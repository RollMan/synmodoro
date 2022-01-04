var timer = document.getElementById('timer');
var type = document.getElementById('type');
var type_button = document.getElementById('timerbutton');
var notify_elemenet = document.getElementById('notify');
var mates_element = document.getElementById('mates');
var editusername_element = document.getElementById('edit-username');
var username = randomUserID();
const wsurl = 'ws://' + window.location.host + '/api/ws';
const stateuri = 'http://' + window.location.host + '/api/state';
const starturi = 'http://' + window.location.host + '/api/start';
const registeruri = 'http://' + window.location.host + '/api/register';
const getmatesuri = 'http://' + window.location.host + '/api/get_mates';
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
  const diff_min = String(Math.floor(diff / 60));
  const diff_sec = String(Math.floor(diff) % 60);
  timer.innerHTML = diff_min.padStart(2, '0') + ':' + diff_sec.padStart(2, '0');
  return diff;
}

function randomUserID(){
  let chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  let rand_str = '';
  const len = 8;
  for ( var i = 0; i < len; i++ ) {
    rand_str += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  return rand_str
}

// function registerName(){
//   fetch(registeruri, {
//     method: 'POST',
//     mode: 'cors',
//     cache: 'no-cache',
//     credentials: 'same-origin',
//     headers: {
//       'Content-Type': 'text/plain',
//     },
//     redirect: 'follow',
//     body: username,
//   }).then(response => {
//     const ok = response.status;
//     if(ok != 200){
//       response.text().then(errmsg => {
//         console.log("Failed to register the username.");
//         console.log(errmsg);
//         showErrorMessage("Failed to register username: " + errmsg);
//       });
//     }
//   });
// 
// }

function onOpen(event){
  console.log("Connected.");
  notify_elemenet.classList.remove('active');
}

function showMates(mates){
    let html = ""
    html += "<table>";
    for (let i = 0; i < mates.length; i++) {
      let mate = mates[i];
      if (mate == username) {
        html += "<tr><td>" + String(mate) + "<span id=\"edit-username\">&#9999</span></td></tr>";
      }else{
        html += "<tr><td>" + String(mate) + "</td></tr>";
      }
    }
    html += "</table>";
    mates_element.innerHTML = html;

    editusername_element = document.getElementById('edit-username');
    editusername_element.onclick = registerUsername;
}

function  onMessage(event) {
  notify_elemenet.classList.remove('active');
  const timerinfo = JSON.parse(event.data);
  if ('EndTimeUnixSec' in timerinfo){
    const endtimeunixsec = timerinfo.timerinfo.EndTimeUnixSec;
    const state = timerinfo.timerinfo.Status;

    var interval = setInterval(function(){
      const diff = draw(state, endtimeunixsec);
      if(diff <= 0) {
        clearInterval(interval);
        const notification = new Notification("Finish " + state);
        const sound = new Audio("/sound/chaim.mp3");
        sound.play();
      }
    }, 1000);
  }else if ('Mates' in timerinfo){
    showMates(timerinfo.Mates);
  }else{
    console.log("Invalid websocket message.");
    console.log(timerinfo);
    showErrorMessage("Invalid websocket message: " + timerinfo);
  }
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
  notify_elemenet.innerHTML += message;

}

function onError(event){
  console.log("WS Connection Error.");
  console.log(event);
  showErrorMessage("Websocket connection error: " + event);
}

function onClose(event){
  console.log("WS Connection Closed.");
  console.log(event);
  showErrorMessage("Websocket connection error: " + event);
}

function registerUsername(event){
  editusername_element.onclick = "";
  editusername_element.innerHTML = "<input type=\"text\" id=\"newusername\"> <input type=\"submit\" id=\"submitusername\">";
  var submitusername_element = document.getElementById("submitusername");
  submitusername_element.onclick = function(event){
    let newusername_element = document.getElementById("newusername");
    let newusername = newusername_element.value;
    let prevname = username;
    let reqJson = {PrevName: prevname, NewName: newusername};
    let reqJsonStr = JSON.stringify(reqJson);
    fetch(registeruri, {
      method: 'POST',
      mode: 'cors',
      cache: 'no-cache',
      credentials: 'same-origin',
      headers: {
        'Content-Type': 'text/plain',
      },
      redirect: 'follow',
      body: reqJsonStr,
    }).then(response => {
      const ok = response.status;
      if (ok != 200){
        response.text().then(errmsg => {
          console.log("Failed to register username.");
          console.log(errmsg);
          showErrorMessage("Failed to register username: " + errmsg);
        });
      }
      username = newusername;
    });
  }
}

window.onload = function(){
  fetch(stateuri).then(response => response.text()).then(response=> {
    const statestr = response;
    const state = parseState(statestr);
    const interval = setInterval(function(){
      const diff = draw(state.Status, state.EndTimeUnixSec);
      if(diff <= 0) {
        clearInterval(interval);
        const notification = new Notification("Finish " + state);
        const sound = new Audio("/sound/chaim.mp3");
        sound.play();
      }
    }, 1000);

    // fetch(getmatesuri).then(response => response.json()).then(responseJson=>{
    //   showMates(responseJson.Mates);
    //   editusername_element = document.getElementById('edit-username');
    //   editusername_element.onclick = registerUsername;
    // });

    const query = new URLSearchParams({Username: username});
    connection = new WebSocket(wsurl + '?' + query);
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


