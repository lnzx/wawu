package internal

import (
	"github.com/shirou/gopsutil/v3/disk"
	"log"
)

const (
	defaultPath = "/"
)

func GetDiskUsed() uint64 {
	usage, err := disk.Usage(defaultPath)
	if err != nil {
		log.Println(err)
		return 0
	}
	return usage.Used
}
