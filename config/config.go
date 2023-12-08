package config

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const fileConfig = "config/config.prop"

type Config struct {
	LogPathFolder    string
	LogFileName      string
	LogMaxFileSizeKB uint16
	RConnect         string
	RTimeout         int
	ServerHost       string
	ServerPort       string
}

func getConfig() *Config {

	return readConfig()
}

func readConfig() *Config {

	conf, err := readConfigFile()
	if err != nil {
		log.Fatalf("Config: %s", err)
	}
	return conf
}

func readConfigFile() (*Config, error) {

	data, err := openConfig()
	if err != nil {
		return nil, err
	}
	records := readDataConfig(data)
	if err = checkFilePath(&records); err != nil {
		return nil, err
	}
	if err = checkHost(records); err != nil {
		return nil, err
	}
	if !checkPortIPIsOK(records["serverPort"]) {
		return nil, fmt.Errorf("incorrect port")
	}
	if !checkFileNameIsOK(records["logFileName"]) {
		return nil, fmt.Errorf(`file name can't include: / \ : * ? " < > |`)
	}

	size, err := strconv.Atoi(records["logMaxFileSizeKB"])
	if err != nil {
		return nil, err
	}
	if 0 > size || size > 2048 {
		size = 2048
		log.Println("max log file size was set up to 2048")
	}

	tOut, err := strconv.Atoi(records["rTimeout"])
	if err != nil {
		return nil, err
	}

	return &Config{
		LogPathFolder:    records["logPathFolder"],
		LogFileName:      records["logFileName"],
		LogMaxFileSizeKB: uint16(size),
		RConnect:         records["rConnect"],
		RTimeout:         tOut,
		ServerHost:       records["serverHost"],
		ServerPort:       records["serverPort"],
	}, nil
}

func openConfig() ([]byte, error) {

	data, err := os.ReadFile(fileConfig)
	if err != nil {
		if newErr := makeConfig(); newErr != nil {
			return nil, newErr
		}
		return nil, fmt.Errorf("%s\n"+
			"[%s] is corrupt and has been replaced\n"+
			"check it and rerun application", err, fileConfig)
	}
	return data, nil
}

func makeConfig() error {

	file, err := os.OpenFile(fileConfig, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Println(err)
		}
	}()

	if _, err = file.Write([]byte(defaultConfig())); err != nil {
		return err
	}

	return nil
}

func defaultConfig() string {

	return `##  logger
logPathFolder:              ./logs
logFileName:                log
logMaxFileSizeKB:           1024

##  rabbitMQ
rConnect:                   amqp://guest:guest@localhost:5672
rTimeout:                   5

##  server
serverHost:                 127.0.0.1
serverPort:                 9000`
}

func readDataConfig(data []byte) map[string]string {

	lines := strings.Split(string(data), "\n")
	records := make(map[string]string)
	for _, line := range lines {
		if line == "\r" || line == "" || line[:1] == "#" {
			continue
		}
		rows := strings.Fields(line)
		if len(rows) != 2 {
			continue
		}
		records[rows[0][:len(rows[0])-1]] = rows[1]
	}
	return records
}

func checkFilePath(records *map[string]string) error {

	var (
		filesPath = []string{
			"logPathFolder"}
	)

	for _, val := range filesPath {
		if (*records)[val] == "" {
			if err := appendConfig(val); err != nil {
				return err
			}
			return fmt.Errorf("not enough data [%s] check \"%s\"\n",
				val, fileConfig)
		}

		(*records)[val] = strings.Join(strings.Split((*records)[val], `/`),
			string(os.PathSeparator))
		(*records)[val] = strings.Join(strings.Split((*records)[val], `\`),
			string(os.PathSeparator))
		if (*records)[val][len((*records)[val])-1] == uint8(os.PathSeparator) {
			(*records)[val] = (*records)[val][:len((*records)[val])-1]
		}

		stat, err := os.Stat((*records)[val])
		if err != nil {
			if strings.Contains(err.Error(), "cannot find") {
				if err = os.Mkdir((*records)[val], 0o660); err != nil {
					return err
				}
				return nil
			}
			return err
		}

		if !stat.IsDir() {
			return fmt.Errorf("wrong path [%s] for [%s]", (*records)[val], val)
		}
		if stat.Mode().Perm()&0o660 != 0o660 {
			return fmt.Errorf("permission [%s] for [%s] denied", (*records)[val], val)
		}
	}
	return nil
}

func appendConfig(s string) error {

	file, err := os.OpenFile(fileConfig, os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if file.Close() != nil {
			log.Println(err)
		}
	}()

	str := strings.Split(defaultConfig(), "\n")
	for _, val := range str {
		if strings.Index(val, s) == 0 {
			_, err = file.Write([]byte("\n" + val))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func checkHost(records map[string]string) error {

	var (
		hosts  = []string{"serverHost"}
		msg, i = "", 0
	)

	for j, val := range hosts {
		i = j
		parts := strings.Split(records[val], ":")
		switch len(parts) {
		case 0: // empty
			msg = "not enough data"
		case 1, 2: // IPv4
			msg = checkIPv4(parts)
		default: // IPv6
			msg = checkIPv6(records[val])
		}
		if msg != "" {
			break
		}
	}
	if msg == "" {
		return nil
	}
	if err := appendConfig(hosts[i]); err != nil {
		return err
	}
	return fmt.Errorf("%s [%s] check \"%s\"\n",
		msg, hosts[i], fileConfig)
}

func checkIPv4(parts []string) string {

	parseIP := net.ParseIP(parts[0])
	if parseIP == nil && strings.ToLower(parts[0]) != "localhost" {
		return "incorrect IPv4"
	}

	if len(parts) == 2 {
		if !checkPortIPIsOK(parts[1]) {
			return "incorrect port"
		}
	}
	return ""
}

func checkIPv6(s string) string {

	if !strings.ContainsAny(s, "[]") {
		parseIP := net.ParseIP(s)
		if parseIP == nil {
			return "incorrect IPv6"
		}
		return ""
	}

	port := s[strings.LastIndex(s, ":")+1:]
	if !checkPortIPIsOK(port) {
		return "incorrect port"
	}

	parseIP := net.ParseIP(s[1:strings.Index(s, "]")])
	if parseIP == nil {
		return "incorrect IPv6"
	}
	return ""
}

func checkPortIPIsOK(s string) bool {

	port, err := strconv.Atoi(s)
	if err != nil || 1<<10 >= port || port >= 1<<16 {
		return false
	}
	return true
}

func checkFileNameIsOK(s string) bool {

	if strings.ContainsAny(s, `/\:*?"<>|`) {
		return false
	}
	return true
}
