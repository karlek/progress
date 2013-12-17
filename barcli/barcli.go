// Package barcli implements a cli frontend for progress bars.
package barcli

import (
	"fmt"

	"github.com/0xC3/progress"
	"github.com/mewkiz/pkg/errutil"
)

// Bar represents a cli frontend of a progress.Bar object.
type Bar struct {
	backend  *progress.Bar
	hasRun   bool
	barCount int
	percent  int
}

// New returns a new Bar object initialized from the given parameters and prints
// the bar.
func New(max int) (bar *Bar, err error) {
	bar = new(Bar)
	bar.backend, err = progress.New(max)
	if err != nil {
		return nil, errutil.Err(err)
	}
	return bar, nil
}

// IncMax increases the Max value of bar by add and prints the bar.
func (bar *Bar) IncMax(add int) (err error) {
	err = bar.backend.IncMax(add)
	if err != nil {
		return errutil.Err(err)
	}
	return nil
}

// IncN increases the Cur value of bar by add and prints the bar.
func (bar *Bar) IncN(add int) (err error) {
	err = bar.backend.IncN(add)
	if err != nil {
		return errutil.Err(err)
	}
	return nil
}

// Inc increases the Cur value of bar by one and prints the bar.
func (bar *Bar) Inc() (err error) {
	err = bar.backend.Inc()
	if err != nil {
		return errutil.Err(err)
	}
	return nil
}

// Print prints the current progress bar.
//
// Note: the terminal has to be at least 14 characters in width.
func (bar *Bar) StringSize(col int) (filled string, unfilled string, err error) {
	const prettyFmt = "%s"

	//    '%s' -> ''  = -2
	//    '%%' -> '%' = -1
	const prettySize = len(prettyFmt) - 3
	barSize := col - prettySize
	if barSize < 2 {
		return "", "", errutil.NewNoPosf("terminal too small (%d) to display progressbar.", col)
	}
	part := bar.backend.Progress()
	barCount := int(part * float64(barSize))
	percent := int(part * 100)
	if bar.hasRun == true && barCount == bar.barCount && percent == bar.percent {
		return "", "", nil
	}
	bar.hasRun = true
	bar.barCount = barCount
	bar.percent = percent
	filledBuf := []byte{}
	unfilledBuf := []byte{}
	for i := 0; i < barSize; i++ {
		if i < barCount {
			filledBuf = append(filledBuf, byte('='))
		} else {
			unfilledBuf = append(unfilledBuf, byte('-'))
		}
	}
	filled += fmt.Sprintf(prettyFmt, string(filledBuf))
	if percent == 100 {
		filled += fmt.Sprintln()
	}
	unfilled += string(unfilledBuf)
	return filled, unfilled, nil
}
