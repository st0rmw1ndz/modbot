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
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"codeberg.org/frosty/modbot/lib/ui"
)

type MemoryInfo struct {
	Total     uint64
	Available uint64
	Used      uint64

	TotalPretty     string
	AvailablePretty string
	UsedPretty      string
}

func ReadMemory() func() (interface{}, error) {
	return func() (interface{}, error) {
		const memoryFile = "/proc/meminfo"

		file, err := os.Open(memoryFile)
		if err != nil {
			return MemoryInfo{}, err
		}
		defer file.Close()

		var memTotal, memAvailable, memUsed uint64
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()
			columns := strings.Fields(line)

			if len(columns) < 2 {
				continue
			}

			value, err := strconv.ParseUint(columns[1], 10, 64)
			if err != nil {
				return MemoryInfo{}, fmt.Errorf("failed to parse memory value: %w", err)
			}

			switch {
			case strings.HasPrefix(line, "MemTotal:"):
				memTotal = value
			case strings.HasPrefix(line, "MemAvailable:"):
				memAvailable = value
			}
		}

		if err := scanner.Err(); err != nil {
			return MemoryInfo{}, err
		}
		if memTotal == 0 {
			return MemoryInfo{}, errors.New("missing MemTotal")
		}
		if memAvailable == 0 {
			return MemoryInfo{}, errors.New("missing MemAvailable")
		}

		memUsed = memTotal - memAvailable
		return MemoryInfo{
			Total:     memTotal,
			Available: memAvailable,
			Used:      memUsed,

			TotalPretty:     ui.PrettifyKib(memTotal, 1),
			AvailablePretty: ui.PrettifyKib(memAvailable, 1),
			UsedPretty:      ui.PrettifyKib(memUsed, 1),
		}, nil
	}
}
