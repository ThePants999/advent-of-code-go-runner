# Advent of Code Runner for Golang
This project provides a simple framework library that makes it easy for you to tackle the Advent of Code in Go while writing as little code as possible that isn't directly solving each day's task. The framework will take care of fetching the input from the AoC website, caching it locally and providing it to your code, as well as timing the execution and outputting the results.

To use it:
* Fetch this module:

`go get github.com/ThePants999/advent-of-code-dotnet-runner`

* Write each day's solution as a pair of functions (part 1 and part 2). Your part 1 function should take an `*slog.Logger` and a string, which is your input, and should return a string (your answer) and any one other parameter, which is some contextual information that will be passed to your part 2 function (e.g. something you calculated in part 1 that will be of use in part 2). Your part 2 function should take the same parameters as part 1 plus your contextual information, and return just your answer as a string.
* Also create a `runner.DayImplementation` for each day that references the part 1 and part 2 functions above. Optionally, it can also encode the example input and part 1/2 answers for that input, as provided on the Advent of Code website, in which case your solution will be tested using that input as well as executed over your real input.
  * It's recommended to put each day in its own file, and you can use [_template.go](blob/main/_template.go) as a starting point.
* Write a `main()` function for your application, which should just call `runner.NewRunner()` - providing a slice of all your `runner.DayImplementation`s - and then call `Run()` on it. Optionally, you can also set up your own `slog.Logger` that controls how logging will work; this logger will be passed to all of your solution functions, as well as used by the framework. But don't worry, you can just pass `nil` if you don't care. An example of what this might look like:

```go
func main() {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelWarn)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: lvl}))

	r := runner.NewRunner(logger, "2024", []runner.DayImplementation{
		Day1,
		Day2,
		Day3,
		Day4})
	r.Run()
}
```

That's it! Now, when you run your program, the framework will handle the behind-the-scenes work. Notably, the first time you run it, it'll ask you to provide your session cookie that your browser is sending to the AoC website so that it can download your inputs on your behalf - that session cookie is stored in a local file called `session`, and the inputs stored in a subdirectory called `inputs/`, so don't forget to add those to your `.gitignore` if you're storing your code in Git.

Provide the `-h` parameter to get help on command-line parameters, but a quick overview:
* By default, it'll attempt to run the current day (i.e. if today's the 10th of the month, it'll attempt to run day 10, assuming that you're doing this during December and keeping up with the puzzles!) You can use `-d <number>` to run a specific day, or `-a` to run all days.
* If you want to get more accurate timings for how quickly your code is executing, use `-s <number>` to run your solution `<number>` times and show the mean/median/fastest/slowest execution.