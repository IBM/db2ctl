package handler

import (
	"fmt"

	"github.com/IBM/db2ctl/internal/config"
	"github.com/IBM/db2ctl/internal/flag"
	"github.com/spf13/pflag"
)

func getFlags(queryParams map[string][]string) *pflag.FlagSet {
	flags := pflag.NewFlagSet("db2ctl-rest", pflag.ContinueOnError)

	var confFile string
	var dryRun bool
	var reRun bool

	flags.StringVarP(&confFile, flag.ConfigurationFile, "c", "db2ctl-sample.yaml", "configuration yaml file needed for application")
	flags.BoolVarP(&dryRun, flag.DryRun, "d", false, "(optional) shows what scripts will run, but does not run the scripts")
	flags.BoolVarP(&reRun, flag.ReRun, "r", false, "(optional) re-run script from initial state, ignoring previously saved state")

	for key, value := range queryParams {
		switch key {
		case "conf":
			flags.Set(flag.ConfigurationFile, value[0])
		case "dryRun":
			flags.Set(flag.DryRun, value[0])
		case "reRun":
			flags.Set(flag.ReRun, value[0])
		}
	}
	return flags
}

func makeAPICallToMagneto(defaultNumOfNodes int, newConfigurationRequest configurationRequest, config *config.Combined) error {
	nvmeList := config.Spec.Nodes.Required.NVMEList
	ipAddresses := config.Spec.Nodes.Required.IPAddresses
	names := config.Spec.Nodes.Required.Names

	var size string
	name := make([]string, defaultNumOfNodes)
	for index, element := range nvmeList[0] {
		for key, value := range element {
			if key == "name" {
				name[index] = value
			} else {
				size = value
			}
		}
	}

	if newConfigurationRequest.Nodes > defaultNumOfNodes {
		count := newConfigurationRequest.Nodes - defaultNumOfNodes //count=12
		index := 0
		for count > 0 {
			nvmeListAdd := make([]map[string]string, defaultNumOfNodes)
			for i := 0; i < defaultNumOfNodes; i++ {
				m := make(map[string]string)
				m["name"] = name[i]
				m["size"] = size
				nvmeListAdd[i] = m
			}
			ipAddressesAdd := ipAddresses[index%defaultNumOfNodes]
			namesAdd := names[index%defaultNumOfNodes]
			config.Spec.Nodes.Required.NVMEList = append(config.Spec.Nodes.Required.NVMEList, nvmeListAdd)
			config.Spec.Nodes.Required.IPAddresses = append(config.Spec.Nodes.Required.IPAddresses, ipAddressesAdd)
			config.Spec.Nodes.Required.Names = append(config.Spec.Nodes.Required.Names, namesAdd)
			count--
			index++
		}
	} else {
		return fmt.Errorf("Minimum number of nodes is %v", 4)
	}
	// fmt.Println("NumNodes: ", config.Spec.Nodes.Required.NumNodes)
	// fmt.Println("NVMEList: ", len(config.Spec.Nodes.Required.NVMEList))
	// fmt.Println("IPAddresses: ", len(config.Spec.Nodes.Required.IPAddresses))
	// fmt.Println("Names: ", len(config.Spec.Nodes.Required.Names))
	return nil
}
