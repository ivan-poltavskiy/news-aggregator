package main

import . "NewsAggregator/command-line-client"

func main() {
	cli := NewCommandLineClient()
	cli.Run()
}
