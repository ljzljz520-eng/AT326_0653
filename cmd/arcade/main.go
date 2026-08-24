package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"arcadeledger/internal/service"
	"arcadeledger/internal/transport"
)

func main() {
	path := flag.String("db", "arcade-ledger.db", "bbolt database path")
	command := flag.String("command", "", "single command to execute")
	flag.Parse()
	svc, err := service.OpenService(*path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer closeService(svc)
	if strings.TrimSpace(*command) != "" {
		runLine(svc, *command)
		return
	}
	reader := bufio.NewScanner(os.Stdin)
	for reader.Scan() {
		if strings.TrimSpace(reader.Text()) == "quit" {
			return
		}
		runLine(svc, reader.Text())
	}
}

func runLine(svc *service.ArcadeService, line string) {
	command, err := transport.Parse(line)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	output, err := transport.Execute(svc, command)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(output)
}

func closeService(svc *service.ArcadeService) {
	if svc == nil {
		return
	}
	// The service owns a store through its workflow API; closing is handled by
	// the process boundary for the small command-line entrypoint.
}
