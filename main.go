package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"gitlab.midas.dev/back/river/cmd"
)

func Ask4confirm() bool {
	var s string

	fmt.Printf("Pay a salary (y/N): ")
	_, err := fmt.Scan(&s)
	if err != nil {
		panic(err)
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s == "y" || s == "yes" {
		return true
	}
	return false
}

func readConfig() {
	var err error

	viper.SetConfigFile("main.env")
	viper.SetConfigType("props")

	err = viper.ReadInConfig()

	if err != nil {
		log.Println(err)
		return
	}

	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		log.Println("WARNING: file .env not found")
	} else {
		viper.SetConfigFile(".env")
		viper.SetConfigType("props")
		err = viper.MergeInConfig()
		viper.SetConfigType("props")
		if err != nil {
			log.Println(err)
			return
		}
	}

	// Override config parameters from environment variables if specified
	for _, key := range viper.AllKeys() {
		err = viper.BindEnv(key)

		if err != nil {
			log.Println(err)
		}
	}

}

func main() {
	if Ask4confirm() {
		readConfig()

		cmd.Execute()
	}
	fmt.Println("Goodbye!")
}
