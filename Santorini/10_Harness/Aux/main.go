package main

import (
	"os"

	admin "github.com/CS4500-F18/dare-rebr/Santorini/Admin/Referee"
	// iplayer "github.com/CS4500-F18/dare-rebr/Santorini/Common/Player"
	sandbox "github.com/CS4500-F18/dare-rebr/Santorini/Admin/Sandbox"
	obs "github.com/CS4500-F18/dare-rebr/Santorini/Observer"
	player "github.com/CS4500-F18/dare-rebr/Santorini/Player/Client"
)

func main() {
	observer := obs.NewJSONObserver("Fido", os.Stdout)

	name1 := "Uno"
	name2 := "Dos"

	p_1 := sandbox.NewTimeoutPlayer(sandbox.TIMEOUT_DEFAULT, player.ValidPlayer(name1))
	p_2 := sandbox.NewTimeoutPlayer(sandbox.TIMEOUT_DEFAULT, player.ValidPlayer(name2))

	referee_2 := admin.NewReferee(name1, p_1, name2, p_2)

	referee_2.AttachObserver(observer)

	referee_2.Play()
}
