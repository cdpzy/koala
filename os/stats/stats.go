package stats

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strconv"
	"time"
)

var _pid int
var _startTime = time.Now()

func init() {
	_pid = os.Getpid()
}

// GetGoroutine ..
func GetGoroutine() *pprof.Profile {
	return pprof.Lookup("goroutine")
}

// GetHeap .
func GetHeap() *pprof.Profile {
	return pprof.Lookup("heap")
}

// GetThreadcreate ..
func GetThreadcreate() *pprof.Profile {
	return pprof.Lookup("threadcreate")
}

// GetBlock ..
func GetBlock() *pprof.Profile {
	return pprof.Lookup("block")
}

// MemProf 创建内存分析文件
func MemProf(path string) (string, error) {
	filename := "mem-" + strconv.Itoa(_pid) + ".memprof"
	filepath, err := filepath.Abs(path + "/" + filename)
	if err != nil {
		return "", err
	}

	f, err := os.Create(filepath)
	if err != nil {
		return "", err
	}

	defer f.Close()

	runtime.GC()
	pprof.WriteHeapProfile(f)
	return filename, nil
}

// CPUProfile CUP分析文件
func CPUProfile(path string) (string, error) {
	filename := "cpu-" + strconv.Itoa(_pid) + ".pprof"
	filepath, err := filepath.Abs(path + "/" + filename)
	if err != nil {
		return "", err
	}

	f, err := os.Create(filepath)
	if err != nil {
		return "", err
	}

	defer f.Close()

	pprof.StartCPUProfile(f)
	time.Sleep(time.Duration(30) * time.Second)
	pprof.StopCPUProfile()

	return filename, nil
}

// GCSummary GC概述
type GCSummary struct {
	NumGC     int64   `json:"numGC"`
	Pause     string  `json:"pause"`
	PauseAvg  string  `json:"pause_avg"`
	Overhead  float64 `json:"overhead"`
	Alloc     string  `json:"alloc"`
	Sys       string  `json:"sys"`
	AllocRate string  `json:"allocrate"`
	Histogram string  `json:"histogram"`
}

// GetGCSummary 获取GC概述
func GetGCSummary() *GCSummary {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	gc := &debug.GCStats{PauseQuantiles: make([]time.Duration, 100)}
	debug.ReadGCStats(gc)

	elapsed := time.Now().Sub(_startTime)
	summary := &GCSummary{}
	if gc.NumGC > 0 {
		summary.NumGC = gc.NumGC
		summary.Pause = toS(gc.Pause[0])
		summary.PauseAvg = toS(avg(gc.Pause))
		summary.Overhead = float64(gc.PauseTotal) / float64(elapsed) * 100
		summary.Alloc = toH(memStats.Alloc)
		summary.Sys = toH(memStats.Sys)
		summary.AllocRate = toH(uint64(float64(memStats.TotalAlloc) / elapsed.Seconds()))
		summary.Histogram = fmt.Sprintf("%s %s %s", gc.PauseQuantiles[94], gc.PauseQuantiles[98], gc.PauseQuantiles[99])
	} else {
		summary.Alloc = toH(memStats.Alloc)
		summary.Sys = toH(memStats.Sys)
		summary.AllocRate = toH(uint64(float64(memStats.TotalAlloc) / elapsed.Seconds()))
	}

	return summary
}

func avg(items []time.Duration) time.Duration {
	var sum time.Duration
	for _, item := range items {
		sum += item
	}
	return time.Duration(int64(sum) / int64(len(items)))
}

// format bytes number friendly
func toH(bytes uint64) string {
	switch {
	case bytes < 1024:
		return fmt.Sprintf("%dB", bytes)
	case bytes < 1024*1024:
		return fmt.Sprintf("%.2fK", float64(bytes)/1024)
	case bytes < 1024*1024*1024:
		return fmt.Sprintf("%.2fM", float64(bytes)/1024/1024)
	default:
		return fmt.Sprintf("%.2fG", float64(bytes)/1024/1024/1024)
	}
}

// short string format
func toS(d time.Duration) string {

	u := uint64(d)
	if u < uint64(time.Second) {
		switch {
		case u == 0:
			return "0"
		case u < uint64(time.Microsecond):
			return fmt.Sprintf("%.2fns", float64(u))
		case u < uint64(time.Millisecond):
			return fmt.Sprintf("%.2fus", float64(u)/1000)
		default:
			return fmt.Sprintf("%.2fms", float64(u)/1000/1000)
		}
	} else {
		switch {
		case u < uint64(time.Minute):
			return fmt.Sprintf("%.2fs", float64(u)/1000/1000/1000)
		case u < uint64(time.Hour):
			return fmt.Sprintf("%.2fm", float64(u)/1000/1000/1000/60)
		default:
			return fmt.Sprintf("%.2fh", float64(u)/1000/1000/1000/60/60)
		}
	}

}
