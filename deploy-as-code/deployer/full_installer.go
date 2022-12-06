package main

import (
	"bytes"
	"container/list"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jcelliott/lumber"
	"github.com/manifoldco/promptui"
	"golang.org/x/crypto/ssh"
	yaml "gopkg.in/yaml.v3"

	//"bufio"
	"deployer/configs"
	"encoding/json"
)

var cloudTemplate string // Which terraform template to choose
var repoDirRoot string
var selectedMod []string
var CloudProvider string
var InfraType string
var db_pswd string
var sshFile string
var cluster_name string

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

//Defining a struct to parse the yaml file
type Digit struct {
	Version string `yaml:"version"`
	Modules []struct {
		Name         string   `yaml:"name"`
		Services     []string `yaml:"services"`
		Dependencies []string `yaml:"dependencies,omitempty"`
	} `yaml:"modules"`
}

type State struct {
	Proceed         string   `yaml:"Proceed"`
	Infra           string   `yaml:"Infra"`
	Product         string   `yaml:"Product"`
	Productversion  string   `yaml:"Productversion"`
	Modules         []string `yaml:"Modules"`
	CloudType       string   `yaml:"CloudType"`
	Awsacess        string   `yaml:"Awsacess"`
	Aws_access_key  string   `yaml:"Aws_access_key"`
	Aws_secret_key  string   `yaml:"Aws_secret_key"`
	Aws_session_key string   `yaml:"Aws_session_key"`
	Aws_profile     string   `yaml:"Aws Profile"`
	Aws_command     string   `yaml:"Aws_command"`
	Clustername     string   `yaml:"Clustername"`
	Db_pass         string   `yaml:"Database password"`
	Git_command     string   `yaml:"Git command"`
	T_init          string   `yaml:"Terraform init"`
	T_plan          string   `yaml:"Terraform plan"`
	T_apply         string   `yaml:"Terraform apply"`
	Domain          string   `yaml:"Domain"`
	Branch          string   `yaml:"Config & mdms branch"`
	Config_git_url  string   `yaml:"Configs Git URL"`
	Mdms_git_url    string   `yaml:"MDMS Git URL"`
	SmsUrl          string   `yaml:"SMS URL"`
	SmsGateway      string   `yaml:"SMS Gateway"`
	SmsSender       string   `yaml:"SMS Sender"`
	SmsUsername     string   `yaml:"SMS Username"`
	Bucket          string   `yaml:"Bucket name"`
	Smsproceed      string   `yaml:"SMS"`
	Fileproceed     string   `yaml:"Filestore"`
	Botproceed      string   `yaml:"Botproceed"`
	SshCreation     string   `yaml:"SSH KEY creation"`
	Hasdomain       string   `yaml:"Hasdomain"`
	Hasgitacc       string   `yaml:"Hasgitacc"`
}

var St State

type Set struct {
	set map[string]bool
}

func NewSet() *Set {
	return &Set{make(map[string]bool)}
}
func (set *Set) Add(i string) bool {
	_, found := set.set[i]
	set.set[i] = true
	return !found //False if it existed already
}
func (set *Set) Get(i string) bool {
	_, found := set.set[i]
	return found
}

