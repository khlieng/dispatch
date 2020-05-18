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
	Port   string `json:"port"`
	Length uint64 `json:"length"`
}

func ParseDCCSend(ctcp *CTCP) *DCCSend {
	params := strings.Split(ctcp.Params, " ")

	if len(params) > 4 {
		ip, err := strconv.Atoi(params[2])
		if err != nil {
			return nil
		}

		length, err := strconv.ParseUint(params[4], 10, 64)
		if err != nil {
			return nil
		}

		filename := path.Base(params[1])
		if filename == "/" || filename == "." {
			filename = ""
		}

		return &DCCSend{
			File:   filename,
			IP:     intToIP(ip),
			Port:   params[3],
			Length: length,
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
		File: pack.File,
	}

	file, err := os.OpenFile(filepath.Join(c.DownloadFolder, pack.File), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		c.downloadFailed(pack, err)
		return
	}
	defer file.Close()

	conn, err := net.Dial("tcp", net.JoinHostPort(pack.IP, pack.Port))
	if err != nil {
		c.downloadFailed(pack, err)
		return
	}
	defer conn.Close()

	totalBytes := uint64(0)
	accBytes := uint64(0)
	averageSpeed := float64(0)
	buf := make([]byte, 4*1024)
	start := time.Now()
	prevUpdate := start

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				c.downloadFailed(pack, err)
				return
			}
			if n == 0 {
				break
			}
		}

		if _, err := file.Write(buf[:n]); err != nil {
			c.downloadFailed(pack, err)
			return
		}

		accBytes += uint64(n)
		totalBytes += uint64(n)

		conn.Write(uint64Bytes(totalBytes))

		if dt := time.Since(prevUpdate); dt >= time.Second {
			prevUpdate = time.Now()

			speed := float64(accBytes) / dt.Seconds()
			if averageSpeed == 0 {
				averageSpeed = speed
			} else {
				averageSpeed = 0.2*speed + 0.8*averageSpeed
			}
			accBytes = 0

			bytesRemaining := float64(pack.Length - totalBytes)
			percentage := 100 * (float64(totalBytes) / float64(pack.Length))

			c.Progress <- DownloadProgress{
				Speed:          humanReadableByteCount(averageSpeed, true),
				PercCompletion: percentage,
				BytesRemaining: humanReadableByteCount(bytesRemaining, false),
				BytesCompleted: humanReadableByteCount(float64(totalBytes), false),
				SecondsElapsed: secondsSince(start),
				SecondsToGo:    bytesRemaining / averageSpeed,
				File:           pack.File,
			}
		}
	}

	// TODO: is this needed?
	conn.Write(uint64Bytes(totalBytes))

	c.Progress <- DownloadProgress{
		PercCompletion: 100,
		BytesCompleted: humanReadableByteCount(float64(totalBytes), false),
		SecondsElapsed: secondsSince(start),
		File:           pack.File,
	}
}

func (c *Client) downloadFailed(pack *DCCSend, err error) {
	c.Progress <- DownloadProgress{
		PercCompletion: -1,
		File:           pack.File,
		Error:          err,
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

func intToIP(n int) string {
	var byte1 = n & 255
	var byte2 = ((n >> 8) & 255)
	var byte3 = ((n >> 16) & 255)
	var byte4 = ((n >> 24) & 255)
	return fmt.Sprintf("%d.%d.%d.%d", byte4, byte3, byte2, byte1)
}

func uint64Bytes(i uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b
}

func secondsSince(t time.Time) int64 {
	return int64(math.Round(time.Since(t).Seconds()))
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
