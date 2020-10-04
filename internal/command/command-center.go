package command

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials/ibmiam"
	"github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/ibm-cos-sdk-go/service/s3/s3manager"

	"github.com/IBM/db2ctl/internal/bash"
	"github.com/IBM/db2ctl/internal/config"
	"github.com/IBM/db2ctl/internal/csv"
	"github.com/IBM/db2ctl/internal/flag"
	"github.com/IBM/db2ctl/internal/template"
	"github.com/IBM/db2ctl/internal/yaml"
	"github.com/IBM/db2ctl/statik"

	"github.com/spf13/pflag"
)

var logDirPathFromEnv string  //This will be set through the build command, see Makefile
var stateDBPathFromEnv string //This will be set through the build command, see Makefile

//constants needed
const (
	SampleConfigFileName = "db2ctl.yaml"

	defaultConfigFileName  = "db2ctl-defaults.yaml"
	stateFileDefaultName   = "db2ctl-state.db"
	mappingFileDefaultName = "mapping.csv"
	binPackingFilename     = "binpacking.csv"

	generatedDir    = "generated"
	templateDir     = "/templates"
	templateFileExt = ".tmpl"
)

func init() {
	if _, err := os.Stat(generatedDir); os.IsNotExist(err) {
		// fmt.Printf("\n%v directory does not exist, creating ...\n\n", generatedDir)
		err := os.Mkdir(generatedDir, 0755)
		if err != nil {
			log.Fatalf("init error while creating %v directory, cannot proceed. err: %v", generatedDir, err)
		}
	}
}

//New creates a new instance for command execution
func New(flags *pflag.FlagSet) *Instance {
	stateFilePath := "./"
	if stateDBPathFromEnv != "" {
		stateDBPathFromEnv = parsePath(stateDBPathFromEnv)
		stateFilePath = stateDBPathFromEnv
	}

	logDir := "logs/"
	if logDirPathFromEnv != "" {
		logDirPathFromEnv = parsePath(logDirPathFromEnv)
		logDir = logDirPathFromEnv
	}

	bashInstance := bash.Instance{
		LogDir:          logDir,
		GeneratedDir:    generatedDir,
		TemplateDir:     templateDir,
		DryRunEnabled:   getBoolFlagValue(flags, flag.DryRun),
		ReRun:           getBoolFlagValue(flags, flag.ReRun),
		TimeoutInterval: time.Hour * 5, //change later
		State: bash.State{
			StateFilePath:        stateFilePath,
			StateFileDefaultname: stateFileDefaultName,
		},
	}
	bashInstance.Init()

	return &Instance{
		CombinedConfig:   &config.Combined{},
		Flags:            flags,
		Instance:         bashInstance,
		StartTime:        time.Now(),
		S3ConfigInstance: &config.S3ConfigStruct{},
	}
}

//CreateSampleConfigFile creates sample config file
func (i *Instance) CreateSampleConfigFile() *Instance {
	if i.Error != nil {
		return i
	}
	err := yaml.CreateSampleConfigFile(SampleConfigFileName)
	if err != nil {
		i.Error = err
	}
	fmt.Println("\nGenerated sample file : ", SampleConfigFileName)
	return i
}

//ParseYaml parses yaml and puts configuration into config struct
func (i *Instance) ParseYaml(confFile string) *Instance {
	if i.Error != nil {
		return i
	}
	err := yaml.Parse(i.CombinedConfig, confFile, defaultConfigFileName)
	if err != nil {
		i.Error = fmt.Errorf("error while parsing YAML file: %v", err)
		return i
	}

	if getBoolFlagValue(i.Flags, flag.Verbose) {
		i.CombinedConfig.PCConfType.Print()
	}
	fmt.Printf("\nConfiguration is valid in file : %v\n", confFile)
	return i
}

//CreateFromConfig creates yaml config
func (i *Instance) CreateFromConfig(confFile string) *Instance {
	if i.Error != nil {
		return i
	}
	err := yaml.CreateFromConfig(i.CombinedConfig, confFile)
	if err != nil {
		i.Error = err
		return i
	}
	return i
}

//GenerateBinPackingCSV generates bin packing csv
func (i *Instance) GenerateBinPackingCSV() *Instance {
	useCustomBinPacking := getBoolFlagValue(i.Flags, flag.CustomBinPacking)
	if i.Error != nil {
		return i
	}
	if !useCustomBinPacking {
		file, err := csv.GenerateBinPackingCsv(
			i.CombinedConfig,
			generatedDir,
			binPackingFilename,
		)
		if err != nil {
			i.Error = fmt.Errorf("error while generating binpacking file: %v", err)
			return i
		}
		if getBoolFlagValue(i.Flags, flag.Verbose) {
			fmt.Println("\n\nBin packing generated successfully : ", file)
		}
	}
	return i
}

