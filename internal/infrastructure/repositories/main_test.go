package repositories

import (
	"os"
	"testing"

	"github.com/kir1l9x/avito-test-assignment/pkg/validators"
)

func TestMain(m *testing.M) {
	_ = validators.InitDomainValidator()

	os.Exit(m.Run())
}
