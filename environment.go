package runner

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

const SESSION_FILENAME string = "session"
const INPUT_DIRNAME string = "inputs"
const AOC_BASE_URL string = "https://adventofcode.com"

type AOCEnvironment struct {
	year       string
	logger     *slog.Logger
	httpClient *http.Client
	inputsDir  string
	inputUrl   string
}

func newAOCEnvironment(year string, baseLogger *slog.Logger) AOCEnvironment {
	if baseLogger == nil {
		baseLogger = slog.Default()
	}
	logger := baseLogger.With(slog.String("year", year))

	thisDir, err := os.Getwd()
	assertOK(err)
	logger.Debug("Found currect directory", slog.String("dir", thisDir))
	sessionFileName := thisDir + string(os.PathSeparator) + SESSION_FILENAME
	inputDirName := thisDir + string(os.PathSeparator) + INPUT_DIRNAME

	// Make sure the inputs directory exists.
	logger.Debug("Checking for input directory", slog.String("dir", inputDirName))
	_, err = os.Stat(inputDirName)
	if err == nil {
		logger.Debug("Input directory exists")
	} else {
		if os.IsNotExist(err) {
			logger.Info("Input directory doesn't exist, creating")
			err = os.Mkdir(inputDirName, 0744)
			assertOK(err)
		} else {
			panic(err)
		}
	}

	// Read the session file.
	var session string
	logger.Debug("Checking for session file", slog.String("file", sessionFileName))
	sessionData, err := os.ReadFile(sessionFileName)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Info("Session file does not exist")
			fmt.Println("In order to download the inputs from the Advent of Code website, this program requires your session cookie.")
			fmt.Println("Please log into the Advent of Code website, then check your browser cookies and enter the value of the 'session' cookie now.")
			fmt.Scanln(&session)
			logger.Info("Session cookie provided", slog.String("cookieValue", session))

			// Write session file
			err = os.WriteFile(sessionFileName, []byte(session), 0644)
			assertOK(err)
			logger.Debug("Session file created")
		} else {
			panic(err)
		}
	} else {
		session = string(sessionData)
		logger.Debug("Session cookie read from file", slog.String("cookieValue", session))
	}

	// Create HTTP client.
	jar, err := cookiejar.New(nil)
	assertOK(err)
	aocUrl, err := url.Parse(AOC_BASE_URL)
	assertOK(err)
	cookie := http.Cookie{
		Name:   "session",
		Value:  session,
		Path:   "/",
		Domain: ".adventofcode.com",
	}
	jar.SetCookies(aocUrl, []*http.Cookie{&cookie})
	client := &http.Client{
		Jar: jar,
	}

	baseUrl := fmt.Sprintf("%s/%s/day/%%d/input", AOC_BASE_URL, year)
	logger.Debug("Constructed base URL for inputs", slog.String("baseUrl", baseUrl))

	return AOCEnvironment{year, logger, client, inputDirName, baseUrl}
}

func assertOK(err error) {
	if err != nil {
		panic(err)
	}
}
