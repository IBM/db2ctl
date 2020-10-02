package config

import (
	"strconv"
	"testing"
)

type mem struct {
	sizeStr string
	size    int
	order   string
}

type sizeCal struct {
	numNodes          int
	numNfsServers     int
	partitions        int
	replicated        bool
	size              int
	nfsSize           string
	db2LocalSizeWOnfs int
	db2LocalSizeWnfs  int
	nodeSizeWOnfs     string
	nodeSizeWnfs      string
	shouldError       bool
}

func TestNFSServerPlacement(t *testing.T) {

	//input 1
	numNFSServers := 4
	numNodes := 4
	var expected = []string{"server0", "server1", "server2", "server3"}
	nfsTest(t, "input1", numNFSServers, numNodes, expected)

	//input 2
	numNFSServers = 4
	numNodes = 8
	expected = []string{"server0", "server4", "server1", "server5"}
	nfsTest(t, "input2", numNFSServers, numNodes, expected)

	//input 3
	numNFSServers = 4
	numNodes = 12
	expected = []string{"server0", "server4", "server8", "server1"}
	nfsTest(t, "input3", numNFSServers, numNodes, expected)

	//input 4
	numNFSServers = 4
	numNodes = 16
	expected = []string{"server0", "server4", "server8", "server12"}
	nfsTest(t, "input4", numNFSServers, numNodes, expected)

	//input 5
	numNFSServers = 5
	numNodes = 12
	expected = []string{"server0", "server4", "server8", "server1", "server5"}
	nfsTest(t, "input5", numNFSServers, numNodes, expected)

	//input 6
	numNFSServers = 5
	numNodes = 16
	expected = []string{"server0", "server4", "server8", "server12", "server1"}
	nfsTest(t, "input6", numNFSServers, numNodes, expected)

	//input 7
	numNFSServers = 6
	numNodes = 24
	expected = []string{"server0", "server4", "server8", "server12", "server16", "server20"}
	nfsTest(t, "input7", numNFSServers, numNodes, expected)

	//input 8
	numNFSServers = 8
	numNodes = 16
	expected = []string{"server0", "server4", "server8", "server12", "server1", "server5", "server9", "server13"}
	nfsTest(t, "input8", numNFSServers, numNodes, expected)

	//input 9
	numNFSServers = 8
	numNodes = 24
	expected = []string{"server0", "server4", "server8", "server12", "server16", "server20", "server1", "server5"}
	nfsTest(t, "input9", numNFSServers, numNodes, expected)
}

func TestNodeInfoMapping(t *testing.T) {
	config := &PCConfType{}

	var nodeNames []string
	for i := 0; i < 2; i++ {
		nodeNames = append(nodeNames, "server"+strconv.Itoa(i))
	}

	config.Spec.Nodes.Required.Names = nodeNames

	m1 := make(map[string]string)
	m2 := make(map[string]string)
	m1["name"] = "/dev/xvdc"
	m1["size"] = "25.5GB"

	m2["name"] = "/dev/xvde"
	m2["size"] = "25GB"

	m3 := make(map[string]string)
	m4 := make(map[string]string)
	m3["name"] = "/dev/xvdc"
	m3["size"] = "1TB"

	m4["name"] = "/dev/xvde"
	m4["size"] = "2TB"

	nvmeL := [][]map[string]string{{m1, m2}, {m3, m4}}
	config.Spec.Nodes.Required.NVMEList = nvmeL
	err := nodeInfoMapping(config)
	if err != nil {
		t.Errorf("error while creating map : %v", err)
	}

	nodeInfoMap := config.Spec.Nodes.Required.NodeInfoMap

	if nodeInfoMap["server0"].Size != 50 || nodeInfoMap["server0"].Order != "GB" {
		t.Errorf("values do not match, expected %v got %v", 50, nodeInfoMap["server0"].Size)
	}
	if nodeInfoMap["server1"].Size != 3 || nodeInfoMap["server1"].Order != "TB" {
		t.Errorf("values do not match, expected %v got %v", 3, nodeInfoMap["server0"].Size)
	}
}

func TestSizingCalculations(t *testing.T) {

	var sizeCalList = []sizeCal{
		{
			numNodes:         4,
			numNfsServers:    4,
			partitions:       24,
			replicated:       false,
			size:             100,
			nfsSize:          "10GB",
			db2LocalSizeWnfs: 90,
		},
		{
			numNodes:          8,
			numNfsServers:     4,
			partitions:        24,
			replicated:        false,
			size:              100,
			nfsSize:           "10GB",
			db2LocalSizeWnfs:  90,
			db2LocalSizeWOnfs: 100,
		},
		{
			numNodes:         4,
			numNfsServers:    4,
			partitions:       24,
			replicated:       true,
			size:             100,
			nfsSize:          "10GB",
			db2LocalSizeWnfs: 18,
			nodeSizeWOnfs:    "6GB",
		},
		{
			numNodes:          8,
			numNfsServers:     4,
			partitions:        24,
			replicated:        true,
			size:              100,
			nfsSize:           "10GB",
			db2LocalSizeWOnfs: 20,
			db2LocalSizeWnfs:  18,
			nodeSizeWOnfs:     "13GB",
			nodeSizeWnfs:      "12GB",
		},
		{
			numNodes:          8,
			numNfsServers:     4,
			partitions:        24,
			replicated:        true,
			size:              50,
			nfsSize:           "10GB",
			db2LocalSizeWOnfs: 10,
			db2LocalSizeWnfs:  8,
			nodeSizeWOnfs:     "6GB",
			nodeSizeWnfs:      "5GB",
		},
		{
			numNodes:          8,
			numNfsServers:     4,
			partitions:        24,
			replicated:        true,
			size:              25,
			nfsSize:           "10GB",
			db2LocalSizeWOnfs: 5,
			db2LocalSizeWnfs:  3,
			nodeSizeWOnfs:     "3GB",
			nodeSizeWnfs:      "2GB",
		},
		{
			numNodes:          8,
			numNfsServers:     4,
			partitions:        24,
			replicated:        true,
			size:              10,
			nfsSize:           "10GB",
			db2LocalSizeWOnfs: 5,
			db2LocalSizeWnfs:  3,
			nodeSizeWOnfs:     "3GB",
			nodeSizeWnfs:      "2GB",
			shouldError:       true,
		},
	}

	for num, sizeCalInstance := range sizeCalList {
		sizeTest(t, "sizing calculation test-"+strconv.Itoa(num+1), sizeCalInstance)
	}
}

