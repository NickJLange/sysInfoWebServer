package main

import "net/http"
import "time"
import "log"
import "io"
import "context"
import "fmt"
import "sync"
import "strconv"


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
  historical map[key][]string
}

//TickDB Master Initialized in Main?
//TickDB Can I make this a locally scoped variable?
var TickDB tTickDB



//start := time.Now()
//End = time.Since(start) (// this is really nice)

func populateTemperature(TickDB tTickDB) {
  start := time.Now()
  smcOpen()
  defer smcClose()
  End := time.Since(start)

  for  true {
   start = time.Now()
   var myLock sync.RWMutex
   myLock.Lock()
   x:=fmt.Sprintf("%v",readTemperature())

   TickDB.current["temperature"]=fmt.Sprintf("%v",x)
   TickDB.current["lastRun"] = fmt.Sprintf("%v", End.Nanoseconds())
//FIXME: I'm not sure this is correct
   TickDB.historical["temperature"] = append(TickDB.historical["temperature"],x)
   TickDB.historical["lastRun"] = append(TickDB.historical["lastRun"],fmt.Sprintf("%v",End.Nanoseconds()))
   myLock.Unlock()
   End = time.Since(start)
   fmt.Printf("It took this long to run %v\n",End)
   time.Sleep(10*time.Second)
 }
}


func main() {
  //Do some hosuekeeping

  TickDB.current=make(map[key]string)
  TickDB.historical=make(map[key][]string)
  // FML
  go populateTemperature(TickDB)

  mux := http.NewServeMux()
  //TODO: Why does nil work?
  mux.HandleFunc("/temperature",reportLatestStatusWrapper(nil))
  mux.HandleFunc("/statistics", reportStatisticsWrapper(nil))
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
    //TODO: Is this passed by value or reference?
    var myLock sync.RWMutex
    myLock.RLock()
    latest := TickDB.current
    myLock.RUnlock()
    ctx := newContextWithLatestEntry(r.Context(), r, latest)
    helloWorldHandler(w, r.WithContext(ctx))
 })
}

func reportStatisticsWrapper(next http.HandlerFunc) http.HandlerFunc {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    //TODO: Is this passed by value or reference?
    var myLock sync.RWMutex
    myLock.RLock()
    latest := TickDB.historical
    myLock.RUnlock()
    w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header)
    w.WriteHeader(http.StatusOK)
    for k,v := range(latest){
      var total,sampleCount float64
      for i := range(v){
            //fmt.Printf("parsing a float %v , %v\n", i, v[i])
            if f,err := strconv.ParseFloat(v[i], 64); err == nil {
             total+= f
             sampleCount++
          }
      }
      fmt.Fprintf(w, "Statistics for %v: Total Samples %v with an Average of %v\n",k,sampleCount,(total/sampleCount))
    }
 })
}



func helloWorldHandler(rw http.ResponseWriter, req *http.Request) {

 latestAndGreatest := pullID(req.Context())
// whoami := latestAndGreatest["fancy"]
 rw.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header)
 rw.WriteHeader(http.StatusOK)
 for k,v := range(latestAndGreatest){
   fmt.Fprintf(rw, "Hello World %v=%v\n",k,v)
 }
 io.WriteString(rw, "Goodbye 世界!\n")
}
