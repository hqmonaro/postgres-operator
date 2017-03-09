/*
 Copyright 2017 Crunchy Data Solutions, Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/crunchydata/operator/tpr"
	"github.com/spf13/cobra"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a database or Cluster",
	Long: `delete allows you to delete a database or cluster
For example:

crunchy delete database mydatabase
crunchy delete cluster mycluster`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			fmt.Println(`You must specify the type of resource to delete.  Valid resource types include:
	* database
	* cluster`)
		} else {
			switch args[0] {
			case "database":
			case "cluster":
				break
			default:
				fmt.Println(`You must specify the type of resource to delete.  Valid resource types include: 
	* database
	* cluster`)
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
	deleteCmd.AddCommand(deleteDatabaseCmd)
	deleteCmd.AddCommand(deleteClusterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// showDatbaseCmd represents the show database command
var deleteDatabaseCmd = &cobra.Command{
	Use:   "database",
	Short: "delete a database",
	Long: `delete a crunchy database. For example:
	crunchy delete database mydatabase`,
	Run: func(cmd *cobra.Command, args []string) {
		deleteDatabase(args)
	},
}

func deleteDatabase(args []string) {
	var err error
	//result := tpr.CrunchyDatabaseList{}
	databaseList := tpr.CrunchyDatabaseList{}
	err = Tprclient.Get().Resource("crunchydatabases").Do().Into(&databaseList)
	if err != nil {
		panic(err)
	}
	for _, arg := range args {
		// delete the crunchydatabase resource instance
		for _, database := range databaseList.Items {
			fmt.Println("database LIST: " + database.Spec.Name)
			if database.Spec.Name == arg {
				err = Tprclient.Delete().
					Resource("crunchydatabases").
					Namespace(api.NamespaceDefault).
					Name(database.Spec.Name).
					//Body(database).
					Do().
					Error()
					//Into(&tpr.CrunchyDatabase{})
				if err != nil {
					fmt.Println("error deleting crunchydatabase " + arg)
					fmt.Println(err.Error())
				}
				fmt.Println("deleted crunchydatabase " + database.Spec.Name)
			}

		}
		// delete the Service
		err = Clientset.Services(v1.NamespaceDefault).Delete(arg,
			&v1.DeleteOptions{})
		if err != nil {
			fmt.Println("error deleting service " + arg)
			fmt.Println(err.Error())
		}
		fmt.Println("deleted service " + arg)

		// delete the Pod
		err = Clientset.Pods(v1.NamespaceDefault).Delete(arg,
			&v1.DeleteOptions{})
		if err != nil {
			fmt.Println("error deleting pod " + arg)
			fmt.Println(err.Error())
		}
		fmt.Println("deleted pod " + arg)

	}
}

// deleteClusterCmd represents the delete cluster command
var deleteClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "delete a cluster",
	Long: `delete a crunchy cluster. For example:
	crunchy delete cluster mycluster`,
	Run: func(cmd *cobra.Command, args []string) {
		deleteCluster(args)
	},
}

func deleteCluster(args []string) {
	for _, arg := range args {
		// Fetch a list of our database TPRs
		fmt.Println("deleting cluster " + arg)
		databaseList := tpr.CrunchyDatabaseList{}
		err := Tprclient.Get().Resource("crunchydatabases").Do().Into(&databaseList)
		if err != nil {
			panic(err)
		}
		for _, database := range databaseList.Items {
			fmt.Println("database LIST: " + database.Spec.Name)
		}
	}
}
