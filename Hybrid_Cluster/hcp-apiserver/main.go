package main

import (
	"Hybrid_Cluster/hcp-apiserver/converter/mappingTable"
	"Hybrid_Cluster/hcp-apiserver/handler"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/service/eks"
)

func parser(w http.ResponseWriter, req *http.Request, input interface{}) {
	jsonDataFromHttp, err := ioutil.ReadAll(req.Body)
	json.Unmarshal(jsonDataFromHttp, input)
	defer req.Body.Close()
	if err != nil {
		log.Println(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
}

func checkErr(w http.ResponseWriter, err error) {
	if err != nil {
		log.Println(err)
	}
}

// hybridctl join <platformName> <ClusterName>
func join(w http.ResponseWriter, req *http.Request) {

	clusterInfo := make(map[string]interface{})
	var info = mappingTable.ClusterInfo{
		PlatformName: clusterInfo["PlatformName"].(string),
		ClusterName:  clusterInfo["ClusterName"].(string),
	}
	parser(w, req, clusterInfo)
	handler.Join(info)
}

func createAddon(w http.ResponseWriter, req *http.Request) {

	var createAddonInput eks.CreateAddonInput

	parser(w, req, &createAddonInput)
	out, err := handler.CreateAddon(createAddonInput)
	// checkErr(w, err)
	var jsonData []byte
	if err != nil {
		log.Println(err)
		jsonData, _ = json.Marshal(&err)
	} else {
		jsonData, _ = json.Marshal(&out)
	}
	w.Write([]byte(jsonData))
}

func deleteAddon(w http.ResponseWriter, req *http.Request) {

	var deleteAddonInput eks.DeleteAddonInput

	parser(w, req, &deleteAddonInput)
	out, err := handler.DeleteAddon(deleteAddonInput)
	checkErr(w, err)
	jsonData, _ := json.Marshal(&out)
	w.Write([]byte(jsonData))
}

func describeAddon(w http.ResponseWriter, req *http.Request) {

	var describeAddonInput eks.DescribeAddonInput

	parser(w, req, &describeAddonInput)
	out, err := handler.DescribeAddon(describeAddonInput)
	checkErr(w, err)
	jsonData, _ := json.Marshal(&out)
	w.Write([]byte(jsonData))
}

func describeAddonVersions(w http.ResponseWriter, req *http.Request) {

	var describeAddonVersionsInput eks.DescribeAddonVersionsInput

	parser(w, req, &describeAddonVersionsInput)
	out, err := handler.DescribeAddonVersions(describeAddonVersionsInput)
	checkErr(w, err)
	jsonData, _ := json.Marshal(&out)
	w.Write(jsonData)

}

func listAddon(w http.ResponseWriter, req *http.Request) {

	var listAddonInput eks.ListAddonsInput

	parser(w, req, &listAddonInput)
	out, err := handler.ListAddon(listAddonInput)
	checkErr(w, err)
	jsonData, _ := json.Marshal(&out)
	w.Write([]byte(jsonData))

}

func updateAddon(w http.ResponseWriter, req *http.Request) {

	var updateAddonInput eks.UpdateAddonInput

	parser(w, req, &updateAddonInput)
	out, err := handler.UpdateAddon(updateAddonInput)
	checkErr(w, err)
	jsonData, _ := json.Marshal(&out)
	w.Write([]byte(jsonData))
}

func main() {
	http.HandleFunc("/join", join)
	http.HandleFunc("/createAddon", createAddon)
	http.HandleFunc("/deleteAddon", deleteAddon)
	http.HandleFunc("/describeAddon", describeAddon)
	http.HandleFunc("/describeAddonVersions", describeAddonVersions)
	http.HandleFunc("/listAddon", listAddon)
	http.HandleFunc("/updateAddon", updateAddon)
	http.ListenAndServe(":8000", nil)
}

/*
*** optionCheck Module ***
- kubernetes Platform Check
*/
// func OptionCheck(w http.ResponseWriter, req *http.Request) {
// 	fmt.Println("---Checking Options start---")

// cli := make(map[string]interface{})

// info, err := ioutil.ReadAll(req.Body)
// json.Unmarshal([]byte(jsonDataFromHttp), &cli)
// defer req.Body.Close()
// if err != nil {
// 	panic(err)
// }
// handler.JoinHandler(info)

// info := mappingTable.CommandInfo{Cmd: cli["Cmd"].(string), Platform: cli["PlatformName"].(string), ClusterName: cli["ClusterName"].(string)}
// platform 이름 기입여부 체크
// switch info.Platform {
// case "gke", "aks", "eks":
// 	printInfo(info.Platform, info.ClusterName)
// 	handler.JoinHandler(info)
// default:
// 	schedulingPlatform()
// 	handler.JoinHandler(info)
// }

// w.Header().Set("Content-Type", "application/json")
// w.WriteHeader(http.StatusOK)
// }

// func printInfo(PlatformName string, clusterName string) {
// 	fmt.Println("kubernetes engine Name : ", PlatformName)
// 	fmt.Printf("Cluster Name : %s\n", clusterName)
// 	fmt.Printf("---Checking Options Done---\n\n")
// }

// func schedulingPlatform() {
// 	fmt.Println("---SchedulingPlatform---")
// }
