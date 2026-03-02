package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/andeya/pholcus/common/config"
	"github.com/andeya/pholcus/runtime/status"
)

// Default configuration values from the config file.
const (
	crawlcap int = 50 // Max spider pool capacity
	logcap                int64  = 10000                       // Log buffer capacity
	loglevel              string = "debug"                     // Global log level (also file output level)
	logconsolelevel       string = "info"                      // Console log level
	logfeedbacklevel      string = "error"                     // Client-to-server feedback log level
	loglineinfo           bool   = false                       // Whether to print line info in logs
	logsave               bool   = true                        // Whether to save all logs to local file
	phantomjs             string = WORK_ROOT + "/phantomjs"    // PhantomJS binary path
	proxylib              string = WORK_ROOT + "/proxy.lib"    // Proxy IP file path
	spiderdir             string = WORK_ROOT + "/spiders"      // Dynamic rule directory
	fileoutdir            string = WORK_ROOT + "/file_out"     // Output dir for files (images, HTML, etc.)
	textoutdir            string = WORK_ROOT + "/text_out"     // Output dir for text (excel/csv)
	dbname                string = TAG                         // Database name
	mgoconnstring         string = "127.0.0.1:27017"           // MongoDB connection string
	mgoconncap            int    = 1024                        // MongoDB connection pool size
	mgoconngcsecond       int64  = 600                         // MongoDB connection pool GC interval (seconds)
	mysqlconnstring       string = "root:@tcp(127.0.0.1:3306)" // MySQL connection string
	mysqlconncap          int    = 2048                        // MySQL connection pool size
	mysqlmaxallowedpacket int    = 1048576                     // MySQL max allowed packet (bytes, default 1MB)
	beanstalkHost         string = "localhost:11300"           // Beanstalkd default host (with port)
	beanstalkTube         string = "pholcus"                   // Beanstalkd default tube
	kafkabrokers          string = "127.0.0.1:9092"            // Kafka brokers (comma-separated)

	mode        int    = status.UNSET // Node role
	port        int    = 2015         // Master node port
	master      string = "127.0.0.1"  // Master node address (no port)
	thread      int    = 20           // Global max concurrency
	pause       int64  = 300          // Pause duration reference (ms, random: Pausetime/2 ~ Pausetime*2)
	outtype     string = "csv"        // Output type
	dockercap   int    = 10000        // Segment dump container capacity
	limit       int64  = 0            // Crawl limit; 0=unlimited; custom if rule sets initial LIMIT
	proxyminute int64  = 0            // Proxy IP rotation interval (minutes)
	success     bool   = true         // Inherit success history
	failure     bool   = true         // Inherit failure history
)

var setting = func() config.Configer {
	mustMkdirAll(HISTORY_DIR)
	mustMkdirAll(CACHE_DIR)
	mustMkdirAll(PHANTOMJS_TEMP)

	iniconf, err := config.NewConfig("ini", CONFIG)
	if err != nil {
		file, err := os.Create(CONFIG)
		if err != nil {
			panic(err)
		}
		if err := file.Close(); err != nil {
			log.Printf("[W] close config file: %v", err)
		}
		iniconf, err = config.NewConfig("ini", CONFIG)
		if err != nil {
			panic(err)
		}
		defaultConfig(iniconf)
		iniconf.SaveConfigFile(CONFIG)
	} else {
		trySet(iniconf)
	}

	mustMkdirAll(iniconf.String("spiderdir"))
	mustMkdirAll(iniconf.String("fileoutdir"))
	mustMkdirAll(iniconf.String("textoutdir"))

	return iniconf
}()

func mustMkdirAll(dir string) {
	if err := os.MkdirAll(filepath.Clean(dir), 0777); err != nil {
		log.Fatalf("[F] create directory %q: %v", dir, err)
	}
}

func defaultConfig(iniconf config.Configer) {
	iniconf.Set("crawlcap", strconv.Itoa(crawlcap))
	// iniconf.Set("datachancap", strconv.Itoa(datachancap))
	iniconf.Set("log::cap", strconv.FormatInt(logcap, 10))
	iniconf.Set("log::level", loglevel)
	iniconf.Set("log::consolelevel", logconsolelevel)
	iniconf.Set("log::feedbacklevel", logfeedbacklevel)
	iniconf.Set("log::lineinfo", fmt.Sprint(loglineinfo))
	iniconf.Set("log::save", fmt.Sprint(logsave))
	iniconf.Set("phantomjs", phantomjs)
	iniconf.Set("proxylib", proxylib)
	iniconf.Set("spiderdir", spiderdir)
	iniconf.Set("fileoutdir", fileoutdir)
	iniconf.Set("textoutdir", textoutdir)
	iniconf.Set("dbname", dbname)
	iniconf.Set("mgo::connstring", mgoconnstring)
	iniconf.Set("mgo::conncap", strconv.Itoa(mgoconncap))
	iniconf.Set("mgo::conngcsecond", strconv.FormatInt(mgoconngcsecond, 10))
	iniconf.Set("mysql::connstring", mysqlconnstring)
	iniconf.Set("mysql::conncap", strconv.Itoa(mysqlconncap))
	iniconf.Set("mysql::maxallowedpacket", strconv.Itoa(mysqlmaxallowedpacket))
	iniconf.Set("kafka::brokers", kafkabrokers)
	iniconf.Set("run::mode", strconv.Itoa(mode))
	iniconf.Set("run::port", strconv.Itoa(port))
	iniconf.Set("run::master", master)
	iniconf.Set("run::thread", strconv.Itoa(thread))
	iniconf.Set("run::pause", strconv.FormatInt(pause, 10))
	iniconf.Set("run::outtype", outtype)
	iniconf.Set("run::dockercap", strconv.Itoa(dockercap))
	iniconf.Set("run::limit", strconv.FormatInt(limit, 10))
	iniconf.Set("run::proxyminute", strconv.FormatInt(proxyminute, 10))
	iniconf.Set("run::success", fmt.Sprint(success))
	iniconf.Set("run::failure", fmt.Sprint(failure))
}

