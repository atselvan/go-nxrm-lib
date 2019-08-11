package nxrm

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
)

type Script struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

type ScriptOutput struct {
	Name   string `json:"name"`
	Result string `json:"result"`
}

type ScriptResult struct {
	Status           string            `json:"status"`
	Name             string            `json:"name"`
	URL              string            `json:"url"`
	Type             string            `json:"type"`
	Format           string            `json:"format"`
	Recipe           string            `json:"recipe"`
	Message          string            `json:"message"`
	Attributes       Attributes        `json:"attributes"`
	ContentSelectors []ContentSelector `json:"contentSelectors"`
	Privileges       []Privilege       `json:"privileges"`
	Roles            []Role            `json:"roles"`
}

func ListScripts(name string) {
	if name != "" {
		script := getScript(name)
		fmt.Println(script)
	} else {
		scriptsList := getScripts()
		sort.Strings(scriptsList)
		printStringSlice(scriptsList)
		if len(scriptsList) == 0 {
			fmt.Println("There are no scripts available in nexus")
		} else {
			fmt.Printf("No of scripts in nexus : %d\n", len(scriptsList))
		}
	}
}

func ScriptsInit() {
	for _, s := range NexusScripts {
		AddOrUpdateScript(s)
	}
}

func AddScript(name string) {
	if name == "" {
		log.Printf("%s : %s", getfuncName(), nameRequiredInfo)
		os.Exit(1)
	}
	url := fmt.Sprintf("%s/%s/%s", NexusURL, apiBase, scriptAPI)
	if !scriptExists(name) {
		payload, err := json.Marshal(Script{Name: name, Type: "groovy", Content: readFile(getScriptPath(name))})
		logJsonMarshalError(err, getfuncName())
		req := createBaseRequest("POST", url, RequestBody{Json: payload})
		_, status := httpRequest(req)
		if status == noContentStatus {
			if Debug {
				log.Printf(scriptAddedInfo, name)
			}
		} else {
			log.Printf("%s : %s", getfuncName(), setVerboseInfo)
			os.Exit(1)
		}
	} else {
		log.Printf(scriptExistsInfo, name)
	}
}

func UpdateScript(name string) {
	if name == "" {
		log.Printf("%s : %s", getfuncName(), nameRequiredInfo)
		os.Exit(1)
	}
	url := fmt.Sprintf("%s/%s/%s/%s", NexusURL, apiBase, scriptAPI, name)
	if scriptExists(name) {
		payload, err := json.Marshal(Script{Name: name, Type: "groovy", Content: readFile(getScriptPath(name))})
		logJsonMarshalError(err, getfuncName())
		req := createBaseRequest("PUT", url, RequestBody{Json: payload})
		_, status := httpRequest(req)
		if status == noContentStatus {
			if Debug {
				log.Printf(scriptUpdatedInfo, name)
			}
		} else {
			log.Printf("%s : %s", getfuncName(), setVerboseInfo)
			os.Exit(1)
		}
	} else {
		log.Printf(scriptNotfoundInfo, name)
	}
}

func AddOrUpdateScript(name string) {
	if name == "" {
		log.Printf("%s : %s", getfuncName(), nameRequiredInfo)
		os.Exit(1)
	}
	if !scriptExists(name) {
		AddScript(name)
	} else {
		UpdateScript(name)
	}
}

func DeleteScript(name string) {
	if name == "" {
		log.Printf("%s : %s", getfuncName(), nameRequiredInfo)
		os.Exit(1)
	}
	url := fmt.Sprintf("%s/%s/%s/%s", NexusURL, apiBase, scriptAPI, name)
	if scriptExists(name) {
		req := createBaseRequest("Delete", url, RequestBody{Json: nil})
		_, status := httpRequest(req)
		if status == noContentStatus {
			if Debug {
				log.Printf(scriptDeletedInfo, name)
			}
		} else {
			log.Printf("%s : %s", getfuncName(), setVerboseInfo)
			os.Exit(1)
		}
	} else {
		log.Printf(scriptNotfoundInfo, name)
	}
}

func RunScript(name, payload string) ScriptResult {
	if name == "" {
		log.Printf("%s : %s", getfuncName(), nameRequiredInfo)
		os.Exit(1)
	}
	var (
		output ScriptOutput
		result ScriptResult
	)
	url := fmt.Sprintf("%s/%s/%s/%s/run", NexusURL, apiBase, scriptAPI, name)
	req := createBaseRequest("POST", url, RequestBody{Text: payload})
	respBody, status := httpRequest(req)
	if status == successStatus {
		if Debug {
			log.Printf(scriptRunSuccessInfo, name)
		}
	} else if status == notFoundStatus {
		log.Printf(scriptRunNotFoundInfo, name)
		os.Exit(1)
	} else {
		log.Printf("%s : %s", getfuncName(), setVerboseInfo)
		os.Exit(1)
	}
	err := json.Unmarshal(respBody, &output)
	logJsonUnmarshalError(err, getfuncName())
	err = json.Unmarshal([]byte(output.Result), &result)
	logJsonUnmarshalError(err, getfuncName())
	return result
}

func getScripts() []string {
	var (
		url         = fmt.Sprintf("%s/%s/%s", NexusURL, apiBase, scriptAPI)
		scripts     []Script
		scriptsList []string
	)
	req := createBaseRequest("GET", url, RequestBody{})
	respBody, status := httpRequest(req)
	if status != successStatus {
		log.Printf("%s : %s", getfuncName(), setVerboseInfo)
		os.Exit(1)
	} else {
		err := json.Unmarshal(respBody, &scripts)
		logJsonUnmarshalError(err, getfuncName())
		for _, s := range scripts {
			scriptsList = append(scriptsList, s.Name)
		}
	}
	return scriptsList
}

func getScript(name string) Script {
	if name == "" {
		log.Printf("%s : %s", getfuncName(), nameRequiredInfo)
		os.Exit(1)
	}
	var (
		url    = fmt.Sprintf("%s/%s/%s/%s", NexusURL, apiBase, scriptAPI, name)
		script Script
	)
	req := createBaseRequest("GET", url, RequestBody{})
	respBody, status := httpRequest(req)
	if status == successStatus {
		err := json.Unmarshal(respBody, &script)
		logJsonUnmarshalError(err, getfuncName())
	} else if status == notFoundStatus {
		log.Printf(scriptNotfoundInfo, name)
		os.Exit(1)
	} else {
		log.Printf("%s : %s", getfuncName(), setVerboseInfo)
		os.Exit(1)
	}
	return script
}

func getScriptPath(name string) string {
	if name == "" {
		log.Printf("%s : %s", getfuncName(), nameRequiredInfo)
		os.Exit(1)
	}
	return fmt.Sprintf("%s/%s.groovy", scriptBasePath, name)
}

func scriptExists(name string) bool {
	if name == "" {
		log.Printf("%s : %s", getfuncName(), nameRequiredInfo)
		os.Exit(1)
	}
	url := fmt.Sprintf("%s/%s/%s/%s", NexusURL, apiBase, scriptAPI, name)
	req := createBaseRequest("GET", url, RequestBody{})
	_, status := httpRequest(req)
	if status == successStatus {
		return true
	}
	return false
}