//ReadBinPackingCSV reads binpacking csv
func (i *Instance) ReadBinPackingCSV() *Instance {
	if i.Error != nil {
		return i
	}
	err := csv.ReadBinPackingCsv(i.CombinedConfig, generatedDir, binPackingFilename)
	if err != nil {
		i.Error = fmt.Errorf("error while parsing binpacking file: %v", err)
		return i
	}

	if getBoolFlagValue(i.Flags, flag.Verbose) {
		fmt.Println("\n\nBin packing parsed successfully")
	}
	return i
}

//GenerateMappingFile generates mapping file
func (i *Instance) GenerateMappingFile() *Instance {
	useCustomMapping := getBoolFlagValue(i.Flags, flag.CustomMap)
	if i.Error != nil {
		return i
	}
	if !useCustomMapping {
		file, csvContent, err := csv.GenerateMappingFile(
			i.CombinedConfig,
			generatedDir,
			mappingFileDefaultName,
		)
		if err != nil {
			i.Error = fmt.Errorf("error generating mapping file : %v", err)
			return i
		}

		if getBoolFlagValue(i.Flags, flag.Verbose) {
			fmt.Println("\n\nMapping file generated successfully: ", file)
			fmt.Println(csvContent)
		}
	}
	return i
}

//ReadMappingCSV reads mapping csv file
func (i *Instance) ReadMappingCSV() *Instance {
	if i.Error != nil {
		return i
	}
	err := csv.ReadMappingCsv(i.CombinedConfig, generatedDir, mappingFileDefaultName)
	if err != nil {
		i.Error = fmt.Errorf("error while parsing mapping file: %v", err)
		return i
	}

	if getBoolFlagValue(i.Flags, flag.Verbose) {
		i.CombinedConfig.BinPacking.Print()
	}
	return i
}

//GenerateConfigFilesFromDir generates all config files
func (i *Instance) GenerateConfigFilesFromDir(dirToGenerateFrom string) *Instance {
	noGenerate := getBoolFlagValue(i.Flags, flag.NoGenerate)
	if i.Error != nil {
		return i
	}
	var configDir string
	var err error
	if dirToGenerateFrom != "" {
		configDir, err = statik.GetActualDirName(dirToGenerateFrom, templateDir)
		if err != nil {
			i.Error = fmt.Errorf("could not get ActualDirName for dir %v, err : %v ", dirToGenerateFrom, err)
			return i
		}
		if configDir == "" {
			i.Error = fmt.Errorf("could not find directory or directory is empty %v", dirToGenerateFrom)
			return i
		}
	}
	//fmt.Println("actual dir : ", configDir)
	i.ConfigDir = configDir
	if !noGenerate {
		//cleaning up all scripts in dir if it exists
		if _, err := os.Stat(generatedDir + configDir); !os.IsNotExist(err) {
			filesDeleted, err := cleanupFilesInDir(generatedDir+configDir, ".sh")
			if err != nil {
				i.Error = fmt.Errorf("could not delete files in %v directory, err: %v", generatedDir+configDir, err)
				return i
			}
			if getBoolFlagValue(i.Flags, flag.Verbose) {
				fmt.Printf("\n%v directory exists, cleaned up %v files inside...\n\n", generatedDir+configDir, filesDeleted)
			}
		}
		err := template.Generate(i.CombinedConfig,
			configDir,
			templateDir,
			templateFileExt,
			generatedDir)
		if err != nil {
			i.Error = fmt.Errorf("error while creating configuration : %v", err)
			return i
		}
	}
	i.PrintSeparator()
	return i
}

//RunBashScripts runs all bash scripts in a directory
func (i *Instance) RunBashScripts() *Instance {
	if i.Error != nil {
		return i
	}
	fullPath := generatedDir + i.ConfigDir
	// fmt.Println("fullPath : ", fullPath)
	if i.DryRunEnabled {
		i.RunScriptsInDir(fullPath)
	} else {
		i.DryRunEnabled = true
		i.RunScriptsInDir(fullPath)
		i.DryRunEnabled = false
		i.RunScriptsInDir(fullPath)
	}
	i.Error = i.Instance.Error
	i.PrintSeparator()
	return i
}

//TimeTaken prints time taken for execution
func (i *Instance) TimeTaken() *Instance {
	if i.Error != nil {
		return i
	}
	fmt.Println("Time taken : ", time.Since(i.StartTime))
	return i
}

//DeleteStateForDir deletes state for given dir
func (i *Instance) DeleteStateForDir(directory string) *Instance {
	if i.Error != nil {
		return i
	}
	i.DeleteState(directory)
	return i
}

//ReturnStateForDir prints state for given dir
func (i *Instance) ReturnStateForDir(directory string) *Instance {
	if i.Error != nil {
		return i
	}
	i.ReturnState(directory)
	return i
}

