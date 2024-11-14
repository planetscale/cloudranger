package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/planetscale/cloudranger"
)

type out struct {
	Cloud  string `json:"cloud"`
	Region string `json:"region"`
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	ranger := cloudranger.New()

	ip := os.Args[1]

	info, found := ranger.GetIP(ip)
	if !found {
		fmt.Fprintf(os.Stderr, "IP not in the database: %q\n", ip)
		os.Exit(0)
	}

	o := out{
		Cloud:  info.Cloud(),
		Region: info.Region(),
	}

	if err := json.NewEncoder(os.Stdout).Encode(o); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("Usage: cloudranger <ip>")
}
