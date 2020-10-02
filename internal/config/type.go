package config

//PCConfType is the main config file needed
type PCConfType struct {
	APIVersion string `yaml:"apiVersion" validate:"required" json:"apiVersion"`
	Kind       string `yaml:"kind" validate:"required" json:"kind"`
	Metadata   struct {
		Name        string `yaml:"name" json:"name"`
		ClusterName string `yaml:"clusterName" json:"clusterName"`
	} `yaml:"metadata" json:"metadata"`
	Spec struct {
		DB2 struct {
			Required struct {
				Role              string `yaml:"role" validate:"oneof=optimized sandbox"`
				Replicated        bool   `yaml:"replicated"`
				DatabaseName      string `yaml:"databaseName" validate:"required"`
				InstancePort      int    `yaml:"instancePort" validate:"required"`
				InstanceName      string `yaml:"instanceName" validate:"required"`
				InstanceID        int    `yaml:"instanceID" validate:"required"`
				InstanceGD        int    `yaml:"instanceGD" validate:"required"`
				FencedID          int    `yaml:"fencedID" validate:"required"`
				FencedGD          int    `yaml:"fencedGD" validate:"required"`
				InstanceSecret    string `yaml:"instanceSecret" validate:"required"`
				InstanceSecretVal string //calculated dynamically, base64 decoded
				FencedName        string `yaml:"fencedName" validate:"required"`
				FencedSecret      string `yaml:"fencedSecret" validate:"required"`
				FencedSecretVal   string //calculated dynamically, base64 decoded
				DB2Version        string `yaml:"db2Version" validate:"required"`
				DB2Binary         string `yaml:"db2Binary" validate:"required"`
				DB2License        string `yaml:"db2License" validate:"required"`
				NumInstances      int    `yaml:"numInstances" validate:"required"`
				NumDB             int    `yaml:"numdb" validate:"required"`
				Organization      string `yaml:"organization" validate:"oneof=row column"`
			} `yaml:"required"`
			Optional struct {
				TopLevelDir string `yaml:"topLevelDir" json:"topLevelDir"`
			} `yaml:"optional" json:"optional"`
		} `yaml:"db2" json:"db2"`
		Nodes struct {
			Required struct {
				NumNodes          int                   `yaml:"numNodes" validate:"node-check,required"`
				NVMEList          [][]map[string]string `yaml:"nvmeList" validate:"nvme-check,gt=0,required" json:"NVMEList"`
				IPAddresses       []string              `yaml:"ipAddresses" validate:"ip-check,gt=0,required,dive,ip" json:"ipAdresses"`
				IPAddressesAsList string                `json:"-"` //generated dynamically
				Names             []string              `yaml:"names" validate:"names-check,gt=0,required" json:"names"`
				NamesAsList       string                `json:"-"` //generated dynamically
				Partitions        int                   `yaml:"partitions" validate:"required" json:"partitions"`
				SizePerPartition  map[string]string     `json:"-"` //generated dynamically
				NodeInfoMap       map[string]*nodeInfo  `json:"-"` //generated dynamically
			} `yaml:"required" json:"required"`
			Optional struct {
				DBPrimitiveNamePrefix string `yaml:"dbPrimitiveNamePrefix" json:"dbPrimitiveNamePrefix"`
			} `yaml:"optional" json:"-"`
		} `yaml:"nodes" json:"nodes"`
		Linbit struct {
			Required struct {
			} `yaml:"required" json:"-"`
			Optional struct {
				NumStripes       int    `yaml:"numStripes" json:"numStripes"`
				StripeSize       string `yaml:"stripeSize" json:"stripeSize"`
				VolumeDefinition struct {
					Nodes struct {
						Size  string `yaml:"size" json:"size"`
						Minor string `json:"-"` //generated dynamically
					} `yaml:"nodes" json:"nodes"`
					NFS struct {
						Size          string `yaml:"size" json:"size"`
						NumNFSServers int    `yaml:"numNFSServers" json:"numNFSServers"`
						Minor         string `json:"-"` //generated dynamically
					} `yaml:"nfs" json:"nfs"`
					DB2Local struct {
						Size string `yaml:"size" json:"size"`
					} `yaml:"db2local" json:"db2local"`
				} `yaml:"volumeDefinition" json:"volumeDefinition"`
			} `yaml:"optional" json:"optional"`
		} `yaml:"linbit" json:"linbit"`
		NFS struct {
			Server struct {
				Required struct {
					VirtualIP         string   `yaml:"nfsVirtualIP" validate:"required,ip" json:"nfsVirtualIP"`
					CIDRNetMask       string   `yaml:"cidrNetmask" validate:"required" json:"cidrNetmask"`
					NIC               string   `yaml:"nic" validate:"required" json:"nic"`
					ClientSpec        string   `json:"-"` //calculated dynamically
					NodesForPlacement []string `json:"-"` //calculated dynamically
				} `yaml:"required" json:"required"`
				Optional struct {
					Path                    string `yaml:"nfsPath" json:"nfsPath"`
					PrimitiveName           string `yaml:"primitiveName" json:"primitiveName"`
					VirtualIPResourceName   string `yaml:"virtualIPResourceName" json:"virtualIPResourceName"`
					ExportResourceName      string `yaml:"exportResourceName" json:"exportResourceName"`
					ExportResourceDirectory string `yaml:"exportResourceDirectory" json:"exportResourceDirectory"`
					MountpointName          string `yaml:"mountpointName" json:"mountpointName"`
					PacemakerOrderName1     string `yaml:"pacemakerOrderName1" json:"pacemakerOrderName1"`
					PacemakerOrderName2     string `yaml:"pacemakerOrderName2" json:"pacemakerOrderName2"`
					PacemakerOrderName3     string `yaml:"pacemakerOrderName3" json:"pacemakerOrderName3"`
					ColocationName1         string `yaml:"colocationName1" json:"colocationName1"`
					ColocationName2         string `yaml:"colocationName2" json:"colocationName2"`
					ColocationName3         string `yaml:"colocationName3" json:"colocationName3"`
				} `yaml:"optional" json:"-"`
			} `yaml:"server" json:"server"`
			Client struct {
				Optional struct {
					MountPoint    string `yaml:"mountPoint" json:"mountPoint"`
					PrimitiveName string `yaml:"primitiveName" json:"primitiveName"`
					CloneName     string `yaml:"cloneName" json:"cloneName"`
					OrderName     string `yaml:"orderName" json:"orderName"`
				} `yaml:"optional" json:"-"`
			} `yaml:"client" json:"-"`
		} `yaml:"nfs" json:"nfs"`
	} `yaml:"spec" json:"spec"`
}

