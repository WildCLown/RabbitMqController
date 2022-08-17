package monitor

import (
	"bufio"
	"fmt"
	"net/http"
	"rabbitmq/shared"
	"strconv"
	"strings"
	"time"
)

type Monitor struct {
	MonitorInterval time.Duration
}

func NewMonitor(t time.Duration) Monitor {
	return Monitor{MonitorInterval: t}
}

func (m Monitor) Monitoring(chStart, chStop chan bool) {
	for {
		chStart <- true
		time.Sleep(m.MonitorInterval) // seconds
		chStop <- true
	}
}

func (m Monitor) GetDeliveryRate() float64 {
	resp, err := http.Get("http://nsr:nsr@localhost:15672/api/vhosts")
	shared.FailOnError(err, err.Error())

	defer resp.Body.Close()

	// Print the HTTP response status.
	fmt.Println("Response status:", resp.Status)

	// Print the first 5 lines of the response body.
	scanner := bufio.NewScanner(resp.Body)

	l := ""
	for i := 0; scanner.Scan() && i < 1; i++ {
		l = scanner.Text()
	}

	sb := "\"deliver_details\":{\"rate\":"
	se := "},\"deliver_get\""
	idx1 := strings.Index(l, sb) + len(sb)
	idx2 := strings.Index(l, se)
	s1 := l[idx1:idx2]

	r, err := strconv.ParseFloat(s1, 64)
	if err != nil {
		shared.FailOnError(err, "Something wrong with getting the Delivery rate")
	}

	return r
}

func (m Monitor) GetSomething() float64 {
	resp, err := http.Get("http://nsr:nsr@localhost:5672/api/vhosts")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Print the HTTP response status.
	//fmt.Println("Response status:", resp.Status)

	// Print the first 5 lines of the response body.
	scanner := bufio.NewScanner(resp.Body)

	for i := 0; scanner.Scan() && i < 1; i++ {
		fmt.Println(scanner.Text())
	}
	return 0.0
}
