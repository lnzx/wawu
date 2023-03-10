package internal

import (
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

const (
	minCPUCores = 6
)

var cpuinfoRegx = regexp.MustCompile(`\d`)

func GetCpuinfo() string {
	info := readCpuinfo()
	numCPU := runtime.NumCPU()
	log.Println("NumCPU:", numCPU)
	if numCPU >= minCPUCores {
		return info
	}

	lines := strings.Split(info, "\n")
	var processor []string
	for _, line := range lines {
		if line == "" || strings.TrimSpace(line) == "" {
			processor = append(processor, line)
			break
		}
		if strings.HasPrefix(line, "siblings") {
			line = replaceProcessorValue(line, minCPUCores)
		}
		if strings.HasPrefix(line, "cpu cores") {
			line = replaceProcessorValue(line, minCPUCores)
		}
		processor = append(processor, line)
	}

	var processors []string
	for i := 0; i < minCPUCores; i++ {
		for _, line := range processor {
			if strings.HasPrefix(line, "processor") {
				line = replaceProcessorValue(line, i)
			}
			if strings.HasPrefix(line, "core id") {
				line = replaceProcessorValue(line, i)
			}
			if strings.HasPrefix(line, "apicid") {
				line = replaceProcessorValue(line, i)
			}
			if strings.HasPrefix(line, "initial apicid") {
				line = replaceProcessorValue(line, i)
			}
			processors = append(processors, line)
		}
	}

	info = strings.Join(processors, "\n")
	return info
}

func readCpuinfo() string {
	bytes, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return ""
	}
	data := string(bytes)
	return data
}

func replaceProcessorValue(line string, value int) string {
	return cpuinfoRegx.ReplaceAllString(line, strconv.Itoa(value))
}
