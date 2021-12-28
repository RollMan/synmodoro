package apiHandlers

import (
  "fmt"
  "io"
  "net/http"
  "github.com/RollMan/synmodoro/app/db"
  "log"
  "errors"
  "database/sql"
)

func StartTimerHandler(w http.ResponseWriter, r *http.Request) {
  var err error
  leftTime := &db.Timer{}

  err = db.DbMap.SelectOne(leftTime, "SELECT * from lefttime")
  if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
      leftTime = &db.Timer{0, 0, false}
      err = db.DbMap.Insert(leftTime)
      if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        return
      }
    }else{
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
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
  return
}
