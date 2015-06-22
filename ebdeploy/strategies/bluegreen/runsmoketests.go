package bluegreen

import (
	"net/http"
	"time"

	"github.com/bernos/go-eb-deployer/ebdeploy/pipeline"
	"log"
	"strings"
)

func runSmokeTest(ctx *pipeline.DeploymentContext, next pipeline.Continue) error {
	if len(ctx.Configuration.SmokeTestUrl) > 0 {
		url := strings.Replace(ctx.Configuration.SmokeTestUrl, "{url}", ctx.TargetEnvironment.Url, -1)

		log.Printf("Running smoke test against %s", url)

		for {

			if resp, err := http.Get(url); err == nil && resp.StatusCode == 200 {
				log.Printf("Smoke test passed!")
				break
			}

			time.Sleep(5 * time.Second)
		}

	} else {
		log.Println("No SmokeTestUrl specified. Skipping smoke tests")
	}

	return next()
}

