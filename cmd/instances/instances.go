// Package instances offers listing all instances for the provided regions.
// This can be used for grepping, etc.
// TODO Add built-in filtering capabilities (Frigga, etc?)
package instances

import (
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/zerocontribution/list_instances/cmd"
	"gopkg.in/alecthomas/kingpin.v2"
)

type command struct {
	cmd.ApplicationOptions
}

// Bind the instances command to the application.
func Bind(app *kingpin.Application, opts *cmd.ApplicationOptions) {
	c := &command{*opts}
	app.Command("instances", "List all instances.").Action(c.run)
}

func (c *command) run(ctx *kingpin.ParseContext) error {
	var insts instances
	instCh := make(chan *instance)
	var wg sync.WaitGroup

	start := time.Now()
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            *aws.NewConfig().WithMaxRetries(c.Retries),
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		logrus.Fatalf("Unable to create AWS session: %v", err)
	}

	// Query EC2 endpoints in each region in parallel
	for _, region := range c.Regions {
		wg.Add(1)
		go func(region string) {
			getInstances(sess, region, instCh)
			wg.Done()
		}(region)
	}
	// Close channel after all EC2 queries have returned
	go func() {
		wg.Wait()
		close(instCh)
	}()
	// Build up a slice of instances as they are returned to us
	for {
		inst, ok := <-instCh
		if !ok {
			break
		}
		insts = append(insts, inst)
	}

	// Sort and print instances that were returned to us
	if len(insts) > 0 {
		insts.Sort()
		insts.PrintTable()
	}

	if c.Debug {
		end := time.Now()
		duration := end.Sub(start)
		logrus.Debugf("Queried %d regions for %d instances in %v", len(c.Regions), len(insts), duration)
	}

	return nil
}

func getInstances(sess *session.Session, region string, instCh chan *instance) {
	ec2Svc := ec2.New(sess, aws.NewConfig().WithRegion(region))

	start := time.Now()
	ec2instances, err := ec2Svc.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running"), aws.String("pending")},
			},
		},
	})
	if err != nil {
		logrus.Errorf("Failed to query EC2 service in region %s: %v", region, err)
	}
	end := time.Now()

	count := 0
	for _, r := range ec2instances.Reservations {
		for _, i := range r.Instances {
			inst := newInstance(i)
			instCh <- &inst
			count++
		}
	}
	duration := end.Sub(start)
	logrus.Debugf("Queried region %s and got %d instances in %v", region, count, duration)
}
