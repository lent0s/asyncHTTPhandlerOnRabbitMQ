package logger

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func StartLog(pathFolder, fileName string, maxSizeKB uint16) {

	startLog(pathFolder, fileName, maxSizeKB)
}

func startLog(folder, name string, maxSize uint16) {

	num := getLastLogFile(folder, name, maxSize)
	fullName := fmt.Sprintf("%s%s%s.%03d", folder, string(os.PathSeparator),
		name, num)
	logFile, err := os.OpenFile(fullName,
		os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("access to log file denied")
	}
	log.SetOutput(logFile)
	go seeker(folder, name, maxSize, num)
}

func seeker(folder, name string, maxSize uint16, num int) {

	fullName := fmt.Sprintf("%s%s%s.%03d", folder, string(os.PathSeparator),
		name, num)
	for range time.Tick(30 * time.Second) {
		stat, err := os.Stat(fullName)
		if err != nil || stat.Size() >= int64(maxSize)*1024 {
			break
		}
	}
	startLog(folder, name, maxSize)
}

func getLastLogFile(folder, name string, maxSize uint16) int {

	files, err := os.ReadDir(folder)
	if err != nil {
		log.Fatal(err)
	}

	last := 0
	for _, file := range files {
		dot := strings.LastIndex(file.Name(), ".")
		if file.IsDir() || dot == -1 ||
			file.Name()[:dot] != name ||
			len(file.Name()) != dot+4 {
			continue
		}
		num, err := strconv.Atoi(file.Name()[dot+1:])
		if err != nil {
			continue
		}
		if last < num {
			last = num
		}
	}

	fullName := fmt.Sprintf("%s%s%s.%03d", folder, string(os.PathSeparator),
		name, last)
	stat, err := os.Stat(fullName)
	if err != nil || stat.Size() > int64(maxSize)*1024 {
		return last + 1
	}
	return last
}
