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
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

type ExecInfo string

func readExec(command string) (interface{}, error) {
	args := []string{"sh", "-c", command}
	var stdout, stderr bytes.Buffer

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return ExecInfo(""), err
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return ExecInfo(""), errors.New("returned non-zero exit code")
	}

	outputLines := strings.Split(stdout.String(), "\n")
	if len(outputLines) == 0 {
		return ExecInfo(""), nil
	}

	outputString := outputLines[0]
	return ExecInfo(outputString), nil

}

func ReadExec(command string) func() (interface{}, error) {
	return func() (interface{}, error) {
		return readExec(command)
	}
}
