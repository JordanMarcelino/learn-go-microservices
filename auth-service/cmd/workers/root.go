package workers

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/config"
	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/log"
	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/provider"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/logger"
	"github.com/spf13/cobra"
)

func Start() {
	cfg := config.InitConfig()
	log.SetLogger(logger.NewZapLogger(cfg.Logger.Level))
	provider.BootstrapGlobal(cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	rootCmd := &cobra.Command{}
	cmd := []*cobra.Command{
		{
			Use:   "serve-all",
			Short: "Run all",
			Run: func(cmd *cobra.Command, _ []string) {
				runHttpWorker(cfg, ctx)
			},
		},
	}

	rootCmd.AddCommand(cmd...)
	if err := rootCmd.Execute(); err != nil {
		log.Logger.Fatal(err)
	}
}
