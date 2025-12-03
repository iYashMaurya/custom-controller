package controller

import (
	"context"
	"fmt"
	"time"

	deploy "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	v1 "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type controller struct {
	clientSet    kubernetes.Interface
	depList      appsv1.DeploymentLister
	depCacheSync cache.InformerSynced
	queue        workqueue.RateLimitingInterface
}

func NewController(clSet kubernetes.Interface, depInformer v1.DeploymentInformer) *controller {
	c := controller{
		clientSet:    clSet,
		depList:      depInformer.Lister(),
		depCacheSync: depInformer.Informer().HasSynced,
		queue:        workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ekpose"),
	}
	depInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.handleAdd,
		DeleteFunc: c.handleDelete,
		//UpdateFunc: handleUpdate,
	})
	return &c
}

func (c *controller) Start(ch <-chan struct{}) {
	fmt.Println("Starting the controller")
	if !cache.WaitForCacheSync(ch, c.depCacheSync) {
		fmt.Println("Error in syncing the cache")
	}
	go wait.Until(c.worker, 1*time.Second, ch)
	<-ch
}

func (c *controller) worker() {
	for c.processItem() {

	}
}

func (c *controller) processItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	defer c.queue.Forget(item)
	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		fmt.Println("Getting key from cache: ", err.Error())
	}
	fmt.Println(key)
	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		fmt.Println("Splitting key into namespace and name: ", err.Error())
		return false
	}
	err = c.syncDeployment(ns, name)
	if err != nil {
		// retry

		fmt.Println("Syncing deployment: ", err.Error())
		return false
	}
	return true
}

func (c *controller) syncDeployment(ns, name string) error {
	ctx := context.Background()
	deployment, err := c.depList.Deployments(ns).Get(name)
	if err != nil {
		fmt.Println("Getting the deployment using informer: ", err.Error())
	}
	if deployment == nil {
		return fmt.Errorf("deployment received is nil")
	}
	// Create service
	svc := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployment.Name,
			Namespace: ns,
		},
		Spec: corev1.ServiceSpec{
			Selector: depLabel(*deployment),
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: 80,
				},
			},
		},
	}
	_, err = c.clientSet.CoreV1().Services(ns).Create(ctx, &svc, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("Creating Service: ", err.Error())
	}
	// Create deployment
	return nil
}

func depLabel(dep deploy.Deployment) map[string]string {
	return dep.Spec.Template.Labels
}

func (c *controller) handleAdd(obj interface{}) {
	fmt.Println("Add was called")
	c.queue.Add(obj)
}

func (c *controller) handleDelete(obj interface{}) {
	fmt.Println("Delete was called")
}
