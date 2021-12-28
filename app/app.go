package main

import (
  "fmt"
  "log"
  "os"
  "net/http"

  "github.com/RollMan/synmodoro/app/apiHandlers"
  "github.com/RollMan/synmodoro/app/db"

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

  r := mux.NewRouter()

  r.HandleFunc("/api/start", apiHandlers.StartTimerHandler).Methods("POST")
  log.Fatal(http.ListenAndServe(":80", r))
}
