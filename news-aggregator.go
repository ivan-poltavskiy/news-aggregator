package main

import (
	. "NewsAggregator/command-line-client"
	. "NewsAggregator/initialization-data"
)

func main() {
	InitializeSource()
	cli := NewCommandLineClient()
	cli.Run()
}
