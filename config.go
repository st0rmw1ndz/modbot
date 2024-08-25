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

package main

import (
	"time"

	"codeberg.org/frosty/modbot/lib/readers"
)

var (
	delim  = []byte("] [")
	prefix = []byte("[")
	suffix = []byte("]")
)

var modules = []Module{
	{
		Func:     readers.ReadCpuUsage(),
		Interval: 5 * time.Second,
		Template: `CPU {{printf "%.0f" .UsagePercent}}%`,
	},
	{
		Func:     readers.ReadBattery("BAT1"),
		Interval: 60 * time.Second,
		Template: "BAT {{.Capacity}}%",
	},
	{
		Func:     readers.ReadDate("15:04:05"),
		Interval: 1 * time.Second,
	},
	{
		Func:     readers.ReadLoad(),
		Interval: 5 * time.Second,
		Template: "{{.OneMinute}}",
	},
}
