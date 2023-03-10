# wawu

```shell
/proc/cpuinfo
/proc/diskstats
/proc/meminfo
/proc/stat
/proc/swaps
/proc/uptime
/proc/slabinfo
/sys/devices/system/cpu/online
```

MemAvailable ≈ MemFree + Buffers + Cached


## 默认设备名
`/dev/sd` + ["a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n"]

## 默认路径
`/var/lib/wawu`

Run:

```shell
./wawu
```

Install
```shell
/usr/bin/wawu
/etc/systemd/system/
```

Usage
```shell
docker run -it \
      -v /var/lib/wawu/cpuinfo:/proc/cpuinfo:rw \
      -v /var/lib/wawu/meminfo:/proc/meminfo:rw \
      ubuntu /bin/bash
```

Fix:

cpu:
```shell
/var/lib/wawu/lscpu
```

disk:

```shell
df -B MB /proc/meminfo
```
