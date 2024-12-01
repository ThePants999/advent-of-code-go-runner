package runner

import (
	"fmt"
	"log/slog"

	"github.com/logrusorgru/aurora/v4"
)

type AOCRunner struct {
	env  *AOCEnvironment
	days []DayImplementation
}

func NewRunner(logger *slog.Logger, year string, days []DayImplementation) AOCRunner {
	env := newAOCEnvironment(year, logger)
	return AOCRunner{&env, days}
}

func (runner AOCRunner) Run() {
	/*lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelDebug)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: lvl}))
	env := newAOCEnvironment("2024", logger)*/

	// Day 1 for now. We'll handle more days tomorrow!
	runner.runDay(1)
}

func (runner AOCRunner) runDay(dayNumber int) {
	day := runner.days[dayNumber-1]

	fmt.Println("-----------------------")
	fmt.Printf("Day %d\n", dayNumber)
	if day.ExampleInput != "" {
		fmt.Print("--Example input--\nPart 1: ")
		res := day.testDay(runner.env)
		if res.part1Correct {
			fmt.Println(aurora.Green("CORRECT"))
		} else {
			fmt.Println(aurora.Red("INCORRECT"))
			fmt.Println("Expected: ", day.ExamplePart1Answer)
			fmt.Println("Received: ", res.result.part1Result)
		}
		fmt.Printf("Part 2: ")
		if res.part2Correct {
			fmt.Println(aurora.Green("CORRECT"))
		} else {
			fmt.Println(aurora.Red("INCORRECT"))
			fmt.Println("Expected: ", day.ExamplePart2Answer)
			fmt.Println("Received: ", res.result.part2Result)
		}
	}
	fmt.Println("--Real input--")
	res := day.executeDay(runner.env)
	fmt.Printf("Part 1: %s (%s)\nPart 2: %s (%s)\n", res.part1Result, res.part1Time, res.part2Result, res.part2Time)
}
