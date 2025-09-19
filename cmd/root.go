package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"gitlab.midas.dev/back/river/db"
	"gitlab.midas.dev/back/river/internal/client/ethereum"
	"gitlab.midas.dev/back/river/internal/config"
	"gitlab.midas.dev/back/river/internal/handler"
	"gitlab.midas.dev/back/river/internal/service/payment"
	"gitlab.midas.dev/back/river/internal/service/salary"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "./river or ./river repay",
	Short: "Salary tools",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}

		// Open database connection
		dbDriver, err := sql.Open("sqlite3", cfg.DatabasePath)
		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			err = dbDriver.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		// Initialize database schema
		_, err = dbDriver.Exec(db.Schema)
		if err != nil {
			log.Printf("%q: %s\n", err, db.Schema)
			return
		}

		// Initialize repositories
		salaryRepository := db.NewSalaryRepository(dbDriver)

		// Initialize Ethereum client
		ethClient, err := ethclient.Dial(cfg.Node)
		if err != nil {
			log.Fatal(err)
		}
		client := ethereum.NewClient(ethClient)

		// Initialize services
		paymentService := payment.New(client)
		salaryService := salary.New(salaryRepository, paymentService)

		// Initialize handler
		h := handler.New(dbDriver, salaryService, cfg)

		// Check if this is a repay command
		isRepay := false
		for _, v := range args {
			if v == "repay" {
				isRepay = true
			}
		}

		// Execute the command
		err = h.Pay(context.Background(), isRepay)
		if err != nil {
			fmt.Println(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
