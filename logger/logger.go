package logger

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"
)

const (
	// Using RFC3339Nano as the timestamp format for logs as prometheus scrape interval can be in milliseconds.
	timestampFormat = time.RFC3339Nano
	// Prefix of the name of the log file to store task ranker logs.
	// This will be suffixed with the timestamp, associated with creating the file, to obtain the log filename.
	taskRankerLogFilePrefix = "task_ranker_logs"
	// Prefix of the name of the log file to store task ranking results.
	// This will be suffixed with the timestamp, associated with creating the file, to obtain the log filename.
	taskRankingResultsLogFilePrefix = "task_ranking_results"
	// Giving everyone read and write permissions to the log files.
	logFilePermissions = 0666

	// log types.
	taskRankerLog = iota
	rankingResultsLog
)

var log = logrus.New()

// createTaskRankerLogFile creates the log file to which task ranker logs are persisted.
func createTaskRankerLogFile(now time.Time) (*os.File, error) {
	filename := fmt.Sprintf("%s_%v.log", taskRankerLogFilePrefix, now.UnixNano())
	taskRankerLogFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, logFilePermissions)
	if err != nil {
		err = errors.Wrap(err, "failed to create task ranker operations log file")
	}
	return taskRankerLogFile, err
}

// createTaskRankingResultsLogFile creates the log file to which task ranking results are persisted.
func createTaskRankingResultsLogFile(now time.Time) (*os.File, error) {
	filename := fmt.Sprintf("%s_%v.log", taskRankingResultsLogFilePrefix, now.UnixNano())
	taskRankingResultsLogFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, logFilePermissions)
	if err != nil {
		err = errors.Wrap(err, "failed to create task ranker log file")
	}
	return taskRankingResultsLogFile, err
}

var taskRankerLogFile *os.File
var taskRankingResultsLogFile *os.File

func Configure() error {
	// Disabling log to stdout.
	log.SetOutput(ioutil.Discard)
	// Setting highest log level.
	log.SetLevel(logrus.InfoLevel)

	// Creating the log files.
	now := time.Now()
	var err error

	taskRankerLogFile, err = createTaskRankerLogFile(now)
	if err != nil {
		return err
	}
	taskRankingResultsLogFile, err = createTaskRankingResultsLogFile(now)
	if err != nil {
		return err
	}

	// Instantiate the hooks.
	// Using JSONFormatter to simplify parsing.
	jsonFormatter := &logrus.JSONFormatter{
		DisableHTMLEscape: true,
		TimestampFormat:   timestampFormat,
	}

	textFormatter := &logrus.TextFormatter{
		DisableColors:   true,
		FullTimestamp:   true,
		TimestampFormat: timestampFormat,
	}

	log.AddHook(newWriterHook(textFormatter, taskRankerLogFile, "stage", "query", "query_result"))
	log.AddHook(newWriterHook(jsonFormatter, taskRankingResultsLogFile, "task_ranking_results"))

	return nil
}

func Done() error {
	var err error
	if taskRankerLogFile != nil {
		err = taskRankerLogFile.Close()
		if err != nil {
			err = errors.Wrap(err, "failed to close task ranker log file")
		}
	}

	if taskRankingResultsLogFile != nil {
		err = taskRankingResultsLogFile.Close()
		if err != nil {
			err = errors.Wrap(err, "failed to close tank ranking results log file")
		}
	}
	return err
}

// Aliasing logrus functions.
var WithField = log.WithField
var WithFields = log.WithFields
var Info = log.Info
var Infof = log.Infof
var Error = log.Error
var Errorf = log.Errorf
var Warn = log.Warn
var Warnf = log.Warnf
