package main

import (
	"flag"
	"os"

	"github.com/exiaohao/deploy-test/pkg/tester"
	"github.com/spf13/cobra"
	"k8s.io/apiserver/pkg/util/logs"
)

var (
	opts tester.InitOptions

	rootCmd = &cobra.Command{
		Use:          "tester",
		Short:        "A tester",
		Long:         "A real tester",
		SilenceUsage: true,
	}

	testerCmd = &cobra.Command{
		Use:   "istio",
		Short: "Test istio",
		RunE: func(*cobra.Command, []string) error {
			logs.InitLogs()
			defer logs.FlushLogs()

			if opts.Namespace == "" {
				opts.Namespace = "istio-system"
			}

			server := new(tester.IstioTest)
			server.Initialize(opts)
			server.Run()
			return nil
		},
	}
)

func init() {
	testerCmd.PersistentFlags().StringVar(&opts.KubeConfig, "kubeconfig", "", "Path to kubeconfig file")
	testerCmd.PersistentFlags().StringVar(&opts.Namespace, "namespace", "", "Target namespace")
	testerCmd.PersistentFlags().BoolVar(&opts.ShowDetail, "detail", false, "Show detail")
}

func main() {
	flag.Parse()

	rootCmd.AddCommand(testerCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
