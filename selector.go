package nxrm

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type ContentSelector struct {
	Name        string                    `json:"name"`
	Type        string                    `json:"type"`
	Description string                    `json:"description"`
	Attributes  ContentSelectorAttributes `json:"attributes"`
}

type ContentSelectorAttributes struct {
	Expression string `json:"expression"`
}

func ListSelectors(name string) {
	if name != "" {
		cs := getSelector(name)
		fmt.Printf("Name: %s\nDescription: %s\nExpression: %s\n",
			cs.Name, cs.Description, cs.Attributes.Expression)
	} else {
		csNames := getSelectorNames()
		printStringSlice(csNames)
		fmt.Printf("Total Number of content selectors : %d\n", len(csNames))
	}
}

func CreateSelector(name, description, expression string) {
	if name == "" || expression == "" {
		log.Printf("%s : %s", getfuncName(), createSelectorRequiredInfo)
		os.Exit(1)
	}
	if !selectorExists(name) {
		attributes := ContentSelectorAttributes{Expression: expression}
		payload, err := json.Marshal(ContentSelector{Name: name, Type: contentSelectorType, Description: getSelectorDescription(description), Attributes: attributes})
		logJsonMarshalError(err, jsonMarshalError)
		result := RunScript("create-content-selector", string(payload))
		if result.Status == successStatus {
			log.Printf(createSelectorSuccessInfo, name)
		} else {
			log.Printf("%s : %s", getfuncName(), setVerboseInfo)
			os.Exit(1)
		}
	} else {
		log.Printf(selectorAlreadyExistsInfo, name)
	}
}

func UpdateSelector(name, description, expression string) {
	if name == "" {
		log.Printf("%s : %s", getfuncName(), nameRequiredInfo)
		os.Exit(1)
	}
	if selectorExists(name) {
		selector := getSelector(name)
		if description != "" {
			selector.Description = description
		}
		if expression != "" {
			selector.Attributes = ContentSelectorAttributes{Expression: expression}
		}
		payload, err := json.Marshal(selector)
		logJsonMarshalError(err, jsonMarshalError)
		result := RunScript("update-content-selector", string(payload))
		if result.Status == successStatus {
			log.Printf(updateSelectorSuccessInfo, name)
		} else {
			log.Printf("%s : %s", getfuncName(), setVerboseInfo)
			os.Exit(1)
		}
	} else {
		log.Printf(selectorNotFoundInfo, name)
	}
}

func DeleteSelector(name string) {
	if name == "" {
		log.Printf("%s : %s", getfuncName(), nameRequiredInfo)
		os.Exit(1)
	}
	if selectorExists(name) {
		selector := getSelector(name)
		payload, err := json.Marshal(selector)
		logJsonMarshalError(err, jsonMarshalError)
		result := RunScript("delete-content-selector", string(payload))
		if result.Status == successStatus {
			log.Printf(deleteSelectorSuccessInfo, name)
		} else {
			log.Printf("%s : %s", getfuncName(), setVerboseInfo)
			os.Exit(1)
		}
	} else {
		log.Printf(selectorNotFoundInfo, name)
	}
}

func getSelectors() []ContentSelector {
	payload, err := json.Marshal(Repository{})
	logJsonMarshalError(err, getfuncName())
	result := RunScript("get-content-selectors", string(payload))
	return result.ContentSelectors
}

func getSelector(name string) ContentSelector {
	if name == "" {
		log.Printf("%s : %s", getfuncName(), nameRequiredInfo)
		os.Exit(1)
	}
	var contentSelector ContentSelector
	contentSelectors := getSelectors()
	for _, cs := range contentSelectors {
		if cs.Name == name {
			contentSelector = cs
		}
	}
	if contentSelector.Name == "" {
		log.Printf(selectorNotFoundInfo, name)
		os.Exit(1)
	}
	return contentSelector
}

func getSelectorNames() []string {
	var csNames []string
	contentSelectors := getSelectors()
	for _, cs := range contentSelectors {
		csNames = append(csNames, cs.Name)
	}
	return csNames
}

func selectorExists(name string) bool {
	csNames := getSelectorNames()
	if entryExists(csNames, name) {
		return true
	}
	return false
}

func getSelectorDescription(description string) string {
	if description == "" {
		return defaultContentSelectorDescription
	}
	return description
}
