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
	"fmt"
	"os"
	"strconv"
)

type CpuTemperatureInfo float32

func ReadCpuTemperature(hwmonName, tempName string) (interface{}, error) {
	tempPath := fmt.Sprintf("/sys/class/hwmon/%s/%s_input", hwmonName, tempName)

	file, err := os.Open(tempPath)
	if err != nil {
		return 0, fmt.Errorf("failed to open %s: %w", tempPath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return 0, fmt.Errorf("failed to read from %s: %w", tempPath, scanner.Err())
	}

	line := scanner.Text()
	cpuTemperatureMdeg, err := strconv.ParseUint(line, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse cpu temperature from %s: %w", tempPath, err)
	}

	cpuTemperature := float32(cpuTemperatureMdeg) / 1000

	return CpuTemperatureInfo(cpuTemperature), nil
}
