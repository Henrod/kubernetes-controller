// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"log"
	"time"

	"github.com/Henrod/kube-controller/controller"
	"github.com/Henrod/kube-controller/models"
	"github.com/spf13/cobra"
)

var (
	namespace    string
	image        string
	numberOfPods int
	servicePort  int
	ticker       time.Duration
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts watcher",
	Long:  `creates namespace, pods and service and watch for pod events`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("connecting to kubernetes")
		kubernetes, err := models.NewKubernetes()
		checkErr(err)

		log.Print("configuring watcher")
		watcher := controller.NewWatcher(
			kubernetes, namespace, image, numberOfPods, ticker)
		checkErr(err)

		err = watcher.Create(servicePort)
		checkErr(err)

		err = watcher.Watch()
		checkErr(err)
	},
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	RootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVar(
		&namespace, "namespace",
		"default", "namespace to run watcher",
	)
	startCmd.Flags().StringVar(
		&image, "image",
		"redis", "image of the pods",
	)
	startCmd.Flags().IntVar(
		&numberOfPods, "numberOfPods",
		3, "number of pods to run on namespace",
	)
	startCmd.Flags().IntVar(
		&servicePort, "port",
		6379, "port to run service",
	)
	startCmd.Flags().DurationVar(
		&ticker, "ticker",
		10*time.Second, "period to watch pods",
	)
}
