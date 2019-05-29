package requestworker

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

type job struct {
	ctx     context.Context
	req     *http.Request
	handler func(resp *http.Response, err error) error

	end chan error
}

type result struct {
	resp *http.Response
	err  error
}

//DefaultMaxIdleConnPerHost max idle
const DefaultMaxIdleConnPerHost = 20

// NewWithClient return new Worker with transport
func NewWithClient(threads int, workerTimeout int, tr *http.Transport) *Worker {

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * time.Duration(workerTimeout),
	}

	w := &Worker{
		jobQuene: make(chan *job),
		threads:  threads,
		tr:       tr,
		client:   client,
	}

	return w
}

//New return http worker
func New(threads int) *Worker {

	tr := &http.Transport{
		Proxy:               NoProxyAllowed,
		MaxIdleConnsPerHost: threads * DefaultMaxIdleConnPerHost,
		// DialContext: (&net.Dialer{
		// 	Timeout:   10 * time.Second,
		// 	KeepAlive: 10 * time.Second,
		// 	DualStack: true,
		// }).DialContext,
		IdleConnTimeout: 300 * time.Second,
		// TLSHandshakeTimeout: 10 * time.Second,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 60,
	}
	w := &Worker{
		jobQuene: make(chan *job),
		threads:  threads,
		tr:       tr,
		client:   client,
	}

	go w.start()
	return w

}

//NoProxyAllowed no proxy
func NoProxyAllowed(request *http.Request) (*url.URL, error) {
	return nil, nil
}

//Worker instance
type Worker struct {
	jobQuene chan *job
	threads  int
	tr       *http.Transport
	client   *http.Client
}

//Execute exec http request
func (w *Worker) Execute(ctx context.Context, req *http.Request, h func(resp *http.Response, err error) error) (err error) {

	j := &job{ctx, req, h, make(chan error)}
	w.jobQuene <- j
	return <-j.end

}

func (w *Worker) run() {
	for j := range w.jobQuene {
		c := make(chan error, 1)
		go func() {
			c <- j.handler(w.client.Do(j.req))
		}()
		select {
		case <-j.ctx.Done():
			w.tr.CancelRequest(j.req)
			j.end <- j.ctx.Err()
		case err := <-c:
			j.end <- err
		}
	}

}

func (w *Worker) start() {

	for i := 0; i < w.threads; i++ {
		go w.run()
	}

}
