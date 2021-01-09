package config

import (
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

func Init() {
	yamlFile, unMarshallError := ioutil.ReadFile("config.yaml")
	log.Printf("Could not load config.yaml, Error = %v", zap.Error(unMarshallError))
	configMap := make(map[interface{}]interface{})
	uError := yaml.Unmarshal(yamlFile, &configMap)
	if uError != nil {
		log.Printf("Could not Unmarshal config.yaml. Error = %v", zap.Error(uError))
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
		} else if os.Getenv(nextRoot) == "" {
			os.Setenv(nextRoot, fmt.Sprintf("%v", v))
		}
	}
}
