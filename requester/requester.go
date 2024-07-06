package requester

import (
	"fmt"
	"github.com/JonecoBoy/stress-test/appContext"
	"net/http"
	"sync"
)

type Requester struct {
	Context *appContext.Context
}

func NewRequester(ctx *appContext.Context) *Requester {
	return &Requester{
		Context: ctx,
	}
}

func (r *Requester) DoRequest(ctx *appContext.Context) {
	resp, err := http.Get(r.Context.URL)
	if err != nil {
		if ctx.Quiet == false {
			fmt.Printf("Status Code: %d - Error: %s", resp.StatusCode, err)
		}
		//return
	}
	defer resp.Body.Close()

	//r.Context.StatusCodes[resp.StatusCode]++
	//if resp.StatusCode >= 200 && resp.StatusCode < 300 {
	//	r.Context.SuccessfulRequests++
	//}
	if resp.StatusCode != http.StatusOK {
		ctx.Errors = append(ctx.Errors, resp.StatusCode)
	} else {
		r.Context.SuccessfulRequests++
	}
	if ctx.Quiet == false {
		fmt.Println("Response status:", resp.Status)
	}
}

func (r *Requester) Start(ctx *appContext.Context) {
	var wg sync.WaitGroup

	for i := 0; i < r.Context.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < r.Context.TotalRequests/r.Context.Concurrency; j++ {
				r.DoRequest(ctx)
			}
		}()
	}

	wg.Wait()
}