func trySet(iniconf config.Configer) {
	if v, e := iniconf.Int("crawlcap"); v <= 0 || e != nil {
		iniconf.Set("crawlcap", strconv.Itoa(crawlcap))
	}

	// if v, e := iniconf.Int("datachancap"); v <= 0 || e != nil {
	// 	iniconf.Set("datachancap", strconv.Itoa(datachancap))
	// }

	if v, e := iniconf.Int64("log::cap"); v <= 0 || e != nil {
		iniconf.Set("log::cap", strconv.FormatInt(logcap, 10))
	}

	level := iniconf.String("log::level")
	if logLevel(level) == -10 {
		level = loglevel
	}
	iniconf.Set("log::level", level)

	consolelevel := iniconf.String("log::consolelevel")
	if logLevel(consolelevel) == -10 {
		consolelevel = logconsolelevel
	}
	iniconf.Set("log::consolelevel", logLevel2(consolelevel, level))

	feedbacklevel := iniconf.String("log::feedbacklevel")
	if logLevel(feedbacklevel) == -10 {
		feedbacklevel = logfeedbacklevel
	}
	iniconf.Set("log::feedbacklevel", logLevel2(feedbacklevel, level))

	if _, e := iniconf.Bool("log::lineinfo"); e != nil {
		iniconf.Set("log::lineinfo", fmt.Sprint(loglineinfo))
	}

	if _, e := iniconf.Bool("log::save"); e != nil {
		iniconf.Set("log::save", fmt.Sprint(logsave))
	}

	if v := iniconf.String("phantomjs"); v == "" {
		iniconf.Set("phantomjs", phantomjs)
	}

	if v := iniconf.String("proxylib"); v == "" {
		iniconf.Set("proxylib", proxylib)
	}

	if v := iniconf.String("spiderdir"); v == "" {
		iniconf.Set("spiderdir", spiderdir)
	}

	if v := iniconf.String("fileoutdir"); v == "" {
		iniconf.Set("fileoutdir", fileoutdir)
	}

	if v := iniconf.String("textoutdir"); v == "" {
		iniconf.Set("textoutdir", textoutdir)
	}

	if v := iniconf.String("dbname"); v == "" {
		iniconf.Set("dbname", dbname)
	}

	if v := iniconf.String("mgo::connstring"); v == "" {
		iniconf.Set("mgo::connstring", mgoconnstring)
	}

	if v, e := iniconf.Int("mgo::conncap"); v <= 0 || e != nil {
		iniconf.Set("mgo::conncap", strconv.Itoa(mgoconncap))
	}

	if v, e := iniconf.Int64("mgo::conngcsecond"); v <= 0 || e != nil {
		iniconf.Set("mgo::conngcsecond", strconv.FormatInt(mgoconngcsecond, 10))
	}

	if v := iniconf.String("mysql::connstring"); v == "" {
		iniconf.Set("mysql::connstring", mysqlconnstring)
	}

	if v, e := iniconf.Int("mysql::conncap"); v <= 0 || e != nil {
		iniconf.Set("mysql::conncap", strconv.Itoa(mysqlconncap))
	}

	if v, e := iniconf.Int("mysql::maxallowedpacket"); v <= 0 || e != nil {
		iniconf.Set("mysql::maxallowedpacket", strconv.Itoa(mysqlmaxallowedpacket))
	}

	if v := iniconf.String("kafka::brokers"); v == "" {
		iniconf.Set("kafka::brokers", kafkabrokers)
	}

	if v, e := iniconf.Int("run::mode"); v < status.UNSET || v > status.CLIENT || e != nil {
		iniconf.Set("run::mode", strconv.Itoa(mode))
	}

	if v, e := iniconf.Int("run::port"); v <= 0 || e != nil {
		iniconf.Set("run::port", strconv.Itoa(port))
	}

	if v := iniconf.String("run::master"); v == "" {
		iniconf.Set("run::master", master)
	}

	if v, e := iniconf.Int("run::thread"); v <= 0 || e != nil {
		iniconf.Set("run::thread", strconv.Itoa(thread))
	}

	if v, e := iniconf.Int64("run::pause"); v < 0 || e != nil {
		iniconf.Set("run::pause", strconv.FormatInt(pause, 10))
	}

	if v := iniconf.String("run::outtype"); v == "" {
		iniconf.Set("run::outtype", outtype)
	}

	if v, e := iniconf.Int("run::dockercap"); v <= 0 || e != nil {
		iniconf.Set("run::dockercap", strconv.Itoa(dockercap))
	}

	if v, e := iniconf.Int64("run::limit"); v < 0 || e != nil {
		iniconf.Set("run::limit", strconv.FormatInt(limit, 10))
	}

	if v, e := iniconf.Int64("run::proxyminute"); v <= 0 || e != nil {
		iniconf.Set("run::proxyminute", strconv.FormatInt(proxyminute, 10))
	}

	if _, e := iniconf.Bool("run::success"); e != nil {
		iniconf.Set("run::success", fmt.Sprint(success))
	}

	if _, e := iniconf.Bool("run::failure"); e != nil {
		iniconf.Set("run::failure", fmt.Sprint(failure))
	}

	iniconf.SaveConfigFile(CONFIG)
}

func logLevel2(l string, g string) string {
	a, b := logLevel(l), logLevel(g)
	if a < b {
		return l
	}
	return g
}
