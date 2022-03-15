package emailer

import (
	"fmt"
	"github.com/OpenCal-FYDP/AsyncCalendarOptimizer/internal/storer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	emailer, err := New()
	require.NoError(t, err)
	e := &storer.EventData{
		CalendarEventID: "dsadas",
		Start:           time.Now().Unix(),
		End:             time.Now().Unix(),
		Attendees:       []string{"jspsun@gmail.com"},
		Location:        "aLoc",
		Summary:         "ASum",
	}
	err = emailer.SendConfirmationEmail("jspsun+test@gmail.com", []string{"jspsun@gmail.com", "jspsun+test@gmail.com"}, e)
	assert.NoError(t, err)
}

func TestTime(t *testing.T) {
	ti := time.Unix(time.Now().Unix(), 0)

	fmt.Println(ti.Format(time.Kitchen))
}
