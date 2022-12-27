package integ

import (
	"flag"
)

var repeatIntegTest = flag.Int("repeat", 1, "the number of repeat")

var repeatInterval = flag.Int("intervalMs", 0, "the interval between test in milliseconds")

var cadenceIntegTest = flag.Bool("cadence", true, "run integ test against cadence")

var temporalIntegTest = flag.Bool("temporal", true, "run integ test against temporal")

var testSearchIntegTest = flag.Bool("search", true, "run search integ test against temporal/Cadence")

var searchWaitTimeIntegTest = flag.Int("searchWaitMs", 100, "the amount of time to wait for ElasticSearch being able to search in milliseconds")
