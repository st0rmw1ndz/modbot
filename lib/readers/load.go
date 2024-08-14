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
	"strings"
)

type LoadInfo struct {
	OneMinute     string
	FiveMinute    string
	FifteenMinute string
}

func (l LoadInfo) String() string {
	return fmt.Sprintf("%v %v %v", l.OneMinute, l.FiveMinute, l.FifteenMinute)
}

func ReadLoad() (interface{}, error) {
	const loadPath = "/proc/loadavg"

	file, err := os.Open(loadPath)
	if err != nil {
		return LoadInfo{}, fmt.Errorf("failed to open %s: %w", loadPath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return LoadInfo{}, fmt.Errorf("failed to read from %s: %w", loadPath, err)
	}

	line := scanner.Text()
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return LoadInfo{}, fmt.Errorf("unexpected format in %s: %s", loadPath, line)
	}

	return LoadInfo{
		OneMinute:     fields[0],
		FiveMinute:    fields[1],
		FifteenMinute: fields[2],
	}, nil
}
