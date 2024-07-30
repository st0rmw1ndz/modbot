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

package main

import (
	"fmt"
	"log"

	"codeberg.org/frosty/sysinfo/lib/readers"
	"codeberg.org/frosty/sysinfo/lib/ui"
)

func main() {
	osInfo, err := readers.ReadOs()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	fmt.Printf("OS: %v %v\n", osInfo.Name, osInfo.Version)

	memoryInfo, err := readers.ReadMemory()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	fmt.Printf("RAM: %v / %v\n", ui.PrettifySize(memoryInfo.Used, 2), ui.PrettifySize(memoryInfo.Total, 2))

	cpuUsageInfo, err := readers.ReadCpuUsage()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	fmt.Printf("CPU: %.*f%%\n", 2, cpuUsageInfo.UsagePercent)
}
