# Santorini Test Harness

Project 6's deliverable is the `xboard` test harness meant to handle simple requests against the game board implementation.
The interface for the test harness and valid requests is listed here: [Santorini](http://www.ccs.neu.edu/home/matthias/4500-f18/6.html#%28tech._board%29).

The implementation of Santorini's board and logic can be found in `Santorini/`.

`board-tests/` contains the `.json` files representing sample player inputs (`n-in.json`) and outputs (`n-out.json`) that can be fed into `xboard` to simulate player actions during a game, and the respective outputs.

To run tests against our `xboard` implementation, pipe the `n-in.json` file of your choosing from `board-tests/` into the `xboard` binary, and compare the output to the corresponding `n-out.json` file inside `board-tests/`

To run tests:
```./TESTME```