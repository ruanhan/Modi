package repo

import (
	"fmt"
	"sort"
	"sync"

	corev1 "k8s.io/api/core/v1"
)

type CoreV1Node []*corev1.Node

func (this CoreV1Node) Len() int {
	return len(this)
}

func (this CoreV1Node) Less(i, j int) bool {
	//根据时间排序    倒排序
	fmt.Println("this[i].CreationTimestamp.Time)===", this[i].CreationTimestamp.Time)
	return this[i].CreationTimestamp.Time.After(this[j].CreationTimestamp.Time)
}
func (this CoreV1Node) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

type NodeRepo struct {
	nodeMap sync.Map
}

func ProviderNodeRepo() *NodeRepo {
	return &NodeRepo{nodeMap: sync.Map{}}
}

func (this *NodeRepo) OnAdd(obj interface{}, isInInitialList bool) {
	this.Add(obj.(*corev1.Node))
}

func (this *NodeRepo) OnUpdate(oldObj, newObj interface{}) {
	err := this.Update(newObj.(*corev1.Node))
	if err != nil {
		return
	}
}

func (this *NodeRepo) OnDelete(obj interface{}) {
	this.Delete(obj.(*corev1.Node))
}

func (this *NodeRepo) Get(ns string, name string) *corev1.Node {
	if items, ok := this.nodeMap.Load(ns); ok {
		for _, item := range items.([]*corev1.Node) {
			if item.Name == name {
				return item
			}
		}
	}
	return nil
}
func (this *NodeRepo) Add(item *corev1.Node) {
	if list, ok := this.nodeMap.Load(item.Namespace); ok {
		list = append(list.([]*corev1.Node), item)
		this.nodeMap.Store(item.Namespace, list)
	} else {
		this.nodeMap.Store(item.Namespace, []*corev1.Node{item})
	}
}
func (this *NodeRepo) Update(item *corev1.Node) error {
	if list, ok := this.nodeMap.Load(item.Namespace); ok {
		for i, range_item := range list.([]*corev1.Node) {
			if range_item.Name == item.Name {
				list.([]*corev1.Node)[i] = item
			}
		}
		return nil
	}
	return fmt.Errorf("Service-%s not found", item.Name)
}
func (this *NodeRepo) Delete(svc *corev1.Node) {
	if list, ok := this.nodeMap.Load(svc.Namespace); ok {
		for i, range_item := range list.([]*corev1.Node) {
			if range_item.Name == svc.Name {
				newList := append(list.([]*corev1.Node)[:i], list.([]*corev1.Node)[i+1:]...)
				this.nodeMap.Store(svc.Namespace, newList)
				break
			}
		}
	}
}

func (this *NodeRepo) List(namespace string) ([]*corev1.Node, error) {
	// 如果指定了namespace，只返回该namespace下的services
	if namespace != "" {
		if services, ok := this.nodeMap.Load(namespace); ok {
			svcs := services.([]*corev1.Node)
			sort.Sort(CoreV1Node(svcs))
			return svcs, nil
		}
		return []*corev1.Node{}, nil
	}

	// 如果没有指定namespace，返回所有namespace下的services
	var allServices []*corev1.Node
	this.nodeMap.Range(func(key, value interface{}) bool {
		services := value.([]*corev1.Node)
		allServices = append(allServices, services...)
		return true
	})
	sort.Sort(CoreV1Node(allServices))
	return allServices, nil
}