func main() {
	if _, err := os.Stat("state.yaml"); err == nil {
		state, err := ioutil.ReadFile("state.yaml")
		if err != nil {
			log.Printf("%v", err)
		}
		err = yaml.Unmarshal(state, &St)
		if err != nil {
			log.Printf("%v", err)
		}
	}
	var optedInfraType string          // Infra types supported to deploy DIGIT
	var servicesToDeploy string        // Modules to be deployed
	var number_of_worker_nodes int = 1 // No of VMs for the k8s worker nodes
	var optedCloud string              // Desired InfraType to deploy
	var cloudLoginCredentials bool     // Is there a valid cloud account and credentials
	var isProductionSetup bool = false

	infraType := []string{
		"0. You have an existing kubernetes Cluster ready, you would like to leverage it to setup DIGIT on that",
		"1. Pilot/POC (Just for a POC to Quickstart and explore)",
		"2. DevTest Setup (To setup and build/customize and test)",
		"3. Production: Bare Minimal (90% reliability), 10 gov services, 10 concurrent users/sec",
		"4. Production: Medium (95% reliability), 50+ concurrent gov services 100 concurrent users/sec",
		"5. Production: HA/DRS Setup (99.99% reliability), 50+ concurrent gov services 1000 concurrent users/sec",
		"6. For custom options, use this calcualtor to determine the required nodes (https://docs.digit.org/Infra-calculator)"}

	Platforms := []string{
		"0. Local machine/Your Existing VM",
		"1. AWS-EC2 - Quickstart with a Single EC2 Instace on AWS"}
	cloudPlatforms := []string{
		"0. On-prem/Private Cloud - Quickstart with Single VM",
		"1. AWS-EKS - Production grade Elastic Kubernetes Service (EKS)",
		"2. AZURE-AKS - Production grade Azure Kubernetes Service (AKS)",
		"3. GOOGLE CLOUD - Production grade Google Kubernetes Engine (GKE)",
		"4. On-prem/Privare Cloud - Production grade Kubernetes Cluster Setup"}

	fmt.Println(string(Green), "\n*******  Welcome to DIGIT Server setup & Deployment !!! ******** \n\n *********\n https://docs.digit.org/Infra-calculator\n")
	const sPreReq = "Pre-requsites (Please Read Carefully):\n\tDIGIT comprises of many microservices that are packaged as docker containers that can be run on any container supported platforms like dockercompose, kubernetes, etc. Here we'll have a setup a kubernetes.\nHence the following are mandatory to have it before you proceed.\n\t1. Kubernetes(K8s) Cluster.\n\t\t[Option a] Local/VM: If you do not have k8s, using this link you can create k8s cluster on your local or on a VM.\n\t\t[b] Cloud: If you have your cloud account like AWS, Azure, GCP, SDC or NIC you can follow this link to create k8s.\n\t2. Post the k8s cluster creation you should get the Kubeconfig file, which you have saved in your local machine.\n\t\n\n Well! Let's get started with the DIGIT Setup process, if you want to abort any time press (Ctl+c), you can always come back and rerun the script."
	fmt.Println(string(Cyan), sPreReq)

	preReqConfirm := []string{"Yes", "No"}
	var proceed string = ""
	if St.Proceed == "" {
		proceed, _ = sel(preReqConfirm, "Are you good to proceed?")
		St.Proceed = proceed
	} else {
		proceed = St.Proceed
	}
	if proceed == "Yes" {
		if St.Infra == "" {
			optedInfraType, _ = sel(infraType, "Select the below suitable infra option for your usecase")
			St.Infra = optedInfraType
		} else {
			optedInfraType = St.Infra
		}
		switch optedInfraType {
		case infraType[0]:
			number_of_worker_nodes = 0
		case infraType[1]:
			InfraType = "quickstart"
			number_of_worker_nodes = 1
		case infraType[2]:
			number_of_worker_nodes = 1
		case infraType[3]:
			number_of_worker_nodes = 3 //TBD
			isProductionSetup = true
		case infraType[4]:
			number_of_worker_nodes = 4 //TBD
			isProductionSetup = true
		case infraType[5]:
			number_of_worker_nodes = 5 //TBD
		case infraType[6]:
			number_of_worker_nodes, _ = strconv.Atoi(enterValue(nil, "How many VM/nodes are required based on the calculation"))
			isProductionSetup = true
		default:
			number_of_worker_nodes = 0
		}

		servicesToDeploy = selectGovServicesToInstall()
		if InfraType == "quickstart" {
			if St.CloudType == "" {
				optedCloud, _ = sel(Platforms, "Choose the Platform type to provision the required resources for the selectd gov stack services?")
				St.CloudType = optedCloud
				// fmt.Println(optedCloud)
			} else {
				optedCloud = St.CloudType
			}
		} else {
			if St.CloudType == "" {
				optedCloud, _ = sel(cloudPlatforms, "Choose the cloud type to provision the required servers for the selectdd gov stack services?")
				St.CloudType = optedCloud
				// fmt.Println(optedCloud)
			} else {
				optedCloud = St.CloudType
			}
		}
		switch optedCloud {
		case Platforms[0]:
			// TBD
		case Platforms[1]:
			var optedAccessType string
			var aws_access_key string
			var aws_secret_key string
			var aws_session_key string

			cloudTemplate = "quickstart-aws-ec2"

			accessTypes := []string{"Root Admin", "Temprory Admin", "Already configured"}
			optedAccessType, _ = sel(accessTypes, "Choose your AWS access type? eg: If your access is session based unlike root admin")

			fmt.Println("\n Great, you need to input your " + optedCloud + "credentials to provision the cloud resources ..\n")

			if optedAccessType == "Temprory Admin" {

				fmt.Println("Input the AWS access key id")
				fmt.Scanln(&aws_access_key)

				fmt.Println("\nInput the AWS secret key")
				fmt.Scanln(&aws_secret_key)

				fmt.Println("\nInput the AWS Session Token")
				fmt.Scanln(&aws_session_key)

				cloudLoginCredentials = awslogin(aws_access_key, aws_secret_key, aws_session_key, "")
			} else if optedAccessType == "Root Admin" {

				fmt.Println("Input the AWS access key id")
				fmt.Scanln(&aws_access_key)

				fmt.Println("\nInput the AWS secret key")
				fmt.Scanln(&aws_secret_key)

				cloudLoginCredentials = awslogin(aws_access_key, aws_secret_key, "", "")
			} else {
				cloudLoginCredentials = awslogin("", "", "", "")
				fmt.Println("Proceeding with the existing AWS profile configured")
			}
		case cloudPlatforms[0]:
			//TBD

		case cloudPlatforms[1]:
			var optedAccessType string
			var aws_access_key string
			var aws_secret_key string
			var aws_session_key string
			CloudProvider = "aws"
			cloudTemplate = "sample-aws"

			accessTypes := []string{"Root Admin", "Temprory Admin", "Already configured"}
			if St.Awsacess == "" {
				optedAccessType, _ = sel(accessTypes, "Choose your AWS access type? eg: If your access is session based unlike root admin")
				St.Awsacess = optedAccessType
			} else {
				optedAccessType = St.Awsacess
			}
			fmt.Println("\n Great, you need to input your " + optedCloud + "credentials to provision the cloud resources ..\n")

			if optedAccessType == "Temprory Admin" {
				if St.Aws_access_key == "" {
					fmt.Println("Input the AWS access key id")
					fmt.Scanln(&aws_access_key)
					St.Aws_access_key = aws_access_key
				} else {
					aws_access_key = St.Aws_access_key
				}
				if St.Aws_secret_key == "" {
					fmt.Println("\nInput the AWS secret key")
					fmt.Scanln(&aws_secret_key)
					St.Aws_secret_key = aws_secret_key
				} else {
					aws_secret_key = St.Aws_secret_key
				}
				if St.Aws_session_key == "" {
					fmt.Println("\nInput the AWS Session Token")
					fmt.Scanln(&aws_session_key)
					St.Aws_session_key = aws_session_key
					writeState()
				} else {
					aws_session_key = St.Aws_session_key
				}

				cloudLoginCredentials = awslogin(aws_access_key, aws_secret_key, aws_session_key, "")
			} else if optedAccessType == "Root Admin" {

				if St.Aws_access_key == "" {
					fmt.Println("Input the AWS access key id")
					fmt.Scanln(&aws_access_key)
					St.Aws_access_key = aws_access_key
				} else {
					aws_access_key = St.Aws_access_key
				}

				if St.Aws_secret_key == "" {
					fmt.Println("\nInput the AWS secret key")
					fmt.Scanln(&aws_secret_key)
					St.Aws_secret_key = aws_secret_key
				} else {
					aws_secret_key = St.Aws_secret_key
				}

				cloudLoginCredentials = awslogin(aws_access_key, aws_secret_key, "", "")
			} else {
				cloudLoginCredentials = awslogin("", "", "", "")
				fmt.Println("Proceeding with the existing AWS profile configured")
			}

		case cloudPlatforms[2]:
			cloudTemplate = "sample-azure"
			fmt.Println("\n Great, you need to input your " + optedCloud + "credentials to provision the cloud resources ..\n")
			azure_username := enterValue(nil, "Please enter your AZURE UserName")
			azure_password := enterValue(nil, "Enter your AZURE Password")
			cloudLoginCredentials = azurelogin(azure_username, azure_password)

		case cloudPlatforms[3]:
			cloudTemplate = "sample-gcp"
			fmt.Println("\n Great, you need to input your " + optedCloud + "credentials to provision the cloud resources ..\n")
			fmt.Println("Support for the " + optedCloud + "is still underway ... you need to wait")

		case cloudPlatforms[4]:
			cloudTemplate = "sample-private-cloud"
			fmt.Println("\n Great, you need to input your " + optedCloud + "credentials to provision the cloud resources ..\n")
			fmt.Println("Support for the " + optedCloud + "is still underway ... you need to wait")

		default:
			//fmt.Println("\n Great, you need to input your " + optedCloud + "credentials to provision the cloud resources ..\n")
			//fmt.Println("Support for the " + optedCloud + "is still underway ... you need to wait")
		}
	}

	if cloudLoginCredentials {
		fmt.Println(string(Green), "\n*******  Let's proceed with cluster creation, please input the requested details below *********\n")
		fmt.Println(string(Green), "Make sure that the cluster name is unique if you are trying consecutively, duplicate DNS/hosts file entry under digit.org domain could have been mapped already\n")
		if St.Clustername == "" {
			cluster_name = enterValue(nil, "How do you want to name the Cluster? eg: your-name_dev or your-name_poc")
			St.Clustername = cluster_name
		} else {
			cluster_name = St.Clustername
		}
		// fmt.Println("How do you want to name the Cluster? \n eg: your-name_dev or your-name_poc")
		// fmt.Scanln(&cluster_name)

		repoDirRoot = "DIGIT-DevOps"
		gitCmd := ""
		_, err := os.Stat(repoDirRoot)
		if os.IsNotExist(err) {
			gitCmd = fmt.Sprintf("git clone -b release https://github.com/egovernments/DIGIT-DevOps.git %s", repoDirRoot)
		} else {
			gitCmd = fmt.Sprintf("git -C %s pull", repoDirRoot)
		}
		err1 := execCommand(gitCmd)
		if err1 == nil {
			St.Git_command = "success"
		}
		if err1 != nil {
			St.Git_command = "failure"
			writeState()
			log.Printf("%v", err)
		}

		if !isProductionSetup {

			sshFile = "./digit-ssh.pem"
			var keyName string = "digit-aws-vm"
			pubKey, _, err := GetKeyPair(sshFile)
			// to pick public ip and private ip from terraform state

			if err != nil {
				log.Fatalf("Failed to generate SSH Key %s\n", err)
			} else {
				execSingleCommand(fmt.Sprintf("terraform -chdir=%s/infra-as-code/terraform/%s init", repoDirRoot, cloudTemplate))

				execSingleCommand(fmt.Sprintf("terraform -chdir=%s/infra-as-code/terraform/%s plan -var=\"public_key=%s\" -var=\"key_name=%s\"", repoDirRoot, cloudTemplate, pubKey, keyName))

				execSingleCommand(fmt.Sprintf("terraform  -chdir=%s/infra-as-code/terraform/%s apply -auto-approve -var=\"public_key=%s\" -var=\"key_name=%s\"", repoDirRoot, cloudTemplate, pubKey, keyName))
				//taking public ip and private ip from terraform.tfstate
				quickState, err := ioutil.ReadFile("DIGIT-DevOps/infra-as-code/terraform/quickstart-aws-ec2/terraform.tfstate")
				if err != nil {
					log.Printf("%v", err)
				}
				var quick configs.Quickstart
				err = json.Unmarshal(quickState, &quick)
				//publicip
				ip := quick.Outputs.PublicIP.Value
				//privateip
				privateip := quick.Resources[0].Instances[0].Attributes.PrivateIP
				createK3d(cluster_name, ip, keyName, privateip)
				changePrivateIp(cluster_name, privateip)

			}

		} else {
			if St.Db_pass == "" {
				db_pswd = enterValue(nil, "What should be the database password to be created, it should be 8 char min")
				St.Db_pass = db_pswd
			} else {
				db_pswd = St.Db_pass
			}
			writeState()
			if St.T_init == "failure" || St.T_init == "" {
				err = execSingleCommand(fmt.Sprintf("terraform -chdir=%s/infra-as-code/terraform/%s init", repoDirRoot, cloudTemplate))
			}
			if err == nil {
				St.T_init = "success"
				writeState()
			} else {
				St.T_init = "failure"
				writeState()
			}
			if St.T_plan == "failure" || St.T_plan == "" {
				err = execSingleCommand(fmt.Sprintf("terraform -chdir=%s/infra-as-code/terraform/%s plan -var=\"cluster_name=%s\" -var=\"db_password=%s\" -var=\"number_of_worker_nodes=%d\"", repoDirRoot, cloudTemplate, cluster_name, db_pswd, number_of_worker_nodes))
			}
			if err == nil {
				St.T_plan = "success"
				writeState()
			} else {
				St.T_plan = "failure"
				writeState()
			}
			if St.T_apply == "failure" || St.T_apply == "" {
				err = execSingleCommand(fmt.Sprintf("terraform -chdir=%s/infra-as-code/terraform/%s apply -auto-approve -var=\"cluster_name=%s\" -var=\"db_password=%s\" -var=\"number_of_worker_nodes=%d\"", repoDirRoot, cloudTemplate, cluster_name, db_pswd, number_of_worker_nodes))
			}
			if err == nil {
				St.T_apply = "success"
				writeState()
			} else {
				St.T_apply = "failure"
				writeState()
			}
			//calling funtion to write config file
			Configsfile()
			//calling function to create secret file
			envSecretsFile()
			writeState()

		}
	}
	contextset := setClusterContext()
	if contextset {
		deployCharts(servicesToDeploy, cluster_name)
	}
	//terraform output to a file
	//replace the env values with the tf output
	//save the kubetconfig and set the currentcontext
	//set dns in godaddy using the api's
	endScript()
}

