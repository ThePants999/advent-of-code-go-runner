package runner

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/pprof"
	"sort"
	"strconv"
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
const USAGE_TEXT = "Usage: %s [-d <day number> | --day <day number>] [-a | --allDays] [-k | --skipTests] [-t | --testsOnly] [-p | --profiling] [-s <num runs> | --stats <num runs>]\n  -d, --day        Run a specific day (mutually exclusive with -a)\n  -a, --allDays    Run all days sequentially (mutually exclusive with -d)\n  -k, --skipTests  Execute only the real inputs, not the examples (mutually exclusive with -t)\n  -t, --testsOnly  Execute only the examples, not the real inputs (mutually exclusive with -k)\n  -p, --profiling  Run with profiling enabled (output to profile.prof)\n  -s, --stats      Perform multiple runs and calculate statistics\nThe default behaviour if run with neither -a nor -d is to attempt to execute the present day.\n"

func printUsage() {
	fmt.Printf(USAGE_TEXT, os.Args[0])
}

func (runner AOCRunner) Run() {
	var allDays, skipTests, testsOnly, profile bool
	var specificDay, numRuns int
	flag.BoolVar(&skipTests, "k", false, "Skip tests")
	flag.BoolVar(&skipTests, "skipTests", false, "Skip tests")
	flag.BoolVar(&testsOnly, "t", false, "Only run tests")
	flag.BoolVar(&testsOnly, "testsOnly", false, "Only run tests")
	flag.BoolVar(&allDays, "a", false, "Run all days")
	flag.BoolVar(&allDays, "allDays", false, "Run all days")
	flag.BoolVar(&profile, "p", false, "Run with profiling enabled")
	flag.BoolVar(&profile, "profiling", false, "Run with profiling enabled")
	flag.IntVar(&specificDay, "d", 0, "Specify a day number to run")
	flag.IntVar(&specificDay, "day", 0, "Specify a day number to run")
	flag.IntVar(&numRuns, "s", 0, "Calculate statistics over multiple runs")
	flag.IntVar(&numRuns, "stats", 0, "Calculate statistics over multiple runs")
	flag.Usage = printUsage
	flag.Parse()

	if skipTests && testsOnly {
		fmt.Println("The -k and -t arguments are mutually exclusive. Specify one or the other.")
		printUsage()
		os.Exit(1)
	}

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
		var totals runStats
		var maxTime time.Duration
		var axisBuilder, labelBuilder1, labelBuilder2 strings.Builder
		for ix := range runner.days {
			result := runner.runDay(ix+1, numRuns, skipTests, testsOnly)
			times[ix] = result.median
			totals.min += result.min
			totals.max += result.max
			totals.mean += result.mean
			totals.median += result.median
			if result.median > maxTime {
				maxTime = result.median
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
		if !testsOnly {
			if numRuns < 2 {
				fmt.Printf("Total time: %s\n\n", totals.median)
			} else {
				fmt.Printf("Total time: %s median, %s mean, %s min, %s max\n\n", totals.median, totals.mean, totals.min, totals.max)
			}

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
		}
	} else {
		if specificDay == 0 {
			_, _, specificDay = time.Now().Date()
		}
		runner.runDay(specificDay, numRuns, skipTests, testsOnly)
	}
}

type runStats struct {
	min    time.Duration
	max    time.Duration
	median time.Duration
	mean   time.Duration
}

func (runner AOCRunner) runDay(dayNumber int, numRuns int, skipTests bool, testsOnly bool) runStats {
	fmt.Println(DAY_SEPARATOR)
	fmt.Printf("Day %d\n", dayNumber)

	results := make([]dayResult, numRuns)

	nowYear, nowMonth, nowDay := time.Now().Date()
	day := runner.days[dayNumber-1]
	year, _ := strconv.Atoi(runner.env.year)
	if nowYear == year && nowMonth == 12 && nowDay < day.DayNumber {
		runner.env.logger.Info("Skipping day - not ready yet")
		fmt.Printf("Skipping - day not published yet\n")
		return runStats{}
	}

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
	}
	if testsOnly {
		return runStats{}
	}

	fmt.Println("--Real input--")
	if numRuns < 2 {
		res := day.executeDay(runner.env)
		total := res.part1Time + res.part2Time
		fmt.Printf("Part 1: %s (%s)\nPart 2: %s (%s)\nTotal time: %s\n", res.part1Result, res.part1Time, res.part2Result, res.part2Time, total)
		return runStats{total, total, total, total}
	} else {
		var p1Aggregate, p2Aggregate runStats
		var p1Total, p2Total time.Duration
		for i := range numRuns {
			results[i] = day.executeDay(runner.env)
			if results[i].part1Time > p1Aggregate.max {
				p1Aggregate.max = results[i].part1Time
			}
			if p1Aggregate.min == 0 || results[i].part1Time < p1Aggregate.min {
				p1Aggregate.min = results[i].part1Time
			}
			if results[i].part2Time > p2Aggregate.max {
				p2Aggregate.max = results[i].part2Time
			}
			if p2Aggregate.min == 0 || results[i].part2Time < p2Aggregate.min {
				p2Aggregate.min = results[i].part2Time
			}
			p1Total += results[i].part1Time
			p2Total += results[i].part2Time
		}
		p1Aggregate.mean = p1Total / time.Duration(numRuns)
		p2Aggregate.mean = p2Total / time.Duration(numRuns)
		sort.Slice(results, func(i, j int) bool {
			return results[i].part1Result < results[j].part1Result
		})
		p1Aggregate.median = results[numRuns/2].part1Time
		if numRuns%2 == 0 {
			p1Aggregate.median += results[numRuns/2+1].part1Time
			p1Aggregate.median /= 2
		}
		sort.Slice(results, func(i, j int) bool {
			return results[i].part2Result < results[j].part2Result
		})
		p2Aggregate.median = results[numRuns/2].part2Time
		if numRuns%2 == 0 {
			p2Aggregate.median += results[numRuns/2+1].part2Time
			p2Aggregate.median /= 2
		}
		totals := runStats{p1Aggregate.min + p2Aggregate.min, p1Aggregate.max + p2Aggregate.max, p1Aggregate.median + p2Aggregate.median, p1Aggregate.mean + p2Aggregate.mean}
		fmt.Printf("Part 1: %s (median %s, mean %s, min %s, max %s)\nPart 2: %s (median %s, mean %s, min %s, max %s)\nTotal time: median %s, mean %s, min %s, max %s\n", results[0].part1Result, p1Aggregate.median, p1Aggregate.mean, p1Aggregate.min, p1Aggregate.max, results[0].part2Result, p2Aggregate.median, p2Aggregate.mean, p2Aggregate.min, p2Aggregate.max, totals.median, totals.mean, totals.min, totals.max)
		return totals
	}
}
