package repo

import (
	"fmt"
	"sort"
	"sync"

	corev1 "k8s.io/api/core/v1"
)

type CoreV1Service []*corev1.Service

func (this CoreV1Service) Len() int {
	return len(this)
}

func (this CoreV1Service) Less(i, j int) bool {
	//根据时间排序    倒排序
	return this[i].CreationTimestamp.Time.After(this[j].CreationTimestamp.Time)
}
func (this CoreV1Service) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

type ServiceRepo struct {
	serviceMap sync.Map
}

func ProviderServiceRepo() *ServiceRepo {
	return &ServiceRepo{serviceMap: sync.Map{}}
}

func (this *ServiceRepo) OnAdd(obj interface{}, isInInitialList bool) {
	this.Add(obj.(*corev1.Service))
}

func (this *ServiceRepo) OnUpdate(oldObj, newObj interface{}) {
	err := this.Update(newObj.(*corev1.Service))
	if err != nil {
		return
	}
}

func (this *ServiceRepo) OnDelete(obj interface{}) {
	this.Delete(obj.(*corev1.Service))
}

func (this *ServiceRepo) Get(ns string, name string) *corev1.Service {
	if items, ok := this.serviceMap.Load(ns); ok {
		for _, item := range items.([]*corev1.Service) {
			if item.Name == name {
				return item
			}
		}
	}
	return nil
}
func (this *ServiceRepo) Add(item *corev1.Service) {
	if list, ok := this.serviceMap.Load(item.Namespace); ok {
		list = append(list.([]*corev1.Service), item)
		this.serviceMap.Store(item.Namespace, list)
	} else {
		this.serviceMap.Store(item.Namespace, []*corev1.Service{item})
	}
}
func (this *ServiceRepo) Update(item *corev1.Service) error {
	if list, ok := this.serviceMap.Load(item.Namespace); ok {
		for i, range_item := range list.([]*corev1.Service) {
			if range_item.Name == item.Name {
				list.([]*corev1.Service)[i] = item
			}
		}
		return nil
	}
	return fmt.Errorf("Service-%s not found", item.Name)
}
func (this *ServiceRepo) Delete(svc *corev1.Service) {
	if list, ok := this.serviceMap.Load(svc.Namespace); ok {
		for i, range_item := range list.([]*corev1.Service) {
			if range_item.Name == svc.Name {
				newList := append(list.([]*corev1.Service)[:i], list.([]*corev1.Service)[i+1:]...)
				this.serviceMap.Store(svc.Namespace, newList)
				break
			}
		}
	}
}

func (this *ServiceRepo) List(namespace string) ([]*corev1.Service, error) {
	// 如果指定了namespace，只返回该namespace下的services
	if namespace != "" {
		if services, ok := this.serviceMap.Load(namespace); ok {
			svcs := services.([]*corev1.Service)
			sort.Sort(CoreV1Service(svcs))
			return svcs, nil
		}
		return []*corev1.Service{}, nil
	}

	// 如果没有指定namespace，返回所有namespace下的services
	var allServices []*corev1.Service
	this.serviceMap.Range(func(key, value interface{}) bool {
		services := value.([]*corev1.Service)
		allServices = append(allServices, services...)
		return true
	})
	sort.Sort(CoreV1Service(allServices))
	return allServices, nil
}
