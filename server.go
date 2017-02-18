package main

import "net/http"
import "time"
import "log"
import "io"
import "context"
import "fmt"



type key string
const (
  //LatestKey Used for fun and games
  LatestKey key = "Latest"
  //MaxHistory should cap the memory size of the program - need to play with this a bit
   MaxHistory=10*60*60
)

//TickDBEntry  - Wrapper to do what I want to do
type TickDBEntry map[key]string
//tTickDB contains the latest and greatest time value
type tTickDB struct {
  current TickDBEntry
  historical []TickDBEntry
}

//TickDB Master Initialized in Main?
var TickDB tTickDB



//start := time.Now()
//End = time.Since(start) (// this is really nice)

func main() {
  //Do some hosuekeeping

  TickDB.current=make(map[key]string)
  TickDB.current["fancy"]="pansy"
  TickDB.historical = append(TickDB.historical,TickDB.current)


  mux := http.NewServeMux()
  mux.HandleFunc("/helloWorld",reportLatestStatusWrapper(nil))
  mainServer:=http.Server {
  Addr: ":8080",
  Handler: mux,
  ReadTimeout: 10 * time.Second,
  WriteTimeout: 10 * time.Second,
 }
  log.Fatal(mainServer.ListenAndServe())
}



func newContextWithLatestEntry(ctx context.Context, r *http.Request, val TickDBEntry) context.Context {
  if val == nil {
   return context.TODO()
  }
 return context.WithValue(ctx, LatestKey, val)
}

func pullID(ctx context.Context) TickDBEntry {

  return ctx.Value(LatestKey).(TickDBEntry)
}


func reportLatestStatusWrapper(next http.HandlerFunc) http.HandlerFunc {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    //TODO: Add Mutex
    //TODO: Is this passed by value or reference?
    latest := TickDB.current
    ctx := newContextWithLatestEntry(r.Context(), r, latest)
    helloWorldHandler(w, r.WithContext(ctx))
 })
}


func helloWorldHandler(rw http.ResponseWriter, req *http.Request) {

 latestAndGreatest := pullID(req.Context())
// whoami := latestAndGreatest["fancy"]
 rw.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header)
 rw.WriteHeader(http.StatusOK)
 fmt.Fprintf(rw, "Hello World %v\n",latestAndGreatest["fancy"])
 io.WriteString(rw, "Hello World!")
}
