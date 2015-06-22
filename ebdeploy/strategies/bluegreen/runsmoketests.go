package bluegreen

import (
	"fmt"
	"net/http"
	"time"

	"log"
	"strings"

	"github.com/bernos/go-eb-deployer/ebdeploy/pipeline"
)

func runSmokeTest(ctx *pipeline.DeploymentContext, next pipeline.Continue) error {
	if len(ctx.Configuration.SmokeTestUrl) > 0 {
		url := strings.Replace(ctx.Configuration.SmokeTestUrl, "{url}", ctx.TargetEnvironment.Url, -1)
		t := ctx.Configuration.SmokeTestTimeout

		if t == 0 {
			t = 30
		}

		timeout := time.Now().Add(time.Second * t)
		log.Printf("Running smoke test against %s", url)

		for {

			if resp, err := http.Get(url); err == nil && resp.StatusCode == 200 {
				log.Printf("Smoke test passed!")
				break
			}

			time.Sleep(5 * time.Second)

			if time.Now().After(timeout) {
				return fmt.Errorf("Smoke test timed out after %d seconds", t)
			}
		}

	} else {
		log.Println("No SmokeTestUrl specified. Skipping smoke tests")
	}

	return next()
}
