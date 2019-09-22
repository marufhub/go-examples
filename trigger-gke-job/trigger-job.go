package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"cloud.google.com/go/datastore"
	uuid "github.com/satori/go.uuid"
	v1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)


func main() {
	// check the file status in datastore
	ctx := context.Background()



		if err := triggerGKEJob(); err != nil {
			log.Fatal(err)
			fmt.Println("Error in trigger gke job")
		} else {
				log.Fatal(err)
			}
		}

	}

}


func triggerGKEJob() error {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" { // check if machine has home directory.
		// read kubeconfig flag. if not provided use config file $HOME/.kube/config
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")

	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		fmt.Println(kubeconfig)
	}
	flag.Parse()

	// build configuration from the config file.
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return err
	}
	// create kubernetes clientset. this clientset can be used to create,delete,patch,list etc for the kubernetes resources
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	jobsClient := clientset.BatchV1().Jobs(apiv1.NamespaceDefault)

	// job, err := jobsClient.Get("test-job", metav1.GetOptions{})

	decode := scheme.Codecs.UniversalDeserializer().Decode
	jobManifests, err := ioutil.ReadFile("job.yaml")
	obj, _, err := decode([]byte(jobManifests), nil, nil)
	if err != nil {
		return err
	}
	job := obj.(*v1.Job)
	if _, err := jobsClient.Create(job); err != nil {
		return err
	}
	return nil
}

