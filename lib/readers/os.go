// sysinfo queries information about your system
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
	"os"
	"strings"
)

type OsInfo struct {
	Name       string
	PrettyName string
	Version    string
}

func ReadOs() (OsInfo, error) {
	file, err := os.Open("/lib/os-release")
	if err != nil {
		return OsInfo{}, err
	}
	defer file.Close()

	var osName, osPrettyName, osVersion string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.ReplaceAll(line, "\"", "")
		columns := strings.Split(line, "=")

		switch {
		case strings.HasPrefix(line, "NAME="):
			osName = columns[1]
		case strings.HasPrefix(line, "PRETTY_NAME="):
			osPrettyName = columns[1]
		case strings.HasPrefix(line, "VERSION_ID="):
			osVersion = columns[1]
		}
	}

	return OsInfo{
		Name:       osName,
		PrettyName: osPrettyName,
		Version:    osVersion,
	}, nil
}
