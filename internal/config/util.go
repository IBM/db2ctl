package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("nvme-check", nvmeCheck)
	validate.RegisterValidation("ip-check", ipsCheck)
	validate.RegisterValidation("names-check", namesCheck)
	validate.RegisterValidation("node-check", numNodesCheck)
}

func ipsCheck(fl validator.FieldLevel) bool {
	numNodes := fl.Parent().FieldByName("NumNodes")
	//fmt.Println("numNodes : ", numNodes)
	numIPs := fl.Field().Len()
	//fmt.Println("numIPs : ", numIPs)
	if int64(numIPs) == numNodes.Int() {
		return true
	}
	return false
}

func nvmeCheck(fl validator.FieldLevel) bool {
	numNodes := fl.Parent().FieldByName("NumNodes")
	//fmt.Println("numNodes : ", numNodes)
	numNVME := fl.Field().Len()
	//fmt.Println("numNVME : ", numNVME)
	if int64(numNVME) == numNodes.Int() {
		return true
	}
	return false
}

func namesCheck(fl validator.FieldLevel) bool {
	numNodes := fl.Parent().FieldByName("NumNodes")
	//fmt.Println("numNodes : ", numNodes)
	numNames := fl.Field().Len()
	//fmt.Println("numNames : ", numIPs)
	if int64(numNames) == numNodes.Int() {
		return true
	}
	return false
}

func numNodesCheck(fl validator.FieldLevel) bool {
	numNodes := fl.Field().Int()
	role := fl.Top().Elem().FieldByName("Spec").FieldByName("DB2").FieldByName("Required").FieldByName("Role")
	if role.String() == "optimized" && numNodes%4 == 0 {
		return true
	} else if role.String() == "sandbox" && numNodes > 0 && numNodes < 16 {
		return true
	}
	return false
}

//Validate the struct
func (config *PCConfType) Validate() error {
	err := validate.Struct(config)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			// fmt.Println("err.Nam", err.Namespace())
			// fmt.Println("err.Fie", err.Field())
			// fmt.Println("err.Str", err.StructNamespace())
			// fmt.Println("err.Str", err.StructField())
			// fmt.Println("err.Tag", err.Tag())
			// fmt.Println("err.Act", err.ActualTag())
			// fmt.Println("err.Kin", err.Kind())
			// fmt.Println("err.Typ", err.Type())
			// fmt.Println("err.Val", err.Value())
			// fmt.Println("err.Par", err.Param())
			// fmt.Println()
			switch err.Tag() {
			case "node-check":
				return fmt.Errorf("err: number of nodes not correct according to role selected")
			case "nvme-check":
				return fmt.Errorf("err: nvme list does not match the number of nodes")
			case "ip-check":
				return fmt.Errorf("err: number of IP addresses do not match the number of nodes")
			case "names-check":
				return fmt.Errorf("err: number of node names do not match the number of nodes")
			case "oneof":
				return fmt.Errorf("%v - possible values: [%v]", err, strings.Join(strings.Split(err.Param(), " "), ","))
			case "min":
				return fmt.Errorf("%v - possible values: [%v]", err, strings.Join(strings.Split(err.Param(), " "), ","))
			}
		}
		return err
	}
	return nil
}

//Print the struct
func (config *PCConfType) Print() error {
	return printPretty(config)
}

//Print the struct
func (config *BinPacking) Print() error {
	return printPretty(config)
}

func printPretty(c interface{}) error {
	//Marshal
	json, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("error printing config : %v", err)
	}
	fmt.Printf("\nParsed configuration: %s\n", string(json))
	return nil
}

//PreconfigureFields preconfigures some fields for config struct
func (config *PCConfType) PreconfigureFields() error {

	//clientspec config calculations
	virtualIP := config.Spec.NFS.Server.Required.VirtualIP
	CIDRNetMask := config.Spec.NFS.Server.Required.CIDRNetMask

	_, clientSpec, err := net.ParseCIDR(virtualIP + "/" + CIDRNetMask)
	if err != nil {
		return fmt.Errorf("error parsing virtual IP for clientspec, err : %v", err)
	}
	config.Spec.NFS.Server.Required.ClientSpec = clientSpec.String()

	//NamesAsList
	config.Spec.Nodes.Required.NamesAsList = strings.Join(config.Spec.Nodes.Required.Names, " ")

	//IP Addresses as list
	config.Spec.Nodes.Required.IPAddressesAsList = strings.Join(config.Spec.Nodes.Required.IPAddresses, " ")

	//Node info mapping
	err = nodeInfoMapping(config)
	if err != nil {
		return err
	}
	//Minor list
	config.Spec.Linbit.Optional.VolumeDefinition.Nodes.Minor = "1000"
	config.Spec.Linbit.Optional.VolumeDefinition.NFS.Minor = "3000"

	err = nfsServerPlacement(config)
	if err != nil {
		return err
	}

	err = sizingCalculations(config)
	if err != nil {
		return err
	}

	err = secretBase64DecodeHandling(config)
	if err != nil {
		return err
	}
	return nil
}

