package controller

import (
	"os"

	"github.com/Reggles44/apperator/cmd/utils"
	"github.com/Reggles44/apperator/pkg/controller"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
)

func NewCommand() *cobra.Command {
	opts := utils.AppOptions{}
	command := cobra.Command{
		Use:   "application-controller",
		Short: "Run Application Controller",
		Long:  "TODO",
		Run: func(cmd *cobra.Command, args []string) {
			scheme := runtime.NewScheme()
			mgr, err := opts.GetManager(scheme)
			if err != nil {
				opts.Log.Error(err, "Failed to create manager", "controller", "application")
				os.Exit(1)
			}

			if err := (&controller.ApplicationReconciler{
				Client: mgr.GetClient(),
				Scheme: mgr.GetScheme(),
			}).SetupWithManager(mgr); err != nil {
				opts.Log.Error(err, "Failed to create controller", "controller", "application")
				os.Exit(1)
			}
		},
	}

	utils.AddCommandFlag(&command, &opts)

	return &command
}
