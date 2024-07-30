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

import "fmt"

const (
	Kibibyte = 1
	Mebibyte = 1024 * Kibibyte
	Gibibyte = 1024 * Mebibyte
)

func PrettifySize(sizeKiB uint64, decimalPlaces uint8) string {
	var size float64
	var unit string

	if sizeKiB < Kibibyte {
		size = float64(sizeKiB)
		unit = "B"
	} else if sizeKiB < Mebibyte {
		size = float64(sizeKiB)
		unit = "KiB"
	} else if sizeKiB < Gibibyte {
		size = float64(sizeKiB) / float64(Mebibyte)
		unit = "MiB"
	} else {
		size = float64(sizeKiB) / float64(Gibibyte)
		unit = "GiB"
	}

	format := fmt.Sprintf("%%.%df %s", decimalPlaces, unit)
	return fmt.Sprintf(format, size)
}
