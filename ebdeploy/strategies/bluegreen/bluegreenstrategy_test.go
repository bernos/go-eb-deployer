package bluegreen

import (
	"github.com/bernos/go-eb-deployer/ebdeploy/pipeline"
	"testing"
)

func TestBlueGreenStrategy(t *testing.T) {
	strategy := NewBlueGreenStrategy()
	strategy.Run(new(pipeline.DeploymentContext))

}
