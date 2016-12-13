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
	"io"
	"log"
	"net/url"

	trello "github.com/VojtechVitek/go-trello"
	"github.com/mfojtik/trello-tools/pkg/util"
	"github.com/spf13/cobra"
)

type RelabelOptions struct {
	Board  string
	List   string
	From   string
	To     string
	Remove bool
}

func NewCmdRelabelCards(fullName string, out io.Writer) *cobra.Command {
	opts := &RelabelOptions{}
	cmd := &cobra.Command{
		Use:   "relabel --board BOARD --list LIST --from proposed-3.1 --to proposed-3.2",
		Short: "Replace labels in given board list",
		Run: func(cmd *cobra.Command, args []string) {
			opts.Complete()
			opts.Run()
		},
	}
	cmd.Flags().StringVar(&opts.Board, "board", opts.Board, "The name of Trello board to operate with.")
	cmd.Flags().StringVar(&opts.List, "list", opts.List, "The name of Trello board list to operate with.")
	cmd.Flags().StringVar(&opts.From, "from", opts.From, "The name of the label to replace.")
	cmd.Flags().StringVar(&opts.To, "to", opts.To, "The name of target label.")
	cmd.Flags().BoolVar(&opts.Remove, "rm", opts.Remove, "Remove the --from label")
	return cmd
}

func (o *RelabelOptions) Complete() {
	if len(o.Board) == 0 {
		log.Fatal("--board must be set")
		return
	}
	if len(o.List) == 0 {
		log.Fatal("--list must be set")
		return
	}
	if !o.Remove && len(o.To) == 0 && len(o.From) == 0 {
		log.Fatal("--to must be set when --from is set and no --rm")
		return
	}
	if o.Remove && len(o.From) == 0 {
		log.Fatal("--from must be set when removing")
		return
	}
}

func (o *RelabelOptions) Run() {
	boards, err := Trello.User.Boards()
	if err != nil {
		log.Fatalf("error listing boards: %v", err)
	}
	found := false
	var board trello.Board
	for _, b := range boards {
		if b.Name == o.Board {
			board = b
			found = true
			break
		}
	}
	if !found {
		log.Fatalf("unable to find board %q", o.Board)
	}
	lists, err := board.Lists()
	if err != nil {
		log.Fatalf("error getting lists for board %q: %v", o.Board, err)
	}
	found = false
	var list trello.List
	for _, l := range lists {
		if l.Name == o.List {
			list = l
			found = true
			break
		}
	}
	if !found {
		log.Fatalf("unable to find list %q", o.List)
	}

	cards, err := list.Cards()
	if err != nil {
		log.Fatalf("error listing cards for %q: %v", list.Name, err)
	}

	labelMap, err := util.GetBoardLabelsIds(Trello.Client, board.Id)
	if err != nil {
		log.Fatalf("error listing board labels: %v", err)
	}

	for _, c := range cards {
		for _, l := range c.Labels {
			if l.Name == o.From {
				if o.Remove || len(o.To) > 0 {
					log.Printf("Removing %q label for %q ...", o.From, c.Name)
					if _, err := Trello.Client.Delete("/cards/" + c.Id + "/idLabels/" + labelMap[o.From]); err != nil {
						log.Printf("error removing label from %q: %v", c.Name, err)
					}
				}
				if len(o.To) > 0 {
					log.Printf("Adding %q label to %q ...", o.To, c.Name)
					if _, err := Trello.Client.Post("/cards/"+c.Id+"/idLabels", url.Values{"value": []string{labelMap[o.To]}}); err != nil {
						log.Printf("error adding label to %q: %v", c.Name, err)
					}
				}
			}
		}
	}
}
