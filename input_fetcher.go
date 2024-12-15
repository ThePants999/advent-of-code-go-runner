package runner

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

func fetchInput(env *AOCEnvironment, logger *slog.Logger, day int) string {
	inputFileName := fmt.Sprintf("%s%s%d", env.inputsDir, string(os.PathSeparator), day)

	logger.Debug("Attempting to read input file", slog.String("fileName", inputFileName))
	var input string
	inputData, err := os.ReadFile(inputFileName)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Info("Input file does not exist")
			url := fmt.Sprintf(env.inputUrl, day)
			logger.Info("Sending request", slog.String("url", url))
			res, err := env.httpClient.Get(url)
			assertOK(err)

			if res.StatusCode != 200 {
				bodyString := "none"
				bodyData, err := io.ReadAll(res.Body)
				if err == nil {
					bodyString = string(bodyData)
				}
				logger.Error(
					"Failed to fetch input data",
					slog.Int("statusCode", res.StatusCode),
					slog.String("statusText", res.Status),
					slog.String("body", bodyString))
				os.Exit(1)
			}

			inputData, err = io.ReadAll(res.Body)
			assertOK(err)
			input = string(inputData)
			logger.Info("Retrieved input", slog.String("input", input))

			// Write input file
			err = os.WriteFile(inputFileName, inputData, 0644)
			assertOK(err)
			logger.Debug("Input file created")
		} else {
			panic(err)
		}
	} else {
		input = string(inputData)
		logger.Debug("Input data read from file", slog.String("input", input))
	}

	return input
}