//PrintStateForDir prints state for given dir
func (i *Instance) PrintStateForDir(directory string) *Instance {
	if i.Error != nil {
		return i
	}
	i.PrintState(directory)
	return i
}

//StopRunningCommand stops currently running command
func (i *Instance) StopRunningCommand() *Instance {
	if i.Error != nil {
		return i
	}
	i.StopRunningCmd()
	i.Error = i.Instance.Error
	return i
}

//GetAllDirsInsideTmpl gets all directories inside template folder
func GetAllDirsInsideTmpl() ([]string, error) {
	dirs, err := statik.GetAllDirsInDir(templateDir)
	if err != nil {
		return nil, err
	}
	// fmt.Println("dirs : ", dirs)
	return dirs, nil
}

//returns number of files cleaned up, along with error (if nil)
func cleanupFilesInDir(directory, fileExt string) (int, error) {
	filesDeleted := 0
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == fileExt {
			err = os.Remove(path)
			if err != nil {
				return fmt.Errorf("could not delete file %v, err : %v", info.Name(), err)
			}
			filesDeleted++
		}
		return nil
	})
	return filesDeleted, err
}

func getBoolFlagValue(flags *pflag.FlagSet, flagname string) bool {
	if value, err := flags.GetBool(flagname); err == nil {
		return value
	}
	return false
}

func getStringFlagValue(flags *pflag.FlagSet, flagname string) string {
	if value, err := flags.GetString(flagname); err == nil {
		return value
	}
	return ""
}

func parsePath(path string) string {
	lastChar := path[len(path)-1:]

	if lastChar != "/" {
		path += "/"
	}
	return path
}

// Log file to display part size and elapsed time
func Log(file string, msg string) {
	var logpath = file
	f, err := os.OpenFile(logpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println(err)
		return
	}

	defer f.Close()

	//set output of logs to f
	log.SetOutput(f)

	log.Println(msg)
}

// ParseS3Config parse S3 Config file
func (i *Instance) ParseS3Config(confFile string) *Instance {
	if i.Error != nil {
		return i
	}
	err := yaml.ParseS3(i.S3ConfigInstance, confFile)
	if err != nil {
		i.Error = fmt.Errorf("error while parsing YAML file: %v", err)
		return i
	}
	return i
}

// UploadToS3 uploads file to the IBM S3 Storage
func (i *Instance) UploadToS3() *Instance {

	Log(i.S3ConfigInstance.Spec.UploadFile.Log, "\n")
	Log(i.S3ConfigInstance.Spec.UploadFile.Log, "-----------------------------------------------------------")
	Log(i.S3ConfigInstance.Spec.UploadFile.Log, "\n")
	// Start time
	start := time.Now()

	// Fetch config data from yaml file
	name := i.S3ConfigInstance.Spec.UploadFile.Name
	serviceEndpoint := i.S3ConfigInstance.Spec.S3.ServiceEndpoint
	authEndpoint := i.S3ConfigInstance.Spec.S3.AuthEndpoint
	apiKey := i.S3ConfigInstance.Spec.S3.APIKey
	serviceInstanceID := i.S3ConfigInstance.Spec.S3.ServiceInstanceID
	bucketLocation := i.S3ConfigInstance.Spec.S3.BucketLocation
	bucketName := i.S3ConfigInstance.Spec.S3.BucketName
	keyName := i.S3ConfigInstance.Spec.UploadFile.KeyName

	partSize, err := strconv.Atoi(i.S3ConfigInstance.Spec.S3.PartSize)
	if err != nil {
		fmt.Println("Conversion failed")
	}

	// Open file to be uploaded
	file, err := os.Open(name)
	if err != nil {
		i.Error = fmt.Errorf("unable to open file , err : %v", err)
		return i
	}

	// Get the file size
	fi, err := file.Stat()
	if err != nil {
		fmt.Println("Cannot obtain file stats")
	}
	fileSize := fi.Size()

	// Create config
	conf := aws.NewConfig().
		WithRegion(bucketLocation).
		WithEndpoint(serviceEndpoint).
		WithCredentials(ibmiam.NewStaticCredentials(aws.NewConfig(), authEndpoint, apiKey, serviceInstanceID)).
		WithS3ForcePathStyle(true)

	// Create session
	sess, err := session.NewSession(conf)
	if err != nil {
		i.Error = fmt.Errorf("unable to create session , err : %v", err)
		return i
	}

	// uploader := s3manager.NewUploader(sess)
	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.PartSize = int64(partSize) // Min chunk size - 5MB per part
		chunks := fileSize / u.PartSize
		Log(i.S3ConfigInstance.Spec.UploadFile.Log, "MaxUploadParts: "+strconv.Itoa(u.MaxUploadParts))
		Log(i.S3ConfigInstance.Spec.UploadFile.Log, "PartSize: "+strconv.FormatInt(u.PartSize/(1024*1024), 10)+"MB")
		Log(i.S3ConfigInstance.Spec.UploadFile.Log, "Upload file size: "+strconv.FormatInt(fileSize/(1024*1024*1024), 10)+"GB")
		Log(i.S3ConfigInstance.Spec.UploadFile.Log, "No of chunks: "+strconv.FormatInt(chunks, 10))
	})

	// upload input paramaeters
	upParams := &s3manager.UploadInput{
		Bucket: &bucketName,
		Key:    &keyName,
		Body:   file,
	}

	// call Upload() function
	_, err = uploader.Upload(upParams)
	if err != nil {
		i.Error = fmt.Errorf("unable to upload file , err : %v", err)
		return i
	}

	// Total elapsed time
	elapsed := time.Since(start)
	fmt.Println("Elapsed: " + elapsed.String())

	Log(i.S3ConfigInstance.Spec.UploadFile.Log, "Elapsed time: "+elapsed.String())
	Log(i.S3ConfigInstance.Spec.UploadFile.Log, "\n")

	return i
}

