package show

import (
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/wnc"
)

// Log field names. They are the vocabulary TROUBLESHOOTING.md indexes, so they are
// named here rather than repeated as literals at each call site.
const (
	fieldController = "controller"
	fieldEndpoint   = "endpoint"
	fieldCause      = "cause"
	fieldStatus     = "status"
)

// Reporter collects what went wrong while reading one controller, so the reads finish before
// anything is written: a line per failure interleaved across concurrent controllers would come out
// in a different order every run.
type Reporter struct {
	target config.Target
	faults []fault
	notes  []string
}

type fault struct {
	endpoint string
	err      error
}

// Degraded records a read that failed while leaving the rows intact. The cells it would have
// filled stay unreported, because a zero in their place would be a reading nothing gave.
func (r *Reporter) Degraded(endpoint string, err error) {
	if err == nil {
		return
	}

	r.faults = append(r.faults, fault{endpoint: endpoint, err: err})
}

// Excluded records rows a filter dropped because the leaf it filters on was not reported, which
// is the difference between "no client on that band" and "the band was never reported".
func (r *Reporter) Excluded(n int, reason string) {
	if n <= 0 {
		return
	}

	r.notes = append(r.notes, pluralRows(n)+" excluded: "+reason)
}

// Note records something an operator should see that is not a failure.
func (r *Reporter) Note(text string) {
	if text == "" {
		return
	}

	r.notes = append(r.notes, text)
}

func (r *Reporter) Degradations() int {
	return len(r.faults)
}

// logFatal writes the failure that cost this controller its rows. The fields stay separate from
// the message because logrus quotes a message holding spaces, which would bury the cause token.
func (r *Reporter) logFatal(logger *logrus.Logger, err error) {
	cause, status := wnc.Classify(err)

	logger.WithFields(logrus.Fields{
		fieldController: r.target.Name,
		fieldCause:      string(cause),
		fieldStatus:     status,
	}).Error(wnc.Message(err))
}

// logDegraded writes this controller's non-fatal faults and notes, and reports whether there were
// any faults.
func (r *Reporter) logDegraded(logger *logrus.Logger) bool {
	for _, f := range r.faults {
		cause, status := wnc.Classify(f.err)

		logger.WithFields(logrus.Fields{
			fieldController: r.target.Name,
			fieldEndpoint:   f.endpoint,
			fieldCause:      string(cause),
			fieldStatus:     status,
		}).Error(wnc.Message(f.err))
	}

	for _, n := range r.notes {
		logger.WithField(fieldController, r.target.Name).Warn(n)
	}

	return len(r.faults) > 0
}

func pluralRows(n int) string {
	if n == 1 {
		return "1 row"
	}

	return strconv.Itoa(n) + " rows"
}
