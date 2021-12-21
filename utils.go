package main

import (
	"github.com/docopt/docopt-go"
)

var usage = `gitlist

Usage:
  gitlist --outputcsv=<file> --orgs=<file>
  gitlist --orgs=<file> --rawjson
  gitlist -h | --help

Options:
  -h --help      	 	 Show this screen.
  --outputcsv=<file> output csv file with processed fields
  --orgs=<file>      path to file with list of github organizations
  --rawjson          dump on stdout raw unprocessed json    
`

///parse flags

func getOrgsFile() string {
	opts, err := docopt.ParseDoc(usage)
	if err != nil {
		panic(err)
	}
	filepath, _ := opts.String("--orgs")
	return filepath
}

func getCSVFile() string {
	opts, err := docopt.ParseDoc(usage)
	if err != nil {
		panic(err)
	}
	filepath, _ := opts.String("--outputcsv")
	return filepath
}

func getRAWJSONFlag() bool {
	opts, err := docopt.ParseDoc(usage)
	if err != nil {
		panic(err)
	}
	rawjson, _ := opts.Bool("--rawjson")
	return rawjson
}
