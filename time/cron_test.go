package time

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/robfig/cron/v3"
)

func TestCronExpr(t *testing.T) {
	croj := cron.New(cron.WithSeconds())

	_, err := croj.AddJob("0/1 * * * * ?", NewMyJob())
	assert.NoError(t, err)

	_, err = croj.AddFunc("0/3 * * * * ?", func() {
		log.Printf("long work started!\n")
		time.Sleep(9 * time.Second)
		log.Printf("long work finished!\n")
	})
	assert.NoError(t, err)

	croj.Start()
	time.Sleep(8 * time.Second)

	finishStop := croj.Stop()
	t.Logf("send finish signal!\n")
	<-finishStop.Done()
	t.Logf("all working job are finished!\n")
}

type MyJob struct{}

func NewMyJob() *MyJob {
	return &MyJob{}
}

func (job *MyJob) Run() {
	log.Printf("my job executed!")
}
