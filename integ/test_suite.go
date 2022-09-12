package integTests

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// NOTE: the following definitions can't be defined in *_test.go
// since they need to be exported and used by our internal tests
type (
	BasicIntegrationTestSuite struct {
		// override suite.Suite.Assertions with require.Assertions; this means that s.NotNil(nil) will stop the test,
		// not merely log an error
		*require.Assertions
		IntegrationBase
	}

	IntegrationBase struct {
		suite.Suite
	}
)
