package main

import (
	"./configuration"
	"./twitter_client"
	"./commands"
	"flag"
)

func main() {

	configurationFile := flag.String("c", "", "Configuration file (default ~/.twitter_cli)")
	count := flag.Int("n", 20, "Count of results")

	flag.Parse()

	config := configuration.ReadConfiguration(*configurationFile)
	client := twitter_client.TwitterClient{Config: config}
	client.DefaultCount = *count

	client.Login()

	commands.PrintUsersAsList(client.GetFollowers("BashfulGeek"))
}