func TestSplitMemorySize(t *testing.T) {

	memAll := []mem{
		{
			sizeStr: "20GB",
			size:    20,
			order:   "GB",
		},
		{
			sizeStr: "20TB",
			size:    20,
			order:   "TB",
		},
		{
			sizeStr: "20.5GB",
			size:    20,
			order:   "GB",
		},
		{
			sizeStr: "20gb",
			size:    20,
			order:   "gb",
		},
		{
			sizeStr: "20Gb",
			size:    20,
			order:   "Gb",
		},
	}

	for _, memInstance := range memAll {
		memTest(t, memInstance)
	}
}

func memTest(t *testing.T, memInstance mem) {
	size, order, err := splitMemorySize(memInstance.sizeStr)

	if size != memInstance.size || order != memInstance.order || err != nil {
		t.Errorf("not matching for %v, got %v %v %v", memInstance, size, order, err)
	}
}

func nfsTest(t *testing.T, testName string, numNFSServers, numNodes int, expected []string) {

	t.Run(testName, func(t *testing.T) {
		config := &PCConfType{}
		var nodeNames []string
		for i := 0; i < numNodes; i++ {
			nodeNames = append(nodeNames, "server"+strconv.Itoa(i))
		}

		config.Spec.Linbit.Optional.VolumeDefinition.NFS.NumNFSServers = numNFSServers
		config.Spec.Nodes.Required.NumNodes = numNodes
		config.Spec.Nodes.Required.Names = nodeNames

		nodeInfoMapping(config)
		nfsServerPlacement(config)

		output := config.Spec.NFS.Server.Required.NodesForPlacement
		result := true
		for i := range expected {
			if expected[i] != output[i] {
				//t.Log("match : ", expected[i], output[i])
				result = false
			}
		}
		if !result {
			t.Errorf("does not match with numNodes %v, numServers %v, wanting %v, got %v", numNodes, numNFSServers, expected, output)
		}
	})

}

func sizeTest(t *testing.T, testName string, sizeCalInstance sizeCal) {

	t.Run(testName, func(t *testing.T) {
		config := &PCConfType{}
		var nodeNames []string

		nodeInfoMap := make(map[string]*nodeInfo)
		for i := 0; i < sizeCalInstance.numNodes; i++ {
			nodeName := "server" + strconv.Itoa(i)
			nodeNames = append(nodeNames, nodeName)
			nodeInfoMap[nodeName] = &nodeInfo{
				Size:  sizeCalInstance.size,
				Order: "GB",
			}
		}

		config.Spec.Nodes.Required.NumNodes = sizeCalInstance.numNodes
		config.Spec.Linbit.Optional.VolumeDefinition.NFS.NumNFSServers = sizeCalInstance.numNfsServers
		config.Spec.Nodes.Required.Partitions = sizeCalInstance.partitions
		config.Spec.Nodes.Required.Names = nodeNames
		config.Spec.Nodes.Required.NodeInfoMap = nodeInfoMap
		config.Spec.DB2.Required.Replicated = sizeCalInstance.replicated

		nfsServerPlacement(config)

		err := sizingCalculations(config)
		if err != nil {
			if sizeCalInstance.shouldError {
				return
			}
			t.Errorf("error: %v", err)
			return
		}

		// config.Print()

		if config.Spec.Linbit.Optional.VolumeDefinition.NFS.Size != sizeCalInstance.nfsSize {
			t.Errorf("not matching for NFS Size, want %v got %v", sizeCalInstance.nfsSize, config.Spec.Linbit.Optional.VolumeDefinition.NFS.Size)
		}

		for _, nodeInfo := range config.Spec.Nodes.Required.NodeInfoMap {
			if nodeInfo.HasNfsServer {
				if nodeInfo.DB2LocalSize != sizeCalInstance.db2LocalSizeWnfs {
					t.Errorf("not matching for Db2local size w nfs, want %v got %v", sizeCalInstance.db2LocalSizeWnfs, nodeInfo.DB2LocalSize)
				}
			} else {
				if nodeInfo.DB2LocalSize != sizeCalInstance.db2LocalSizeWOnfs {
					t.Errorf("not matching for Db2local size wo nfs, want %v got %v", sizeCalInstance.db2LocalSizeWOnfs, nodeInfo.DB2LocalSize)
				}
			}
		}
		for _, size := range config.Spec.Nodes.Required.SizePerPartition {
			if size != sizeCalInstance.nodeSizeWOnfs && size != sizeCalInstance.nodeSizeWnfs {
				t.Errorf("not matching for node size, want %v got %v", sizeCalInstance.nodeSizeWOnfs, size)
			}
		}
	})

}
