package configs

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

var region = "ap-south-1b"

// Quickstart kubeconfig struct
type Config struct {
	APIVersion string `yaml:"apiVersion"`
	Clusters   []struct {
		Cluster struct {
			CertificateAuthorityData string `yaml:"certificate-authority-data"`
			Server                   string `yaml:"server"`
		} `yaml:"cluster"`
		Name string `yaml:"name"`
	} `yaml:"clusters"`
	Contexts []struct {
		Context struct {
			Cluster string `yaml:"cluster"`
			User    string `yaml:"user"`
		} `yaml:"context"`
		Name string `yaml:"name"`
	} `yaml:"contexts"`
	CurrentContext string `yaml:"current-context"`
	Kind           string `yaml:"kind"`
	Preferences    struct {
	} `yaml:"preferences"`
	Users []struct {
		Name string `yaml:"name"`
		User struct {
			ClientCertificateData string `yaml:"client-certificate-data"`
			ClientKeyData         string `yaml:"client-key-data"`
		} `yaml:"user"`
	} `yaml:"users"`
}

// environment secret struct
type Secret struct {
	ClusterConfigs struct {
		Secrets struct {
			Db struct {
				Username       string `yaml:"username"`
				Password       string `yaml:"password"`
				FlywayUsername string `yaml:"flywayUsername"`
				FlywayPassword string `yaml:"flywayPassword"`
			} `yaml:"db"`
			EgovNotificationSms struct {
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			} `yaml:"egov-notification-sms"`
			EgovFilestore struct {
				AwsKey       string `yaml:"aws-key"`
				AwsSecretKey string `yaml:"aws-secret-key"`
			} `yaml:"egov-filestore"`
			EgovLocation struct {
				Gmapskey string `yaml:"gmapskey"`
			} `yaml:"egov-location"`
			EgovPgService struct {
				AxisMerchantID         string `yaml:"axis-merchant-id"`
				AxisMerchantSecretKey  string `yaml:"axis-merchant-secret-key"`
				AxisMerchantUser       string `yaml:"axis-merchant-user"`
				AxisMerchantPwd        string `yaml:"axis-merchant-pwd"`
				AxisMerchantAccessCode string `yaml:"axis-merchant-access-code"`
				PayuMerchantKey        string `yaml:"payu-merchant-key"`
				PayuMerchantSalt       string `yaml:"payu-merchant-salt"`
			} `yaml:"egov-pg-service"`
			Pgadmin struct {
				AdminEmail    string `yaml:"admin-email"`
				AdminPassword string `yaml:"admin-password"`
				ReadEmail     string `yaml:"read-email"`
				ReadPassword  string `yaml:"read-password"`
			} `yaml:"pgadmin"`
			EgovEncService struct {
				MasterPassword      string `yaml:"master-password"`
				MasterSalt          string `yaml:"master-salt"`
				MasterInitialvector string `yaml:"master-initialvector"`
			} `yaml:"egov-enc-service"`
			EgovNotificationMail struct {
				Mailsenderusername string `yaml:"mailsenderusername"`
				Mailsenderpassword string `yaml:"mailsenderpassword"`
			} `yaml:"egov-notification-mail"`
			GitSync struct {
				SSH        string `yaml:"ssh"`
				KnownHosts string `yaml:"known-hosts"`
			} `yaml:"git-sync"`
			Kibana struct {
				Namespace   string `yaml:"namespace"`
				Credentials string `yaml:"credentials"`
			} `yaml:"kibana"`
			EgovSiMicroservice struct {
				SiMicroserviceUser     string `yaml:"si-microservice-user"`
				SiMicroservicePassword string `yaml:"si-microservice-password"`
				MailSenderPassword     string `yaml:"mail-sender-password"`
			} `yaml:"egov-si-microservice"`
			EgovEdcrNotification struct {
				EdcrMailUsername string `yaml:"edcr-mail-username"`
				EdcrMailPassword string `yaml:"edcr-mail-password"`
				EdcrSmsUsername  string `yaml:"edcr-sms-username"`
				EdcrSmsPassword  string `yaml:"edcr-sms-password"`
			} `yaml:"egov-edcr-notification"`
			Chatbot struct {
				ValuefirstUsername string `yaml:"valuefirst-username"`
				ValuefirstPassword string `yaml:"valuefirst-password"`
			} `yaml:"chatbot"`
			EgovUserChatbot struct {
				CitizenLoginPasswordOtpFixedValue string `yaml:"citizen-login-password-otp-fixed-value"`
			} `yaml:"egov-user-chatbot"`
			Oauth2Proxy struct {
				ClientID     string `yaml:"clientID"`
				ClientSecret string `yaml:"clientSecret"`
				CookieSecret string `yaml:"cookieSecret"`
			} `yaml:"oauth2-proxy"`
		} `yaml:"secrets"`
	} `yaml:"cluster-configs"`
}

//terrafrom struct
type Output struct {
	Outputs struct {
		ClusterEndpoint struct {
			Value string `json:"value"`
		} `json:"cluster_endpoint"`
		DbInstanceEndpoint struct {
			Value string `json:"value"`
		} `json:"db_instance_endpoint"`
		DbInstanceName struct {
			Value string `json:"value"`
		} `json:"db_instance_name"`
		DbInstancePort struct {
			Value int `json:"value"`
		} `json:"db_instance_port"`
		DbInstanceUsername struct {
			Value string `json:"value"`
		} `json:"db_instance_username"`
		EsDataVolumeIds struct {
			Value []string `json:"value"`
		} `json:"es_data_volume_ids"`
		EsMasterVolumeIds struct {
			Value []string `json:"value"`
		} `json:"es_master_volume_ids"`
		KafkaVolIds struct {
			Value []string `json:"value"`
		} `json:"kafka_vol_ids"`
		KubectlConfig struct {
			Value string `json:"value"`
		} `json:"kubectl_config"`
		PrivateSubnets struct {
			Value []string `json:"value"`
		} `json:"private_subnets"`
		PublicSubnets struct {
			Value []string `json:"value"`
		} `json:"public_subnets"`
		VpcID struct {
			Value string `json:"value"`
		} `json:"vpc_id"`
		ZookeeperVolumeIds struct {
			Value []string `json:"value"`
		} `json:"zookeeper_volume_ids"`
	} `json:"outputs"`
}

//quickstart tfstate struct
type Quickstart struct {
	Outputs struct {
		PublicIP struct {
			Value string `json:"value"`
		} `json:"public_ip"`
	} `json:"outputs"`
	Resources []struct {
		Instances []struct {
			Attributes struct {
				PrivateIP string `json:"private_ip"`
			} `json:"attributes"`
		} `json:"instances"`
	} `json:"resources"`
}

