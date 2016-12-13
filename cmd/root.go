// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"

	trello "github.com/VojtechVitek/go-trello"
	"github.com/spf13/cobra"
)

type config struct {
	AppKey string
	Token  string
	Member string

	User   *trello.Member
	Client *trello.Client
}

var Trello = &config{}

var RootCmd = &cobra.Command{
	Use:   "trello-tools",
	Short: "A set of handy tools to work with Trello",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	Trello.AppKey = os.Getenv("TRELLO_APP_KEY")
	Trello.Token = os.Getenv("TRELLO_TOKEN")
	Trello.Member = os.Getenv("TRELLO_MEMBER")

	if len(Trello.AppKey) == 0 || len(Trello.Token) == 0 || len(Trello.Member) == 0 {
		fmt.Printf("Must set TRELLO_APP_KEY, TRELLO_TOKEN and TRELLO_MEMBER environment variables.")
		os.Exit(1)
	}

	client, err := trello.NewAuthClient(Trello.AppKey, &Trello.Token)
	if err != nil {
		log.Fatalf("unable to get trello client: %v", err)
		os.Exit(1)
	}
	Trello.Client = client

	user, err := client.Member(Trello.Member)
	if err != nil {
		log.Fatalf("unable to get member %s: %v", Trello.Member, err)
		os.Exit(1)
	}
	Trello.User = user

	RootCmd.AddCommand(NewCmdRelabelCards("trello-tools", os.Stdout))
}