func getService(fullChart Digit, service string, set Set, svclist *list.List) {
	for _, s := range fullChart.Modules {
		if s.Name == service {
			if set.Add(service) {
				svclist.PushFront(service) //Add services into the list
				if s.Dependencies != nil {
					for _, deps := range s.Dependencies {
						getService(fullChart, deps, set, svclist)
					}
				}
			}
		}
	}
}

// create a cluster in vm
func createK3d(clusterName string, publicIp string, keyName string, privateIp string) {
	commands := []string{
		"mkdir ~/kube && sudo chmod 777 ~/kube",
		"sudo k3d kubeconfig get k3s-default > " + clusterName + "_k3dconfig",
	}
	createClusterCmd := fmt.Sprintf("sudo k3d cluster create --api-port %s:6550 --k3s-server-arg --no-deploy=traefik --agents 2 -v /home/ubuntu/kube:/kube@agent[0,1] -v /home/ubuntu/kube:/kube@server[0] --port 8333:9000@loadbalancer --k3s-server-arg --tls-san=%s", privateIp, publicIp)
	command := fmt.Sprintf("%s&&%s&&%s", commands[0], createClusterCmd, commands[1])
	execRemoteCommand("ubuntu", publicIp, sshFile, command)
	copyConfig := fmt.Sprintf("scp ubuntu@%s:%s_k3dconfig  .", publicIp, clusterName)
	execCommand(copyConfig)
}

