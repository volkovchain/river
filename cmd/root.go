package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.midas.dev/back/river/db"
	"gitlab.midas.dev/back/river/internal/payment"
	"gitlab.midas.dev/back/river/internal/salary"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "./river or ./river repay",
	Short: "Salary tools",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {

		dbDriver, err := sql.Open("sqlite3", "./main.db")

		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			err = dbDriver.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		_, err = dbDriver.Exec(db.Schema)
		if err != nil {
			log.Printf("%q: %s\n", err, db.Schema)
			return
		}

		node := viper.GetString("NODE")
		if node == "" {
			log.Fatal("Environment variable NODE is required")
		}

		privateKeys := viper.GetStringSlice("PRIVATE_KEYS")
		if len(privateKeys) == 0 {
			log.Fatal("Environment variable PRIVATE_KEYS is required")
		}

		paymentService := payment.New(node, privateKeys)
		salaryRepository := db.NewSalaryRepository(dbDriver)
		salaryService := salary.New(salaryRepository, paymentService)

		isrepay := false
		for _, v := range args {
			if v == "repay" {
				isrepay = true
			}
		}

		if isrepay {
			err = salaryService.Repay(context.Background())
		} else {
			err = salaryService.Pay(context.Background())
		}

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
