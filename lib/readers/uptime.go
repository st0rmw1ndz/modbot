// modbot is a system information agregator
// Copyright (C) 2024 frosty <inthishouseofcards@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package readers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	SecondsInMinute = 60
	SecondsInHour   = 3600
	SecondsInDay    = 86400
)

type UptimeInfo uint64

func (u UptimeInfo) Days() int {
	return int(u) / SecondsInDay
}

func (u UptimeInfo) Hours() int {
	return (int(u) % SecondsInDay) / SecondsInHour
}

func (u UptimeInfo) Minutes() int {
	return (int(u) % SecondsInHour) / SecondsInMinute
}

func (u UptimeInfo) Seconds() int {
	return int(u) % SecondsInMinute
}

func (u UptimeInfo) String() string {
	var builder strings.Builder

	if u.Days() > 0 {
		builder.WriteString(fmt.Sprintf("%d days, ", u.Days()))
	}
	if u.Hours() > 0 {
		builder.WriteString(fmt.Sprintf("%d hours, ", u.Hours()))
	}
	if u.Minutes() > 0 {
		builder.WriteString(fmt.Sprintf("%d minutes, ", u.Minutes()))
	}
	builder.WriteString(fmt.Sprintf("%d seconds", u.Seconds()))

	return builder.String()
}

func ReadUptime() func() (interface{}, error) {
	return func() (interface{}, error) {
		const uptimePath = "/proc/uptime"

		file, err := os.Open(uptimePath)
		if err != nil {
			return 0, fmt.Errorf("failed to open %s: %w", uptimePath, err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		if !scanner.Scan() {
			return 0, fmt.Errorf("failed to read from %s: %w", uptimePath, scanner.Err())
		}

		line := scanner.Text()
		fields := strings.SplitN(line, ".", 2)
		if len(fields) < 1 {
			return 0, fmt.Errorf("unexpected format in %s: %s", uptimePath, line)
		}

		uptime, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse uptime from %s: %w", uptimePath, err)
		}

		return UptimeInfo(uptime), nil
	}
}
