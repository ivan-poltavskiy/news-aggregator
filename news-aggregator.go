package main

import (
	. "NewsAggregator/command-line-client"
	. "NewsAggregator/initialization-data"
	"NewsAggregator/parser"
)

func main() {
	InitializeSource()
	parser.InitializeParserMap()
	cli := NewCommandLineClient()
	cli.Run()
}
