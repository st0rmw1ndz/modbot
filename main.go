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

package main

import (
	"fmt"
	"log"

	"codeberg.org/frosty/modbot/lib/readers"
	"codeberg.org/frosty/modbot/lib/ui"
)

const (
	MemoryDecimalPlaces     = 2
	CpuUsageDecimalPlaces   = 2
	CpuTemperatureHwmonName = "hwmon6"
	CpuTemperatureTempName  = "temp1"
	BatteryName             = "BAT1"
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
	fmt.Printf("RAM: %v / %v\n", ui.PrettifyKib(memoryInfo.Used, MemoryDecimalPlaces), ui.PrettifyKib(memoryInfo.Total, MemoryDecimalPlaces))

	cpuUsageInfo, err := readers.ReadCpuUsage()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	cpuTemperatureInfo, err := readers.ReadCpuTemperature(CpuTemperatureHwmonName, CpuTemperatureTempName)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	fmt.Printf("CPU: %.1f°C (%.*f%%)\n", cpuTemperatureInfo, CpuUsageDecimalPlaces, cpuUsageInfo.UsagePercent)

	loadInfo, err := readers.ReadLoad()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	fmt.Printf("Load: %v %v %v\n", loadInfo.OneMinute, loadInfo.FiveMinute, loadInfo.FifteenMinute)

	// batteryInfo, err := readers.ReadBattery(BatteryName)
	// if err != nil {
	// 	log.Fatalf("%v\n", err)
	// }
	// fmt.Printf("Battery: %v%% (%v) w/ %v\n", batteryInfo.Capacity, batteryInfo.Status, batteryInfo.Technology)

	uptimeInfo, err := readers.ReadUptime()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	fmt.Printf("Uptime: %v\n", uptimeInfo)
}
