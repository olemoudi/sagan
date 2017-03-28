package main

import (
	"flag"
	"fmt"
	//"gopkg.in/libgit2/git2go.v24"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const reposPath = "repos"

var (
	infoLog   *log.Logger
	debugLog  *log.Logger
	debugMode bool
	exiting   chan struct{}
	wg        sync.WaitGroup
	pm        ProjectManager
)

func main() {
	fmt.Println(banner)
	flag.BoolVar(&debugMode, "debug", false, "log additional debug traces")

	flag.Parse()

	LogInit(debugMode)
	//initSignals()
	exiting = make(chan struct{})

	wg.Add(1)
	go func() {
		<-time.After(time.Second * 2)
		wg.Done()
	}()

	pm = ProjectManager{make(chan *Project), make(map[string]*Project)}
	go pm.Run()
	p := makeProject("pin8", "git://github.com/olemoudi/pin8.git")
	wg.Add(1)
	pm.add <- &p
	p = makeProject("sqlmap", "git://github.com/sqlmapproject/sqlmap.git")
	pm.add <- &p
	//<-time.After(time.Second * 100)
	//p = makeProject("oportuno", "git://github.com/olemoudi/oportuno.git")
	wg.Add(1)
	webServer()

	// sync workers
	broadcastExit("regular")
	wg.Wait()

	/*

		//repo, err := git.Clone("git://github.com/sqlmapproject/sqlmap.git", "sqlmap", &git.CloneOptions{})
		repo, err := git.OpenRepository("sqlmap")
		if err != nil {
			panic(err)
		}

		iter, err := repo.NewBranchIterator(git.BranchAll)
		if err != nil {
			panic(err)
		}

		iter.ForEach(printBranch)
	*/

}

func shutDown(msg string) {
	info("Shutting down: ", msg)
	close(exiting)
}

func LogInit(debug_flag bool) {
	logfile, err := os.OpenFile("/tmp/sagan.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file")
	}
	infowriter := io.MultiWriter(logfile, os.Stdout)

	if debug_flag {
		debuglogfile, err := os.OpenFile("/tmp/sagan.debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("Error opening debug log file")
		}

		infowriter = io.MultiWriter(logfile, os.Stdout, debuglogfile)

		debugwriter := io.MultiWriter(debuglogfile, os.Stdout)
		debugLog = log.New(debugwriter, "[DEBUG] ", log.Ldate|log.Ltime)

	} else {
		debugLog = log.New(ioutil.Discard, "", 0)
	}

	infoLog = log.New(infowriter, "", log.Ldate|log.Ltime)

}

func info(msg ...string) {
	s := make([]interface{}, len(msg))
	for i, v := range msg {
		s[i] = v
	}
	infoLog.Println(s...)
	//infoLog.Println(msg)
}

func debug(msg ...string) {
	s := make([]interface{}, len(msg))
	for i, v := range msg {
		s[i] = v
	}
	debugLog.Println(s...)
}

func broadcastExit(msg string) {
	debug("broadcasting exit from", msg)
	close(exiting)
}

func initSignals() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	wg.Add(1)
	go func() {
		select {
		case <-c:
			debug("close exiting because of signal")
			broadcastExit("Interrupt/SIGTERM Signal")
		case <-exiting:
		}
		wg.Done()
		return
	}()
}
