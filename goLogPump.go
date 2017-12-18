package main

import (
	"bufio"
	"os"
	"strconv"
	"time"
)

type pump struct {
	from  string
	to    string
	delay int
}

func populatePumpers(osargs []string) []pump {
	var pmps []pump
	j := 0
	for i := 0; i < len(osargs); i += j {
		if (i + 2) < len(osargs) {
			pumpDelay, err := strconv.Atoi(osargs[i+2])
			if err == nil {
				j = 3
			} else {
				j = 2
				pumpDelay = 5000
			}
			pmps = append(pmps, pump{from: osargs[i], to: osargs[i+1], delay: pumpDelay})
		} else {
			j = 2
			pumpDelay := 5000
			pmps = append(pmps, pump{from: osargs[i], to: osargs[i+1], delay: pumpDelay})
		}
	}
	return pmps
}

func pumpLogs(ps pump, done chan bool) {
	for {
		fIn, _ := os.Open(ps.from)
		defer fIn.Close()
		fOut, _ := os.Create(ps.to)
		defer fOut.Close()
		scanner := bufio.NewScanner(fIn)

		for scanner.Scan() {
			//fmt.Println("read log :", scanner.Text())
			//fmt.Fprintf(fOut, scanner.Text()+"\n")
			time.Sleep(time.Duration(ps.delay) * time.Millisecond)
		}
	}
	done <- true
}

func main() {
	pumps := populatePumpers(os.Args[1:]) // Take off 0th element to remove application name form args
	done := make(chan bool)

	for _, p := range pumps {
		go pumpLogs(p, done)
	}

	for i := 0; i < len(pumps); i++ {
		<-done
	}
}
