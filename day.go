package runner

import (
	"log/slog"
	"time"
)

type dayResult struct {
	part1Result string
	part1Time   time.Duration
	part2Result string
	part2Time   time.Duration
}

type testResult struct {
	result       dayResult
	part1Correct bool
	part2Correct bool
}

type DayImplementation struct {
	DayNumber          int
	ExecutePart1       func(*slog.Logger, string) (string, any)
	ExecutePart2       func(*slog.Logger, string, any) string
	ExampleInput       string
	ExamplePart1Answer string
	ExamplePart2Answer string
}

func (day *DayImplementation) testDay(env *AOCEnvironment) testResult {
	logger := env.logger.With(slog.Int("day", day.DayNumber))

	if day.ExampleInput != "" {
		logger.Debug("Example input provided")
		result := day.runDayWithInput(logger, day.ExampleInput)
		part1Correct := result.part1Result == day.ExamplePart1Answer
		part2Correct := result.part2Result == day.ExamplePart2Answer
		logger.Info("Tests completed", slog.Bool("part1Correct", part1Correct), slog.Bool("part2Correct", part2Correct))
		if !part1Correct {
			logger.Warn("Incorrect answer for part 1!", slog.String("expected", day.ExamplePart1Answer), slog.String("received", result.part1Result))
		}
		if !part2Correct {
			logger.Warn("Incorrect answer for part 2!", slog.String("expected", day.ExamplePart2Answer), slog.String("received", result.part2Result))
		}
		return testResult{result, part1Correct, part2Correct}
	} else {
		logger.Info("Skipping tests as no example input provided")
		return testResult{}
	}
}

func (day *DayImplementation) executeDay(env *AOCEnvironment) dayResult {
	logger := env.logger.With(slog.Int("day", day.DayNumber))
	input := fetchInput(env, logger, day.DayNumber)
	result := day.runDayWithInput(logger, input)
	return result
}

func (day *DayImplementation) runDayWithInput(logger *slog.Logger, input string) dayResult {
	start := time.Now()
	logger.Debug("Starting part 1", slog.Time("startTime", start))
	part1Result, part1Context := day.ExecutePart1(logger, input)
	part1Done := time.Now()
	part1Time := part1Done.Sub(start)
	part2Result := day.ExecutePart2(logger, input, part1Context)
	part2Done := time.Now()
	part2Time := part2Done.Sub(part1Done)
	logger.Info("Part 1 results", slog.String("result", part1Result), slog.Time("endTime", part1Done), slog.Duration("duration", part1Time))
	logger.Info("Part 2 results", slog.String("result", part2Result), slog.Time("endTime", part2Done), slog.Duration("duration", part2Time))
	return dayResult{part1Result, part1Time, part2Result, part2Time}
}
