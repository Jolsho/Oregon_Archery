package network

import (
	"fmt"
	"os"
	"server/utils"
	"time"
)

type LogEntry struct {
    Time    string `json:"time"`
    Level   string `json:"level"`
    Message string `json:"message"`
}

type Logger struct {
	file 		*os.File
	file_size 	int
	logChan		chan *LogEntry
	logStack 	chan *LogEntry
}

const MAX_LOG_FILE_SIZE = 1 * 1024 * 1024 * 1024

func (logger *Logger) start_writer(group *utils.WorkGroup) {
	group.WG.Add(1)
	for {
		select {
		case entry := <- logger.logChan: {
			str := fmt.Sprintf(
				"%s [%s] %s\n", 
				entry.Time, entry.Level, entry.Message,
			)
			n, _ := logger.file.WriteString(str)
			logger.file_size += n;

			if logger.file_size >= MAX_LOG_FILE_SIZE {
				// TASK_6
			}

			select {
			case logger.logStack <- entry:
			default:
			}
		}

		case <-group.Ctx.Done(): {
			for {
				select {
				case entry := <-logger.logChan:
					str := fmt.Sprintf(
						"%s [%s] %s\n", 
						entry.Time, entry.Level, entry.Message,
					)
					logger.file.WriteString(str)
					logger.logStack <- entry
				default:
					logger.file.Close()
					group.WG.Done();
					return
				}
			}
		}
		}
	}
}

func New_logger(path string) *Logger {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644);
	if err != nil { panic(err.Error()); }

	const LOG_CHAN_SIZE = 256;
	const LOG_STACK_SIZE = 128;

	logger := &Logger {
		file: file,
		logChan: make(chan *LogEntry, LOG_CHAN_SIZE),
		logStack: make(chan *LogEntry, LOG_STACK_SIZE),
	}
	for range LOG_STACK_SIZE { logger.logStack <- &LogEntry{}; }

	return logger;
}


const ERROR_LEVEL = "ERROR";
const WARNING_LEVEL = "WARNING";
const INFO_LEVEL = "INFO";
const DEBUG_LEVEL = "DEBUG";

func (logger *Logger) Log(level, msg string) {

	var log *LogEntry;
	select {
	case log = <-logger.logStack:
	default: log = &LogEntry{}
	}

	log.Time = time.Now().Local().Format("01/02/2006 03:04:05 PM")
	log.Level  = level;
	log.Message = msg;

	logger.logChan <- log;
}