//changes the private ip in k3dconfig
func changePrivateIp(clusterName string, privateIp string) {
	path := fmt.Sprintf("%s_k3dconfig", clusterName)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("%v", err)
	}
	var con configs.Config
	err = yaml.Unmarshal(file, &con)
	if err != nil {
		log.Printf("%v", err)
	}
	server := fmt.Sprintf("https://%s:6550", privateIp)
	con.Clusters[0].Cluster.Server = server
	newfile, err := yaml.Marshal(&con)
	if err != nil {
		log.Printf("%v", err)

	}
	err = ioutil.WriteFile("new_k3dconfig", newfile, 0644)
	if err != nil {
		log.Printf("%v", err)
	}

}

func execCommand(command string) error {
	var err error
	parts := strings.Fields(command)
	//	The first part is the command, the rest are the args:
	head := parts[0]
	args := parts[1:len(parts)]
	//	Format the command

	log.Println(string(Blue), " ==> "+command)
	cmd := exec.Command(head, args...)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err = cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return err
}

func setClusterContext() bool {

	validatepath := func(input string) error {
		_, err := os.Stat(input)
		if os.IsNotExist(err) {
			return errors.New("The File does not exist in the given path")
		}
		return nil
	}

	var kubeconfig string
	kubeconfig = enterValue(validatepath, "Please enter the fully qualified path of your kubeconfig file")

	if kubeconfig != "" {
		getcontextcmd := fmt.Sprintf("kubectl config get-contexts --kubeconfig=%s", kubeconfig)
		err := execCommand(getcontextcmd)
		if err == nil {
			context := enterValue(nil, "Please enter the cluster context to be used from the avaliable contexts")
			if context != "" {
				usecontextcmd := fmt.Sprintf("kubectl config use-context %s --kubeconfig=%s", context, kubeconfig)
				err := execCommand(usecontextcmd)
				if err == nil {
					return true
				}
			}
		}
	}
	return false
}

