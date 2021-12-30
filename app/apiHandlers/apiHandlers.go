package apiHandlers

import (
  "fmt"
  "io"
  "net/http"
  "github.com/RollMan/synmodoro/app/db"
  "github.com/RollMan/synmodoro/app/ws"
  "log"
  "errors"
  "database/sql"
  "encoding/json"
  "strconv"
)

func currentTimerStatus() (*db.Timer, error){
  var err error
  leftTime := &db.Timer{}

  err = db.DbMap.SelectOne(leftTime, "SELECT * from lefttime")
  if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
      leftTime = &db.Timer{0, 0, false}
      err = db.DbMap.Insert(leftTime)
      if err != nil {
        return nil, err
      }
    }else{
      return nil, err
    }
  }
  return leftTime, nil
}

func stateResponse() ([]byte, error) {
  var err error
  leftTime, err := currentTimerStatus()
  if err != nil {
    return nil, err
  }

  var currentPhaseName string
  if leftTime.IsBreak {
    if leftTime.IsEnded(){
      currentPhaseName = "work"
    }else{
      currentPhaseName = "break"
    }
  }else{
    if leftTime.IsEnded(){
      currentPhaseName = "break"
    }else{
      currentPhaseName = "work"
    }
  }

  response_map := map[string]string {
    "Id"               : strconv.FormatUint(leftTime.Id, 10),
    "EndTimeUnixSec"   : strconv.FormatInt(leftTime.EndTimeUnixSec, 10),
    "Status"          : currentPhaseName,
  }

  response, err := json.Marshal(response_map)
  return response, err
}

func StartTimerHandler(hub *ws.Hub, w http.ResponseWriter, r *http.Request) {
  var err error

  leftTime, err := currentTimerStatus()
  if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        return
  }

  if !leftTime.IsEnded() {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusBadRequest)

    io.WriteString(w, "Timer is already started.\n")
    return
  }

  next := leftTime.Start()
  leftTime.EndTimeUnixSec = next.EndTimeUnixSec
  leftTime.IsBreak = next.IsBreak

  count, err := db.DbMap.Update(leftTime)
  if count == 0 {
    log.Printf("No entries updated.\n", count)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  if count > 1 {
    log.Printf("Unexpected error: more than one data (%d) updated.\n", count)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  if err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  w.WriteHeader(http.StatusOK)
  io.WriteString(w, fmt.Sprintln(leftTime))

  // Send timer start message via websocket
  response, err := stateResponse()
  if err != nil {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  hub.Broadcast <- response
  return
}

func StateHandler(w http.ResponseWriter, r *http.Request) {
  var err error
  response, err := stateResponse()
  if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        return
  }

  w.Header().Set("Content-Type", "text/json; charset=utf-8")
  w.WriteHeader(http.StatusOK)
  io.WriteString(w, string(response))
}
