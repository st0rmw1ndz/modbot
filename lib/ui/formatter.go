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

package ui

import "fmt"

const (
	kibibyte = 1
	mebibyte = 1024 * kibibyte
	gibibyte = 1024 * mebibyte
)

func PrettifySize(sizeKibibytes uint64, decimalPlaces uint8) string {
	var size float64
	var unit string

	if sizeKibibytes < mebibyte {
		size = float64(sizeKibibytes)
		unit = "KiB"
	} else if sizeKibibytes < gibibyte {
		size = float64(sizeKibibytes) / float64(mebibyte)
		unit = "MiB"
	} else {
		size = float64(sizeKibibytes) / float64(gibibyte)
		unit = "GiB"
	}

	format := fmt.Sprintf("%%.%df %s", decimalPlaces, unit)
	return fmt.Sprintf(format, size)
}
