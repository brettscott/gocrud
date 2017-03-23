package tools

import (
	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogger_withFileAndLine(t *testing.T) {
	fields := logrus.Fields{}
	updatedFields := withFileAndLine(fields)
	assert.NotEmpty(t, updatedFields["file"])
	assert.NotEmpty(t, updatedFields["line"])
}