func selectGovServicesToInstall() string {

	var versionfiles []string
	var modules []string
	svclist := list.New()
	set := NewSet()
	var argStr string = ""
	var releaseChartDir string = "../../config-as-code/product-release-charts/"

	// Get the versions from the chart and display it to user to select
	file, err := os.Open(releaseChartDir)
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	prodList, _ := file.Readdirnames(0) // 0 to read all files and folders

	var optedProduct string = ""
	if St.Product == "" {
		optedProduct, _ = sel(prodList, "Choose the Gov stack services that you would you like to install")
		St.Product = optedProduct
	} else {
		optedProduct = St.Product
	}
	if optedProduct != "" {
		files, err := ioutil.ReadDir(releaseChartDir + optedProduct)
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			name := f.Name()
			versionfiles = append(versionfiles, name[strings.Index(name, "-")+1:strings.Index(name, ".y")])
		}
		var version string = ""
		if St.Productversion == "" {
			version, _ = sel(versionfiles, "Which version of the selected product would like to install?")
			St.Productversion = version
		} else {
			version = St.Productversion
		}
		if version != "" {
			argFile := releaseChartDir + optedProduct + "/dependancy_chart-" + version + ".yaml"

			// Decode the yaml file and assigning the values to a map
			chartFile, err := ioutil.ReadFile(argFile)
			if err != nil {
				fmt.Println("\n\tERROR: Preparing required services details =>", argFile, err)
				return ""
			}

			// Parse the yaml values
			fullChart := Digit{}
			err = yaml.Unmarshal(chartFile, &fullChart)
			if err != nil {
				fmt.Println("\n\tERROR: Sourcing the the gov services matrix for your requirement => ", argFile, err)
				return ""
			}

			// Mapping the images to servicename
			var m = make(map[string][]string)
			for _, s := range fullChart.Modules {
				m[s.Name] = s.Services
				if strings.Contains(s.Name, "m_") {
					modules = append(modules, s.Name)
				}
			}
			modules = append(modules, "Exit")
			if len(St.Modules) == 0 {
				result, err := sel(modules, "Select the DIGIT's Gov services that you want to install, choose Exit to complete selection")

				for result != "Exit" && err == nil {
					selectedMod = append(selectedMod, result)
					St.Modules = append(St.Modules, result)
					result, err = sel(modules, "Select the modules you want to install, you can select multiple if you wish, choose Exit to complete selection")
				}
			} else {
				selectedMod = St.Modules
			}
			if selectedMod != nil {
				for _, mod := range selectedMod {
					getService(fullChart, mod, *set, svclist)
				}
				for element := svclist.Front(); element != nil; element = element.Next() {
					imglist := m[element.Value.(string)]
					imglistsize := len(imglist)
					for i, service := range imglist {
						argStr = argStr + service
						if !(element.Next() == nil && i == imglistsize-1) {
							argStr = argStr + ","
						}

					}
				}
			}
		}
	}
	return argStr
}

// To deploy all the services selected

