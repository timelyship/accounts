package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"

	//"reflect"
	"gopkg.in/yaml.v2"
)

/*

GOOGLE_OAUTH_SCOPES: abc def
GOOGLE_OAUTH_CLIENT_ID: 1s4rf
*/

func Init() {

	yamlFile, unMarshallError := ioutil.ReadFile("config.yaml")
	fmt.Println("Error", unMarshallError)
	configMap := make(map[interface{}]interface{})
	uError := yaml.Unmarshal(yamlFile, &configMap)
	if uError != nil {
		log.Fatalf("error: %v", uError)
	}

	flatten("", configMap)
}

func flatten(root string, configMap map[interface{}]interface{}) {
	for k, v := range configMap {
		var nextRoot string
		if root == "" {
			nextRoot = k.(string)
		} else {
			nextRoot = fmt.Sprintf("%s.%s", root, k)
		}
		if reflect.ValueOf(v).Kind() == reflect.Map {
			flatten(nextRoot, v.(map[interface{}]interface{}))
		} else {
			//fmt.Printf("%v %v\n", nextRoot, v)
			if os.Getenv(nextRoot) == "" {
				os.Setenv(nextRoot, fmt.Sprintf("%v", v))
			}
		}
	}
}
