// Copyright © 2021 NAME HERE <EMAIL ADDRESS>
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var flagvar int

type Cluster_info struct {
	Project_id    string `json:"project_id"`
	Cluster_name  string `json:"cluster_name"`
	Region        string `json:"region"`
	Gke_num_nodes uint64 `json:"gke_num_nodes"`
}

type TF struct {
	Resource *Resource            `json:"resource"`
	Provider *map[string]Platform `json:"provider"`
}

type Platform struct {
	Project string `json:"project"`
	Region  string `json:"region"`
}

type Resource struct {
	Google_container_cluster   *map[string]Cluster_type   `json:"google_container_cluster"`
	Google_container_node_pool *map[string]Node_pool_type `json:"google_container_node_pool"`
}

type Cluster_type struct {
	Name                     string `json:"name"`
	Location                 string `json:"location"`
	Remove_default_node_pool string `json:"remove_default_node_pool"`
	Initial_node_count       int    `json:"initial_node_count"`
}

type Node_pool_type struct {
	Name        string       `json:"name"`
	Location    string       `json:"location"`
	Cluster     string       `json:"cluster"`
	Node_count  int          `json:"node_count"`
	Node_config *Node_config `json:"node_config"`
}

type Labels struct {
	Env string `json:"env"`
}
type Node_config struct {
	Oauth_scopes []string  `json:"oauth_scopes"`
	Labels       *Labels   `json:"labels"`
	Machine_type string    `json:"machine_type"`
	Tags         []string  `json:"tags"`
	Metadata     *Metadata `json:"metadata"`
}

type Metadata struct {
	Disable_legacy_endpoints string `json:"disable-legacy-endpoints"`
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("create called ㅡㅡ")
		cli := Cli{args[0], args[1]}
		fmt.Println(cli)
		if len(args) == 0 {
			fmt.Println("Run 'hybridctl create --help' to view all commands")
		} else if args[0] == "gke" {
			if args[1] == "" {

				create_gke(cli)
				fmt.Println("Run 'hybridctl create --help' to view all commands")
			} else {
				fmt.Println("kubernetes engine Name : ", args[0])
				fmt.Printf("Cluster Name : %s\n", args[1])

				create_gke(cli)
			}
		}
	},
}

func create_gke(info Cli) {

	cluster := "cluster_1"
	num := 1
	data := make([]Cluster_info, 1)
	platform := "google"

	data[0].Project_id = "keti-container"
	data[0].Cluster_name = info.ClusterName
	data[0].Region = "us-central1-a"
	data[0].Gke_num_nodes = uint64(num)

	doc, _ := json.Marshal(data)

	fmt.Println(strings.Trim(string(doc), "[]"))

	err := ioutil.WriteFile("/root/go/src/Hybrid_Cluster/terraform/gke/"+cluster+"/"+info.ClusterName+".tfvars.json", []byte(strings.Trim(string(doc), "[]")), os.FileMode(0644))

	if err != nil {
		panic(err)
	}

	send_js_cluster := TF{
		Provider: &map[string]Platform{
			platform: {
				Project: "keti-container",
				Region:  "us-central1-a",
			},
		},
		Resource: &Resource{
			Google_container_cluster: &map[string]Cluster_type{
				info.ClusterName: {
					Name:                     info.ClusterName,
					Location:                 "us-central1-a",
					Remove_default_node_pool: "true",
					Initial_node_count:       num,
				},
			},
			Google_container_node_pool: &map[string]Node_pool_type{
				info.ClusterName + "_node_pool": {
					Name:       info.ClusterName + "-nodes",
					Location:   "us-central1-a",
					Cluster:    info.ClusterName,
					Node_count: num,
					Node_config: &Node_config{
						Labels: &Labels{
							Env: "keti-container",
						},
						Metadata: &Metadata{
							Disable_legacy_endpoints: "true",
						},
						Tags:         []string{"gke-node", "keti-container-gke"},
						Machine_type: "n1-standard-1",
						Oauth_scopes: []string{"https://www.googleapis.com/auth/logging.write", "https://www.googleapis.com/auth/monitoring"},
					},
				},
			},
		},
	}

	send, err := json.MarshalIndent(send_js_cluster, "", " ")
	if err != nil {
		panic(err)
	}

	// src, _ := json.Marshal([]byte(string(resource)))

	err = ioutil.WriteFile("/root/go/src/Hybrid_Cluster/terraform/gke/"+cluster+"/"+info.ClusterName+".tf.json", []byte(string(send)), os.FileMode(0644))
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("terraform", "apply", "-auto-approve")
	// cmd := exec.Command("terraform", "plan")
	cmd.Dir = "../terraform/gke/" + cluster

	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(output))
	}
}

func init() {
	RootCmd.AddCommand(createCmd)

	// flag.IntVar(&flagvar, "node", 1, "set node count")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
