package models

import (
	"testing"

	"github.com/kir1l9x/avito-test-assignment/pkg/validators"
)

func TestMain(m *testing.M) {
	_ = validators.InitDomainValidator()

	m.Run()
}
