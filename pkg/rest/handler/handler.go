package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/IBM/db2ctl/internal/command"
	"github.com/IBM/db2ctl/internal/flag"
	"github.com/IBM/db2ctl/internal/ws"
)

//Index page of application
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

//WebSocket starts the websocket
func WebSocket(w http.ResponseWriter, r *http.Request) {
	ws.ServeWs(w, r)
}

//Configuration endpoint
func Configuration(w http.ResponseWriter, r *http.Request) {

	flags := getFlags(r.URL.Query())
	commandCenter := command.New(flags)
	commandCenter.CreateSampleConfigFile().
		ParseYaml(command.SampleConfigFileName)

		//make changes
	reqBody, _ := ioutil.ReadAll(r.Body)
	var newConfigurationRequest configurationRequest
	json.Unmarshal(reqBody, &newConfigurationRequest)

	config := commandCenter.CombinedConfig
	defaultNumOfNodes := config.Spec.Nodes.Required.NumNodes
	config.Spec.Nodes.Required.NumNodes = newConfigurationRequest.Nodes
	config.Spec.DB2.Required.Role = newConfigurationRequest.Type

	//TODO - replace with actual Magneto call later
	makeAPICallToMagneto(defaultNumOfNodes, newConfigurationRequest, config)

	err := config.Validate()
	if err != nil {
		http.Error(w, "error while validating values, err: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = config.PreconfigureFields()
	if err != nil {
		http.Error(w, "error while Preconfiguring Fields, err: "+err.Error(), http.StatusBadRequest)
		return
	}

	confFile, _ := flags.GetString(flag.ConfigurationFile)
	commandCenter.CreateFromConfig(confFile)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)

}

//InstallModule installs a module
func InstallModule(w http.ResponseWriter, r *http.Request) {

	commandCenter := command.New(getFlags(r.URL.Query()))
	err := commandCenter.ParseYaml(command.SampleConfigFileName).Error
	if err != nil {
		http.Error(w, "error with configuration, err: "+err.Error(), http.StatusBadRequest)
		return
	}

	var moduleRequest moduleRequest
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &moduleRequest)
	if moduleRequest.ModuleName == "" {
		http.Error(w, "invalid module", http.StatusBadRequest)
		return
	}

	moduleName := strings.Join(strings.Split(moduleRequest.ModuleName, " "), "/")
	err = commandCenter.GenerateConfigFilesFromDir(moduleName).Error
	if err != nil {
		http.Error(w, "invalid module", http.StatusBadRequest)
		return
	}

	go func() {
		commandCenter.
			RunBashScripts()
	}()

	fmt.Fprint(w, "Success!")
}

//GetState gets the state of a module
func GetState(w http.ResponseWriter, r *http.Request) {

	var moduleRequest moduleRequest
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &moduleRequest)

	moduleName := strings.Join(strings.Split(moduleRequest.ModuleName, " "), "/")
	commandCenter := command.New(getFlags(r.URL.Query())).ReturnStateForDir(moduleName)
	stateMap := commandCenter.ReturnStateMap

	regexFile := regexp.MustCompile("([0-9]+-)")
	regexDir := regexp.MustCompile("([a-zA-Z0-9/]+/[0-9-]*)")
	stateJSON := stateJSON{}
	stepInstance := step{}
	stepList := []step{}
	var stateOfExecution string

	var sortedDirkeys []string
	for key := range stateMap {
		sortedDirkeys = append(sortedDirkeys, key)
	}
	sort.Strings(sortedDirkeys)

	for _, dir := range sortedDirkeys {
		stepInstance.Module = regexDir.ReplaceAllString(dir, "")
		taskInstance := task{}
		taskList := []task{}

		mapDir := stateMap[dir]
		var sortedFileKeys []string
		for key := range mapDir {
			sortedFileKeys = append(sortedFileKeys, key)
		}
		sort.Strings(sortedFileKeys)

		for _, file := range sortedFileKeys {
			taskInstance.TaskName = regexFile.ReplaceAllString(strings.TrimSuffix(file, ".sh"), "")
			taskInstance.FileExecStatus = mapDir[file]
			if taskInstance.FileExecStatus.State == "running" { //update running time
				taskInstance.FileExecStatus.TimeTaken = time.Since(taskInstance.FileExecStatus.StartTime).String()
			}
			if taskInstance.FileExecStatus.State != "" {
				stateOfExecution = string(taskInstance.FileExecStatus.State)
			}
			taskList = append(taskList, taskInstance)
		}
		stepInstance.Tasks = taskList
		stepList = append(stepList, stepInstance)
	}
	stateJSON.Steps = stepList
	stateJSON.State = stateOfExecution

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stateJSON)
}

//Cancel cancels running command
func Cancel(w http.ResponseWriter, r *http.Request) {
	err := command.New(getFlags(r.URL.Query())).
		StopRunningCommand().
		Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