func sizingCalculations(config *PCConfType) error {

	var nfsSize int
	if config.Spec.Linbit.Optional.VolumeDefinition.NFS.Size == "" {
		minSize := int(math.Inf(0))
		var order string
		for _, nodeName := range config.Spec.Nodes.Required.Names {
			size := config.Spec.Nodes.Required.NodeInfoMap[nodeName].Size
			if size < minSize {
				minSize = size
			}
			order = config.Spec.Nodes.Required.NodeInfoMap[nodeName].Order
		}
		// fmt.Println("min : ", minSize)
		nfsSizeF := math.Min(math.Max(10, float64((minSize*1)/10)), 100)
		nfsSize = int(nfsSizeF)
		// fmt.Println("nfsSize : ", nfsSize)
		// config.Spec.Linbit.Optional.VolumeDefinition.NFS.Size = strconv.FormatFloat(nfsSize, 'f', 0, 64) + order
		config.Spec.Linbit.Optional.VolumeDefinition.NFS.Size = strconv.Itoa(nfsSize) + order
	} else {
		var err error
		nfsSize, _, err = splitMemorySize(config.Spec.Linbit.Optional.VolumeDefinition.NFS.Size)
		if err != nil {
			return err
		}
	}

	//db2local
	for _, nodeName := range config.Spec.Nodes.Required.Names {
		nodeInfo := config.Spec.Nodes.Required.NodeInfoMap[nodeName]

		if config.Spec.Linbit.Optional.VolumeDefinition.DB2Local.Size == "" {

			//start with node size, more calculations to follow
			db2LocalSize := nodeInfo.Size
			if nodeInfo.HasNfsServer {
				db2LocalSize -= nfsSize
			}
			if config.Spec.DB2.Required.Replicated {
				db2LocalSize = (db2LocalSize * 2) / 10
			}
			if db2LocalSize <= 0 {
				return fmt.Errorf("db2LocalSize size <= 0. Please check values")
			}
			nodeInfo.DB2LocalSize = db2LocalSize
		} else {
			db2LocalSize, _, err := splitMemorySize(config.Spec.Linbit.Optional.VolumeDefinition.DB2Local.Size)
			if err != nil {
				return err
			}
			nodeInfo.DB2LocalSize = db2LocalSize
		}
	}

	//node sizing
	if config.Spec.DB2.Required.Replicated {
		partitionsPerNode := config.Spec.Nodes.Required.Partitions / config.Spec.Nodes.Required.NumNodes
		sizePerPartition := make(map[string]string)
		var partitionSize int
		for nodeNum, nodeName := range config.Spec.Nodes.Required.Names {
			nodeInfo := config.Spec.Nodes.Required.NodeInfoMap[nodeName]

			if config.Spec.Linbit.Optional.VolumeDefinition.Nodes.Size == "" {

				nodeSize := nodeInfo.Size
				db2LocalSize := nodeInfo.DB2LocalSize
				step1 := (nodeSize - db2LocalSize)
				if nodeInfo.HasNfsServer {
					step1 -= nfsSize
				}
				// fmt.Println("step1 : ", step1)
				// partitionSize := step1 / float64(partitionsPerNode)
				partitionSize = step1 / (partitionsPerNode * 2)
				if partitionSize <= 0 {
					return fmt.Errorf("partition size <= 0. Please check values")
				}
			} else {
				var err error
				partitionSize, _, err = splitMemorySize(config.Spec.Linbit.Optional.VolumeDefinition.Nodes.Size)
				if err != nil {
					return err
				}
			}
			for numpartition := nodeNum * partitionsPerNode; numpartition < (nodeNum*partitionsPerNode)+partitionsPerNode; numpartition++ {
				keyName := config.Spec.Nodes.Optional.DBPrimitiveNamePrefix + fmt.Sprintf("%04d", numpartition)
				sizePerPartition[keyName] = strconv.Itoa(partitionSize) + config.Spec.Nodes.Required.NodeInfoMap[nodeName].Order
			}
		}
		config.Spec.Nodes.Required.SizePerPartition = sizePerPartition
	}
	return nil
}

