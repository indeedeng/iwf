package integ

import (
	"flag"
)

var repeatIntegTest = flag.Int("repeat", 1, "the number of repeat")

var repeatInterval = flag.Int("interval", 1, "the interval between test in seconds")

var cadenceIntegTest = flag.Bool("cadence", true, "run integ test against cadence")

var temporalIntegTest = flag.Bool("temporal", true, "run integ test against temporal")
