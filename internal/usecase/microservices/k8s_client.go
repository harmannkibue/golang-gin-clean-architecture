package microservices

import (
	"context"
	"log"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/harmannkibue/actsml-jobs-orchestrator/internal/entity/intfaces"
)

type K8sClient struct {
	clientset *kubernetes.Clientset
}

func NewK8sClient() (intfaces.KubernetesClient, error) {
	config, err := rest.InClusterConfig()
	if err == nil {
		clientset, e := kubernetes.NewForConfig(config)
		if e != nil {
			return nil, e
		}
		log.Println("K8S CLIENT: Using InClusterConfig")
		return &K8sClient{clientset: clientset}, nil
	}

	kubePath := os.Getenv("KUBECONFIG")
	if kubePath == "" {
		kubePath = os.ExpandEnv("$HOME/.kube/config")
	}
	config, err = clientcmd.BuildConfigFromFlags("", kubePath)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	log.Println("K8S CLIENT: Using local kubeconfig:", kubePath)
	return &K8sClient{clientset: clientset}, nil
}

func (c *K8sClient) CreateJob(ctx context.Context, namespace string, job *batchv1.Job) (*batchv1.Job, error) {
	return c.clientset.BatchV1().Jobs(namespace).Create(ctx, job, metav1.CreateOptions{})
}

func (c *K8sClient) DeleteJob(ctx context.Context, namespace, name string) error {
	return c.clientset.BatchV1().Jobs(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (c *K8sClient) GetJob(ctx context.Context, namespace, name string) (*batchv1.Job, error) {
	return c.clientset.BatchV1().Jobs(namespace).Get(ctx, name, metav1.GetOptions{})
}
