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
	"time"

	"codeberg.org/frosty/modbot/lib/readers"
)

var (
	delim  = " ][ "
	prefix = "[ "
	suffix = " ]"
)

var modules = []Module{
	{
		Func:     readers.ReadUptime,
		Interval: 1 * time.Second,
	},
	{
		Func:     readers.ReadLoad,
		Interval: 5 * time.Second,
	},
	{
		Func:     readers.ReadMemory,
		Interval: 5 * time.Second,
		Template: "{{.UsedPretty}} / {{.TotalPretty}}",
	},
}