//BinPacking struct
type BinPacking struct {
	DataToNode []DataToNode  `json:"-"`
	Mapping    []DBMapStruct `json:"-"`
}

type nodeInfo struct {
	NVMEList     string
	Size         int
	Order        string
	HasNfsServer bool
	DB2LocalSize int
}

//DataToNode struct holds the relationship between a disk and its nodes
type DataToNode struct {
	DBPrimitiveName string
	PrimaryServer   string
	ReplicaServer   string
}

//DBMapStruct struct holds the mapping between mount point and device names
type DBMapStruct struct {
	DBMountPoint string `csv:"DB Mount Point"`
	DBDeviceName string `csv:"DB Device Name"`
}

//Combined is a combination of the 2 main structs
type Combined struct {
	PCConfType
	BinPacking
}

// S3Config struct for S3 config yaml file
type S3Config struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string
	} `yaml:"metadata"`
	Spec struct {
		UploadFile struct {
			KeyName string `yaml:"keyName"`
			Name    string `yaml:"name"`
			Log     string `yaml:"log"`
		} `yaml:"uploadFile"`
		S3 struct {
			APIKey            string `yaml:"apiKey"`
			ServiceInstanceID string `yaml:"serviceInstanceID"`
			AuthEndpoint      string `yaml:"authEndpoint"`
			ServiceEndpoint   string `yaml:"serviceEndpoint"`
			BucketLocation    string `yaml:"bucketLocation"`
			BucketName        string `yaml:"bucketName"`
			PartSize          string `yaml:"partSize"`
		} `yaml:"s3"`
		DownloadFile struct {
			Prefix string `yaml:"prefix"`
			Name   string `yaml:"name"`
			Log    string `yaml:"log"`
		} `yaml:"downloadFile"`
	} `yaml:"spec"`
}

// S3ConfigStruct struct
type S3ConfigStruct struct {
	S3Config
}
