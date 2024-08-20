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
// along with this prograb.  If not, see <https://www.gnu.org/licenses/>.

// Credit to gocaudices (https://github.com/LordRusk/gocaudices) for the general outline of how to create the goroutines necessary, and for the X connection code.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

var (
	updateChan    = make(chan int)
	moduleOutputs = make([]string, len(modules))

	mutex      sync.Mutex
	lastOutput string

	// X connection data
	x    *xgb.Conn
	root xproto.Window
)

type Flags struct {
	SetXRootName bool
}

type Module struct {
	Func     func() (interface{}, error)
	Interval time.Duration
	Template string
	Signal   int

	// Internal
	pos int
}

func (m *Module) Run() {
	if m.pos < 0 || m.pos >= len(modules) {
		log.Printf("invalid module index %d\n", m.pos)
		return
	}

	info, err := m.Func()
	if err != nil {
		moduleOutputs[m.pos] = "failed"
		return
	}

	var output string
	if m.Template != "" {
		// Parse the output and apply the provided template
		tmpl, err := template.New("module").Parse(m.Template)
		if err != nil {
			log.Printf("template parsing error: %v\n", err)
			return
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, info); err != nil {
			log.Printf("template execution error: %v\n", err)
			return
		}

		output = buf.String()
	} else {
		// Use the output as is
		output = fmt.Sprintf("%v", info)
	}

	mutex.Lock()
	moduleOutputs[m.pos] = output
	updateChan <- 1
	mutex.Unlock()
}

func parseFlags() Flags {
	var flags Flags
	flag.BoolVar(&flags.SetXRootName, "x", false, "set x root window name")

	flag.Parse()

	return flags
}

func main() {
	flags := parseFlags()

	// Connect to X and get the root window if requested
	if flags.SetXRootName {
		var err error
		x, err = xgb.NewConn()
		if err != nil {
			log.Fatalf("X connection failed: %s\n", err.Error())
		}
		root = xproto.Setup(x).DefaultScreen(x).Root
	}

	sigChan := make(chan os.Signal, 1024)
	signalMap := make(map[os.Signal][]*Module)

	// Initialize modules
	for i := range modules {
		go func(m *Module, i int) {
			m.pos = i

			if m.Signal != 0 {
				sig := syscall.Signal(34 + m.Signal)
				if _, exists := signalMap[sig]; !exists {
					signal.Notify(sigChan, sig)
				}
				signalMap[sig] = append(signalMap[sig], m)
			}

			m.Run()
			if m.Interval > 0 {
				for {
					time.Sleep(m.Interval)
					m.Run()
				}
			}
		}(&modules[i], i)
	}

	// Update output on difference
	go func() {
		for range updateChan {
			mutex.Lock()
			var combinedOutput string
			for i, output := range moduleOutputs {
				if i > 0 {
					combinedOutput += delim
				}
				combinedOutput += output
			}
			combinedOutput = prefix + combinedOutput + suffix
			mutex.Unlock()

			// Output to either X root window name or stdout based on flags
			if combinedOutput != lastOutput {
				if flags.SetXRootName {
					// Set the X root window name
					outputBytes := []byte(combinedOutput)
					xproto.ChangeProperty(x, xproto.PropModeReplace, root, xproto.AtomWmName, xproto.AtomString, 8, uint32(len(outputBytes)), outputBytes)
				} else {
					// Print to stdout
					fmt.Printf("%v\n", combinedOutput)
				}
				lastOutput = combinedOutput
			}
		}
	}()

	// Handle module signals
	for sig := range sigChan {
		go func(sig *os.Signal) {
			ms := signalMap[*sig]
			for _, m := range ms {
				go m.Run()
			}
		}(&sig)
	}
}
