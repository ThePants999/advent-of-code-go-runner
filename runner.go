package runner

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/pprof"
	"strings"
	"time"

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

const DAY_SEPARATOR = "-----------------------"
const USAGE_TEXT = "Usage: %s [-d <day number> | --day <day number>] [-a | --allDays] [-s | --skipTests] [-p | --profiling]\n  -d, --day        Run a specific day\n  -a, --allDays    Run all days sequentially\n  -s, --skipTests  Execute only the real inputs, not the examples\n  -p, --profiling  Run with profiling enabled (output to profile.prof)\nThe -a and -d arguments are mutually exclusive.\nThe default behaviour if run with no arguments is to attempt to execute the present day.\n"

func printUsage() {
	fmt.Printf(USAGE_TEXT, os.Args[0])
}

func (runner AOCRunner) Run() {
	var allDays, skipTests, profile bool
	var specificDay int
	flag.BoolVar(&skipTests, "s", false, "Skip tests")
	flag.BoolVar(&skipTests, "skipTests", false, "Skip tests")
	flag.BoolVar(&allDays, "a", false, "Run all days")
	flag.BoolVar(&allDays, "allDays", false, "Run all days")
	flag.BoolVar(&profile, "p", false, "Run with profiling enabled")
	flag.BoolVar(&profile, "profiling", false, "Run with profiling enabled")
	flag.IntVar(&specificDay, "d", 0, "Specify a day number to run")
	flag.IntVar(&specificDay, "day", 0, "Specify a day number to run")
	flag.Usage = printUsage
	flag.Parse()

	if allDays && specificDay > 0 {
		fmt.Println("The -a and -d arguments are mutually exclusive. Specify one or the other.")
		printUsage()
		os.Exit(1)
	}

	if profile {
		f, _ := os.Create("profile.prof")
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if allDays {
		times := make([]time.Duration, len(runner.days))
		var totalTime time.Duration
		var maxTime time.Duration
		var axisBuilder, labelBuilder1, labelBuilder2 strings.Builder
		for ix := range runner.days {
			times[ix] = runner.runDay(ix+1, skipTests)
			totalTime += times[ix]
			if times[ix] > maxTime {
				maxTime = times[ix]
			}
			axisBuilder.WriteByte('-')
			tensDigit := byte((ix + 1) / 10)
			unitsDigit := byte((ix + 1) % 10)
			if tensDigit > 0 {
				labelBuilder1.WriteByte(tensDigit + '0')
				labelBuilder2.WriteByte(unitsDigit + '0')
			} else {
				labelBuilder1.WriteByte(unitsDigit + '0')
				labelBuilder2.WriteByte(' ')
			}
		}
		println(DAY_SEPARATOR)
		fmt.Printf("Total time: %s\n\n", totalTime)

		for threshold := 1.0; threshold > 0.01; threshold -= 0.1 {
			fmt.Print("| ")
			for _, time := range times {
				if float64(time)/float64(maxTime) >= threshold {
					fmt.Print("#")
				} else {
					fmt.Print(" ")
				}
			}
			fmt.Println()
		}
		fmt.Print("|-")
		fmt.Println(axisBuilder.String())
		fmt.Print("  ")
		fmt.Println(labelBuilder1.String())
		fmt.Print("  ")
		fmt.Println(labelBuilder2.String())
	} else {
		if specificDay == 0 {
			_, _, specificDay = time.Now().Date()
		}
		runner.runDay(specificDay, skipTests)
	}
}

func (runner AOCRunner) runDay(dayNumber int, skipTests bool) time.Duration {
	day := runner.days[dayNumber-1]

	fmt.Println(DAY_SEPARATOR)
	fmt.Printf("Day %d\n", dayNumber)
	if day.ExampleInput != "" && !skipTests {
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
		fmt.Println("--Real input--")
	}
	res := day.executeDay(runner.env)
	fmt.Printf("Part 1: %s (%s)\nPart 2: %s (%s)\n", res.part1Result, res.part1Time, res.part2Result, res.part2Time)

	return res.part1Time + res.part2Time
}
