package main

import (
  "fmt"
  "log"
  "os"
  "net/http"

  "github.com/RollMan/synmodoro/app/apiHandlers"
  "github.com/RollMan/synmodoro/app/db"
  "github.com/RollMan/synmodoro/app/ws"

  "github.com/gorilla/mux"
)

func main(){
  log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
  log.Print("Server started.")
  {
    dsn := fmt.Sprintf("%s:%s@tcp(db:3306)/synmodoro?charset=utf8&parseTime=true", "root", os.Getenv("MYSQL_ROOT_PASSWORD"))
    log.Printf("dsn: %s", dsn)
    db.InitDB(dsn)
    log.Print("DB OK.")
  }

  hub := ws.NewHub()
  go hub.Run()

  r := mux.NewRouter()

  r.HandleFunc("/api/start", func(w http.ResponseWriter, r *http.Request){
    apiHandlers.StartTimerHandler(hub, w, r)
  }).Methods("POST")
  r.HandleFunc("/api/state", apiHandlers.StateHandler).Methods("GET")
  r.HandleFunc("/api/ws", func(w http.ResponseWriter, r *http.Request){
    ws.WsHandler(hub, w, r)
  })
  log.Fatal(http.ListenAndServe(":80", r))
}
