/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"pesu-cli/cli/config"
	"strings"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Username: ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)

		fmt.Print("Password: ")
		// Note: basic read for now. In production use terminal/password for masking
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)

		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		loginData := map[string]string{
			"username": username,
			"password": password,
		}
		jsonData, _ := json.Marshal(loginData)

		resp, err := http.Post(cfg.ApiURL+"/auth/login", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error connecting to server:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Login failed! Check credentials.")
			return
		}

		var result map[string]string
		json.NewDecoder(resp.Body).Decode(&result)

		cfg.Token = result["token"]
		if err := config.SaveConfig(cfg); err != nil {
			fmt.Println("Error saving config:", err)
			return
		}

		fmt.Println("Login successful!")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
