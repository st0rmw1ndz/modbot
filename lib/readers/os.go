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

type OsInfo struct {
	Name       string
	PrettyName string
	Version    string
}

func readOs() (interface{}, error) {
	const osPath = "/lib/os-release"

	file, err := os.Open(osPath)
	if err != nil {
		return OsInfo{}, fmt.Errorf("failed to open %s: %w", osPath, err)
	}
	defer file.Close()

	var osName, osPrettyName, osVersion string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.ReplaceAll(line, "\"", "")
		fields := strings.SplitN(line, "=", 2)

		if len(fields) < 2 {
			return OsInfo{}, fmt.Errorf("unexpected format in %s: %s", osPath, line)
		}

		switch {
		case strings.HasPrefix(line, "NAME"):
			osName = fields[1]
		case strings.HasPrefix(line, "PRETTY_NAME"):
			osPrettyName = fields[1]
		case strings.HasPrefix(line, "VERSION_ID"):
			osVersion = fields[1]
		}
	}

	if err := scanner.Err(); err != nil {
		return OsInfo{}, fmt.Errorf("failed to read from %s: %w", osPath, err)
	}

	return OsInfo{
		Name:       osName,
		PrettyName: osPrettyName,
		Version:    osVersion,
	}, nil
}

func ReadOs() func() (interface{}, error) {
	return func() (interface{}, error) {
		return readOs()
	}
}