func DeployConfig(Config map[string]interface{}, kafkaVolumes []string, ZookeeperVolumes []string, ElasticDataVolumes []string, ElasticMasterVolumes []string, modules []string, smsproceed string, fileproceed string, botproceed string, flag string) {

	file, err := ioutil.ReadFile("DIGIT-DevOps/config-as-code/environments/egov-demo.yaml")
	if err != nil {
		log.Printf("%v", err)
	}
	var data map[string]interface{}
	DataMap := make(map[string]interface{})
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		log.Printf("%v", err)
	}
	for i := range data {
		if i == "global" {
			Global := data[i].(map[string]interface{})
			for j := range Global {
				if j == "domain" {
					Global[j] = Config["Domain"]
				}
			}
		}
		if i == "cluster-configs" {
			// fmt.Println("found cluster-configs")
			ClusterConfigs := data[i].(map[string]interface{})
			for j := range ClusterConfigs {
				if j == "configmaps" {
					// fmt.Println("found configmaps")
					Configmaps := ClusterConfigs[j].(map[string]interface{})
					for k := range Configmaps {
						if k == "egov-config" {
							// fmt.Println("found egov-config")
							EgovConfig := Configmaps[k].(map[string]interface{})
							for l := range EgovConfig {
								if l == "data" {
									// fmt.Println("found data")
									Data := EgovConfig[l].(map[string]interface{})
									for m := range Data {
										if m == "db-host" {
											Host := Config["db-host"].(string)
											provider := Host[:strings.IndexByte(Host, ':')]
											Data[m] = provider
										}
										if m == "db-name" {
											Data[m] = Config["db_name"]
										}
										if m == "db-url" {
											url := fmt.Sprintf("jdbc:postgresql://%s/%s", Config["db-host"], Config["db_name"])
											Data[m] = url
										}
										if m == "domain" {
											Data[m] = Config["Domain"]
										}
										if m == "egov-services-fqdn-name" {
											fqdn := fmt.Sprintf("https://%s/", Config["Domain"])
											Data[m] = fqdn
										}
										if m == "s3-assets-bucket" {

										}
										if m == "es-host" {

										}
										if m == "es-indexer-host" {

										}
										if m == "flyway-locations" {

										}
										if m == "kafka-brokers" {

										}
										if m == "kafka-infra-brokers" {

										}
										if m == "logging-level-jdbc" {

										}
										if m == "mobile-validation-workaround" {

										}
										if m == "serializers-timezone-in-ist" {

										}
										if m == "server-tomcat-max-connections" {

										}
										if m == "server-tomcat-max-threads" {

										}
										if m == "sms-enabled" {

										}
										if m == "spring-datasource-tomcat-initialSize" {

										}
										if m == "spring-datasource-tomcat-max-active" {

										}
										if m == "spring-jpa-show-sql" {

										}
										if m == "timezone" {

										}
										if m == "tracer-errors-provideexceptionindetails" {

										}
									}
								}
							}
						}

					}
				}
			}
		}
		if i == "egov-mdms-service" || i == "egov-indexer" || i == "egov-persister" || i == "egov-data-uploader" || i == "egov-searcher" || i == "dashboard-analytics" || i == "dashboard-ingest" || i == "report" || i == "pdf-service" {
			// fmt.Println("in mdms")
			Service := data[i].(map[string]interface{})
			for j := range Service {
				if j == "search-yaml-path" {

				}
				if j == "config-schema-paths" {

				}
				if j == "replicas" {

				}
				if j == "mdms-path" {

				}
				if j == "heap" {

				}
				if j == "memory_limits" {

				}
				if j == "mdms-path" {

				}
				if j == "persist-yml-path" {

				}
				if j == "initContainers" {
					// fmt.Println("in init")
					InitContainers := Service[j].(map[string]interface{})
					for k := range InitContainers {
						if k == "gitSync" {
							// fmt.Println("in git sync")
							GitSync := InitContainers[k].(map[string]interface{})
							for l := range GitSync {
								if l == "branch" {
									GitSync[l] = Config["BranchName"]
								}
								if l == "repo" {
									if data[i] == "egov-mdms-service" {
										GitSync[l] = Config["mdms-branch"]
									} else {
										GitSync[l] = Config["configs-branch"]
									}
								}
							}
						}
					}
				}
				if j == "mdms-folder" {

				}
				if j == "masters-config-url" {

				}
				if j == "java-args" {

				}
				if j == "egov-indexer-yaml-repo-path" {

				}
			}
		}
		if i == "cert-manager" {
			CertManager := data[i].(map[string]interface{})
			for j := range CertManager {
				if j == "email" {
					CertManager[j] = ""
				}
			}
		}
		if i == "kafka-v2" {
			KafkaV2 := data[i].(map[string]interface{})
			for j := range KafkaV2 {
				if j == "persistence" {
					Persistence := KafkaV2[j].(map[string]interface{})
					for k := range Persistence {
						if k == "aws" {
							Aws := Persistence[k].([]interface{})
							N := 0
							for l := range Aws {
								Volume := Aws[l].(map[string]interface{})
								for m := range Volume {
									if m == "volumeId" && N == l {
										Volume[m] = kafkaVolumes[l]
									}
									if m == "zone" {
										Volume[m] = region
									}
								}
								N++

							}
						}
					}
				}
			}
		}
		if i == "zookeeper-v2" {
			ZookeeperV2 := data[i].(map[string]interface{})
			for j := range ZookeeperV2 {
				if j == "persistence" {
					Persistence := ZookeeperV2[j].(map[string]interface{})
					for k := range Persistence {
						if k == "aws" {
							Aws := Persistence[k].([]interface{})
							N := 0
							for l := range Aws {
								Volume := Aws[l].(map[string]interface{})
								for m := range Volume {
									if m == "volumeId" && N == l {
										Volume[m] = ZookeeperVolumes[l]
									}
									if m == "zone" {
										Volume[m] = region
									}
								}
								N++

							}
						}
					}
				}
			}
		}
		if i == "elasticsearch-data-v1" {
			ElasticsearchDataV := data[i].(map[string]interface{})
			for j := range ElasticsearchDataV {
				if j == "persistence" {
					Persistence := ElasticsearchDataV[j].(map[string]interface{})
					for k := range Persistence {
						if k == "aws" {
							Aws := Persistence[k].([]interface{})
							N := 0
							for l := range Aws {
								ElasticDVol := Aws[l].(map[string]interface{})
								for m := range ElasticDVol {
									if m == "volumeId" && N == l {
										ElasticDVol[m] = ElasticDataVolumes[l]
									}
									if m == "zone" {
										ElasticDVol[m] = region
									}
								}
								N++

							}
						}
					}
				}
			}
		}
		if i == "elasticsearch-master-v1" {
			ElasticsearchMasterV := data[i].(map[string]interface{})
			for j := range ElasticsearchMasterV {
				if j == "persistence" {
					persistence := ElasticsearchMasterV[j].(map[string]interface{})
					for k := range persistence {
						if k == "aws" {
							Aws := persistence[k].([]interface{})
							N := 0
							for l := range Aws {
								ElasticMVol := Aws[l].(map[string]interface{})
								for m := range ElasticMVol {
									if m == "volumeId" && N == l {
										ElasticMVol[m] = ElasticMasterVolumes[l]
									}
									if m == "zone" {
										ElasticMVol[m] = region
									}
								}
								N++

							}
						}
					}
				}
			}
		}
		if i == "employee" {
			Employee := data[i].(map[string]interface{})
			for j := range Employee {
				if j == "dashboard-url" {

				}
				if j == "custom-js-injection" {

				}
			}
		}
		if i == "citizen" {
			Citizen := data[i].(map[string]interface{})
			for j := range Citizen {
				if j == "custom-js-injection" {

				}
			}
		}
		if i == "digit-ui" {
			DigitUi := data[i].(map[string]interface{})
			for j := range DigitUi {
				if j == "custom-js-injection" {
				}
			}
		}
		if i == "egov-filestore" && fileproceed == "yes" {
			Filestore := data[i].(map[string]interface{})
			for j := range Filestore {
				if j == "volume" {

				}
				if j == "is-bucket-fixed" {

				}
				if j == "minio.url" {

				}
				if j == "aws.s3.url" {

				}
				if j == "is-s3-enabled" {

				}
				if j == "minio-enabled" {

				}
				if j == "allowed-file-formats-map" {

				}
				if j == "llowed-file-formats" {

				}
				if j == "filestore-url-validity" {

				}
				if j == "fixed-bucketname" {
					Filestore[j] = Config["fixed-bucket"]
				}
			}

		}
		if i == "egov-notification-sms" && smsproceed == "yes" {
			Notification := data[i].(map[string]interface{})
			for j := range Notification {
				if j == "sms-provider-url" {
					Notification[j] = Config["sms-provider-url"]
				}
				if j == "sms.provider.class" {

				}
				if j == "sms.provider.contentType" {

				}
				if j == "sms-config-map" {

				}
				if j == "sms-gateway-to-use" {
					Notification[j] = Config["sms-gateway-to-use"]
				}
				if j == "sms-sender" {
					Notification[j] = Config["sms-sender"]
				}
				if j == "sms-sender-requesttype" {

				}
				if j == "sms-custom-config" {

				}
				if j == "sms-extra-req-params" {

				}
				if j == "sms-sender-req-param-name" {

				}
				if j == "sms-sender-username-req-param-name" {

				}
				if j == "sms-sender-password-req-param-name" {

				}
				if j == "sms-destination-mobile-req-param-name" {

				}
				if j == "sms-message-req-param-name" {

				}
				if j == "sms-error-codes" {

				}
			}
			DataMap["egov-notification-sms"] = data["egov-notification-sms"]
		}
		if i == "egov-user" {
			EgovUser := data[i].(map[string]interface{})
			for j := range EgovUser {
				if j == "heap" {

				}
				if j == "memory_limits" {

				}
				if j == "otp-validation" {

				}
				if j == "citizen-otp-enabled" {

				}
				if j == "employee-otp-enabled" {

				}
				if j == "access-token-validity" {

				}
				if j == "refresh-token-validity" {

				}
				if j == "default-password-expiry" {

				}
				if j == "mobile-number-validation" {

				}
				if j == "roles-state-level" {

				}
				if j == "zen-registration-withlogin" {

				}
				if j == "citizen-otp-fixed" {

				}
				if j == "citizen-otp-fixed-enabled" {

				}
				if j == "egov-state-level-tenant-id" {

				}
				if j == "decryption-abac-enabled" {

				}
			}
		}
		if i == "chatbot" && botproceed == "yes" {
			Chtbot := data[i].(map[string]interface{})
			for j := range Chtbot {
				if j == "kafka-topics-partition-count" {

				}
				if j == "kafka-topics-replication-factor" {

				}
				if j == "kafka-consumer-poll-ms" {

				}
				if j == "kafka-producer-linger-ms" {

				}
				if j == "contact-card-whatsapp-number" {

				}
				if j == "contact-card-whatsapp-name" {

				}
				if j == "valuefirst-whatsapp-number" {

				}
				if j == "valuefirst-notification-assigned-templateid" {

				}
				if j == "valuefirst-notification-resolved-templateid" {

				}
				if j == "valuefirst-notification-rejected-templateid" {

				}
				if j == "valuefirst-notification-reassigned-templateid" {

				}
				if j == "valuefirst-notification-commented-templateid" {

				}
				if j == "valuefirst-notification-welcome-templateid" {

				}
				if j == "valuefirst-notification-root-templateid" {

				}
				if j == "valuefirst-send-message-url" {

				}
				if j == "user-service-chatbot-citizen-passwrord" {

				}
			}
			DataMap["chatbot"] = data["chatbot"]
		}
		if i == "bpa-services" {
			Bpa := data[i].(map[string]interface{})
			for j := range Bpa {
				if j == "memory_limits" {

				}
				if j == "java-args" {

				}
				if j == "java-debug" {

				}
				if j == "tracing-enabled" {

				}
				if j == "egov.idgen.bpa.applicationNum.format" {

				}
			}
		}
		if i == "bpa-calculator" {
			BpaCalc := data[i].(map[string]interface{})
			for j := range BpaCalc {
				if j == "memory_limits" {

				}
				if j == "java-args" {

				}
				if j == "java-debug" {

				}
				if j == "tracing-enabled" {

				}
			}
		}
		if i == "ws-services" {
			WsService := data[i].(map[string]interface{})
			for j := range WsService {
				if j == "wcid-format" {

				}
			}
		}
		if i == "sw-services" {
			SwSvc := data[i].(map[string]interface{})
			for j := range SwSvc {
				if j == "scid-format" {

				}
			}
		}
		if i == "egov-pg-service" {
			PgSvc := data[i].(map[string]interface{})
			for j := range PgSvc {
				if j == "axis" {

				}
			}
		}
		if i == "report" {
			Report := data[i].(map[string]interface{})
			for j := range Report {
				if j == "heap" {

				}
				if j == "tracing-enabled" {

				}
				if j == "spring-datasource-tomcat-max-active" {

				}
				if j == "initContainers" {
					Init := Report[j].(map[string]interface{})
					for k := range Init {
						if k == "gitSync" {
							GSync := Init[k].(map[string]interface{})
							for l := range GSync {
								if l == "repo" {

								}
								if l == "branch" {
									GSync[l] = Config["BranchName"]
								}
							}
						}
					}
				}
				if j == "report-locationsfile-path" {

				}
			}
		}
		if i == "pdf-service" {
			PdfSvc := data[i].(map[string]interface{})
			for j := range PdfSvc {
				if j == "initContainers" {
					InitContainer := PdfSvc[j].(map[string]interface{})
					for k := range InitContainer {
						if k == "gitSync" {
							Git := InitContainer[k].(map[string]interface{})
							for l := range Git {
								if l == "repo" {

								}
								if l == "branch" {
									Git[l] = Config["BranchName"]
								}
							}
						}
					}
				}
				if j == "data-config-urls" {

				}
				if j == "format-config-urls" {

				}

			}
		}
		if i == "egf-master" {
			EgfMaster := data[i].(map[string]interface{})
			for j := range EgfMaster {
				if j == "db-url" {

				}
				if j == "memory_limits" {

				}
				if j == "heap" {

				}

			}
		}
		if i == "egov-custom-consumer" {
			EgovConsumer := data[i].(map[string]interface{})
			for j := range EgovConsumer {
				if j == "erp-host" {

				}
			}
		}
		if i == "egov-apportion-service" {
			Apportion := data[i].(map[string]interface{})
			for j := range Apportion {
				if j == "memory_limits" {

				}
				if j == "heap" {

				}
			}
		}
		if i == "redoc" {
			Redoc := data[i].(map[string]interface{})
			for j := range Redoc {
				if j == "replicas" {

				}
				if j == "images" {

				}
				if j == "service_type" {

				}
			}
		}
		if i == "nginx-ingress" {
			Nginx := data[i].(map[string]interface{})
			for j := range Nginx {
				if j == "images" {

				}
				if j == "replicas" {

				}
				if j == "default-backend-service" {

				}
				if j == "namespace" {

				}
				if j == "cert-issuer" {

				}
				if j == "ssl-protocols" {

				}
				if j == "ssl-ciphers" {

				}
				if j == "ssl-ecdh-curve" {

				}
			}
		}
		if i == "cert-manager" {
			CertManager := data[i].(map[string]interface{})
			for j := range CertManager {
				if j == "email" {

				}
			}
		}
		if i == "zuul" {
			Zuul := data[i].(map[string]interface{})
			for j := range Zuul {
				if j == "replicas" {

				}
				if j == "custom-filter-property" {

				}
				if j == "tracing-enabled" {

				}
				if j == "heap" {

				}
				if j == "server-tomcat-max-threads" {

				}
				if j == "server-tomcat-max-connections" {

				}
				if j == "egov-open-endpoints-whitelist" {

				}
				if j == "egov-mixed-mode-endpoints-whitelist" {

				}
			}
		}
		if i == "collection-services" {
			CollectionService := data[i].(map[string]interface{})
			for j := range CollectionService {
				if j == "receiptnumber-servicebased" {

				}
				if j == "receipt-search-paginate" {

				}
				if j == "receipt-search-defaultsize" {

				}
				if j == "user-create-enabled" {

				}
			}
		}
		if i == "collection-receipt-voucher-consumer" {
			Voucher := data[i].(map[string]interface{})
			for j := range Voucher {
				if j == "jalandhar-erp-host" {

				}
				if j == "mohali-erp-host" {

				}
				if j == "nayagaon-erp-host" {

				}
				if j == "amritsar-erp-host" {

				}
				if j == "kharar-erp-host" {

				}
				if j == "zirakpur-erp-host" {

				}
			}
		}
		if i == "finance-collections-voucher-consumer" {
			FinanceCollection := data[i].(map[string]interface{})
			for j := range FinanceCollection {
				if j == "erp-env-name" {

				}
				if j == "erp-domain-name" {

				}
			}
		}
		if i == "rainmaker-pgr" {
			Rainmaker := data[i].(map[string]interface{})
			for j := range Rainmaker {
				if j == "notification-sms-enabled" {

				}
				if j == "notification-email-enabled" {

				}
				if j == "new-complaint-enabled" {

				}
				if j == "reassign-complaint-enabled" {

				}
				if j == "reopen-complaint-enabled" {

				}
				if j == "comment-by-employee-notif-enabled" {

				}
				if j == "notification-allowed-status" {

				}
			}
		}
		if i == "pt-services-v2" {
			PropertyServices := data[i].(map[string]interface{})
			for j := range PropertyServices {
				if j == "pt-userevents-pay-link" {

				}
			}
		}
		if i == "pt-calculator-v2" {
			Calculator := data[i].(map[string]interface{})
			for j := range Calculator {
				if j == "logging-level" {

				}
			}
		}
		if i == "tl-services" {
			TlServices := data[i].(map[string]interface{})
			for j := range TlServices {
				if j == "heap" {

				}
				if j == "memory_limits" {

				}
				if j == "java-args" {

				}
				if j == "tl-application-num-format" {

				}
				if j == "tl-license-num-format" {

				}
				if j == "tl-userevents-pay-link" {

				}
				if j == "tl-payment-topic-name" {

				}
				if j == "host-link" {

				}
				if j == "pdf-link" {

				}
				if j == "tl-search-default-limit" {

				}
			}
		}
		if i == "egov-hrms" {
			EgovHrms := data[i].(map[string]interface{})
			for j := range EgovHrms {
				if j == "java-args" {

				}
				if j == "heap" {

				}
				if j == "employee-applink" {

				}
			}
		}
		if i == "egov-weekly-impact-notifier" {
			EgovNotifier := data[i].(map[string]interface{})
			for j := range EgovNotifier {
				if j == "mail-to-address" {

				}
				if j == "mail-interval-in-secs" {

				}
				if j == "schedule" {

				}
			}
		}
		if i == "kafka-config" {
			Kafka := data[i].(map[string]interface{})
			for j := range Kafka {
				if j == "topics" {

				}
				if j == "zookeeper-connect" {

				}
				if j == "kafka-brokers" {

				}
			}
		}
		if i == "logging-config" {
			Logging := data[i].(map[string]interface{})
			for j := range Logging {
				if j == "es-host" {

				}
				if j == "es-port" {

				}
			}
		}
		if i == "jaeger-config" {
			Jaeger := data[i].(map[string]interface{})
			for j := range Jaeger {
				if j == "host" {

				}
				if j == "port" {

				}
				if j == "sampler-type" {

				}
				if j == "sampler-param" {

				}
				if j == "sampling-strategies" {

				}
			}
		}
		if i == "redis" {
			Redis := data[i].(map[string]interface{})
			for j := range Redis {
				if j == "replicas" {

				}
				if j == "images" {

				}
			}
		}
		if i == "playground" {
			playground := data[i].(map[string]interface{})
			for j := range playground {
				if j == "replicas" {

				}
				if j == "images" {

				}
			}
		}
		if i == "fluent-bit" {
			FuentBit := data[i].(map[string]interface{})
			for j := range FuentBit {
				if j == "images" {

				}
				if j == "egov-services-log-topic" {

				}
				if j == "egov-infra-log-topic" {

				}
			}
		}
		if i == "egov-workflow-v2" {
			EgovWorkflow := data[i].(map[string]interface{})
			for j := range EgovWorkflow {
				if j == "logging-level" {

				}
				if j == "java-args" {

				}
				if j == "heap" {

				}
				if j == "workflow-statelevel" {

				}
				if j == "host-link" {

				}
				if j == "pdf-link" {

				}
			}
		}
	}
	DataMap["global"] = data["global"]
	DataMap["cluster-configs"] = data["cluster-configs"]
	DataMap["employee"] = data["employee"]
	DataMap["citizen"] = data["citizen"]
	DataMap["digit-ui"] = data["digit-ui"]
	DataMap["egov-filestore"] = data["egov-filestore"]
	DataMap["egov-idgen"] = data["egov-idgen"]
	DataMap["egov-user"] = data["egov-user"]
	DataMap["egov-indexer"] = data["egov-indexer"]
	DataMap["egov-persister"] = data["egov-persister"]
	DataMap["egov-data-uploader"] = data["egov-data-uploader"]
	DataMap["egov-searcher"] = data["egov-searcher"]
	DataMap["report"] = data["report"]
	DataMap["pdf-service"] = data["pdf-service"]
	DataMap["egf-master"] = data["egf-master"]
	DataMap["egov-custom-consumer"] = data["egov-custom-consumer"]
	DataMap["egov-apportion-service"] = data["egov-apportion-service"]
	DataMap["redoc"] = data["redoc"]
	DataMap["nginx-ingress"] = data["nginx-ingress"]
	DataMap["cert-manager"] = data["cert-manager"]
	DataMap["zuul"] = data["zuul"]
	DataMap["collection-services"] = data["collection-services"]
	DataMap["collection-receipt-voucher-consumer"] = data["collection-receipt-voucher-consumer"]
	DataMap["finance-collections-voucher-consumer"] = data["finance-collections-voucher-consumer"]
	DataMap["egov-workflow-v2"] = data["egov-workflow-v2"]
	DataMap["egov-hrms"] = data["egov-hrms"]
	DataMap["egov-weekly-impact-notifier"] = data["egov-weekly-impact-notifier"]
	DataMap["kafka-config"] = data["kafka-config"]
	DataMap["logging-config"] = data["logging-config"]
	DataMap["jaeger-config"] = data["jaeger-config"]
	DataMap["redis"] = data["redis"]
	DataMap["playground"] = data["playground"]
	DataMap["fluent-bit"] = data["fluent-bit"]
	DataMap["kafka-v2"] = data["kafka-v2"]
	DataMap["zookeeper-v2"] = data["zookeeper-v2"]
	DataMap["elasticsearch-data-v1"] = data["elasticsearch-data-v1"]
	DataMap["elasticsearch-master-v1"] = data["elasticsearch-master-v1"]
	DataMap["es-curator"] = data["es-curator"]
	for i := range modules {
		if modules[i] == "m_pgr" {
			DataMap["egov-pg-service"] = data["egov-pg-service"]
			DataMap["rainmaker-pgr"] = data["rainmaker-pgr"]
		}
		if modules[i] == "m_property-tax" {
			DataMap["pt-services-v2"] = data["pt-services-v2"]
			DataMap["pt-calculator-v2"] = data["pt-calculator-v2"]
		}
		if modules[i] == "m_sewerage" {
			DataMap["sw-services"] = data["sw-services"]
		}
		if modules[i] == "m_bpa" {
			DataMap["bpa-services"] = data["bpa-services"]
			DataMap["bpa-calculator"] = data["bpa-calculator"]
		}
		if modules[i] == "m_trade-license" {
			DataMap["tl-services"] = data["tl-services"]
		}
		if modules[i] == "m_firenoc" {

		}
		if modules[i] == "m_water-service" {
			DataMap["ws-services"] = data["ws-services"]
		}
		if modules[i] == "m_dss" {
			DataMap["dashboard-analytics"] = data["dashboard-analytics"]
			DataMap["dashboard-ingest"] = data["dashboard-ingest"]
		}
		if modules[i] == "m_fsm" {

		}
		if modules[i] == "m_echallan" {

		}
		if modules[i] == "m_edcr" {

		}
		if modules[i] == "m_finance" {

		}
	}
	newfile, err := yaml.Marshal(&DataMap)
	if err != nil {
		log.Printf("%v", err)

	}
	filename := fmt.Sprintf("../../config-as-code/environments/%s.yaml", Config["file_name"])
	err = ioutil.WriteFile(filename, newfile, 0644)
	if err != nil {
		log.Printf("%v", err)
	}
}

