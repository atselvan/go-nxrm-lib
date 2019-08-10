package nxrm

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	regexp2 "regexp"
	"strings"
)

type Repository struct {
	Name       string     `json:"name"`
	URL        string     `json:"url"`
	Type       string     `json:"type"`
	Format     string     `json:"format"`
	Recipe     string     `json:"recipe"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	Storage       Storage       `json:"storage"`
	Maven         Maven         `json:"maven"`
	Proxy         Proxy         `json:"proxy"`
	Httpclient    HttpClient    `json:"httpclient"`
	Group         Group         `json:"group"`
	NegativeCache NegetiveCache `json:"negativeCache"`
	Docker        Docker        `json:"docker"`
	DockerProxy   DockerProxy   `json:"dockerProxy"`
	Cleanup       Cleanup       `json:"cleanup"`
}

type Storage struct {
	BlobStoreName               string `json:"blobStoreName"`
	WritePolicy                 string `json:"writePolicy"`
	StrictContentTypeValidation bool   `json:"strictContentTypeValidation"`
}

type Maven struct {
	VersionPolicy string `json:"versionPolicy"`
	LayoutPolicy  string `json:"layoutPolicy"`
}

type Proxy struct {
	RemoteURL      string `json:"remoteUrl"`
	ContentMaxAge  int    `json:"contentMaxAge"`
	MetadataMaxAge int    `json:"metadataMaxAge"`
}

type HttpClient struct {
	Blocked        bool           `json:"blocked"`
	AutoBlock      bool           `json:"autoBlock"`
	Authentication HttpClientAuth `json:"authentication"`
}

type HttpClientAuth struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Docker struct {
	HTTPPort       float64 `json:"httpPort"`
	HTTPSPort      float64 `json:"httpsPort"`
	ForceBasicAuth bool    `json:"forceBasicAuth"`
	V1Enabled      bool    `json:"v1Enabled"`
}

type DockerProxy struct {
	IndexType string `json:"indexType"`
}

type Group struct {
	MemberNames []string `json:"memberNames"`
}

type NegetiveCache struct {
	Enabled    bool    `json:"enabled"`
	TimeToLive float64 `json:"timeToLive"`
}

type Cleanup struct {
	PolicyName string `json:"policyName"`
}

func ListRepositories(name, format string) {
	var repositoryList []string

	if name != "" {
		repository := getRepository(name)
		fmt.Printf("Name: %s\nRecipe: %s\nURL: %s\n", repository.Name, repository.Recipe, repository.URL)
	} else if name == "" && format == "" {
		repositoryList = getRepositoryList()
	} else {
		repositoryList = getRepositoryListByFormat(validateRepositoryFormat(format))
	}
	if name == "" {
		printStringSlice(repositoryList)
		fmt.Printf("Number of repositories : %d\n", len(repositoryList))
	}
}

func CreateHosted(name, blobStoreName, format string, dockerHttpPort, dockerHttpsPort float64, releases bool) {
	if name == "" || format == "" {
		log.Printf("%s : %s", getfuncName(), repoNameFormatRequiredInfo)
		os.Exit(1)
	}
	format = validateRepositoryFormat(format)

	var attributes Attributes
	recipe := fmt.Sprintf("%s-hosted", format)
	storage := Storage{BlobStoreName: getBlobStoreName(blobStoreName), StrictContentTypeValidation: true, WritePolicy: getWritePolicy(releases)}

	if format == "maven2" {
		maven := Maven{VersionPolicy: getVersionPolicy(releases), LayoutPolicy: "STRICT"}
		attributes = Attributes{Storage: storage, Maven: maven}
	} else if format == "docker" {
		if dockerHttpPort == 0 && dockerHttpsPort == 0 {
			log.Printf("%s : %s", getfuncName(), dockerPortsInfo)
			os.Exit(1)
		}
		docker := Docker{HTTPPort: dockerHttpPort, HTTPSPort: dockerHttpsPort, ForceBasicAuth: true, V1Enabled: false}
		attributes = Attributes{Storage: storage, Docker: docker}
	} else {
		attributes = Attributes{Storage: storage}
	}

	repository := Repository{Name: name, Format: format, Recipe: recipe, Attributes: attributes}
	payload, err := json.Marshal(repository)
	logJsonMarshalError(err, getfuncName())
	result := RunScript(createHostedRepoScript, string(payload))
	printCreateRepoStatus(name, result.Status)
}

func CreateProxy(name, blobStoreName, format, remoteURL, proxyUsername, proxyPassword string, dockerHttpPort, dockerHttpsPort float64, releases bool) {
	if name == "" || format == "" {
		log.Printf("%s : %s", getfuncName(), repoNameFormatRequiredInfo)
		os.Exit(1)
	} else if remoteURL == "" {
		log.Printf("%s : %s", getfuncName(), proxyRepoRequiredInfo)
		os.Exit(1)
	}
	format = validateRepositoryFormat(format)
	validateProxyAuthInfo(proxyUsername, proxyPassword)
	validateRemoteURL(remoteURL)

	var attributes Attributes
	recipe := fmt.Sprintf("%s-proxy", format)
	storage := Storage{BlobStoreName: getBlobStoreName(blobStoreName), StrictContentTypeValidation: true, WritePolicy: getWritePolicy(releases)}
	proxy := Proxy{RemoteURL: remoteURL, ContentMaxAge: -1, MetadataMaxAge: 1440}
	proxyAuth := HttpClientAuth{Username: proxyUsername, Password: proxyPassword}
	proxyHttpClient := HttpClient{Blocked: false, AutoBlock: true, Authentication: proxyAuth}
	negetiveCache := NegetiveCache{Enabled: true, TimeToLive: 1440}

	if format == "maven2" {
		maven := Maven{VersionPolicy: getVersionPolicy(releases), LayoutPolicy: "STRICT"}
		attributes = Attributes{Storage: storage, Maven: maven, Proxy: proxy, Httpclient: proxyHttpClient, NegativeCache: negetiveCache}
	} else if format == "docker" {
		if dockerHttpPort == 0 && dockerHttpsPort == 0 {
			log.Printf("%s : %s", getfuncName(), dockerPortsInfo)
			os.Exit(1)
		}
		docker := Docker{HTTPPort: dockerHttpPort, HTTPSPort: dockerHttpsPort, ForceBasicAuth: true, V1Enabled: false}
		dockerProxy := DockerProxy{IndexType: "REGISTRY"}
		attributes = Attributes{Storage: storage, Docker: docker, Proxy: proxy, DockerProxy: dockerProxy, Httpclient: proxyHttpClient, NegativeCache: negetiveCache}
	} else {
		storage := Storage{BlobStoreName: getBlobStoreName(blobStoreName), StrictContentTypeValidation: true, WritePolicy: getWritePolicy(releases)}
		attributes = Attributes{Storage: storage, Proxy: proxy, Httpclient: proxyHttpClient, NegativeCache: negetiveCache}
	}

	repository := Repository{Name: name, Format: format, Recipe: recipe, Attributes: attributes}
	payload, err := json.Marshal(repository)
	logJsonMarshalError(err, getfuncName())
	result := RunScript(createProxyRepoScript, string(payload))
	printCreateRepoStatus(name, result.Status)
}

func CreateGroup(name, blobStoreName, format, repoMembers string, dockerHttpPort, dockerHttpsPort float64, releases bool) {
	if name == "" || format == "" {
		log.Printf("%s : %s", getfuncName(), repoNameFormatRequiredInfo)
		os.Exit(1)
	} else if repoMembers == "" {
		log.Printf("%s : %s", getfuncName(), groupRepoRequiredInfo)
		os.Exit(1)
	}
	format = validateRepositoryFormat(format)
	validList := validateGroupMembers(repoMembers, format)

	var attributes Attributes
	recipe := fmt.Sprintf("%s-group", format)
	storage := Storage{BlobStoreName: getBlobStoreName(blobStoreName), StrictContentTypeValidation: true, WritePolicy: getWritePolicy(releases)}
	group := Group{MemberNames: validList}

	if format == "maven2" {
		maven := Maven{VersionPolicy: getVersionPolicy(releases), LayoutPolicy: "STRICT"}
		attributes = Attributes{Storage: storage, Maven: maven, Group: group}
	} else if format == "docker" {
		if dockerHttpPort == 0 && dockerHttpsPort == 0 {
			log.Printf("%s : %s", getfuncName(), dockerPortsInfo)
			os.Exit(1)
		}
		docker := Docker{HTTPPort: dockerHttpPort, HTTPSPort: dockerHttpsPort, ForceBasicAuth: true, V1Enabled: false}
		attributes = Attributes{Storage: storage, Docker: docker, Group: group}
	} else {
		storage := Storage{BlobStoreName: getBlobStoreName(blobStoreName), StrictContentTypeValidation: true, WritePolicy: getWritePolicy(releases)}
		attributes = Attributes{Storage: storage, Group: group}
	}

	repository := Repository{Name: name, Format: format, Recipe: recipe, Attributes: attributes}
	payload, err := json.Marshal(repository)
	logJsonMarshalError(err, getfuncName())
	result := RunScript(createGroupRepoScript, string(payload))
	printCreateRepoStatus(name, result.Status)
}

func AddMembersToGroup(name, repoMembers string) {
	if name == "" {
		log.Printf("%s : %s", getfuncName(), name)
		os.Exit(1)
	} else if repoMembers == "" {
		log.Printf("%s : %s", getfuncName(), groupRepoRequiredInfo)
		os.Exit(1)
	}
	if repositoryExists(name) {
		repo := getRepository(name)
		validateGroupRepo(repo)
		format := repo.Format
		validList := validateGroupMembers(repoMembers, format)
		currentMembers := repo.Attributes.Group.MemberNames
		for _, newMember := range validList {
			if entryExists(currentMembers, newMember) {
				log.Printf(groupMemberAlreadyExistsInfo, newMember, name)
			} else if newMember == name {
				log.Printf(cannotBeSameRepoInfo, newMember, name)
			} else {
				log.Printf(groupMemberAddSuccessInfo, newMember, name)
				currentMembers = append(currentMembers, newMember)
			}
		}
		repo.Attributes.Group = Group{MemberNames: currentMembers}
		repository := Repository{Name: name, Format: format, Attributes: repo.Attributes}
		payload, err := json.Marshal(repository)
		logJsonMarshalError(err, getfuncName())
		result := RunScript(updateGroupMembersScript, string(payload))
		printUpdateRepoStatus(name, result.Status)
	} else {
		log.Printf(repositoryNotFoundInfo, name)
	}
}

func RemoveMembersFromGroup(name, repoMembers string) {
	if name == "" {
		log.Printf("%s : %s", getfuncName(), nameRequiredInfo)
		os.Exit(1)
	} else if repoMembers == "" {
		log.Printf("%s : %s", getfuncName(), groupRepoRequiredInfo)
		os.Exit(1)
	}
	if repositoryExists(name) {
		repo := getRepository(name)
		validateGroupRepo(repo)
		format := repo.Format
		validList := validateGroupMembers(repoMembers, format)
		currentMembers := repo.Attributes.Group.MemberNames
		for _, newMember := range validList {
			if !entryExists(currentMembers, newMember) {
				log.Printf(groupMemberRemoveNotFoundInfo, newMember, name)
			} else if newMember == name {
				log.Printf(cannotBeSameRepoInfo, newMember, name)
			} else {
				log.Printf(groupMemberRemoveSuccessInfo, newMember, name)
				currentMembers = removeEntryFromSlice(currentMembers, newMember)
			}
		}
		repo.Attributes.Group = Group{MemberNames: currentMembers}
		repository := Repository{Name: name, Format: format, Attributes: repo.Attributes}
		payload, err := json.Marshal(repository)
		logJsonMarshalError(err, getfuncName())
		result := RunScript(updateGroupMembersScript, string(payload))
		printUpdateRepoStatus(name, result.Status)
	} else {
		log.Printf(repositoryNotFoundInfo, name)
	}
}

func DeleteRepository(name string) {
	if name == "" {
		log.Printf("%s : %s", getfuncName(), nameRequiredInfo)
		os.Exit(1)
	}
	payload, err := json.Marshal(Repository{Name: name})
	logJsonMarshalError(err, getfuncName())
	result := RunScript(deleteRepoScript, string(payload))
	printDeleteRepoStatus(name, result.Status)
}

func getRepository(name string) Repository {
	if name == "" {
		log.Printf("%s : %s", getfuncName(), nameRequiredInfo)
		os.Exit(1)
	}
	var result ScriptResult
	if repositoryExists(name) {
		payload, err := json.Marshal(Repository{Name: name})
		logJsonMarshalError(err, getfuncName())
		result = RunScript(getRepoScript, string(payload))
		if result.Status != successStatus {
			log.Printf("%s : %s", getfuncName(), setVerboseInfo)
			os.Exit(1)
		}
	} else {
		log.Printf("%s : %s", getfuncName(), fmt.Sprintf(repositoryNotFoundInfo, name))
		os.Exit(1)
	}
	return Repository{Name: result.Name, URL: result.URL, Type: result.Type, Format: result.Format, Recipe: result.Recipe, Attributes: result.Attributes}
}

func getRepositories() []Repository {
	url := fmt.Sprintf("%s/%s/%s", NexusURL, apiBase, repositoryPath)
	var repositories []Repository
	req := createBaseRequest("GET", url, RequestBody{})
	respBody, status := httpRequest(req)
	if status != successStatus {
		log.Printf("%s : %s", getfuncName(), setVerboseInfo)
		os.Exit(1)
	} else {
		err := json.Unmarshal(respBody, &repositories)
		logJsonUnmarshalError(err, getfuncName())
	}
	return repositories
}

func getRepositoryList() []string {
	var repositoryList []string
	repositories := getRepositories()
	for _, r := range repositories {
		repositoryList = append(repositoryList, r.Name)
	}
	return repositoryList
}

func getRepositoryListByFormat(format string) []string {
	var repositoryList []string
	repositories := getRepositories()
	for _, r := range repositories {
		if format == r.Format {
			repositoryList = append(repositoryList, r.Name)
		}
	}
	return repositoryList
}

func repositoryExists(name string) bool {
	if name == "" {
		log.Printf("%s : %s", getfuncName(), nameRequiredInfo)
		os.Exit(1)
	}
	var isExists bool
	payload, err := json.Marshal(Repository{Name: name})
	logJsonMarshalError(err, getfuncName())
	result := RunScript(getRepoScript, string(payload))
	if result.Status == successStatus {
		isExists = true
	} else if result.Status == notFoundStatus {
		isExists = false
	} else {
		log.Printf("%s : %s", getfuncName(), setVerboseInfo)
		os.Exit(1)
	}
	return isExists
}

func getBlobStoreName(blobStoreName string) string {
	if blobStoreName == "" {
		blobStoreName = "default"
	}
	return blobStoreName
}

func getVersionPolicy(release bool) string {
	var versionPolicy string
	if release {
		versionPolicy = "RELEASE"
	} else {
		versionPolicy = "SNAPSHOT"
	}
	return versionPolicy
}

func getWritePolicy(releases bool) string {
	var writePolicy string
	if releases {
		writePolicy = "ALLOW_ONCE"
	} else {
		writePolicy = "ALLOW"
	}
	return writePolicy
}

func validateRepositoryFormat(format string) string {
	if format == "" {
		log.Printf("%s : %s", getfuncName(), repoFormatRequiredInfo)
		os.Exit(1)
	}
	formatChoice := map[string]bool{"": true}
	for _, repoFormat := range RepoFormats {
		formatChoice[repoFormat] = true
	}
	if _, validChoice := formatChoice[format]; !validChoice {
		log.Printf("%s : %s", getfuncName(), fmt.Sprintf(RepoFormatNotValidInfo, format, RepoFormats))
		os.Exit(1)
	}
	if format == "maven" {
		return "maven2"
	}
	return format
}

func validateProxyAuthInfo(proxyUsername, proxyPassword string) {
	if proxyUsername == "" && proxyPassword == "" {
		return
	} else if proxyUsername != "" && proxyPassword != "" {
		return
	} else {
		log.Printf("%s : %s\n", getfuncName(), proxyCredsNotValidInfo)
		os.Exit(1)
	}
}

func validateRemoteURL(url string) {
	httpRegex, _ := regexp2.Compile(`^(http://).*`)
	httpsRegex, _ := regexp2.Compile(`^(https://).*`)
	if httpRegex.MatchString(url) || httpsRegex.MatchString(url) {
		return
	} else {
		log.Printf("%s : %s", getfuncName(), fmt.Sprintf(remoteURLNotValidInfo, url))
		os.Exit(1)
	}
}

func validateGroupRepo(repo Repository) {
	if strings.Contains(repo.Recipe, "group") {
		return
	} else {
		log.Printf(notAGroupRepoInfo, repo.Name)
		os.Exit(1)
	}
}

func validateGroupMembers(repoMembers, format string) []string {
	var validList []string
	repoMembersList := strings.Split(strings.Replace(repoMembers, " ", "", -1), ",")
	for _, repoMember := range repoMembersList {
		if repositoryExists(repoMember) {
			repoDetails := getRepository(repoMember)
			if strings.Contains(repoDetails.Recipe, format) {
				validList = append(validList, repoMember)
			} else {
				log.Printf(groupMemberInvalidFormatInfo, repoMember, format)
			}
		} else {
			log.Printf(groupMemberNotFoundInfo, repoMember)
		}
	}
	if len(validList) < 1 {
		log.Printf("%s : %s\n", getfuncName(), groupMemberRequiredInfo)
		os.Exit(1)
	}
	return validList
}

func printCreateRepoStatus(name, status string) {
	if status == successStatus {
		log.Printf(repoCreatedInfo, name)
	} else if status == foundStatus {
		log.Printf(repoExistsInfo, name)
	} else {
		log.Printf(repoCreateErrorInfo, setVerboseInfo)
	}
}

func printUpdateRepoStatus(name, status string) {
	if status == successStatus {
		log.Printf(repoUpdatedStatus, name)
	} else if status == notFoundStatus {
		log.Printf(repositoryNotFoundInfo, name)
	} else {
		log.Printf(repoUpdateErrorInfo, setVerboseInfo)
	}
}

func printDeleteRepoStatus(name, status string) {
	if status == successStatus {
		log.Printf(repoDeletedInfo, name)
	} else if status == notFoundStatus {
		log.Printf(repositoryNotFoundInfo, name)
	} else {
		log.Printf(repoDeleteErrorInfo, setVerboseInfo)
	}
}
