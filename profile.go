package helpers

import (
	"fmt"
	"github.com/bocajim/helpers/log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"
)

var activeCpuProfiling = false
var activeMemProfiling = false
var activeWebProfiling = false

func memProfileSaver(minuteInterval int, fileName string) {
	i := 1
	for {
		select {
		case <-time.After(time.Duration(minuteInterval) * time.Minute):
			f, err := os.Create(fileName + fmt.Sprintf(".%d", i*minuteInterval))
			if err != nil {
				log.Printf(log.Warn, "Could not write memory profile: "+err.Error())
				break
			}
			pprof.WriteHeapProfile(f)
			log.Printf(log.Info, "Saved memory profile.")
			i++
			f.Close()
			break
		}
	}
}

func ConfigureProfiling(config string) {

	s := strings.Split(config, ",")
	if len(s) >= 2 {
		switch s[0] {
		case "cpu", "CPU":
			activeCpuProfiling = true
			f, err := os.Create(s[1])
			if err != nil {
				log.Printf(log.Warn, "Could not initialize CPU profiler: "+err.Error())
				break
			}
			pprof.StartCPUProfile(f)
			break
		case "mem", "MEM":
			activeMemProfiling = true
			interval := 10
			if len(s) == 3 {
				v, err := strconv.Atoi(s[2])
				if err == nil {
					interval = v
				}
			}
			go memProfileSaver(interval, s[1])
			break
		case "web", "WEB":
			activeWebProfiling = true
			go func() {
				err := http.ListenAndServe(s[1], nil)
				if err != nil {
					log.Printf(log.Warn, "Could not initialize debug web server: "+err.Error())
				}
			}()
			break
		}
	}
}

func DeferredProfiling() {
	if activeCpuProfiling {
		pprof.StopCPUProfile()
	}
}
