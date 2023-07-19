package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
	"log"
	"net/http"
	httpPprof "net/http/pprof"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strconv"
	"time"
)

func main() {
	envProfileHTTP := os.Getenv("TF_PROVIDER_PROFILING_HTTP")
	if envProfileHTTP == "true" {
		http.HandleFunc("/debug/pprof/", httpPprof.Index)
		http.HandleFunc("/debug/pprof/cmdline", httpPprof.Cmdline)
		http.HandleFunc("/debug/pprof/profile", httpPprof.Profile)
		http.HandleFunc("/debug/pprof/symbol", httpPprof.Symbol)
		http.HandleFunc("/debug/pprof/trace", httpPprof.Trace)

	}

	envProfilingPath := os.Getenv("TF_PROVIDER_PROFILING_PATH")
	if envProfilingPath != "" {
		pFile := "terraform-provider-aws-" + strconv.Itoa(os.Getpid())

		envProfilePeriod := os.Getenv("TF_PROVIDER_PROFILE_PERIOD")
		if envProfilePeriod == "" {
			envProfilePeriod = "1h"
		}
		d, err := time.ParseDuration(envProfilePeriod)
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			i := 0
			ticker := time.NewTicker(d)
			var f *os.File
			var err error
			for range ticker.C {
				if f != nil {
					pprof.StopCPUProfile()
					f.Close()
				}
				f, err = os.Create(filepath.Join(envProfilingPath, fmt.Sprintf("%s-cpu-%d.prof", pFile, i)))
				if err != nil {
					log.Fatal(err)
				}
				i++
				if err := pprof.StartCPUProfile(f); err != nil {
					log.Fatal(err)
				}
			}
		}()

		go func() {
			i := 0
			ticker := time.NewTicker(d)
			for range ticker.C {
				f, err := os.Create(filepath.Join(envProfilingPath, fmt.Sprintf("%s-mem-%d.prof", pFile, i)))
				if err != nil {
					log.Fatal(err)
				}
				i++
				err = pprof.WriteHeapProfile(f)
				f.Close()
				if err != nil {
					log.Fatal(err)
				}
			}
		}()
	}

	debugFlag := flag.Bool("debug", false, "Start provider in debug mode.")
	flag.Parse()

	serverFactory, _, err := provider.ProtoV5ProviderServerFactory(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf5server.ServeOpt

	if *debugFlag {
		serveOpts = append(serveOpts, tf5server.WithManagedDebug())
	}

	logFlags := log.Flags()
	logFlags = logFlags &^ (log.Ldate | log.Ltime)
	log.SetFlags(logFlags)

	err = tf5server.Serve(
		"registry.terraform.io/hashicorp/aws",
		serverFactory,
		serveOpts...,
	)

	if err != nil {
		log.Fatal(err)
	}
}
