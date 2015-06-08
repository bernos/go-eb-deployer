package ebdeploy

import (
	"testing"
)

func TestCalculateBucketName(t *testing.T) {
	expected := "my-application-name-packages"
	actual := calculateBucketName("My Application name")

	if expected != actual {
		t.Errorf("Expected %s, but got %s", expected, actual)
	}
}
