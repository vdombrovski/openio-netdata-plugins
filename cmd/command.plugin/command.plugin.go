// OpenIO netdata collectors
// Copyright (C) 2019 OpenIO SAS
//
// This library is free software; you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3.0 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Lesser General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public
// License along with this program. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"flag"
	"log"
	"os"
	"time"

	"oionetdata/collector"
	"oionetdata/command"
	"oionetdata/netdata"
	"oionetdata/util"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("argument required")
	}
	var conf string
	var cmdInterval int64
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.Int64Var(&cmdInterval, "interval", 3600, "Interval between commands in seconds")
	fs.StringVar(&conf, "conf", "/etc/netdata/commands.conf", "Command configuration file")
	fs.Parse(os.Args[2:])
	intervalSeconds := collector.ParseIntervalSeconds(os.Args[1])

	cmds := make(map[string]command.Command)

	out, err := util.Commands(conf)
	if err != nil {
		log.Fatalln("ERROR: Command plugin: Could not load commands", err)
	}

	for name, cmd := range out {
		cmds[name] = command.Command{Cmd: cmd, Desc: "OpenIO command", Family: "command"}
	}

	log.Printf("INFO: Command plugin: Loaded %d commands", len(cmds))

	writer := netdata.NewDefaultWriter()
	worker := netdata.NewWorker(time.Duration(intervalSeconds)*time.Second, writer)
	collector := command.NewCollector(cmds, cmdInterval, worker)
	worker.SetCollector(collector)

	worker.Run()
}
