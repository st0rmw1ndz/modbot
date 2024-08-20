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

type BatteryStatus int
type BatteryTechnology int

const (
	StatusUnknown BatteryStatus = iota
	StatusCharging
	StatusDischarging
	StatusNotCharging
	StatusFull
)

const (
	TechnologyUnknown BatteryTechnology = iota
	TechnologyNimh
	TechnologyLiion
	TechnologyLipoly
	TechnologyLife
	TechnologyNicd
	TechnologyLimn
)

type BatteryInfo struct {
	Capacity   uint8
	Status     BatteryStatus
	Technology BatteryTechnology
}

func (bs BatteryStatus) String() string {
	switch bs {
	case StatusCharging:
		return "Charging"
	case StatusDischarging:
		return "Discharging"
	case StatusNotCharging:
		return "Not charging"
	case StatusFull:
		return "Full"
	default:
		return "Unknown"
	}
}

func (bt BatteryTechnology) String() string {
	switch bt {
	case TechnologyNimh:
		return "NiMH"
	case TechnologyLiion:
		return "Li-ion"
	case TechnologyLipoly:
		return "Li-poly"
	case TechnologyLife:
		return "LiFe"
	case TechnologyNicd:
		return "NiCd"
	case TechnologyLimn:
		return "LiMn"
	default:
		return "Unknown"
	}
}

var batteryStatusMap = map[string]BatteryStatus{
	"Unknown":      StatusUnknown,
	"Charging":     StatusCharging,
	"Discharging":  StatusDischarging,
	"Not charging": StatusNotCharging,
	"Full":         StatusFull,
}

var batteryTechnologyMap = map[string]BatteryTechnology{
	"Unknown": TechnologyUnknown,
	"NiMH":    TechnologyNimh,
	"Li-ion":  TechnologyLiion,
	"Li-poly": TechnologyLipoly,
	"LiFe":    TechnologyLife,
	"NiCd":    TechnologyNicd,
	"LiMn":    TechnologyLimn,
}

func BatteryStatusFromStr(statusStr string) BatteryStatus {
	if status, exists := batteryStatusMap[statusStr]; exists {
		return status
	}
	return StatusUnknown
}

func BatteryTechnologyFromStr(technologyStr string) BatteryTechnology {
	if technology, exists := batteryTechnologyMap[technologyStr]; exists {
		return technology
	}
	return TechnologyUnknown
}

func readBattery(batteryName string) (interface{}, error) {
	capacityPath := fmt.Sprintf("/sys/class/power_supply/%s/capacity", batteryName)
	statusPath := fmt.Sprintf("/sys/class/power_supply/%s/status", batteryName)
	technologyPath := fmt.Sprintf("/sys/class/power_supply/%s/technology", batteryName)

	capacityFile, err := os.Open(capacityPath)
	if err != nil {
		return BatteryInfo{}, fmt.Errorf("failed to open %s: %w", capacityPath, err)
	}
	defer capacityFile.Close()

	capacityScanner := bufio.NewScanner(capacityFile)
	if !capacityScanner.Scan() {
		return BatteryInfo{}, fmt.Errorf("failed to read from %s: %w", capacityPath, capacityScanner.Err())
	}

	statusFile, err := os.Open(statusPath)
	if err != nil {
		return BatteryInfo{}, fmt.Errorf("failed to open %s: %w", statusPath, err)
	}
	defer statusFile.Close()

	statusScanner := bufio.NewScanner(statusFile)
	if !statusScanner.Scan() {
		return BatteryInfo{}, fmt.Errorf("failed to read from %s: %w", statusPath, statusScanner.Err())
	}

	technologyFile, err := os.Open(technologyPath)
	if err != nil {
		return BatteryInfo{}, fmt.Errorf("failed to open %s: %w", technologyPath, err)
	}
	defer technologyFile.Close()

	technologyScanner := bufio.NewScanner(technologyFile)
	if !technologyScanner.Scan() {
		return BatteryInfo{}, fmt.Errorf("failed to read from %s: %w", technologyPath, technologyScanner.Err())
	}

	batteryCapacityStr := capacityScanner.Text()
	batteryStatus := statusScanner.Text()
	batteryTechnology := technologyScanner.Text()

	batteryCapacity, err := strconv.ParseUint(batteryCapacityStr, 10, 8)
	if err != nil {
		return BatteryInfo{}, fmt.Errorf("failed to parse capacity from %s: %w", capacityPath, err)
	}

	return BatteryInfo{
		Capacity:   uint8(batteryCapacity),
		Status:     BatteryStatusFromStr(batteryStatus),
		Technology: BatteryTechnologyFromStr(batteryTechnology),
	}, nil
}

func ReadBattery(batteryName string) func() (interface{}, error) {
	return func() (interface{}, error) {
		return readBattery(batteryName)
	}
}
