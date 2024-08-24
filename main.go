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
	"syscall"
	"time"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

var (
	updateChan    = make(chan int)
	moduleOutputs = make([][]byte, len(modules))

	sigChan   = make(chan os.Signal, 1024)
	signalMap = make(map[os.Signal][]*Module)

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
		log.Printf("invalid module index: %d\n", m.pos)
		return
	}

	var output bytes.Buffer

	info, err := m.Func()
	if err != nil {
		output.WriteString("failed")
	} else {
		if m.Template != "" {
			// Parse the output and apply the provided template
			tmpl, err := template.New("module").Parse(m.Template)
			if err != nil {
				log.Printf("template parsing error: %v\n", err)
				return
			}

			if err := tmpl.Execute(&output, info); err != nil {
				log.Printf("template execution error: %v\n", err)
				return
			}
		} else {
			// Use the output as is
			fmt.Fprintf(&output, "%v", info)
		}
	}

	moduleOutputs[m.pos] = output.Bytes()
	updateChan <- 1
}

func (m *Module) Init(pos int) {
	m.pos = pos

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
}

func grabXRootWindow() (*xgb.Conn, xproto.Window, error) {
	conn, err := xgb.NewConn()
	if err != nil {
		return nil, 0, err
	}

	root := xproto.Setup(conn).DefaultScreen(conn).Root
	if conn == nil {
		return nil, 0, fmt.Errorf("failed to create X connection")
	}

	return conn, root, nil
}

func createOutput(b *bytes.Buffer) {
	b.Write(prefix)
	first := true
	for _, output := range moduleOutputs {
		if output == nil {
			continue
		}
		if !first {
			b.Write(delim)
		} else {
			first = false
		}
		b.Write(output)
	}
	b.Write(suffix)
}

func monitorUpdates(setXRootName bool) {
	var lastOutput []byte
	var combinedOutput bytes.Buffer

	for range updateChan {
		combinedOutput.Reset()
		createOutput(&combinedOutput)
		combinedOutputBytes := combinedOutput.Bytes()

		if !bytes.Equal(combinedOutputBytes, lastOutput) {
			if setXRootName {
				// Set X root window name
				xproto.ChangeProperty(x, xproto.PropModeReplace, root, xproto.AtomWmName, xproto.AtomString, 8, uint32(len(combinedOutputBytes)), combinedOutputBytes)
			} else {
				// Send to stdout
				fmt.Printf("%s\n", combinedOutputBytes)
			}
			lastOutput = append([]byte(nil), combinedOutputBytes...)
		}
	}
}

func handleSignal(sig os.Signal) {
	ms := signalMap[sig]
	for _, m := range ms {
		go m.Run()
	}
}

func main() {
	// Parse flags
	var flags Flags
	flag.BoolVar(&flags.SetXRootName, "x", false, "set x root window name")
	flag.Parse()

	// Grab X root window if requested
	if flags.SetXRootName {
		var err error
		x, root, err = grabXRootWindow()
		if err != nil {
			log.Fatalf("error grabbing X root window: %v\n", err)
		}
		defer x.Close()
	}

	// Initialize modules
	for i := range modules {
		go modules[i].Init(i)
	}

	// Monitor changes to the combined output
	go monitorUpdates(flags.SetXRootName)

	// Handle signals sent for modules
	for sig := range sigChan {
		go handleSignal(sig)
	}
}
