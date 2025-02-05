package informer

import (
	"github.com/bigartists/Modi/client"
	"github.com/bigartists/Modi/src/repo"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
)

type InformerManager struct {
	secretRepo    *repo.SecretRepo
	configmapRepo *repo.ConfigMapRepo
	fac           informers.SharedInformerFactory
	eventRepo     *repo.EventRepo
	namespaceRepo *repo.NamespaceRepo
	rsRep         *repo.RsRep
	deployRepo    *repo.DeploymentRepo
	podRepo       *repo.PodRepo
	serviceRepo   *repo.ServiceRepo
	nodeRepo      *repo.NodeRepo
}

func NewInformerManager(
	secretRepo *repo.SecretRepo,
	configmapRepo *repo.ConfigMapRepo,
	eventRepo *repo.EventRepo,
	namespaceRepo *repo.NamespaceRepo,
	rsRep *repo.RsRep,
	deployRepo *repo.DeploymentRepo,
	podRepo *repo.PodRepo,
	serviceRepo *repo.ServiceRepo,
	nodeRepo *repo.NodeRepo,
) *InformerManager {
	return &InformerManager{
		secretRepo:    secretRepo,
		configmapRepo: configmapRepo,
		fac:           informers.NewSharedInformerFactory(client.K8sClient, 0),
		eventRepo:     eventRepo,
		namespaceRepo: namespaceRepo,
		rsRep:         rsRep,
		deployRepo:    deployRepo,
		podRepo:       podRepo,
		serviceRepo:   serviceRepo,
		nodeRepo:      nodeRepo,
	}
}

func (this *InformerManager) InitInformers() error {
	//fact := informers.NewSharedInformerFactory(client.K8sClient, 0)
	fact := this.fac
	nsInformer := fact.Core().V1().Namespaces().Informer()
	_, err := nsInformer.AddEventHandler(this.namespaceRepo)
	if err != nil {
		return err
	}

	// 初始化 deployment监听
	depInformer := fact.Apps().V1().Deployments().Informer()
	_, err = depInformer.AddEventHandler(this.deployRepo)
	if err != nil {
		return err
	}

	// 初始化 replicaSet监听
	rsInformer := fact.Apps().V1().ReplicaSets().Informer()
	_, err = rsInformer.AddEventHandler(this.rsRep)
	if err != nil {
		return err
	}

	// 初始化 pod监听
	podInformer := fact.Core().V1().Pods().Informer()
	_, err = podInformer.AddEventHandler(this.podRepo)
	if err != nil {
		return err
	}

	// 初始化 event 监听
	eventInformer := fact.Core().V1().Events().Informer()
	_, err = eventInformer.AddEventHandler(this.eventRepo)
	if err != nil {
		return err
	}

	// 初始化 secret 监听
	secretInformer := fact.Core().V1().Secrets().Informer()
	_, err = secretInformer.AddEventHandler(this.secretRepo)
	if err != nil {
		return err
	}

	// 初始化configMap 监听
	configMapInformer := fact.Core().V1().ConfigMaps().Informer()
	_, err = configMapInformer.AddEventHandler(this.configmapRepo)
	if err != nil {
		return err
	}

	// 初始化 service 的监听
	serviceInformer := fact.Core().V1().Services().Informer()
	_, err = serviceInformer.AddEventHandler(this.serviceRepo)
	if err != nil {
		return err
	}

	// 初始化 node的监听

	nodeInformer := fact.Core().V1().Nodes().Informer()
	_, err = nodeInformer.AddEventHandler(this.nodeRepo)

	fact.Start(wait.NeverStop)
	return nil
}
