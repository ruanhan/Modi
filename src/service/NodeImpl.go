package service

import (
	"github.com/bigartists/Modi/src/helpers"
	"github.com/bigartists/Modi/src/model/NodeModel"
	"github.com/bigartists/Modi/src/repo"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

type INode interface {
	ListAllNodes(ns string) []*NodeModel.NodeModel
}

type NodeService struct {
	nodeRepo *repo.NodeRepo
	podRepo  *repo.PodRepo
	Metric   *versioned.Clientset
}

func ProviderNodeService(nodeRepo *repo.NodeRepo, podRepo *repo.PodRepo, metric *versioned.Clientset) INode {
	return &NodeService{nodeRepo: nodeRepo, podRepo: podRepo, Metric: metric}
}

func (this *NodeService) ListAllNodes(ns string) []*NodeModel.NodeModel {
	list, err := this.nodeRepo.List(ns)
	if err != nil {
		return nil
	}
	ret := make([]*NodeModel.NodeModel, len(list))
	for i, item := range list {
		nodeUsage := helpers.GetNodeUsage(this.Metric, item)

		ret[i] = &NodeModel.NodeModel{
			Name:     item.Name,
			IP:       item.Status.Addresses[0].Address,
			HostName: item.Status.Addresses[1].Address,
			Labels:   helpers.FilterLables(item.Labels),
			Taints:   helpers.FilterTaints(item.Spec.Taints),
			Capacity: NodeModel.NewNodeCapacity(item.Status.Capacity.Cpu().Value(),
				item.Status.Capacity.Memory().Value(), item.Status.Capacity.Pods().Value()),
			Usage:      NodeModel.NewNodeUsage(this.podRepo.GetNumByNode(item.Name), nodeUsage[0], nodeUsage[1]),
			CreateTime: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
	}
	return ret
}
