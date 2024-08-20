// modbot is a system information agregator
// Copyright (C) 2024 frosty <inthishouseofcards@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package readers

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type CpuUsageInfo struct {
	InUse        uint64
	Total        uint64
	UsagePercent float64
}

func (cu CpuUsageInfo) String() string {
	return fmt.Sprintf("%d%%", uint(cu.UsagePercent))
}

func readCpuUsage() (interface{}, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return CpuUsageInfo{}, err
	}
	defer file.Close()

	var cpuInUse, cpuTotal uint64
	scanner := bufio.NewScanner(file)

	scanner.Scan()
	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "cpu") {
			break
		}

		var cpuN string
		var cpuUser, cpuNice, cpuSystem, cpuIdle, cpuIoWait, cpuIrq, cpuSoftIrq uint64

		_, err := fmt.Sscanf(line, "cpu%s %d %d %d %d %d %d %d",
			&cpuN, &cpuUser, &cpuNice, &cpuSystem, &cpuIdle, &cpuIoWait, &cpuIrq, &cpuSoftIrq)
		if err != nil {
			return CpuUsageInfo{}, fmt.Errorf("failed to parse CPU stats: %w", err)
		}

		inUse := cpuUser + cpuNice + cpuSystem
		total := inUse + cpuIdle + cpuIoWait + cpuIrq + cpuSoftIrq

		cpuInUse += inUse
		cpuTotal += total
	}

	if err := scanner.Err(); err != nil {
		return CpuUsageInfo{}, err
	}
	if cpuTotal == 0 {
		return CpuUsageInfo{}, errors.New("no CPU stats found")
	}

	return CpuUsageInfo{
		InUse:        cpuInUse,
		Total:        cpuTotal,
		UsagePercent: float64(cpuInUse) * 100 / float64(cpuTotal),
	}, nil
}

func ReadCpuUsage() func() (interface{}, error) {
	return func() (interface{}, error) {
		return readCpuUsage()
	}
}
