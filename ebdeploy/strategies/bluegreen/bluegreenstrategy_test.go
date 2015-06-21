package bluegreen

import (
	"github.com/bernos/go-eb-deployer/ebdeploy"
	"testing"
)

func TestBlueGreenStrategy(t *testing.T) {
	strategy := NewBlueGreenStrategy()
	strategy.Run(new(ebdeploy.DeploymentContext))

}