func deployCharts(argStr string, configFile string) {

	var goDeployCmd string = fmt.Sprintf("go run main.go deploy -c -e %s %s", configFile, argStr)
	var previewDeployCmd string = fmt.Sprintf("%s -p", goDeployCmd)

	confirm := []string{"Yes", "No"}
	preview, _ := sel(confirm, "Do you want to preview the k8s manifests before the actual Deployment")
	if preview == "Yes" {
		fmt.Println("That's cool... preview is getting loaded. Please review it and cross check the kubernetes manifests before the deployment")
		err := execCommand(previewDeployCmd)
		if err == nil {
			fmt.Println("You can now start actual deployment")
			err := execCommand(goDeployCmd)
			if err == nil {
				fmt.Println("We are done with the deployment. You can start using the services. Thank You!!!")
				return
			} else {
				fmt.Println("Something went wrong, refer the error\n")
				fmt.Println(err)
			}
			return
		} else {
			fmt.Println("Something went wrong, refer the error\n")
			fmt.Println(err)
		}
	} else {
		consent, _ := sel(confirm, "Are we good to proceed with the actual deployment?")
		if consent == "Yes" {
			fmt.Println("Whola!, That's great... Sit back and wait for the deployment to complete in about 10 min")
			err := execCommand(goDeployCmd)
			if err == nil {
				fmt.Println("We are done with the deployment. You can start using the services. Thank You!!!")
				fmt.Println("Hope I made your life easy with the deployment ... Have a goodd day !!!")
				return
			} else {
				fmt.Println("Something went wrong, refer the error\n")
				fmt.Println(err)
			}
		} else {
			endScript()
		}

	}

}

func execRemoteCommand(user string, ip string, sshFileLocation string, command string) error {
	var err error
	sshPreFix := fmt.Sprintf("ssh %s@%s -i %s \"%s\" ", user, ip, sshFileLocation, command)

	cmd := exec.Command("sh", "-c", sshPreFix)

	log.Println(string(Blue), " ==> "+sshPreFix)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err = cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return err
}
func execSingleCommand(command string) error {
	var err error

	cmd := exec.Command("sh", "-c", command)

	log.Println(string(Blue), " ==> "+command)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err = cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	// fmt.Println(err)
	return err
}

// Cloud cloudLoginCredentials functions
func awslogin(accessKey string, secretKey string, sessionToken string, profile string) bool {

	var cloudLoginCredentials bool = false
	var awslogincommand string = ""

	if accessKey != "" && secretKey != "" && sessionToken == "" {
		awslogincommand = fmt.Sprintf("aws configure --profile digit-infra-aws set aws_access_key_id \"%s\" && aws configure --profile digit-infra-aws set aws_secret_access_key \"%s\" && aws configure --profile digit-infra-aws set region \"ap-south-1\"", accessKey, secretKey)
	} else if sessionToken != "" {
		awslogincommand = fmt.Sprintf("aws configure --profile digit-infra-aws set aws_access_key_id \"%s\" && aws configure --profile digit-infra-aws set aws_secret_access_key \"%s\" && aws configure --profile digit-infra-aws set aws_session_token \"%s\"  && aws configure --profile digit-infra-aws set region \"ap-south-1\"", accessKey, secretKey, sessionToken)
	} else {
		awsProf := ""
		profile := ""
		awsProf = fmt.Sprintf("aws configure list-profiles")
		out, err := execCommandWithOutput(awsProf)
		if err != nil {
			log.Printf("%s", err)
		}
		profList := strings.Fields(out)
		if St.Aws_profile == "" {
			profile, _ = sel(profList, "choose the profile with right access")
			St.Aws_profile = profile
		} else {
			profile = St.Aws_profile
		}
		awslogincommand = fmt.Sprintf("aws configure --profile %s set region \"ap-south-1\"", profile)
		// execCommand(fmt.Sprintf("aws configure list"))

	}

	log.Println(awslogincommand)
	err := execSingleCommand(awslogincommand)
	if err == nil {
		St.Aws_command = "success"
		cloudLoginCredentials = true
	} else {
		St.Aws_command = "failure"
		writeState()
	}
	return cloudLoginCredentials
}

func azurelogin(userName string, password string) bool {

	var cloudLoginCredentials bool = false
	if userName != "" && password != "" {
		azurelogincommand := fmt.Sprintf("az cloudLoginCredentials -u %s -p %s", userName, password)
		err := execCommand(azurelogincommand)
		if err == nil {
			cloudLoginCredentials = true
		}
	}
	return cloudLoginCredentials
}

// Input functions

func sel(items []string, label string) (string, error) {
	var result string
	var err error
	prompt := promptui.Select{
		Label: label,
		Items: items,
		Size:  30,
	}
	_, result, err = prompt.Run()

	//if err != nil {
	//	fmt.Printf("Invalid Selection %v\n", err)
	//}
	return result, err
}

func enterValue(validate promptui.ValidateFunc, label string) string {
	var result string
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}
	result, _ = prompt.Run()

	//if err != nil {
	//	fmt.Printf("Invalid Selection %v\n", err)
	//}
	return result
}

