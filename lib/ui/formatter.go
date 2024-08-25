// modbot is a system information agregator
// Copyright (C) 2024 frosty <inthishouseofcards@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package ui

import "fmt"

const (
	Kibibyte = 1
	Mebibyte = 1024 * Kibibyte
	Gibibyte = 1024 * Mebibyte
	Tebibyte = 1024 * Gibibyte
)

func PrettifyKib(sizeKibibytes uint64, decimalPlaces uint8) string {
	var size float64
	var unit string

	if sizeKibibytes < Mebibyte {
		size = float64(sizeKibibytes)
		unit = "K"
	} else if sizeKibibytes < Gibibyte {
		size = float64(sizeKibibytes) / float64(Mebibyte)
		unit = "M"
	} else if sizeKibibytes < Tebibyte {
		size = float64(sizeKibibytes) / float64(Gibibyte)
		unit = "G"
	} else {
		size = float64(sizeKibibytes) / float64(Tebibyte)
		unit = "T"
	}

	return fmt.Sprintf(fmt.Sprintf("%%.%df%s", decimalPlaces, unit), size)
}
