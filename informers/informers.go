package informers

import (
	"context"
	"log"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// SharedInformer code
func SharedInformer(clientset *kubernetes.Clientset) {
	// informerFactory := informers.NewSharedInformerFactory(clientset, time.Second*30)
	// podInformer := informerFactory.Core().V1().Pods().Informer()
	// podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
	// 	AddFunc:    func(new interface{}) { log.Println("pod added") },
	// 	UpdateFunc: func(old, new interface{}) { log.Println("pod updated") },
	// 	DeleteFunc: func(obj interface{}) { log.Println("pod deleted") },
	// })
	// podInformer.Run(wait.NeverStop)
	// pod, err := podInformer.Lister().Pods("programming-kubernetes").Get("client-go")

	// listerWatcher(clientset)
	// reflector(clientset)
	controller(clientset)
}

func controller(clientset *kubernetes.Clientset) {
	store, controller := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return clientset.CoreV1().Pods("default").List(
					context.Background(),
					metav1.ListOptions{},
				)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				// Direct call to the API server, using the job client
				return clientset.CoreV1().Pods("default").Watch(
					context.Background(),
					metav1.ListOptions{},
				)
			},
		},
		&v1.Pod{},
		5*time.Second,
		cache.ResourceEventHandlerFuncs{
			DeleteFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				log.Println(pod.Name)
			},
		},
	)
	log.Println(store.ListKeys())
	controller.Run(wait.NeverStop)
}

func reflector(clientset *kubernetes.Clientset) {
	store := cache.NewStore(
		func(obj interface{}) (string, error) {
			pod := obj.(*v1.Pod)
			return pod.Name, nil
		},
	)
	reflector := cache.NewReflector(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return clientset.CoreV1().Pods("default").List(
					context.Background(),
					metav1.ListOptions{},
				)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				// Direct call to the API server, using the job client
				return clientset.CoreV1().Pods("default").Watch(
					context.Background(),
					metav1.ListOptions{},
				)
			},
		},
		&v1.Pod{},
		store,
		30*time.Second,
	)

	go reflector.Run(wait.NeverStop)
	for {
		log.Println(store.ListKeys())
	}
}

// lister watcher just list and watch for resources.
func listerWatcher(clientset *kubernetes.Clientset) {
	listwatch := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return clientset.CoreV1().Pods("default").List(
				context.Background(),
				metav1.ListOptions{},
			)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			// Direct call to the API server, using the job client
			return clientset.CoreV1().Pods("default").Watch(
				context.Background(),
				metav1.ListOptions{},
			)
		},
	}

	object, _ := listwatch.List(metav1.ListOptions{})
	log.Println(object.GetObjectKind().GroupVersionKind())

	watcher, _ := listwatch.Watch(metav1.ListOptions{})
	for {
		channel := watcher.ResultChan()
		log.Println((<-channel).Type)
	}
}

// event handler function constraints?
// which event occurred? metadata?
// sharedinformer vs sharedindexinformer
// look into different way to initialilze informers.
// controllers
// sample programs
// dynamic informers
