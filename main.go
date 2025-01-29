package main

import (
	"log"
	"time"

	"github.com/Stephen10121/planningcenterbackend/endpoints"
	"github.com/Stephen10121/planningcenterbackend/initializers"
	"github.com/Stephen10121/planningcenterbackend/updaters"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	initializers.SetupEnv()
	base := pocketbase.New()

	base.OnServe().BindFunc(func(e *core.ServeEvent) error {
		endpoints.WebhookTest(e)
		endpoints.TestEndpoint(e)
		endpoints.GetEvents(e, base)

		return e.Next()
	})

	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if updaters.UpdateEvents(base) {
					continue
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	defer close(quit)

	tickerday := time.NewTicker(24 * time.Hour)
	quitday := make(chan struct{})
	go func() {
		for {
			select {
			case <-tickerday.C:
				if updaters.CheckForExpiredTokens(base) {
					continue
				}
			case <-quitday:
				tickerday.Stop()
				return
			}
		}
	}()

	defer close(quitday)

	if err := base.Start(); err != nil {
		log.Fatal(err)
	}
}
