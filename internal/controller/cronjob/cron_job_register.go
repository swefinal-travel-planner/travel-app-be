package cronjob

import (
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
)

type CronJobRegister struct {
	tripService service.TripService
	cron        *cron.Cron
}

func NewCronJobRegister(tripService service.TripService) *CronJobRegister {
	return &CronJobRegister{
		tripService: tripService,
		cron:        cron.New(),
	}
}

func (c *CronJobRegister) RegisterJobs() {
	c.cron.AddFunc("0 0 * * *", func() {
		ctx := &gin.Context{}
		err := c.tripService.UpdateStatusTripStart(ctx)
		if err != nil {
			log.Error("CronJobRegister.RegisterJobs - UpdateStatusTripStart Error: " + err.Error())
		}
	})
	c.cron.AddFunc("0 0 * * *", func() {
		ctx := &gin.Context{}
		err := c.tripService.UpdateStatusTripEnd(ctx)
		if err != nil {
			log.Error("CronJobRegister.RegisterJobs - UpdateStatusTripEnd Error: " + err.Error())
		}
	})
	c.cron.AddFunc("0 0 * * *", func() {
		ctx := &gin.Context{}
		err := c.tripService.SendTripStartReminders(ctx)
		if err != nil {
			log.Error("CronJobRegister.RegisterJobs - SendTripStartReminders Error: " + err.Error())
		}
	})
}

func (c *CronJobRegister) Start() {
	c.RegisterJobs()
	c.cron.Start()
}
