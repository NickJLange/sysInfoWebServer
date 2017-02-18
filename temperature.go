package main


// #include <stdio.h>
// #include <string.h>
// #include <smc.h>
//#include <CoreFoundation/CoreFoundation.h>
//#include <CoreFoundation/CFArray.h>
//#include <IOKit/IOKitLib.h>
//#include <IOKit/ps/IOPSKeys.h>
//#include <IOKit/ps/IOPowerSources.h>
// #cgo CFLAGS: -framework IOKit -framework CoreFoundation  -stdlib=libstdc++ -Wno-deprecated-declarations
// #cgo LDFLAGS: -framework IOKit -framework CoreFoundation  -stdlib=libstdc++ -Wno-deprecated-declarations
import "C"


//CPUKeyString is a magic flag from the .h file - I don't think I can steal it.
const CPUKeyString string  = "TC0P"

func smcOpen() {
   C.SMCOpen();
}

func smcClose() {
   C.SMCClose()
}

func readTemperature() float64 {
   var tossMe C.double
   tossMe = C.SMCGetTemperature(C.CString(CPUKeyString))
   //fmt.Printf("ECHO WORLD %v", tossMe)
   return float64(tossMe)
}
