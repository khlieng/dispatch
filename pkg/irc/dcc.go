package irc

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type DCCSend struct {
	File   string `json:"file"`
	IP     string `json:"ip"`
	Port   uint16 `json:"port"`
	Length uint64 `json:"length"`
}

func ParseDCCSend(ctcp *CTCP) *DCCSend {
	params := strings.Split(ctcp.Params, " ")

	if len(params) > 4 {
		ip, err := strconv.Atoi(params[2])
		port, err := strconv.Atoi(params[3])
		length, err := strconv.Atoi(params[4])

		if err != nil {
			return nil
		}

		ip3 := uint32ToIP(ip)

		filename := path.Base(params[1])
		if filename == "/" || filename == "." {
			filename = ""
		}

		return &DCCSend{
			File:   filename,
			IP:     ip3,
			Port:   uint16(port),
			Length: uint64(length),
		}
	}

	return nil
}

func (c *Client) Download(pack *DCCSend) {
	if !c.Autoget {
		// TODO: ask user if he/she wants to download the file
		return
	}
	c.Progress <- DownloadProgress{
		PercCompletion: 0,
		File:           pack.File,
	}
	file, err := os.OpenFile(filepath.Join(c.DownloadFolder, pack.File), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		c.Progress <- DownloadProgress{
			PercCompletion: -1,
			File:           pack.File,
			Error:          err,
		}
		return
	}
	defer file.Close()

	con, err := net.Dial("tcp", fmt.Sprintf("%s:%d", pack.IP, pack.Port))
	if err != nil {
		c.Progress <- DownloadProgress{
			PercCompletion: -1,
			File:           pack.File,
			Error:          err,
		}
		return
	}

	defer con.Close()

	var speed float64
	var prevUpdate time.Time
	secondsElapsed := int64(0)
	totalBytes := uint64(0)
	buf := make([]byte, 0, 4*1024)
	start := time.Now().UnixNano()
	for {
		n, err := con.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
		}

		if _, err := file.Write(buf); err != nil {
			c.Progress <- DownloadProgress{
				PercCompletion: -1,
				File:           pack.File,
				Error:          err,
			}
			return
		}

		cycleBytes := uint64(len(buf))
		totalBytes += cycleBytes
		percentage := round2(100 * float64(totalBytes) / float64(pack.Length))

		now := time.Now().UnixNano()
		secondsElapsed = (now - start) / 1e9
		speed = round2(float64(totalBytes) / (float64(secondsElapsed)))
		secondsToGo := round2((float64(pack.Length) - float64(totalBytes)) / speed)

		con.Write(byteRead(totalBytes))

		if time.Since(prevUpdate) >= time.Second {
			prevUpdate = time.Now()

			c.Progress <- DownloadProgress{
				Speed:          humanReadableByteCount(speed, true),
				PercCompletion: percentage,
				BytesRemaining: humanReadableByteCount(float64(pack.Length-totalBytes), false),
				BytesCompleted: humanReadableByteCount(float64(totalBytes), false),
				SecondsElapsed: secondsElapsed,
				SecondsToGo:    secondsToGo,
				File:           pack.File,
			}
		}
	}

	con.Write(byteRead(totalBytes))

	c.Progress <- DownloadProgress{
		Speed:          humanReadableByteCount(speed, true),
		PercCompletion: 100,
		BytesCompleted: humanReadableByteCount(float64(totalBytes), false),
		SecondsElapsed: secondsElapsed,
		File:           pack.File,
	}
}

type DownloadProgress struct {
	File           string  `json:"file"`
	Error          error   `json:"error"`
	BytesCompleted string  `json:"bytes_completed"`
	BytesRemaining string  `json:"bytes_remaining"`
	PercCompletion float64 `json:"perc_completion"`
	Speed          string  `json:"speed"`
	SecondsElapsed int64   `json:"elapsed"`
	SecondsToGo    float64 `json:"eta"`
}

func (p DownloadProgress) ToJSON() string {
	progress, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(progress)
}

func uint32ToIP(n int) string {
	var byte1 = n & 255
	var byte2 = ((n >> 8) & 255)
	var byte3 = ((n >> 16) & 255)
	var byte4 = ((n >> 24) & 255)
	return fmt.Sprintf("%d.%d.%d.%d", byte4, byte3, byte2, byte1)
}

func byteRead(totalBytes uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, totalBytes)
	return b
}

func round2(source float64) float64 {
	return math.Round(100*source) / 100
}

const (
	_ = 1.0 << (10 * iota)
	kibibyte
	mebibyte
	gibibyte
)

func humanReadableByteCount(b float64, speed bool) string {
	unit := ""
	value := b

	switch {
	case b >= gibibyte:
		unit = "GiB"
		value = value / gibibyte
	case b >= mebibyte:
		unit = "MiB"
		value = value / mebibyte
	case b >= kibibyte:
		unit = "KiB"
		value = value / kibibyte
	case b > 1 || b == 0:
		unit = "bytes"
	case b == 1:
		unit = "byte"
	}

	if speed {
		unit = unit + "/s"
	}

	stringValue := strings.TrimSuffix(
		fmt.Sprintf("%.2f", value), ".00",
	)

	return fmt.Sprintf("%s %s", stringValue, unit)
}