func nodeInfoMapping(config *PCConfType) error {
	nvmeList := config.Spec.Nodes.Required.NVMEList
	nodeInfoMap := make(map[string]*nodeInfo)

	for index, nodeName := range config.Spec.Nodes.Required.Names {

		var nvmeNames []string
		var nvmeSizeForNode int
		var order string
		//check for tests
		if len(nvmeList) == len(config.Spec.Nodes.Required.Names) {
			nvmeListForNode := nvmeList[index]
			//set stripe size equal to nvmeListSize of node
			if config.Spec.Linbit.Optional.NumStripes == 0 {
				config.Spec.Linbit.Optional.NumStripes = len(nvmeListForNode)
			}
			for _, nvme := range nvmeListForNode {
				nvmeNames = append(nvmeNames, nvme["name"])
				nvmeSize, orderUsed, err := splitMemorySize(nvme["size"])
				if err != nil {
					return err
				}
				order = orderUsed
				nvmeSizeForNode += nvmeSize
			}
		}
		nvmeListAsString := strings.Join(nvmeNames, " ")
		nodeInfoInstance := &nodeInfo{
			NVMEList: strings.TrimSpace(nvmeListAsString),
			Size:     nvmeSizeForNode,
			Order:    order,
		}
		nodeInfoMap[nodeName] = nodeInfoInstance
	}
	config.Spec.Nodes.Required.NodeInfoMap = nodeInfoMap
	return nil
}

func nfsServerPlacement(config *PCConfType) error {
	numNFSServers := config.Spec.Linbit.Optional.VolumeDefinition.NFS.NumNFSServers
	numNodes := config.Spec.Nodes.Required.NumNodes
	nodeNames := config.Spec.Nodes.Required.Names

	if numNFSServers < 4 || numNFSServers > 8 {
		return fmt.Errorf("%v NFS servers not supported", numNFSServers)
	}
	var nodesForNFS []string
	//hitsPerUnit records the hits per unit (4 servers)
	// example: map[0:2 1:2 2:2] -> 0th unit needs 2 placements, 1st unit needs 1, etc.
	hitsPerUnit := make(map[int]int)

	for numServer := 0; numServer < numNFSServers; numServer++ {

		//we find which unit the server should be placed on
		unitForServer := numServer % (numNodes / 4)

		//if unit had a hit earlier, we need to give the next available server
		if hitForUnit := hitsPerUnit[unitForServer]; hitForUnit != 0 {
			hitsPerUnit[unitForServer] = hitForUnit + 1
			nodeName := nodeNames[(4*unitForServer)+hitForUnit]
			infoMapForNode := config.Spec.Nodes.Required.NodeInfoMap[nodeName]
			// fmt.Println("info map : ", infoMapForNode)
			infoMapForNode.HasNfsServer = true
			// fmt.Println("info map : ", infoMapForNode)
			nodesForNFS = append(nodesForNFS, nodeName)
		} else { //else, we can give the first server on that unit
			hitsPerUnit[unitForServer] = 1
			nodeName := nodeNames[4*unitForServer]
			infoMapForNode := config.Spec.Nodes.Required.NodeInfoMap[nodeName]
			infoMapForNode.HasNfsServer = true
			nodesForNFS = append(nodesForNFS, nodeName)
		}
	}
	//fmt.Println(nodesForNFS)
	config.Spec.NFS.Server.Required.NodesForPlacement = nodesForNFS
	return nil
}
func secretBase64DecodeHandling(config *PCConfType) error {

	//instance secret
	base64InstanceSecret := config.Spec.DB2.Required.InstanceSecret
	decodedInstanceVal, err := base64.StdEncoding.DecodeString(strings.TrimSpace(base64InstanceSecret))
	if err != nil {
		return fmt.Errorf("could not decode instanceSecret %v to base64, err : %v", base64InstanceSecret, err)
	}
	config.Spec.DB2.Required.InstanceSecretVal = string(decodedInstanceVal)

	//fenced secret
	base64FencedSecret := config.Spec.DB2.Required.FencedSecret
	decodedFencedVal, err := base64.StdEncoding.DecodeString(strings.TrimSpace(base64FencedSecret))
	if err != nil {
		return fmt.Errorf("could not decode fencedSecret %v to base64, err : %v", base64FencedSecret, err)
	}
	config.Spec.DB2.Required.FencedSecretVal = string(decodedFencedVal)
	return nil
}

func splitMemorySize(memory string) (size int, order string, err error) {

	memory = strings.TrimSpace(memory)

	size, err = strconv.Atoi(regexp.MustCompile("[0-9]+").FindString(memory))
	if err != nil {
		return size, order, fmt.Errorf("cannot parse memory size : %v", memory)
	}

	order = regexp.MustCompile("[a-zA-Z]+").FindString(memory)
	return size, order, err
}