// DownloadFromS3 downloads file from the IBM S3 Storage
func (i *Instance) DownloadFromS3() *Instance {

	Log(i.S3ConfigInstance.Spec.DownloadFile.Log, "\n")
	Log(i.S3ConfigInstance.Spec.DownloadFile.Log, "-----------------------------------------------------------")
	Log(i.S3ConfigInstance.Spec.DownloadFile.Log, "\n")

	// Download Start time
	start := time.Now()

	// fetch data from yaml file
	downloadFileName := i.S3ConfigInstance.Spec.DownloadFile.Name
	serviceEndpoint := i.S3ConfigInstance.Spec.S3.ServiceEndpoint
	authEndpoint := i.S3ConfigInstance.Spec.S3.AuthEndpoint
	apiKey := i.S3ConfigInstance.Spec.S3.APIKey
	serviceInstanceID := i.S3ConfigInstance.Spec.S3.ServiceInstanceID
	bucketLocation := i.S3ConfigInstance.Spec.S3.BucketLocation
	bucketName := i.S3ConfigInstance.Spec.S3.BucketName
	prefix := i.S3ConfigInstance.Spec.DownloadFile.Prefix

	partSize, err := strconv.Atoi(i.S3ConfigInstance.Spec.S3.PartSize)
	if err != nil {
		fmt.Println("Conversion failed")
	}

	// Create Config
	conf := aws.NewConfig().
		WithRegion(bucketLocation).
		WithEndpoint(serviceEndpoint).
		WithCredentials(ibmiam.NewStaticCredentials(aws.NewConfig(), authEndpoint, apiKey, serviceInstanceID)).
		WithS3ForcePathStyle(true)

	// Create Session
	sess, err := session.NewSession(conf)
	client := s3.New(sess, conf)
	if err != nil {
		i.Error = fmt.Errorf("unable to create session , err : %v", err)
		return i
	}

	downloader := s3manager.NewDownloader(sess, func(u *s3manager.Downloader) {
		u.PartSize = int64(partSize) // Min chunk size - 5MB per part
		Log(i.S3ConfigInstance.Spec.DownloadFile.Log, "PartSize: "+strconv.FormatInt(u.PartSize/(1024*1024), 10)+"MB")
	})

	// create  a directory to download files from S3 bucket
	err = os.Mkdir(prefix, 0755)
	if err != nil {
		fmt.Println(err)
	}

	// List all objects in the bucket
	listParams := s3.ListObjectsInput{
		Bucket: &bucketName,
	}
	buckets, err := client.ListObjects(&listParams)
	contents := buckets.Contents
	for i := 0; i < len(contents); i++ {
		key := contents[i].Key
		if strings.Contains(*key, prefix) == true {
			keyName := *key
			// Download Parameters
			getParams := s3.GetObjectInput{
				Bucket: &bucketName,
				Key:    &keyName,
			}
			// Create file for download
			name := downloadFileName + "_" + strings.Split(keyName, "/")[1]
			file, err := os.Create(filepath.Join(prefix, name))
			if err != nil {
				fmt.Println("Error creating file")
				fmt.Println(err)
			}
			// Download file from s3
			_, err = downloader.Download(file, &getParams)
			if err != nil {
				fmt.Println("Error downloading file")
				fmt.Println(err)
			}
		}
	}

	// Total elapsed time
	elapsed := time.Since(start)
	fmt.Println("Elapsed: " + elapsed.String())
	Log(i.S3ConfigInstance.Spec.DownloadFile.Log, "Elapsed time: "+elapsed.String())
	Log(i.S3ConfigInstance.Spec.DownloadFile.Log, "\n")

	return i
}
