package service

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
	"sync"
)

type RsMapStruct struct {
	data sync.Map
}

func (this *RsMapStruct) Add(obj *v1.ReplicaSet) {
	if list, ok := this.data.Load(obj.Namespace); ok {
		list = append(list.([]*v1.ReplicaSet), obj)
		this.data.Store(obj.Namespace, list)
	} else {
		this.data.Store(obj.Namespace, []*v1.ReplicaSet{obj})
	}
}

func (this *RsMapStruct) Update(obj *v1.ReplicaSet) error {
	if list, ok := this.data.Load(obj.Namespace); ok {
		for i, rangeObj := range list.([]*v1.ReplicaSet) {
			if rangeObj.Name == obj.Name {
				list.([]*v1.ReplicaSet)[i] = obj
			}
		}
		return nil
	}
	return fmt.Errorf("rs-%s not found", obj.Name)
}

func (this *RsMapStruct) Delete(obj *v1.ReplicaSet) {
	if list, ok := this.data.Load(obj.Namespace); ok {
		for i, rangeObj := range list.([]*v1.ReplicaSet) {
			if rangeObj.Name == obj.Name {
				newList := append(list.([]*v1.ReplicaSet)[:i], list.([]*v1.ReplicaSet)[i+1:]...)
				this.data.Store(obj.Namespace, newList)
				break
			}
		}
	}
}

func (this *RsMapStruct) RsByNs(ns string) ([]*v1.ReplicaSet, error) {
	if list, ok := this.data.Load(ns); ok {
		return list.([]*v1.ReplicaSet), nil
	} else {
		return nil, fmt.Errorf("rs  not found")
	}
}

var RsMapInstance *RsMapStruct

func init() {
	RsMapInstance = &RsMapStruct{}
}
