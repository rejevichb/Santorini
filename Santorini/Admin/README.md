# Admin
Admin contains all code specific to the Administration of Santorini

Within Admin, the interface used is the `WrappedPlayer` interface, instead of the `IPlayer` interface, meaning that all calls to the "Players" within the course of a Game/Sequence of Games/Tournament are sandboxes beforehand (see )

## Referee
Contains code required to referee a game between 2 players
* `referee.go` -- Code for a Referee component, which handles Game sets
  - `referee_test.go` -- tests on the Referee component

## Tournament
Contains code required to run a tournament between any number of players, alongside configuration information to set up a tournament structure
* `tournament_manager.go` -- tournament manager component

### Config
Code for accepting `IPlayers` into a Tournament, and for wrapping those IPlayers in config-specific `WrappedPlayer` implementations depending on what method of communication is desired for the given Tournament (internal code-loading? TCP? etc.).
* `config.go` -- configuration interface
* `static_config.go` -- configuration code for dynamically-loaded players (currently not dynamically loaded, as we encountered compilation issues when trying to target Linux, so instead we currently switch over the `Kind` provided by the JSON configuration)
* `remote_config.go` -- configuration code for loading remote players over TCP