# Off Week Rework

On our off week, we decided to look for bugs within our different components, and we decided to rework our Board implementation to better accommodate the networking we anticipate in the future. Similarly, we've set ourselves up better to implement a single struct for each type of action to be applied on a Board (placement, turn [move + build]) so that we have less code duplication.

While it's not here in the commit listing, we worked for 20 hours over the weekend to try and figure out two seperate bugs that, as of yet, remain in our codebase:
1. dynamic loading of player components _should_ work, but we cannot compile for a Linux environment
2. A Player that holds priority infinitely (`for {}`) causes the goroutine it is executing within to remain around for the rest of the program's lifetime

Both of these problems require different solutions, and we attempted to solve both, but failed at both endeavors.

Our attempt to solve (1) was to dockerize our build process, as it would unify the environment within which we were creating the different binaries and plugins, but unfortunately the issues faced with that solution remained the same as those we faced running builds locally (for some reason, some flags being sent to `linker` from the `go build` toolchain were invalid). Given the depth of the issue, and given that Go Plugins are a relatively new feature such that there was practically no information freely available within the community regarding our issue, we decided to leave it as-is, since networked players will require a different implementation regardless.

To address (2), we tried multiple methods of stopping a running Goroutine using different stop channels, goroutine contexts, and runtime function calls, to no avail. Given that our infinite player's interface is that it returns a `struct` for its Placement/Turn, we can't deal with promises and timeouts as we would most likely be able to in other languages. Because the current assignments have had us loading untrusted code into our codebase, and because this is something you are not meant to do within Go, the language does not allow for good handling of this specific error case.

If we continue to face these issues post-networking, we will continue to debug and search for solutions, but as it is now, the time spent to implement a fix is far greater than the time spent implementing networking such that these issues aren't in our codebase.

# Changes (by commit)

Here is a complete list of (relevant) git commits we have included, along with a commit-specific annotation explaining some of the changes in more detail. We've excluded annotations for compilation-, merge-, and comment-only commits, as they don't add much of value to the codebase.

* `commit b5e9c2: Comment out infinite tests, much better referee coverage`

Testing for an infinitely-looping player impacts (for an unknown reason as of yet) other tests within our Referee tests, so they have been removed for the time being, since it’s likely that bugs stemming from Infinite Player accesses will be removed when we implement networking

* `commit bd61be -- Tests should have Ron win, remove observer`

Update the xrun test fourth test to make sure that the second player given consistently wins in a 1v1 against another valid player (though, this is dependant on our implementation of the strategy, and on which strategy each player picks for where they place their workers)

* `commit 13a01d -- messed up tournament signature`

In our Tournament Manager, the signature for starting a set of games WAS (timeout, games). It is not (games, timeout), which is what the xrun harness was expected. Due to this, a set between two valid players within the harness was taking much longer than wanted, since it was playing 10 games instead of 3 (with the intended configuration being best-of-3 on a 10-second timeout)

* `commit 1423d3 -- Copy instead of iteration`

Simplify the `deepCopy` so instead of iterating over players on a Board, just copy them into a new array (which works because a Board only ever has 2 players)

* `commit 22d6bc -- Pull from master, remove Valid tests`

Merge commit, but removed tests on a Valid player within Referee tests on timeouts, since we don’t care if the turn we get is valid in that context, only that we got one back from the Player.

* `commit 456ed5 -- Merge commit, ignored`
* `commit 6e07d0 -- Merge commit, ignored`

* `commit d202cd -- Duplicate player array on deepcopy`

To keep track of player order on a Board, we maintain an array of player names in addition to the map of player names to Workers. However, we initially forgot to include this addition to our deepCopy method, so we were getting failures given that any move or build essentially voided out player names.

* `commit d4c4c6 -- increase code coverage from referee tests`

Tested more methods on referee. Brought code coverage from 13% to 86% of all statements in referee.go. 

* `commit fa5b4f5 -- Remove single comment, ignored`

* `commit 1a480b: [Compile…], fix parsing to use new move syntax`

Given the previous work on board, our interface for Move and Build changed to incorporate player name and worker ID instead of Worker struct

* `commit b5c1b4 -- Compilation commit, ignored`

* `commit 5121a1 -- [Merge…], work on board`

We changed our board implementation to deal with a “fuzzy grid” of cells, instead of a hard-set 6x6 grid (that is, a 2D map from (row, col) to Cell).

* `commit 77b8b6 -- Reworked.md`

Added the markdown file for this assignment to git. 

* `commit a0ffaa9 -- remove prints.`

Removed print statements we had used for debugging purposes. 

* `commit 84f7f7 -- board strings go in board`

Removed a circular dependancy when reformatting error strings. 

* `commit 5fb7df -- TODO comment, ignored`

* `commit d8ee1: added comments, converted to error string constants`

Added comments to method that were missing them. 

