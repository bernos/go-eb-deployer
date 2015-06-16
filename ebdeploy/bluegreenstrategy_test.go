package ebdeploy

import (
	"testing"
)

func TestBlueGreenStrategy(t *testing.T) {
	strategy := NewBlueGreenStrategy()
	strategy.Run(new(DeploymentContext))

}