func addDNS(dnsDomain string, dnsType string, dnsName string, dnsValue string) bool {

	var headers string = "Authorization: sso-key 3mM44UcBKoVvB2_Xspi4jKZqJSQUkdouMV4Ck:3pzZiuUPNxzZKu2FfUD9Sm"

	dnsCommand := fmt.Sprintf("curl -X PATCH \"https://api.godaddy.com/v1/domains/%s/records -H %s -H Content-Type: application/json --data-raw [{\"data\":\"%s\",\"name\":\"%s\",\"type\":\"%s\"}]", dnsDomain, headers, dnsValue, dnsName, dnsType)
	fmt.Println(dnsCommand)
	err := execSingleCommand(dnsCommand)
	if err == nil {
		return true
	} else {
		return false
	}
}

func GetKeyPair(file string) (string, string, error) {
	// read keys from file
	_, err := os.Stat(file)
	if err == nil {
		priv, err := ioutil.ReadFile(file)
		if err != nil {
			lumber.Debug("Failed to read file - %s", err)
			goto genKeys
		}
		pub, err := ioutil.ReadFile(file + ".pub")
		if err != nil {
			lumber.Debug("Failed to read pub file - %s", err)
			goto genKeys
		}
		return string(pub), string(priv), nil
	}

	// generate keys and save to file
genKeys:
	pub, priv, err := GenKeyPair()
	err = ioutil.WriteFile(file, []byte(priv), 0600)
	if err != nil {
		return "", "", fmt.Errorf("Failed to write file - %s", err)
	}
	err = ioutil.WriteFile(file+".pub", []byte(pub), 0644)
	if err != nil {
		return "", "", fmt.Errorf("Failed to write pub file - %s", err)
	}

	return pub, priv, nil
}

func GenKeyPair() (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	var private bytes.Buffer
	if err := pem.Encode(&private, privateKeyPEM); err != nil {
		return "", "", err
	}

	// generate public key
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}

	public := ssh.MarshalAuthorizedKey(pub)
	return string(public), private.String(), nil
}

// below function can be used to store output of command to variable
func execCommandWithOutput(command string) (string, error) {

	parts := strings.Fields(command)
	//	The first part is the command, the rest are the args:
	head := parts[0]
	args := parts[1:len(parts)]
	//	Format the command

	log.Println(string(Blue), " ==> "+command)
	cmd := exec.Command(head, args...)
	out, err := cmd.Output()
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	if err != nil {
		log.Fatalf("%s", err)
	}
	return string(out), err
}

