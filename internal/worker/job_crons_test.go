package worker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJobCronsSkipsTicketBookings(t *testing.T) {
	s := &Scheduler{tasks: map[string]*Task{}}

	booking := NewTask("订单表加索引", "halocli", []string{"execute", "-T", "ticket-uuid"})
	booking.UUID = "booking-1"
	booking.NextRun = time.Now().Add(time.Hour)
	s.tasks[booking.UUID] = booking

	job := NewTask("metrics rollup", "halocli", []string{"metrics", "rollup"})
	job.UUID = "job-1"
	job.NextRun = time.Now().Add(time.Minute)
	s.tasks[job.UUID] = job

	crons, err := s.JobCrons()
	require.NoError(t, err)
	require.Len(t, crons, 1)
	require.Equal(t, "job-1", crons[0].UUID)

	require.True(t, booking.IsTicketBooking())
	require.False(t, job.IsTicketBooking())
}
