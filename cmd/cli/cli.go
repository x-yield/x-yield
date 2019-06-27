// inspired by telegraf collector, yandex-tank, pandora and dstat

package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kardianos/service"

	"x_yield/internal"
	"x_yield/internal/core"
)

var fDebug = flag.Bool("debug", false, "turn on debug logging")
var fQuiet = flag.Bool("quiet", false, "run in quiet mode")

var fConfig = flag.String("config", "", "configuration file to load")
var fVersion = flag.Bool("version", false, "display the version and exit")

var fPidfile = flag.String("pidfile", "", "file to write our pid to")
var fUsage = flag.String("usage", "", "print usage for a plugin, ie, 'cli --usage pandora'")

var (
	version string
	commit  string
	branch  string
)

var stop chan struct{}

func ReloadLoop(stop chan struct{}, providers []string, processors []string, listeners []string) {
	reload := make(chan bool, 1)
	reload <- true
	for <-reload {
		reload <- false

		ctx, cancel := context.WithCancel(context.Background())

		signals := make(chan os.Signal)
		signal.Notify(signals, os.Interrupt, syscall.SIGHUP,
			syscall.SIGTERM, syscall.SIGINT)
		go func() {
			select {
			case sig := <-signals:
				if sig == syscall.SIGHUP {
					log.Printf("I! Reloading x-yield config")
					<-reload
					reload <- true
				}
				cancel()
			case <-stop:
				cancel()
			}
		}()

		err := Run(ctx, providers, processors, listeners)
		if err != nil {
			log.Fatalf("E! [telegraf] Error running agent: %v", err)
		}
	}
}

func Run(ctx context.Context, providers []string, processors []string, listeners []string) error {
	log.Print("Starting cli")

	c := core.NewCore()

	c.Providers = providers
	c.Processors = processors
	c.Listeners = listeners
	err := c.LoadConfig(*fConfig)
	if err != nil {
		return err
	}

	if len(c.Inputs) == 0 {
		return errors.New("Error: no inputs found, did you provide a valid config file?")
	}

	ag, err := agent.NewAgent(c)
	if err != nil {
		return err
	}

	// Setup logging as configured.
	logger.SetupLogging(
		ag.Config.Agent.Debug || *fDebug,
		ag.Config.Agent.Quiet || *fQuiet,
		ag.Config.Agent.Logfile,
	)

	if *fTest {
		return ag.Test(ctx)
	}

	log.Printf("I! Loaded inputs: %s", strings.Join(c.InputNames(), " "))
	log.Printf("I! Loaded aggregators: %s", strings.Join(c.AggregatorNames(), " "))
	log.Printf("I! Loaded processors: %s", strings.Join(c.ProcessorNames(), " "))
	log.Printf("I! Loaded outputs: %s", strings.Join(c.OutputNames(), " "))
	log.Printf("I! Tags enabled: %s", c.ListTags())

	if *fPidfile != "" {
		f, err := os.OpenFile(*fPidfile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("E! Unable to create pidfile: %s", err)
		} else {
			fmt.Fprintf(f, "%d\n", os.Getpid())

			f.Close()

			defer func() {
				err := os.Remove(*fPidfile)
				if err != nil {
					log.Printf("E! Unable to remove pidfile: %s", err)
				}
			}()
		}
	}

	return ag.Run(ctx)
}

func UsageExit(rc int) {
	fmt.Println(internal.Usage)
	os.Exit(rc)
}

type program struct {
	inputFilters      []string
	outputFilters     []string
	aggregatorFilters []string
	processorFilters  []string
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}
func (p *program) run() {
	stop = make(chan struct{})
	ReloadLoop(
		stop,
		p.inputFilters,
		p.outputFilters,
		p.aggregatorFilters,
		p.processorFilters,
	)
}
func (p *program) Stop(s service.Service) error {
	close(stop)
	return nil
}

func formatFullVersion() string {
	var parts = []string{"cli"}

	if version != "" {
		parts = append(parts, version)
	} else {
		parts = append(parts, "unknown")
	}

	if branch != "" || commit != "" {
		if branch == "" {
			branch = "unknown"
		}
		if commit == "" {
			commit = "unknown"
		}
		git := fmt.Sprintf("(git: %s %s)", branch, commit)
		parts = append(parts, git)
	}

	return strings.Join(parts, " ")
}
