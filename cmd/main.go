package main

import revolt "github.com/WelcomerTeam/Revolt/internal"

func main() {

	bot := revolt.NewRevoltBot("tokenHere")

	err := bot.Start()
	if err != nil {
		print(err.Error())
	}

}