// write configs to environment file
func Configsfile() {
	Confirm := []string{"Yes", "No"}
	var out configs.Output
	State, err := ioutil.ReadFile("DIGIT-DevOps/infra-as-code/terraform/sample-aws/terraform.tfstate")
	if err != nil {
		log.Printf("%v", err)
	}
	err = json.Unmarshal(State, &out)
	Config := make(map[string]interface{})
	var HasDomain string
	if St.Hasdomain != "" {
		HasDomain = St.Hasdomain
	} else {
		HasDomain, _ = sel(Confirm, "Do you have a URL which can be used after DIGIT installation to access actual site ?")
		St.Hasdomain = HasDomain
	}
	var Domain string
	if HasDomain == "Yes" {
		if St.Domain != "" {
			Domain = St.Domain
		} else {
			Domain = enterValue(nil, "Enter Domain name")
			St.Domain = Domain
			writeState()
		}
	} else {
		fmt.Println("Create a domain From Godaddy or any other DNS providers and come back.")
		Domain = enterValue(nil, "Enter Domain name")
		Domain = St.Domain
	}
	var HasGitacc string
	if St.Hasgitacc != "" {
		HasGitacc = St.Hasgitacc
	} else {
		HasGitacc, _ = sel(Confirm, "Do you have a Github Account ?")
		St.Hasgitacc = HasGitacc
	}
	var BranchName string
	if HasGitacc == "Yes" {
		if St.Branch != "" {
			BranchName = St.Branch
		} else {
			BranchName = enterValue(nil, "Enter Branch name for Configs and MDMS eg: UAT,QA")
			St.Branch = BranchName
		}
	} else {
		fmt.Println("click on the URL https://github.com/ to create a github account \n\tNote: Create two github account one will be used as Org account and other as user.")
	}
	KafkaVolumes := out.Outputs.KafkaVolIds.Value
	ZookeeperVolumes := out.Outputs.ZookeeperVolumeIds.Value
	ElasticDataVolumes := out.Outputs.EsDataVolumeIds.Value
	ElasticMasterVolumes := out.Outputs.EsMasterVolumeIds.Value
	fmt.Println("Fork the Configs and egov-mdms-data repositories to Org account and give User account permission over the repos")
	var ConfigsBranch string
	if St.Config_git_url != "" {
		ConfigsBranch = St.Config_git_url
	} else {
		ConfigsBranch = enterValue(nil, "Enter your configs repo ssh url")
		St.Config_git_url = ConfigsBranch
	}
	var MdmsBranch string
	if St.Mdms_git_url != "" {
		MdmsBranch = St.Mdms_git_url
	} else {
		MdmsBranch = enterValue(nil, "Enter your mdms repo ssh url")
		St.Mdms_git_url = MdmsBranch
	}
	Config["Domain"] = Domain
	Config["BranchName"] = BranchName
	Config["db-host"] = out.Outputs.DbInstanceEndpoint.Value
	Config["db_name"] = out.Outputs.DbInstanceName.Value
	Config["configs-branch"] = ConfigsBranch
	Config["mdms-branch"] = MdmsBranch
	Config["file_name"] = cluster_name
	var HasSMSGateway string
	if St.SmsGateway != "" {
		HasSMSGateway = St.Smsproceed
	} else {
		HasSMSGateway, _ = sel(Confirm, "Do You have your sms Gateway?")
		St.Smsproceed = HasSMSGateway
	}
	if HasSMSGateway == "Yes" {
		var SmsUrl string
		if St.SmsUrl != "" {
			SmsUrl = St.SmsUrl
		} else {
			SmsUrl = enterValue(nil, "Enter your SMS provider url")
			St.SmsUrl = SmsUrl
		}
		var SmsGateway string
		if St.SmsGateway != "" {
			SmsGateway = St.SmsGateway
		} else {
			SmsGateway = enterValue(nil, "Enter your SMS Gateway")
			St.SmsGateway = SmsGateway
		}
		var SmsSender string
		if St.SmsSender != "" {
			SmsSender = St.SmsSender
		} else {
			SmsSender = enterValue(nil, "Enter your SMS sender")
			St.SmsSender = SmsSender
		}

		Config["sms-provider-url"] = SmsUrl
		Config["sms-gateway-to-use"] = SmsGateway
		Config["sms-sender"] = SmsSender

	}
	var IsFilestore string
	if St.Fileproceed != "" {
		IsFilestore = St.Fileproceed
	} else {
		IsFilestore, _ = sel(Confirm, "Do You need filestore which is used to to store files?")
		St.Fileproceed = IsFilestore
	}
	if IsFilestore == "Yes" {
		if CloudProvider == "aws" {
			var bucket string
			if St.Bucket != "" {
				bucket = St.Bucket
			} else {
				bucket := enterValue(nil, "Enter the filestore bucket name")
				St.Bucket = bucket
			}
			Config["fixed-bucket"] = bucket
		}
		if CloudProvider == "sdc" {
			var bucket string
			if St.Bucket != "" {
				bucket = St.Bucket
			} else {
				bucket := enterValue(nil, "Enter the filestore bucket name")
				St.Bucket = bucket
			}
			Config["fixed-bucket"] = bucket
		}
	}
	var IsChatbot string
	if St.Botproceed != "" {
		IsChatbot = St.Botproceed
	} else {
		IsChatbot, _ = sel(Confirm, "Do You need chatbot?")
		St.Botproceed = IsChatbot
	}
	writeState()
	configs.DeployConfig(Config, KafkaVolumes, ZookeeperVolumes, ElasticDataVolumes, ElasticMasterVolumes, selectedMod, HasSMSGateway, IsFilestore, IsChatbot, CloudProvider)

}

// writes secrets to file
func envSecretsFile() {
	generateSsh()
	ssh := ""
	ssh = fmt.Sprintf("cat private.pem")
	Out, err := execCommandWithOutput(ssh)
	if err != nil {
		log.Printf("%s", err)
	}
	configs.SecretFile(cluster_name, Out, selectedMod)
}

// generate ssh key to configs file
func generateSsh() {
	// generate key
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("Cannot generate RSA keyn")
		os.Exit(1)
	}
	publickey := &privatekey.PublicKey

	// dump private key to file
	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(privatekey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privatePem, err := os.Create("private.pem")
	if err != nil {
		fmt.Printf("error when create private.pem: %s n", err)
		os.Exit(1)
	}
	err = pem.Encode(privatePem, privateKeyBlock)
	if err != nil {
		fmt.Printf("error when encode private pem: %s n", err)
		os.Exit(1)
	}

	// dump public key to file
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publickey)
	if err != nil {
		fmt.Printf("error when dumping publickey: %s n", err)
		os.Exit(1)
	}
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicPem, err := os.Create("public.pem")
	if err != nil {
		fmt.Printf("error when create public.pem: %s n", err)
		os.Exit(1)
	}
	err = pem.Encode(publicPem, publicKeyBlock)
	if err != nil {
		fmt.Printf("error when encode public pem: %s n", err)
		os.Exit(1)
	}
}

func endScript() {
	fmt.Println("Take your time, You can come back at any time ... Thank for leveraging me :)!!!")
	fmt.Println("Hope I made your life easy with the deployment ... Have a good day !!!")
	return
}

// Writes the state to file
func writeState() {
	state, err := yaml.Marshal(&St)
	if err != nil {
		log.Printf("%v", err)

	}
	stateFile := fmt.Sprintf("state.yaml")
	err = ioutil.WriteFile(stateFile, state, 0644)
	if err != nil {
		log.Printf("%v", err)
	}
}
