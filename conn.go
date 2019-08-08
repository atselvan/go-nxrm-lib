package nxrm

type ConnDetails struct {
	NexusURL string
	Username string
	Password string
}

func StoreConnectionDetails() {
	configuration := ConnDetails{NexusURL: NexusURL, Username: AuthUser.Username, Password: AuthUser.Password}
	configureJson, err := json.Marshal(configuration)
	logJsonMarshalError(err, jsonMarshalError)
	writeFile(ConfFileName, configureJson)
	log.Printf(connDetailsSuccessInfo, ConfFileName)
}

func getConnectionDetails() ConnDetails {
	var conf ConnDetails
	data := readFile(ConfFileName)
	err := json.Unmarshal([]byte(data), &conf)
	logJsonUnmarshalError(err, jsonUnmarshalError)
	return conf
}

func SetConnectionDetails() {
	if fileExists(ConfFileName) {
		conf := getConnectionDetails()
		NexusURL = conf.NexusURL
		AuthUser.Username = conf.Username
		AuthUser.Password = conf.Password
	} else {
		log.Printf(connDetailsEmptyInfo, "nexus3-repository-cli configure")
		os.Exit(1)
	}
}
