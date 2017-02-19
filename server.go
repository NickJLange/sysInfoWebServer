package main

import "net/http"
import "time"
import "log"
import "context"
import "fmt"
import "sync"
import "strconv"
import "flag"
import "os"
import "github.com/NickJLange/monitoring"
import "github.com/NickJLange/tickdatabase"
import "github.com/bmizerany/perks/quantile"

const (
  //LatestKey Used for fun and games
  LatestKey tickdatabase.Key = "Latest"
  //MaxHistory should cap the memory size of the program - need to play with this a bit
   MaxHistory=10*60*60
)

//TickDB Master Initialized in Main?
//TickDB Can I make this a locally scoped variable?
var TickDB tickdatabase.TTickDB



func main() {
  //Do some hosuekeeping
  portPtr := flag.Int("port",8008,"Listen Port")
  ipPtr := flag.String("ip","","IP to bind on - defaults to INADDR_ANY")
  //ifPtr := flag.String("if","","Interface to try to use - might not work on certain OS (Program will bomb.)")
  flag.Parse()
  //FIXME: Validate Data Quality
  if (*portPtr < 1024 || *portPtr > 32768){
     fmt.Printf("Invalid Port outside Daemon Range: %v", *portPtr)
     os.Exit(-1)
  }
  var bindAddr = fmt.Sprintf("%s:%d",*ipPtr,*portPtr)
  fmt.Printf("Starting up Listening to %v\n", bindAddr)
  TickDB.Current=make(map[tickdatabase.Key]string)
  TickDB.Historical=make(map[tickdatabase.Key][]string)
  // FML
  go monitoring.PopulateTemperature(TickDB)

  mux := http.NewServeMux()
  //TODO: Why does nil work?
  mux.HandleFunc("/now",reportLatestStatusWrapper(nil))
  mux.HandleFunc("/statistics", reportStatisticsWrapper(nil))
  mainServer:=http.Server {
  Addr: bindAddr,
  Handler: mux,
  ReadTimeout: 10 * time.Second,
  WriteTimeout: 10 * time.Second,
 }
  log.Fatal(mainServer.ListenAndServe())
}



func newContextWithLatestEntry(ctx context.Context, r *http.Request, val tickdatabase.TickDBEntry) context.Context {
  if val == nil {
   return context.TODO()
  }
 return context.WithValue(ctx, LatestKey, val)
}

func pullID(ctx context.Context) tickdatabase.TickDBEntry {

  return ctx.Value(LatestKey).(tickdatabase.TickDBEntry)
}


func reportLatestStatusWrapper(next http.HandlerFunc) http.HandlerFunc {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    //TODO: Is this passed by value or reference?
    var myLock sync.RWMutex
    myLock.RLock()
    latest := TickDB.Current
    myLock.RUnlock()
   // whoami := latestAndGreatest["fancy"]
    w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header)
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Current Status of: \n")
    for k,v := range(latest){
      fmt.Fprintf(w, "%v=%v\n",k,v)
    }
 })
}

func reportStatisticsWrapper(next http.HandlerFunc) http.HandlerFunc {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    start := time.Now()
    //TODO: Is this passed by value or reference?
    var myLock sync.RWMutex
    myLock.RLock()
    latest := TickDB.Historical
    myLock.RUnlock()
    w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header)
    w.WriteHeader(http.StatusOK)

    for k,v := range(latest){
      var sampleCount float64
      q := quantile.NewTargeted(0.0001,0.50, 0.90, 0.99,0.9999)
      for i := range(v){
            //fmt.Printf("parsing a float %v , %v\n", i, v[i])
            if f,err := strconv.ParseFloat(v[i], 64); err == nil {
             q.Insert(f)
             sampleCount++
          }
      }
      fmt.Fprintf(w, "Statistics for %v across sample size %v - p1: %v, p50: %v, p90: %v, p99: %v, p99.99: %v \n",k, sampleCount, q.Query(0.0001),q.Query(0.50),q.Query(.90),q.Query(0.99),q.Query(0.9999))
    }
    End := time.Since(start)
    fmt.Fprintf(w, "\n Page Generated in %v\n", End)
 })
}