//secrets config
func SecretFile(cluster_name string, Ssh string, modules []string) {
	var sec Secret
	secret, err := ioutil.ReadFile("DIGIT-DevOps/config-as-code/environments/egov-demo-secrets.yaml")
	if err != nil {
		log.Printf("%v", err)
	}
	err = yaml.Unmarshal(secret, &sec)
	if err != nil {
		log.Printf("%v", err)
	}
	secFilename := fmt.Sprintf("../../config-as-code/environments/%s-secrets.yaml", cluster_name)
	var Sec_state Secret
	if _, err := os.Stat("secrectstate.yaml"); err == nil {
		state, err := ioutil.ReadFile("secrectstate.yaml")
		if err != nil {
			log.Printf("%v", err)
		}
		err = yaml.Unmarshal(state, &Sec_state)
		if err != nil {
			log.Printf("%v", err)
		}
	}
	// fmt.Println("The secret list")
	// fmt.Println(St)
	// eUsername := sec.ClusterConfigs.Secrets.Db.Username
	// fmt.Println(eUsername)
	var Db_Username string
	var Db_Password string
	var Db_FlywayUsername string
	var Db_FlywayPassword string
	var EgovNotificationSms_Password string
	var EgovFilestore_AwsKey string
	var EgovFilestore_AwsSecretKey string
	var EgovLocation_Gmapskey string
	var EgovPgService_AxisMerchantID string
	var EgovPgService_AxisMerchantSecretKey string
	var EgovPgService_AxisMerchantUser string
	var EgovPgService_AxisMerchantPwd string
	var EgovPgService_AxisMerchantAccessCode string
	var EgovPgService_PayuMerchantKey string
	var EgovPgService_PayuMerchantSalt string
	var Pgadmin_AdminEmail string
	var Pgadmin_AdminPassword string
	var Pgadmin_ReadEmail string
	var Pgadmin_ReadPassword string
	var EgovEncService_MasterPassword string
	var EgovEncService_MasterSalt string
	var EgovEncService_MasterInitialvector string
	var EgovNotificationMail_Mailsenderusername string
	var EgovNotificationMail_Mailsenderpassword string
	var Kibana_Namespace string
	var Kibana_Credentials string
	var EgovSiMicroservice_SiMicroserviceUser string
	var EgovSiMicroservice_SiMicroservicePassword string
	var EgovSiMicroservice_MailSenderPassword string
	var EgovEdcrNotification_EdcrMailUsername string
	var EgovEdcrNotification_EdcrMailPassword string
	var EgovEdcrNotification_EdcrSmsUsername string
	var EgovEdcrNotification_EdcrSmsPassword string
	var Chatbot_ValuefirstUsername string
	var Chatbot_ValuefirstPassword string
	var EgovUserChatbot_CitizenLoginPasswordOtpFixedValue string
	var Oauth2Proxy_ClientID string
	var Oauth2Proxy_ClientSecret string
	var Oauth2Proxy_CookieSecret string
	var EgovNotificationSms_Username string

	Username := sec.ClusterConfigs.Secrets.Db.Username
	Password := sec.ClusterConfigs.Secrets.Db.Password
	FlywayUsername := sec.ClusterConfigs.Secrets.Db.FlywayUsername
	FlywayPassword := sec.ClusterConfigs.Secrets.Db.FlywayPassword
	NotUsername := sec.ClusterConfigs.Secrets.EgovNotificationSms.Username
	NotPassword := sec.ClusterConfigs.Secrets.EgovNotificationSms.Password
	AwsKey := sec.ClusterConfigs.Secrets.EgovFilestore.AwsKey
	AwsSecretKey := sec.ClusterConfigs.Secrets.EgovFilestore.AwsSecretKey
	Gmapskey := sec.ClusterConfigs.Secrets.EgovLocation.Gmapskey
	AxisMerchantID := sec.ClusterConfigs.Secrets.EgovPgService.AxisMerchantID
	AxisMerchantSecretKey := sec.ClusterConfigs.Secrets.EgovPgService.AxisMerchantSecretKey
	AxisMerchantUser := sec.ClusterConfigs.Secrets.EgovPgService.AxisMerchantUser
	AxisMerchantPwd := sec.ClusterConfigs.Secrets.EgovPgService.AxisMerchantPwd
	AxisMerchantAccessCode := sec.ClusterConfigs.Secrets.EgovPgService.AxisMerchantAccessCode
	PayuMerchantKey := sec.ClusterConfigs.Secrets.EgovPgService.PayuMerchantKey
	PayuMerchantSalt := sec.ClusterConfigs.Secrets.EgovPgService.PayuMerchantSalt
	AdminEmail := sec.ClusterConfigs.Secrets.Pgadmin.AdminEmail
	AdminPassword := sec.ClusterConfigs.Secrets.Pgadmin.AdminPassword
	ReadEmail := sec.ClusterConfigs.Secrets.Pgadmin.ReadEmail
	ReadPassword := sec.ClusterConfigs.Secrets.Pgadmin.ReadPassword
	MasterPassword := sec.ClusterConfigs.Secrets.EgovEncService.MasterPassword
	MasterSalt := sec.ClusterConfigs.Secrets.EgovEncService.MasterSalt
	MasterInitialvector := sec.ClusterConfigs.Secrets.EgovEncService.MasterInitialvector
	Mailsenderusername := sec.ClusterConfigs.Secrets.EgovNotificationMail.Mailsenderusername
	Mailsenderpassword := sec.ClusterConfigs.Secrets.EgovNotificationMail.Mailsenderpassword
	KnownHosts := sec.ClusterConfigs.Secrets.GitSync.KnownHosts
	Namespace := sec.ClusterConfigs.Secrets.Kibana.Namespace
	Credentials := sec.ClusterConfigs.Secrets.Kibana.Credentials
	SiMicroserviceUser := sec.ClusterConfigs.Secrets.EgovSiMicroservice.SiMicroserviceUser
	SiMicroservicePassword := sec.ClusterConfigs.Secrets.EgovSiMicroservice.SiMicroservicePassword
	MailSenderPassword := sec.ClusterConfigs.Secrets.EgovSiMicroservice.MailSenderPassword
	EdcrMailUsername := sec.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrMailUsername
	EdcrMailPassword := sec.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrMailPassword
	EdcrSmsUsername := sec.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrSmsUsername
	EdcrSmsPassword := sec.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrSmsPassword
	ValuefirstUsername := sec.ClusterConfigs.Secrets.Chatbot.ValuefirstUsername
	ValuefirstPassword := sec.ClusterConfigs.Secrets.Chatbot.ValuefirstPassword
	CitizenLoginPasswordOtpFixedValue := sec.ClusterConfigs.Secrets.EgovUserChatbot.CitizenLoginPasswordOtpFixedValue
	ClientID := sec.ClusterConfigs.Secrets.Oauth2Proxy.ClientID
	ClientSecret := sec.ClusterConfigs.Secrets.Oauth2Proxy.ClientSecret
	CookieSecret := sec.ClusterConfigs.Secrets.Oauth2Proxy.CookieSecret
	if Sec_state.ClusterConfigs.Secrets.Db.Username != "" {
		Sec_state.ClusterConfigs.Secrets.Db.Username = Sec_state.ClusterConfigs.Secrets.Db.Username
	} else {
		fmt.Println("Enter Db_Username:")
		fmt.Scanln(&Db_Username)
		Sec_state.ClusterConfigs.Secrets.Db.Username = Db_Username
	}
	if Db_Username != "" {
		sec.ClusterConfigs.Secrets.Db.Username = Db_Username
	} else {
		sec.ClusterConfigs.Secrets.Db.Username = Username
	}
	if Sec_state.ClusterConfigs.Secrets.Db.Password != "" {
		Sec_state.ClusterConfigs.Secrets.Db.Password = Sec_state.ClusterConfigs.Secrets.Db.Password
	} else {
		fmt.Println("Enter Db_Password:")
		fmt.Scanln(&Db_Password)
		Sec_state.ClusterConfigs.Secrets.Db.Password = Db_Password
	}
	if Db_Password != "" {
		sec.ClusterConfigs.Secrets.Db.Password = Db_Password
	} else {
		sec.ClusterConfigs.Secrets.Db.Password = Password
	}
	if Sec_state.ClusterConfigs.Secrets.Db.FlywayUsername != "" {
		Sec_state.ClusterConfigs.Secrets.Db.FlywayUsername = Sec_state.ClusterConfigs.Secrets.Db.FlywayUsername
	} else {
		fmt.Println("Enter Db_FlywayUsername:")
		fmt.Scanln(&Db_FlywayUsername)
		Sec_state.ClusterConfigs.Secrets.Db.FlywayUsername = Db_FlywayUsername
	}
	if Db_FlywayUsername != "" {
		sec.ClusterConfigs.Secrets.Db.FlywayUsername = Db_FlywayUsername
	} else {
		sec.ClusterConfigs.Secrets.Db.FlywayUsername = FlywayUsername
	}
	if Sec_state.ClusterConfigs.Secrets.Db.FlywayPassword != "" {
		Sec_state.ClusterConfigs.Secrets.Db.FlywayPassword = Sec_state.ClusterConfigs.Secrets.Db.FlywayPassword
	} else {
		fmt.Println("Enter Db_FlywayPassword:")
		fmt.Scanln(&Db_FlywayPassword)
		Sec_state.ClusterConfigs.Secrets.Db.FlywayPassword = Db_FlywayPassword
	}
	if Db_FlywayPassword != "" {
		sec.ClusterConfigs.Secrets.Db.FlywayPassword = Db_FlywayPassword
	} else {
		sec.ClusterConfigs.Secrets.Db.FlywayPassword = FlywayPassword
	}
	for i := range modules {
		if modules[i] == "m_property-tax" {
			if Sec_state.ClusterConfigs.Secrets.EgovNotificationSms.Username != "" {
				Sec_state.ClusterConfigs.Secrets.EgovNotificationSms.Username = Sec_state.ClusterConfigs.Secrets.EgovNotificationSms.Username
			} else {
				fmt.Println("Enter EgovNotificationSms_Username:")
				fmt.Scanln(&EgovNotificationSms_Username)
				Sec_state.ClusterConfigs.Secrets.EgovNotificationSms.Username = EgovNotificationSms_Username
			}
			if EgovNotificationSms_Username != "" {
				sec.ClusterConfigs.Secrets.EgovNotificationSms.Username = EgovNotificationSms_Username
			} else {
				sec.ClusterConfigs.Secrets.EgovNotificationSms.Username = NotUsername
			}
			if Sec_state.ClusterConfigs.Secrets.EgovNotificationSms.Password != "" {
				Sec_state.ClusterConfigs.Secrets.EgovNotificationSms.Password = Sec_state.ClusterConfigs.Secrets.EgovNotificationSms.Password
			} else {
				fmt.Println("Enter EgovNotificationSms_Password:")
				fmt.Scanln(&EgovNotificationSms_Password)
				Sec_state.ClusterConfigs.Secrets.EgovNotificationSms.Password = EgovNotificationSms_Password
			}
			if EgovNotificationSms_Password != "" {
				sec.ClusterConfigs.Secrets.EgovNotificationSms.Password = EgovNotificationSms_Password
			} else {
				sec.ClusterConfigs.Secrets.EgovNotificationSms.Password = NotPassword
			}
			if Sec_state.ClusterConfigs.Secrets.EgovFilestore.AwsKey != "" {
				Sec_state.ClusterConfigs.Secrets.EgovFilestore.AwsKey = Sec_state.ClusterConfigs.Secrets.EgovFilestore.AwsKey
			} else {
				fmt.Println("Enter EgovFilestore_AwsKey:")
				fmt.Scanln(&EgovFilestore_AwsKey)
				Sec_state.ClusterConfigs.Secrets.EgovFilestore.AwsKey = EgovFilestore_AwsKey
			}
			if EgovFilestore_AwsKey != "" {
				sec.ClusterConfigs.Secrets.EgovFilestore.AwsKey = EgovFilestore_AwsKey
			} else {
				sec.ClusterConfigs.Secrets.EgovFilestore.AwsKey = AwsKey
			}
			if Sec_state.ClusterConfigs.Secrets.EgovFilestore.AwsSecretKey != "" {
				Sec_state.ClusterConfigs.Secrets.EgovFilestore.AwsSecretKey = Sec_state.ClusterConfigs.Secrets.EgovFilestore.AwsSecretKey
			} else {
				fmt.Println("Enter EgovFilestore_AwsSecretKey:")
				fmt.Scanln(&EgovFilestore_AwsSecretKey)
				Sec_state.ClusterConfigs.Secrets.EgovFilestore.AwsSecretKey = EgovFilestore_AwsSecretKey
			}
			if EgovFilestore_AwsSecretKey != "" {
				sec.ClusterConfigs.Secrets.EgovFilestore.AwsSecretKey = EgovFilestore_AwsSecretKey
			} else {
				sec.ClusterConfigs.Secrets.EgovFilestore.AwsSecretKey = AwsSecretKey
			}
			if Sec_state.ClusterConfigs.Secrets.EgovLocation.Gmapskey != "" {
				Sec_state.ClusterConfigs.Secrets.EgovLocation.Gmapskey = Sec_state.ClusterConfigs.Secrets.EgovLocation.Gmapskey
			} else {
				fmt.Println("Enter EgovLocation_Gmapskey:")
				fmt.Scanln(&EgovLocation_Gmapskey)
				Sec_state.ClusterConfigs.Secrets.EgovLocation.Gmapskey = EgovLocation_Gmapskey
			}
			if EgovLocation_Gmapskey != "" {
				sec.ClusterConfigs.Secrets.EgovLocation.Gmapskey = EgovLocation_Gmapskey
			} else {
				sec.ClusterConfigs.Secrets.EgovLocation.Gmapskey = Gmapskey
			}
			if Sec_state.ClusterConfigs.Secrets.EgovPgService.AxisMerchantID != "" {
				Sec_state.ClusterConfigs.Secrets.EgovPgService.AxisMerchantID = Sec_state.ClusterConfigs.Secrets.EgovPgService.AxisMerchantID
			} else {
				fmt.Println("Enter EgovPgService_AxisMerchantID:")
				fmt.Scanln(&EgovPgService_AxisMerchantID)
				Sec_state.ClusterConfigs.Secrets.EgovPgService.AxisMerchantID = EgovPgService_AxisMerchantID
			}
			if EgovPgService_AxisMerchantID != "" {
				sec.ClusterConfigs.Secrets.EgovPgService.AxisMerchantID = EgovPgService_AxisMerchantID
			} else {
				sec.ClusterConfigs.Secrets.EgovPgService.AxisMerchantID = AxisMerchantID
			}
			if Sec_state.ClusterConfigs.Secrets.EgovPgService.AxisMerchantSecretKey != "" {
				Sec_state.ClusterConfigs.Secrets.EgovPgService.AxisMerchantSecretKey = Sec_state.ClusterConfigs.Secrets.EgovPgService.AxisMerchantSecretKey
			} else {
				fmt.Println("Enter EgovPgService_AxisMerchantSecretKey:")
				fmt.Scanln(&EgovPgService_AxisMerchantSecretKey)
				Sec_state.ClusterConfigs.Secrets.EgovPgService.AxisMerchantSecretKey = EgovPgService_AxisMerchantSecretKey
			}
			if EgovPgService_AxisMerchantSecretKey != "" {
				sec.ClusterConfigs.Secrets.EgovPgService.AxisMerchantSecretKey = EgovPgService_AxisMerchantSecretKey
			} else {
				sec.ClusterConfigs.Secrets.EgovPgService.AxisMerchantSecretKey = AxisMerchantSecretKey
			}
			if Sec_state.ClusterConfigs.Secrets.EgovPgService.AxisMerchantUser != "" {
				Sec_state.ClusterConfigs.Secrets.EgovPgService.AxisMerchantUser = Sec_state.ClusterConfigs.Secrets.EgovPgService.AxisMerchantUser
			} else {
				fmt.Println("Enter EgovPgService_AxisMerchantUser:")
				fmt.Scanln(&EgovPgService_AxisMerchantUser)
				Sec_state.ClusterConfigs.Secrets.EgovPgService.AxisMerchantUser = EgovPgService_AxisMerchantUser
			}
			if EgovPgService_AxisMerchantUser != "" {
				sec.ClusterConfigs.Secrets.EgovPgService.AxisMerchantUser = EgovPgService_AxisMerchantUser
			} else {
				sec.ClusterConfigs.Secrets.EgovPgService.AxisMerchantUser = AxisMerchantUser
			}
			if Sec_state.ClusterConfigs.Secrets.EgovPgService.AxisMerchantPwd != "" {
				Sec_state.ClusterConfigs.Secrets.EgovPgService.AxisMerchantPwd = Sec_state.ClusterConfigs.Secrets.EgovPgService.AxisMerchantPwd
			} else {
				fmt.Println("Enter EgovPgService_AxisMerchantPwd:")
				fmt.Scanln(&EgovPgService_AxisMerchantPwd)
				Sec_state.ClusterConfigs.Secrets.EgovPgService.AxisMerchantPwd = EgovPgService_AxisMerchantPwd
			}
			if EgovPgService_AxisMerchantPwd != "" {
				sec.ClusterConfigs.Secrets.EgovPgService.AxisMerchantPwd = EgovPgService_AxisMerchantPwd
			} else {
				sec.ClusterConfigs.Secrets.EgovPgService.AxisMerchantPwd = AxisMerchantPwd
			}
			if Sec_state.ClusterConfigs.Secrets.EgovPgService.AxisMerchantAccessCode != "" {
				Sec_state.ClusterConfigs.Secrets.EgovPgService.AxisMerchantAccessCode = Sec_state.ClusterConfigs.Secrets.EgovPgService.AxisMerchantAccessCode
			} else {
				fmt.Println("Enter EgovPgService_AxisMerchantAccessCode:")
				fmt.Scanln(&EgovPgService_AxisMerchantAccessCode)
				Sec_state.ClusterConfigs.Secrets.EgovPgService.AxisMerchantAccessCode = EgovPgService_AxisMerchantAccessCode
			}
			if EgovPgService_AxisMerchantAccessCode != "" {
				sec.ClusterConfigs.Secrets.EgovPgService.AxisMerchantAccessCode = EgovPgService_AxisMerchantAccessCode
			} else {
				sec.ClusterConfigs.Secrets.EgovPgService.AxisMerchantAccessCode = AxisMerchantAccessCode
			}
			if Sec_state.ClusterConfigs.Secrets.EgovPgService.PayuMerchantKey != "" {
				Sec_state.ClusterConfigs.Secrets.EgovPgService.PayuMerchantKey = Sec_state.ClusterConfigs.Secrets.EgovPgService.PayuMerchantKey
			} else {
				fmt.Println("Enter EgovPgService_PayuMerchantKey:")
				fmt.Scanln(&EgovPgService_PayuMerchantKey)
				Sec_state.ClusterConfigs.Secrets.EgovPgService.PayuMerchantKey = EgovPgService_PayuMerchantKey
			}
			if EgovPgService_PayuMerchantKey != "" {
				sec.ClusterConfigs.Secrets.EgovPgService.PayuMerchantKey = EgovPgService_PayuMerchantKey
			} else {
				sec.ClusterConfigs.Secrets.EgovPgService.PayuMerchantKey = PayuMerchantKey
			}
			if Sec_state.ClusterConfigs.Secrets.EgovPgService.PayuMerchantSalt != "" {
				Sec_state.ClusterConfigs.Secrets.EgovPgService.PayuMerchantSalt = Sec_state.ClusterConfigs.Secrets.EgovPgService.PayuMerchantSalt
			} else {
				fmt.Println("Enter EgovPgService_PayuMerchantSalt:")
				fmt.Scanln(&EgovPgService_PayuMerchantSalt)
				Sec_state.ClusterConfigs.Secrets.EgovPgService.PayuMerchantSalt = EgovPgService_PayuMerchantSalt
			}
			if EgovPgService_PayuMerchantSalt != "" {
				sec.ClusterConfigs.Secrets.EgovPgService.PayuMerchantSalt = EgovPgService_PayuMerchantSalt
			} else {
				sec.ClusterConfigs.Secrets.EgovPgService.PayuMerchantSalt = PayuMerchantSalt
			}
			if Sec_state.ClusterConfigs.Secrets.Pgadmin.AdminEmail != "" {
				Sec_state.ClusterConfigs.Secrets.Pgadmin.AdminEmail = Sec_state.ClusterConfigs.Secrets.Pgadmin.AdminEmail
			} else {
				fmt.Println("Enter Pgadmin_AdminEmail:")
				fmt.Scanln(&Pgadmin_AdminEmail)
				Sec_state.ClusterConfigs.Secrets.Pgadmin.AdminEmail = Pgadmin_AdminEmail
			}
			if Pgadmin_AdminEmail != "" {
				sec.ClusterConfigs.Secrets.Pgadmin.AdminEmail = Pgadmin_AdminEmail
			} else {
				sec.ClusterConfigs.Secrets.Pgadmin.AdminEmail = AdminEmail
			}
			if Sec_state.ClusterConfigs.Secrets.Pgadmin.AdminPassword != "" {
				Sec_state.ClusterConfigs.Secrets.Pgadmin.AdminPassword = Sec_state.ClusterConfigs.Secrets.Pgadmin.AdminPassword
			} else {
				fmt.Println("Enter Pgadmin_AdminPassword:")
				fmt.Scanln(&Pgadmin_AdminPassword)
				Sec_state.ClusterConfigs.Secrets.Pgadmin.AdminPassword = Pgadmin_AdminPassword
			}
			if Pgadmin_AdminPassword != "" {
				sec.ClusterConfigs.Secrets.Pgadmin.AdminPassword = Pgadmin_AdminPassword
			} else {
				sec.ClusterConfigs.Secrets.Pgadmin.AdminPassword = AdminPassword
			}
			if Sec_state.ClusterConfigs.Secrets.Pgadmin.ReadEmail != "" {
				Sec_state.ClusterConfigs.Secrets.Pgadmin.ReadEmail = Sec_state.ClusterConfigs.Secrets.Pgadmin.ReadEmail
			} else {
				fmt.Println("Enter Pgadmin_ReadEmail:")
				fmt.Scanln(&Pgadmin_ReadEmail)
				Sec_state.ClusterConfigs.Secrets.Pgadmin.ReadEmail = Pgadmin_ReadEmail
			}
			if Pgadmin_ReadEmail != "" {
				sec.ClusterConfigs.Secrets.Pgadmin.ReadEmail = Pgadmin_ReadEmail
			} else {
				sec.ClusterConfigs.Secrets.Pgadmin.ReadEmail = ReadEmail
			}
			if Sec_state.ClusterConfigs.Secrets.Pgadmin.ReadPassword != "" {
				Sec_state.ClusterConfigs.Secrets.Pgadmin.ReadPassword = Sec_state.ClusterConfigs.Secrets.Pgadmin.ReadPassword
			} else {
				fmt.Println("Enter Pgadmin_ReadPassword:")
				fmt.Scanln(&Pgadmin_ReadPassword)
				Sec_state.ClusterConfigs.Secrets.Pgadmin.ReadPassword = Pgadmin_ReadPassword
			}
			if Pgadmin_ReadPassword != "" {
				sec.ClusterConfigs.Secrets.Pgadmin.ReadPassword = Pgadmin_ReadPassword
			} else {
				sec.ClusterConfigs.Secrets.Pgadmin.ReadPassword = ReadPassword
			}
		}
	}
	if Sec_state.ClusterConfigs.Secrets.EgovEncService.MasterPassword != "" {
		Sec_state.ClusterConfigs.Secrets.EgovEncService.MasterPassword = Sec_state.ClusterConfigs.Secrets.EgovEncService.MasterPassword
	} else {
		fmt.Println("Enter EgovEncService_MasterPassword:")
		fmt.Scanln(&EgovEncService_MasterPassword)
		Sec_state.ClusterConfigs.Secrets.EgovEncService.MasterPassword = EgovEncService_MasterPassword
	}
	if EgovEncService_MasterPassword != "" {
		sec.ClusterConfigs.Secrets.EgovEncService.MasterPassword = EgovEncService_MasterPassword
	} else {
		sec.ClusterConfigs.Secrets.EgovEncService.MasterPassword = MasterPassword
	}
	if Sec_state.ClusterConfigs.Secrets.EgovEncService.MasterSalt != "" {
		Sec_state.ClusterConfigs.Secrets.EgovEncService.MasterSalt = Sec_state.ClusterConfigs.Secrets.EgovEncService.MasterSalt
	} else {
		fmt.Println("Enter EgovEncService_MasterSalt:")
		fmt.Scanln(&EgovEncService_MasterSalt)
		Sec_state.ClusterConfigs.Secrets.EgovEncService.MasterSalt = EgovEncService_MasterSalt
	}
	if EgovEncService_MasterSalt != "" {
		sec.ClusterConfigs.Secrets.EgovEncService.MasterSalt = EgovEncService_MasterSalt
	} else {
		sec.ClusterConfigs.Secrets.EgovEncService.MasterSalt = MasterSalt
	}
	if Sec_state.ClusterConfigs.Secrets.EgovEncService.MasterInitialvector != "" {
		Sec_state.ClusterConfigs.Secrets.EgovEncService.MasterInitialvector = Sec_state.ClusterConfigs.Secrets.EgovEncService.MasterInitialvector
	} else {
		fmt.Println("Enter EgovEncService_MasterInitialvector:")
		fmt.Scanln(&EgovEncService_MasterInitialvector)
		Sec_state.ClusterConfigs.Secrets.EgovEncService.MasterInitialvector = EgovEncService_MasterInitialvector
	}
	if EgovEncService_MasterInitialvector != "" {
		sec.ClusterConfigs.Secrets.EgovEncService.MasterInitialvector = EgovEncService_MasterInitialvector
	} else {
		sec.ClusterConfigs.Secrets.EgovEncService.MasterInitialvector = MasterInitialvector
	}
	if Sec_state.ClusterConfigs.Secrets.EgovNotificationMail.Mailsenderusername != "" {
		Sec_state.ClusterConfigs.Secrets.EgovNotificationMail.Mailsenderusername = Sec_state.ClusterConfigs.Secrets.EgovNotificationMail.Mailsenderusername
	} else {
		fmt.Println("Enter EgovNotificationMail_Mailsenderusername:")
		fmt.Scanln(&EgovNotificationMail_Mailsenderusername)
		Sec_state.ClusterConfigs.Secrets.EgovNotificationMail.Mailsenderusername = EgovNotificationMail_Mailsenderusername
	}
	if EgovNotificationMail_Mailsenderusername != "" {
		sec.ClusterConfigs.Secrets.EgovNotificationMail.Mailsenderusername = EgovNotificationMail_Mailsenderusername
	} else {
		sec.ClusterConfigs.Secrets.EgovNotificationMail.Mailsenderusername = Mailsenderusername
	}
	if Sec_state.ClusterConfigs.Secrets.EgovNotificationMail.Mailsenderpassword != "" {
		Sec_state.ClusterConfigs.Secrets.EgovNotificationMail.Mailsenderpassword = Sec_state.ClusterConfigs.Secrets.EgovNotificationMail.Mailsenderpassword
	} else {
		fmt.Println("Enter EgovNotificationMail_Mailsenderpassword:")
		fmt.Scanln(&EgovNotificationMail_Mailsenderpassword)
		Sec_state.ClusterConfigs.Secrets.EgovNotificationMail.Mailsenderpassword = EgovNotificationMail_Mailsenderpassword
	}
	if EgovNotificationMail_Mailsenderpassword != "" {
		sec.ClusterConfigs.Secrets.EgovNotificationMail.Mailsenderpassword = EgovNotificationMail_Mailsenderpassword
	} else {
		sec.ClusterConfigs.Secrets.EgovNotificationMail.Mailsenderpassword = Mailsenderpassword
	}
	sec.ClusterConfigs.Secrets.GitSync.SSH = Ssh
	sec.ClusterConfigs.Secrets.GitSync.KnownHosts = KnownHosts
	if Sec_state.ClusterConfigs.Secrets.Kibana.Namespace != "" {
		Sec_state.ClusterConfigs.Secrets.Kibana.Namespace = Sec_state.ClusterConfigs.Secrets.Kibana.Namespace
	} else {
		fmt.Println("Enter Kibana_Namespace:")
		fmt.Scanln(&Kibana_Namespace)
		Sec_state.ClusterConfigs.Secrets.Kibana.Namespace = Kibana_Namespace
	}
	if Kibana_Namespace != "" {
		sec.ClusterConfigs.Secrets.Kibana.Namespace = Kibana_Namespace
	} else {
		sec.ClusterConfigs.Secrets.Kibana.Namespace = Namespace
	}
	if Sec_state.ClusterConfigs.Secrets.Kibana.Credentials != "" {
		Sec_state.ClusterConfigs.Secrets.Kibana.Credentials = Sec_state.ClusterConfigs.Secrets.Kibana.Credentials
	} else {
		fmt.Println("Enter Kibana_Credentials:")
		fmt.Scanln(&Kibana_Credentials)
		Sec_state.ClusterConfigs.Secrets.Kibana.Credentials = Kibana_Credentials
	}
	if Kibana_Credentials != "" {
		sec.ClusterConfigs.Secrets.Kibana.Credentials = Kibana_Credentials
	} else {
		sec.ClusterConfigs.Secrets.Kibana.Credentials = Credentials
	}
	if Sec_state.ClusterConfigs.Secrets.EgovSiMicroservice.SiMicroserviceUser != "" {
		Sec_state.ClusterConfigs.Secrets.EgovSiMicroservice.SiMicroserviceUser = Sec_state.ClusterConfigs.Secrets.EgovSiMicroservice.SiMicroserviceUser
	} else {
		fmt.Println("Enter EgovSiMicroservice_SiMicroserviceUser:")
		fmt.Scanln(&EgovSiMicroservice_SiMicroserviceUser)
		Sec_state.ClusterConfigs.Secrets.EgovSiMicroservice.SiMicroserviceUser = EgovSiMicroservice_SiMicroserviceUser
	}
	if EgovSiMicroservice_SiMicroserviceUser != "" {
		sec.ClusterConfigs.Secrets.EgovSiMicroservice.SiMicroserviceUser = EgovSiMicroservice_SiMicroserviceUser
	} else {
		sec.ClusterConfigs.Secrets.EgovSiMicroservice.SiMicroserviceUser = SiMicroserviceUser
	}
	if Sec_state.ClusterConfigs.Secrets.EgovSiMicroservice.SiMicroservicePassword != "" {
		Sec_state.ClusterConfigs.Secrets.EgovSiMicroservice.SiMicroservicePassword = Sec_state.ClusterConfigs.Secrets.EgovSiMicroservice.SiMicroservicePassword
	} else {
		fmt.Println("Enter EgovSiMicroservice_SiMicroservicePassword:")
		fmt.Scanln(&EgovSiMicroservice_SiMicroservicePassword)
		Sec_state.ClusterConfigs.Secrets.EgovSiMicroservice.SiMicroservicePassword = EgovSiMicroservice_SiMicroservicePassword
	}
	if EgovSiMicroservice_SiMicroservicePassword != "" {
		sec.ClusterConfigs.Secrets.EgovSiMicroservice.SiMicroservicePassword = EgovSiMicroservice_SiMicroservicePassword
	} else {
		sec.ClusterConfigs.Secrets.EgovSiMicroservice.SiMicroservicePassword = SiMicroservicePassword
	}
	if Sec_state.ClusterConfigs.Secrets.EgovSiMicroservice.MailSenderPassword != "" {
		Sec_state.ClusterConfigs.Secrets.EgovSiMicroservice.MailSenderPassword = Sec_state.ClusterConfigs.Secrets.EgovSiMicroservice.MailSenderPassword
	} else {
		fmt.Println("Enter EgovSiMicroservice_MailSenderPassword:")
		fmt.Scanln(&EgovSiMicroservice_MailSenderPassword)
		Sec_state.ClusterConfigs.Secrets.EgovSiMicroservice.MailSenderPassword = EgovSiMicroservice_MailSenderPassword
	}
	if EgovSiMicroservice_MailSenderPassword != "" {
		sec.ClusterConfigs.Secrets.EgovSiMicroservice.MailSenderPassword = EgovSiMicroservice_MailSenderPassword
	} else {
		sec.ClusterConfigs.Secrets.EgovSiMicroservice.MailSenderPassword = MailSenderPassword
	}
	if Sec_state.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrMailUsername != "" {
		Sec_state.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrMailUsername = Sec_state.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrMailUsername
	} else {
		fmt.Println("Enter EgovEdcrNotification_EdcrMailUsername:")
		fmt.Scanln(&EgovEdcrNotification_EdcrMailUsername)
		Sec_state.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrMailUsername = EgovEdcrNotification_EdcrMailUsername
	}
	if EgovEdcrNotification_EdcrMailUsername != "" {
		sec.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrMailUsername = EgovEdcrNotification_EdcrMailUsername
	} else {
		sec.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrMailUsername = EdcrMailUsername
	}
	if Sec_state.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrMailPassword != "" {
		Sec_state.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrMailPassword = Sec_state.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrMailPassword
	} else {
		fmt.Println("Enter EgovEdcrNotification_EdcrMailPassword:")
		fmt.Scanln(&EgovEdcrNotification_EdcrMailPassword)
		Sec_state.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrMailPassword = EgovEdcrNotification_EdcrMailPassword
	}
	if EgovEdcrNotification_EdcrMailPassword != "" {
		sec.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrMailPassword = EgovEdcrNotification_EdcrMailPassword
	} else {
		sec.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrMailPassword = EdcrMailPassword
	}
	if Sec_state.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrSmsUsername != "" {
		Sec_state.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrSmsUsername = Sec_state.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrSmsUsername
	} else {
		fmt.Println("Enter EgovEdcrNotification_EdcrSmsUsername:")
		fmt.Scanln(&EgovEdcrNotification_EdcrSmsUsername)
		Sec_state.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrSmsUsername = EgovEdcrNotification_EdcrSmsUsername
	}
	if EgovEdcrNotification_EdcrSmsUsername != "" {
		sec.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrSmsUsername = EgovEdcrNotification_EdcrSmsUsername
	} else {
		sec.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrSmsUsername = EdcrSmsUsername
	}
	if Sec_state.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrSmsPassword != "" {
		Sec_state.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrSmsPassword = Sec_state.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrSmsPassword
	} else {
		fmt.Println("Enter EgovEdcrNotification_EdcrSmsPassword:")
		fmt.Scanln(&EgovEdcrNotification_EdcrSmsPassword)
		Sec_state.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrSmsPassword = EgovEdcrNotification_EdcrSmsPassword
	}
	if EgovEdcrNotification_EdcrSmsPassword != "" {
		sec.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrSmsPassword = EgovEdcrNotification_EdcrSmsPassword
	} else {
		sec.ClusterConfigs.Secrets.EgovEdcrNotification.EdcrSmsPassword = EdcrSmsPassword
	}
	if Sec_state.ClusterConfigs.Secrets.Chatbot.ValuefirstUsername != "" {
		Sec_state.ClusterConfigs.Secrets.Chatbot.ValuefirstUsername = Sec_state.ClusterConfigs.Secrets.Chatbot.ValuefirstUsername
	} else {
		fmt.Println("Enter Chatbot_ValuefirstUsername:")
		fmt.Scanln(&Chatbot_ValuefirstUsername)
		Sec_state.ClusterConfigs.Secrets.Chatbot.ValuefirstUsername = Chatbot_ValuefirstUsername
	}
	if Chatbot_ValuefirstUsername != "" {
		sec.ClusterConfigs.Secrets.Chatbot.ValuefirstUsername = Chatbot_ValuefirstUsername
	} else {
		sec.ClusterConfigs.Secrets.Chatbot.ValuefirstUsername = ValuefirstUsername
	}
	if Sec_state.ClusterConfigs.Secrets.Chatbot.ValuefirstPassword != "" {
		Sec_state.ClusterConfigs.Secrets.Chatbot.ValuefirstPassword = Sec_state.ClusterConfigs.Secrets.Chatbot.ValuefirstPassword
	} else {
		fmt.Println("Enter Chatbot_ValuefirstPassword:")
		fmt.Scanln(&Chatbot_ValuefirstPassword)
		Sec_state.ClusterConfigs.Secrets.Chatbot.ValuefirstPassword = Chatbot_ValuefirstPassword
	}
	if Chatbot_ValuefirstPassword != "" {
		sec.ClusterConfigs.Secrets.Chatbot.ValuefirstPassword = Chatbot_ValuefirstPassword
	} else {
		sec.ClusterConfigs.Secrets.Chatbot.ValuefirstPassword = ValuefirstPassword
	}
	if Sec_state.ClusterConfigs.Secrets.EgovUserChatbot.CitizenLoginPasswordOtpFixedValue != "" {
		Sec_state.ClusterConfigs.Secrets.EgovUserChatbot.CitizenLoginPasswordOtpFixedValue = Sec_state.ClusterConfigs.Secrets.EgovUserChatbot.CitizenLoginPasswordOtpFixedValue
	} else {
		fmt.Println("Enter EgovUserChatbot_CitizenLoginPasswordOtpFixedValue:")
		fmt.Scanln(&EgovUserChatbot_CitizenLoginPasswordOtpFixedValue)
		Sec_state.ClusterConfigs.Secrets.EgovUserChatbot.CitizenLoginPasswordOtpFixedValue = EgovUserChatbot_CitizenLoginPasswordOtpFixedValue
	}
	if EgovUserChatbot_CitizenLoginPasswordOtpFixedValue != "" {
		sec.ClusterConfigs.Secrets.EgovUserChatbot.CitizenLoginPasswordOtpFixedValue = EgovUserChatbot_CitizenLoginPasswordOtpFixedValue
	} else {
		sec.ClusterConfigs.Secrets.EgovUserChatbot.CitizenLoginPasswordOtpFixedValue = CitizenLoginPasswordOtpFixedValue
	}
	if Sec_state.ClusterConfigs.Secrets.Oauth2Proxy.ClientID != "" {
		Sec_state.ClusterConfigs.Secrets.Oauth2Proxy.ClientID = Sec_state.ClusterConfigs.Secrets.Oauth2Proxy.ClientID
	} else {
		fmt.Println("Enter Oauth2Proxy_ClientID:")
		fmt.Scanln(&Oauth2Proxy_ClientID)
		Sec_state.ClusterConfigs.Secrets.Oauth2Proxy.ClientID = Oauth2Proxy_ClientID
	}
	if Oauth2Proxy_ClientID != "" {
		sec.ClusterConfigs.Secrets.Oauth2Proxy.ClientID = Oauth2Proxy_ClientID
	} else {
		sec.ClusterConfigs.Secrets.Oauth2Proxy.ClientID = ClientID
	}
	if Sec_state.ClusterConfigs.Secrets.Oauth2Proxy.ClientSecret != "" {
		Sec_state.ClusterConfigs.Secrets.Oauth2Proxy.ClientSecret = Sec_state.ClusterConfigs.Secrets.Oauth2Proxy.ClientSecret
	} else {
		fmt.Println("Enter Oauth2Proxy_ClientSecret:")
		fmt.Scanln(&Oauth2Proxy_ClientSecret)
		Sec_state.ClusterConfigs.Secrets.Oauth2Proxy.ClientSecret = Oauth2Proxy_ClientSecret
	}
	if Oauth2Proxy_ClientSecret != "" {
		sec.ClusterConfigs.Secrets.Oauth2Proxy.ClientSecret = Oauth2Proxy_ClientSecret
	} else {
		sec.ClusterConfigs.Secrets.Oauth2Proxy.ClientSecret = ClientSecret
	}
	if Sec_state.ClusterConfigs.Secrets.Oauth2Proxy.CookieSecret != "" {
		Sec_state.ClusterConfigs.Secrets.Oauth2Proxy.CookieSecret = Sec_state.ClusterConfigs.Secrets.Oauth2Proxy.CookieSecret
	} else {
		fmt.Println("Enter Oauth2Proxy_CookieSecret:")
		fmt.Scanln(&Oauth2Proxy_CookieSecret)
		Sec_state.ClusterConfigs.Secrets.Oauth2Proxy.CookieSecret = Oauth2Proxy_CookieSecret
	}
	if Oauth2Proxy_CookieSecret != "" {
		sec.ClusterConfigs.Secrets.Oauth2Proxy.CookieSecret = Oauth2Proxy_CookieSecret
	} else {
		sec.ClusterConfigs.Secrets.Oauth2Proxy.CookieSecret = CookieSecret
	}
	secretstate, err := yaml.Marshal(&Sec_state)
	if err != nil {
		log.Printf("%v", err)

	}
	err = ioutil.WriteFile("secrectstate.yaml", secretstate, 0644)
	if err != nil {
		log.Printf("%v", err)
	}
	secretsmar, err := yaml.Marshal(&sec)
	if err != nil {
		log.Printf("%v", err)

	}
	err = ioutil.WriteFile(secFilename, secretsmar, 0644)
	if err != nil {
		log.Printf("%v", err)
	}
}
