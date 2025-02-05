package NodeModel

type NodeModel struct {
	Name       string
	IP         string
	HostName   string
	Labels     []string
	Taints     []string
	Capacity   *NodeCapacity
	Usage      *NodeUsage
	CreateTime string
}

type NodeUsage struct {
	Pods   int
	Cpu    float64
	Memory float64
}

func NewNodeUsage(pods int, cpu float64, memory float64) *NodeUsage {
	return &NodeUsage{Pods: pods, Cpu: cpu, Memory: memory}
}

type NodeCapacity struct {
	Cpu    int64
	Memory int64
	Pods   int64
}

func NewNodeCapacity(cpu int64, memory int64, pods int64) *NodeCapacity {
	return &NodeCapacity{Cpu: cpu, Memory: memory, Pods: pods}
}
