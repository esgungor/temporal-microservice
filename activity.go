package app

import "fmt"

func UpdateKubernetesDeployment(s string) error {
	fmt.Printf("Launch updated :%v", s)
	// Call kubeapps from here
	// Lambda function
	return nil

}

func DeleteKubernetesDeployment(s string) error {
	fmt.Printf("Launch updated :%v", s)
	// Call kubeapps from here
	return nil

}
