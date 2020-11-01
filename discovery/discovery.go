package discovery

import (
	"context"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	rbacv1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
)

func namespaceDiscovery(client v1.CoreV1Interface) {
	// get namespace interface
	namespaceHandler := client.Namespaces()

	// list all namespace
	namespaceList, err := namespaceHandler.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("error getting namespacelist: %s", err)
		return
	}
	log.Println("Namespace list =============================================")
	for _, item := range namespaceList.Items {
		log.Println(item.Name)
	}

	// list a single namespace with name
	log.Println("Single Namespace =============================================")
	namespace, err := namespaceHandler.Get(context.Background(), "linkerd", metav1.GetOptions{})
	if err != nil {
		log.Printf("error getting linkerd namespace: %s", err)
		return
	}
	log.Println(namespace.Name)
}

func nodeDiscovery(client v1.CoreV1Interface) {
	// get node interface
	nodeHandler := client.Nodes()

	// list all node
	nodeList, err := nodeHandler.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("error getting nodelist: %s", err)
		return
	}
	log.Println("Nodes list =============================================")
	for _, item := range nodeList.Items {
		log.Println(item.Name)
	}

	// list a single node with name
	log.Println("Single Node =============================================")
	node, err := nodeHandler.Get(context.Background(), "kube-node-d19d", metav1.GetOptions{})
	if err != nil {
		log.Printf("error getting node: %s", err)
		return
	}
	log.Println(node.Name)
}

func roleDiscovery(client rbacv1.RbacV1Interface) {
	// get cluster role
	clusterRole := client.ClusterRoles()
	list, err := clusterRole.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("error getting cluster role list: %s", err)
		return
	}
	log.Println("Cluster Role List =============================================")
	for _, item := range list.Items {
		log.Println(item.Name)
	}

	// get cluster role bindings
	clusterRoleBindings := client.ClusterRoleBindings()
	bindingList, err := clusterRoleBindings.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("error getting cluster role bindings list: %s", err)
		return
	}
	log.Println("Cluster Role Bindings List =============================================")
	for _, item := range bindingList.Items {
		log.Println(item.Name)
	}

}

func pvDiscovery(client v1.CoreV1Interface) {
	// pv
	pv := client.PersistentVolumes()
	list, err := pv.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("error getting persistent volume list: %s", err)
		return
	}
	log.Println("Persistent Volume List =============================================")
	for _, item := range list.Items {
		log.Println(item.Name)
	}

	// pv claim
	pvClaim := client.PersistentVolumeClaims("linkerd")
	claimList, err := pvClaim.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("error getting persistent volume list: %s", err)
		return
	}
	log.Println("Persistent Volume Claims List =============================================")
	for _, item := range claimList.Items {
		log.Println(item.Name)
	}
}

// Discovery driver
func Discovery(clientset *kubernetes.Clientset) {
	// discovery
	namespaceDiscovery(clientset.CoreV1())
	nodeDiscovery(clientset.CoreV1())
	roleDiscovery(clientset.RbacV1())
	pvDiscovery(clientset.CoreV1())

	// Services - corev1
	// Configmaps - corev1
	// Secrets - corev1
	// Deployments - appv1
	// Statefulset - appv1
	// Daemonsets - appv1
	// Jobs and cronjobs - batchv1
	// Roles/RoleBindings - rbacv1
	// Pods - corev1
	// Replicaset - appv1
}
