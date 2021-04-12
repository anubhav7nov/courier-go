/*
Example:

```go

package main

import (

	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"***REMOVED***"
	"***REMOVED***/metrics"
)

func main() {
	reg := prometheus.NewRegistry()
	m := metrics.NewPrometheus()
	if err := m.AddToRegistry(reg); err != nil {
		panic(err)
	}

	c, err := courier.NewClient(
			courier.WithUsername("username"),
			courier.WithPassword("password"),
			courier.WithTCPAddress("localhost", 1883),
			courier.WithKeepAlive(10*time.Second),
			courier.WithConnectTimeout(10*time.Second),
			courier.WithWriteTimeout(2*time.Second),
			courier.WithMaxReconnectInterval(5*time.Minute),
			courier.WithGracefulShutdownPeriod(time.Minute),
			courier.WithCustomMetrics(m),
			// courier.WithPersistence(s),   // persistence for qos > 0 use-cases
		)

	if err != nil {
		panic(err)
	}

	ctx, disconnect := context.WithCancel(context.Background())
	errCh := make(chan error)

	metricsServer := http.Server{
		Addr:              ":9090",
		Handler:           promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
	}
	go func() {
		_ = metricsServer.ListenAndServe()
	}()
	go func() {
		errCh <- c.Start(ctx)
	}()

	if !courier.WaitForConnection(c, 5*time.Second, 250*time.Millisecond) {
		os.Exit(1)
	}

	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, []os.Signal{os.Interrupt, syscall.SIGTERM}...)

	go func() {
		tick := time.NewTicker(time.Second)
		for {
			select {
			case t := <-tick.C:
				msg := map[string]interface{}{
					"time": t.UnixNano(),
				}
				if err := c.Publish("topic", courier.QOSOne, false, msg); err != nil {
					fmt.Printf("Publish() error = %s\n", err)
				} else {
					fmt.Println("Publish() success")
				}
			case <-stopCh:
				tick.Stop()
				return
			}
		}
	}()

	<-stopCh

	stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = metricsServer.Shutdown(stopCtx)
	disconnect()

	// wait for result on errCh to allow graceful shutdown
	if err := <-errCh; err != nil {
		panic(err)
	}
}

```
*/
package courier
