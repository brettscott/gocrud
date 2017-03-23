package tools

import "testing"
import "github.com/DataDog/datadog-go/statsd"

func TestNewStatsD(t *testing.T) {
	config := StatsDConfig{
		isProduction: true,
		host:         "localhost",
		port:         "8080",
	}
	sd, _ := newMMStatsD(config)
	namespace := sd.ddstatsd.Namespace

	if namespace != "app." {
		t.Error("Was expecting the default namespace, but got ", namespace)
	}
}

func TestReturnsDummyStatsDIfConfigIsPoo(t *testing.T) {
	config := StatsDConfig{
		isProduction: true,
		host:         "",
		port:         "",
	}
	_, err := NewStatsD(config)

	if err != nil {
		t.Error("Was not expecting to get an error, yet error came")
	}
}

func TestGlobalTagging(t *testing.T) {
	config := StatsDConfig{
		isProduction: true,
		host:         "localhost",
		port:         "8080",
	}
	sd, _ := newMMStatsD(config)

	tags := sd.ddstatsd.Tags

	if !contains(tags, "env:local") {
		t.Error("Expected a global tag of 'env:local', but it wasn't in ", tags)
	}

	if !contains(tags, "component:a-service-has-no-name") {
		t.Error("Expected a global tag of 'component:a-service-has-no-name', but it wasn't in ", tags)
	}
}

func TestAddGlobalNamespace(t *testing.T) {
	sd := statsd.Client{}
	addGlobalNamespace(&sd)

	if sd.Namespace != "app." {
		t.Error("Was expecting the default namespace, but got ", sd.Namespace)
	}
}

func TestAddGlobalTags(t *testing.T) {
	sd := statsd.Client{}
	addGlobalTags(&sd)

	if !contains(sd.Tags, "env:local") {
		t.Error("Expected a global tag of 'env:local', but it wasn't in ", sd.Tags)
	}

	if !contains(sd.Tags, "component:a-service-has-no-name") {
		t.Error("Expected a global tag of 'component:a-service-has-no-name', but it wasn't in ", sd.Tags)
	}
}

func contains(strings []string, target string) bool {
	for _, s := range strings {
		if s == target {
			return true
		}
	}
	return false
}
