package main

import "fmt"

func modifyConf(svcName, key, value string) {
	fmt.Printf("Modifying [%v] given key = [%v], value = [%v]\n", svcName, key, value)
}

func restartService(svcName string) {
	fmt.Printf("Restarting [%v]\n", svcName)
}

func applyManifest(resType, image string) {
	fmt.Printf("Applying manifest to [%v] given image = [%v]\n", resType, image)
}
