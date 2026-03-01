package network

import (
	"fmt"
	"os"
	"server/utils"
	"sort"
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

			n, err := logger.file.WriteString(str)
			if err != nil {
				fmt.Fprintln(os.Stderr, "log write failed:", err)
				return
			}

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

const LOG_FILE = "/var/ohsal/ohsal.log"
func New_logger() *Logger {
	file, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644);
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

func (logger *Logger) rotateLog() {
	// Close current file
	logger.file.Close()

	// Rename old log with timestamp
	timestamp := time.Now().UTC().Format("2006-01-02_15-04-05")
	rotatedName := fmt.Sprintf("ohsal.%s.log", timestamp)

	_ = os.Rename(LOG_FILE, rotatedName)

	// Open new log file
	file, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil { 
		panic(err) 
	}

	logger.file = file
	logger.file_size = 0

	logger.cleanupOldLogs(10);
}

func (logger *Logger) cleanupOldLogs(maxFiles int) {
	entries, err := os.ReadDir(".")
	if err != nil { return }

	var logs []os.DirEntry
	for _, e := range entries {
		if len(e.Name()) > 6 && e.Name()[:6] == "ohsal." && 
		e.Name()[len(e.Name())-4:] == ".log" {
			logs = append(logs, e)
		}
	}

	if len(logs) <= maxFiles { return }

	// Oldest first (name contains timestamp)
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].Name() < logs[j].Name()
	})

	for i := 0; i < len(logs)-maxFiles; i++ {
		_ = os.Remove(logs[i].Name())
	}
}
