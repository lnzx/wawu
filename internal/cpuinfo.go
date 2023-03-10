package internal

import (
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

const (
	minCPUCores = 6
)

var numberRegx = regexp.MustCompile(`\d`)
var stringRegx = regexp.MustCompile(`0-\d`)

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
			line = replaceNumberValue(line, minCPUCores)
		}
		if strings.HasPrefix(line, "cpu cores") {
			line = replaceNumberValue(line, minCPUCores)
		}
		processor = append(processor, line)
	}

	var processors []string
	for i := 0; i < minCPUCores; i++ {
		for _, line := range processor {
			if strings.HasPrefix(line, "processor") {
				line = replaceNumberValue(line, i)
			}
			if strings.HasPrefix(line, "core id") {
				line = replaceNumberValue(line, i)
			}
			if strings.HasPrefix(line, "apicid") {
				line = replaceNumberValue(line, i)
			}
			if strings.HasPrefix(line, "initial apicid") {
				line = replaceNumberValue(line, i)
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

func replaceNumberValue(line string, value int) string {
	return numberRegx.ReplaceAllString(line, strconv.Itoa(value))
}

func replaceStringValue(line string, value string) string {
	return stringRegx.ReplaceAllString(line, value)
}

func Lscpu() string {
	cmd := exec.Command("lscpu")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return ""
	}
	numCPU := runtime.NumCPU()
	log.Println("NumCPU:", numCPU)
	if numCPU >= minCPUCores {
		return string(out)
	}
	info := string(out)
	lines := strings.Split(info, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "CPU(s):") {
			line = replaceNumberValue(line, minCPUCores)
			lines[i] = line
		}
		if strings.HasPrefix(strings.TrimSpace(line), "On-line CPU(s) list:") {
			line = replaceStringValue(line, "0-5")
			lines[i] = line
		}
		if strings.HasPrefix(strings.TrimSpace(line), "Core(s) per socket:") ||
			strings.HasPrefix(strings.TrimSpace(line), "Core(s) per cluster:") {
			line = replaceNumberValue(line, minCPUCores)
			lines[i] = line
		}
		if strings.HasPrefix(strings.TrimSpace(line), "NUMA node0 CPU(s):") {
			line = replaceStringValue(line, "0-5")
			lines[i] = line
		}
	}

	info = strings.Join(lines, "\n")
	return info
}
