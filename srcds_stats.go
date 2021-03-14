package srcds

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Stats struct {
	Cpu     *float64
	NetIn   *float64
	NetOut  *float64
	Uptime  *float64
	Maps    *float64
	Fps     *float64
	Players *float64
	Svms    *float64
	Ms      *float64
	Tick    *float64
}

func ParseStats(resp []byte) (Stats, error) {
	stats := Stats{}

	lines := strings.Split(string(resp), "\n")

	if len(lines) < 2 {
		return stats, fmt.Errorf("expected two lines in response")
	}
	re := regexp.MustCompile(`\s*([\d.]+)\s+([\d.]+)\s+([\d.]+)\s+([\d.]+)\s+([\d.]+)\s+([\d.]+)\s+([\d.]+)\s+([\d.]+)\s+([\d.]+)\s+([\d.]+)\s*`)
	m := re.FindStringSubmatch(lines[1])

	{
		x, err := strconv.ParseFloat(m[1], 64)
		if err != nil {
			return stats, err
		}
		stats.Cpu = &x

	}
	{
		x, err := strconv.ParseFloat(m[2], 64)
		if err != nil {
			return stats, err
		}
		stats.NetIn = &x
	}
	{
		x, err := strconv.ParseFloat(m[3], 64)
		if err != nil {
			return stats, err
		}
		stats.NetOut = &x
	}
	{
		x, err := strconv.ParseFloat(m[4], 64)
		if err != nil {
			return stats, err
		}
		stats.Uptime = &x
	}
	{
		x, err := strconv.ParseFloat(m[5], 64)
		if err != nil {
			return stats, err
		}
		stats.Maps = &x
	}
	{
		x, err := strconv.ParseFloat(m[6], 64)
		if err != nil {
			return stats, err
		}
		stats.Fps = &x
	}
	{
		x, err := strconv.ParseFloat(m[7], 64)
		if err != nil {
			return stats, err
		}
		stats.Players = &x
	}
	{
		x, err := strconv.ParseFloat(m[8], 64)
		if err != nil {
			return stats, err
		}
		stats.Svms = &x
	}
	{
		x, err := strconv.ParseFloat(m[9], 64)
		if err != nil {
			return stats, err
		}
		stats.Ms = &x
	}
	{
		x, err := strconv.ParseFloat(m[10], 64)
		if err != nil {
			return stats, err
		}
		stats.Tick = &x
	}

	return stats, nil
}

type Status struct {
	Hostname    *string
	Version     *string
	Address     *string
	PublicIp    *string
	Os          *string
	Type        *string
	Map         *string
	Players     *int32
	Bots        *int32
	MaxPlayers  *int32
	Hibernating *bool
}

func ParseStatus(resp []byte) (Status, error) {
	status := Status{}

	lines := strings.Split(string(resp), "\n")

	for _, line := range lines {

		if m := regexp.MustCompile(`hostname: (.*)`).FindStringSubmatch(line); len(m) > 0 {
			x := strings.TrimSpace(m[1])
			status.Hostname = &x
		} else if m := regexp.MustCompile(`version : (.*)`).FindStringSubmatch(line); len(m) > 0 {
			x := strings.TrimSpace(m[1])
			status.Version = &x
		} else if m := regexp.MustCompile(`udp/ip  : ([\d\.:]+)  \(public ip: (.*)\)`).FindStringSubmatch(line); len(m) > 0 {
			status.Address = &(m[1])
			status.PublicIp = &(m[2])
		} else if m := regexp.MustCompile(`os      : (.*)`).FindStringSubmatch(line); len(m) > 0 {
			x := strings.TrimSpace(m[1])
			status.Os = &x
		} else if m := regexp.MustCompile(`type    : (.*)`).FindStringSubmatch(line); len(m) > 0 {
			x := strings.TrimSpace(m[1])
			status.Type = &x
		} else if m := regexp.MustCompile(`map     : (.*)`).FindStringSubmatch(line); len(m) > 0 {
			x := strings.TrimSpace(m[1])
			status.Map = &x
		} else if m := regexp.MustCompile(`players : (\d+) humans, (\d+) bots \((\d+)/(\d+) max\) (\(hibernating\)|\(not hibernating\))?`).FindStringSubmatch(line); len(m) > 0 {
			{
				x, err := strconv.ParseInt(m[1], 10, 32)
				if err != nil {
					return status, err
				}
				y := int32(x)
				status.Players = &y
			}
			{
				x, err := strconv.ParseInt(m[2], 10, 32)
				if err != nil {
					return status, err
				}
				y := int32(x)
				status.Bots = &y
			}
			{
				x, err := strconv.ParseInt(m[3], 10, 32)
				if err != nil {
					return status, err
				}
				y := int32(x)
				status.MaxPlayers = &y
			}
			x := (m[5] == "(hibernating)")
			status.Hibernating = &x
		} else {
			log.Printf("unexpected line %q", line)
		}
	}

	// TODO(peterklein): parse player list

	return status, nil
}
