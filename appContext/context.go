package appContext

import "time"

type Context struct {
	URL                string
	TotalRequests      int
	Concurrency        int
	SuccessfulRequests int
	Errors             []int
	Quiet              bool
	//StatusCodes        map[int]int
	RequestTimes    []time.Duration
	TotalTime       time.Duration
	LoaderJsContent string
}

// this will be used to be reference as a pointer so everything that i create in my application will be able to share the same info
