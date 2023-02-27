package internal

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	defaultMemGB = 32
)

var numberRegx = regexp.MustCompile(`\d{2,}`)

func GetMeminfo() string {
	info := readMeminfo()
	lines := strings.Split(info, "\n")
	multiplier := 1
	for i, line := range lines {
		if strings.HasPrefix(line, "MemTotal:") {
			memTotal := parseLineValue(line)
			memTotalGB := KBToGB(memTotal)
			if memTotalGB >= defaultMemGB {
				return info
			}
			if defaultMemGB%memTotalGB == 0 {
				multiplier = defaultMemGB / memTotalGB
			} else {
				multiplier = defaultMemGB/memTotalGB + 1
			}
			lines[i] = fixLineValue(line, multiplier)
			continue
		}
		if strings.HasPrefix(line, "MemFree:") {
			lines[i] = fixLineValue(line, multiplier)
			continue
		}
		if strings.HasPrefix(line, "MemAvailable:") {
			lines[i] = fixLineValue(line, multiplier)
			continue
		}
		if strings.HasPrefix(line, "Buffers:") {
			lines[i] = fixLineValue(line, multiplier)
			continue
		}
		if strings.HasPrefix(line, "Cached:") {
			lines[i] = fixLineValue(line, multiplier)
			continue
		}
		if strings.HasPrefix(line, "DirectMap2M:") {
			lines[i] = fixLineValue(line, multiplier)
			continue
		}
	}

	info = strings.Join(lines, "\n")
	return info
}

func fixLineValue(line string, multiplier int) string {
	value := parseLineValue(line)
	value = value * multiplier
	line = replaceLineValue(line, value)
	return line
}

func parseLineValue(line string) int {
	find := numberRegx.FindString(line)
	find = strings.TrimSpace(find)
	number, err := strconv.Atoi(find)
	if err != nil {
		return 0
	}
	return number
}

func replaceLineValue(line string, value int) string {
	return numberRegx.ReplaceAllString(line, strconv.Itoa(value))
}

func readMeminfo() string {
	bytes, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return ""
	}
	data := string(bytes)
	return data
}
