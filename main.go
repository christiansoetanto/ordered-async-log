package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"ordered-async-log/oalog"
	"os"
	"sync"
	"time"
)

func main() {
	fileName := "./log/access.log"
	err := os.Truncate(fileName, 0)
	if err != nil {
		panic(err)
	}
	err = oalog.New()
	if err != nil {
		panic(err)
	}
	wg := sync.WaitGroup{}
	// i is each HTTP request
	process := 2
	logEachProcess := 1000
	for i := 0; i < process; i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			ctx := context.Background()
			ctx = oalog.InitLogContext(ctx)
			ctx = oalog.SetCtxProcessNo(ctx, x)
			//j is each logv2.Debug
			for j := 0; j < logEachProcess; j++ {
				oalog.Debug(ctx, fmt.Sprint(j))
			}
		}(i)
	}
	wg.Wait()
	fmt.Println("done wait")
	//term := make(chan os.Signal, 1)
	//signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	//select {
	//case <-term:
	//}
	//graceful wait x second
	time.Sleep(9 * time.Second)
	/* ALERT @CHRISTIAN ok here's the problem.
	karena process nya go routine, bisa aja pas Close .size udah 0
	tapi masih ada .append() yang belum dijalanin.

	OR
	issuenya adalah q nya di-pop dulu (size udah 0, chan terima data), baru removePII jalan
	makanya bbrp kali coba selalu selisih 1 log, which is yg terakhir

	bisa aja sih abis terima <-qChan  kasih time.Sleep lagi


	will check behaviour if implemented together with socketmaster
	harusnya last process (or http request) AND the remaining log flush can be done in 10s graceful time

	*/
	oalog.Close()

	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	var logs []Log
	for scanner.Scan() {
		log := Log{}
		err = json.Unmarshal(scanner.Bytes(), &log)
		if err != nil {
			panic(err)
		}
		logs = append(logs, log)
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("The file %s has %d lines\n", fileName, lineCount)
	if lineCount != process*logEachProcess {
		fmt.Printf("ninu ninu lineCount != process*logEachProcess. lineCount = %d, should be %d\n", lineCount, process*logEachProcess)
	}

	realData := make(map[string]string)
	for _, log := range logs {
		pid := log.Metadata.ProcessNo
		realData[pid] = fmt.Sprintf("%s%s", realData[pid], log.Msg)
	}

	expectation := make(map[string]string)
	for i := 0; i < process; i++ {
		s := ""
		for j := 0; j < logEachProcess; j++ {
			s += fmt.Sprint(j)
		}
		expectation[fmt.Sprint(i)] = s
	}

	for i := 0; i < process; i++ {
		idx := fmt.Sprint(i)
		if realData[idx] != expectation[idx] {
			panic(fmt.Sprint(realData[idx], expectation[idx]))
		}
	}
	fmt.Println("done")
}

type Log struct {
	Metadata Metadata `json:"metadata"`
	Msg      string   `json:"msg"`
}

type Metadata struct {
	ProcessNo string `json:"process-no"`
}
