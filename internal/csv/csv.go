package csv

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/IBM/db2ctl/internal/config"
	"github.com/gocarina/gocsv"
)

//GenerateMappingFile generates mapping file
func GenerateMappingFile(pcConf *config.Combined, generatedDir, mappingFileDefaultName string) (string, string, error) {
	fileName := generatedDir + "/" + mappingFileDefaultName
	fileToGenerate, err := os.Create(fileName)
	if err != nil {
		return "", "", fmt.Errorf("error while creating file %v: %v", fileName, err)
	}
	defer fileToGenerate.Close()

	mapArray := []*config.DBMapStruct{}

	dbMountPointName := pcConf.Spec.DB2.Optional.TopLevelDir + "/" +
		pcConf.Spec.DB2.Required.InstanceName + "/" + pcConf.Spec.Nodes.Optional.DBPrimitiveNamePrefix

	for i := 0; i < pcConf.Spec.Nodes.Required.Partitions; i++ {

		dbMinorValue, _ := strconv.Atoi(pcConf.Spec.Linbit.Optional.VolumeDefinition.Nodes.Minor)
		loopValue, _ := strconv.Atoi(fmt.Sprintf("%04d", i))
		mapArray = append(mapArray, &config.DBMapStruct{
			DBMountPoint: dbMountPointName + fmt.Sprintf("%04d", i),
			DBDeviceName: "/dev/drbd" + strconv.Itoa(dbMinorValue+loopValue),
		})
	}
	csvContent, err := gocsv.MarshalString(&mapArray)  // Get all clients as CSV string
	err = gocsv.MarshalFile(&mapArray, fileToGenerate) // Use this to save the CSV back to the file
	if err != nil {
		return "", "", fmt.Errorf("error while marshalling in %v file : %v", fileName, err)
	}
	return fileName, csvContent, nil
}

//ReadMappingCsv reads mapping csv
func ReadMappingCsv(conf *config.Combined, generatedDir, mappingFileDefaultName string) error {
	mapArray := []*config.DBMapStruct{}

	fileName := generatedDir + "/" + mappingFileDefaultName
	mappingFile, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("error opening file %v : %v", fileName, err)
	}
	defer mappingFile.Close()

	if err := gocsv.UnmarshalFile(mappingFile, &mapArray); err != nil {
		return fmt.Errorf("error unmarshalling mapping csv %v : %v", fileName, err)
	}

	for _, row := range mapArray {
		mapStructInstance := config.DBMapStruct{
			DBMountPoint: row.DBMountPoint,
			DBDeviceName: row.DBDeviceName,
		}
		conf.Mapping = append(conf.Mapping, mapStructInstance)
	}
	return nil
}

//GenerateBinPackingCsv generates bin packing csv
func GenerateBinPackingCsv(conf *config.Combined, generatedDir, binPackingFilename string) (string, error) {

	dataToNodeMap := make(map[int]map[string]string)
	prepareDataToNodeMapDynamically(conf, dataToNodeMap)

	var dataToNode []config.DataToNode

	//doing this for sorting keys
	var keysForMap []int
	for key := range dataToNodeMap {
		keysForMap = append(keysForMap, key)
	}
	sort.Ints(keysForMap)

	for _, key := range keysForMap {
		//fmt.Println("key : ", key)
		dataToNodeInstance := config.DataToNode{
			DBPrimitiveName: conf.Spec.Nodes.Optional.DBPrimitiveNamePrefix + fmt.Sprintf("%04d", key),
			PrimaryServer:   dataToNodeMap[key]["primary"],
			ReplicaServer:   dataToNodeMap[key]["replicaServer"],
		}
		dataToNode = append(dataToNode, dataToNodeInstance)
	}

	fileName := generatedDir + "/" + binPackingFilename
	binPackingFile, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("error opening file %v : %v", binPackingFilename, err)
	}
	defer binPackingFile.Close()

	_, err = gocsv.MarshalString(dataToNode)
	err = gocsv.MarshalFile(&dataToNode, binPackingFile)
	if err != nil {
		return "", fmt.Errorf("error while marshalling in %v file : %v", fileName, err)
	}
	return fileName, nil
}

//ReadBinPackingCsv reads data from bin packing csv
func ReadBinPackingCsv(conf *config.Combined, generatedDir, binPackingFilename string) error {
	fileName := generatedDir + "/" + binPackingFilename
	binPackingFile, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("error opening file %v : %v", binPackingFilename, err)
	}
	defer binPackingFile.Close()

	var dataToNode []config.DataToNode
	if err := gocsv.UnmarshalFile(binPackingFile, &dataToNode); err != nil {
		return fmt.Errorf("error unmarshalling binpacking csv %v : %v", binPackingFilename, err)
	}

	conf.DataToNode = dataToNode
	return nil
}

func prepareDataToNodeMapDynamically(config *config.Combined, dataToNodeMap map[int]map[string]string) {
	partitions := config.Spec.Nodes.Required.Partitions
	numNodes := config.Spec.Nodes.Required.NumNodes
	numMLNPerServer := partitions / numNodes
	nodes := config.Spec.Nodes.Required.Names

	//fill in the data partitions
	for i := 0; i < partitions; i++ {
		newMap := make(map[string]string)
		newMap["primary"] = nodes[i/numMLNPerServer]
		dataToNodeMap[i] = newMap
	}

	for iteration := 0; iteration < numNodes/4; iteration++ {
		nodeNum := iteration * 4
		partition := (partitions / (numNodes / 4)) * (iteration)
		for {
			// fmt.Println("iteration: ", iteration)
			// fmt.Println("partition: ", partition)
			// fmt.Println("nodeNum: ", nodeNum+1)
			if dataToNodeMap[partition]["primary"] == nodes[nodeNum] {
				nodeNum = (iteration * 4) + (nodeNum+1)%4
			} else {
				dataToNodeMap[partition]["replicaServer"] = nodes[nodeNum]
				//fmt.Println("dataToNodeMap[partition] : ", dataToNodeMap[partition])
				nodeNum = (iteration * 4) + (nodeNum+1)%4
				if partition == ((partitions/(numNodes/4))*(iteration+1))-1 {
					break
				}
				partition = (partition + 1) % ((partitions / (numNodes / 4)) * (iteration + 1))
			}
		}
	}
}
