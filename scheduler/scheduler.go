package scheduler

import (
	"cube/node"
	"cube/task"
	"log"
	"math"
	"time"
)

const (
	// LIEB square ice constant
	// https://en.wikipedia.org/wiki/Lieb%27s_square_ice_constant
	LIEB = 1.53960071783900203869
)

type Scheduler interface {
	SelectCandidateNodes(t task.Task, nodes []*node.Node) []*node.Node
	Score(t task.Task, nodes []*node.Node) map[string]float64
	Pick(scores map[string]float64, candidates []*node.Node) *node.Node
}

type RoundRobin struct {
	Name       string
	lastWorker int
}
type Epvm struct {
	Name string
}

func (r *RoundRobin) SelectCandidateNodes(t task.Task, nodes []*node.Node) []*node.Node {
	return nodes
}

func (r *RoundRobin) Score(t task.Task, nodes []*node.Node) map[string]float64 {
	nodeScores := make(map[string]float64)
	var newWorker int
	if r.lastWorker+1 < len(nodes) {
		newWorker = r.lastWorker + 1
		r.lastWorker++
	} else {
		newWorker = 0
		r.lastWorker = 0
	}
	for idx, node := range nodes {
		if idx == newWorker {
			nodeScores[node.Name] = 0.1
		} else {
			nodeScores[node.Name] = 1.0
		}
	}
	return nodeScores
}

func (r *RoundRobin) Pick(scores map[string]float64, candidates []*node.Node) *node.Node {
	var bestNode *node.Node
	var lowestScore float64
	for idx, node := range candidates {
		if idx == 0 {
			bestNode = node
			lowestScore = scores[node.Name]
			continue
		}
		if scores[node.Name] < lowestScore {
			bestNode = node
			lowestScore = scores[node.Name]
		}
	}
	return bestNode
}

func (e *Epvm) SelectCandidateNodes(t task.Task, nodes []*node.Node) []*node.Node {
	var candidates []*node.Node
	for node := range nodes {
		if checkDisk(t, nodes[node].Disk-nodes[node].DiskAllocated) {
			candidates = append(candidates, nodes[node])
		}
	}
	return candidates
}
func checkDisk(t task.Task, diskAvailable int64) bool {
	return t.Disk <= diskAvailable
}

func calculateLoad(usage float64, capacity float64) float64 {
	return usage / capacity
}

func calculateCpuUsage(node *node.Node) (*float64, error) {
	//stat1 := getNodeStats(node)
	stat1, err := node.GetStats()
	if err != nil {
		return nil, err
	}
	time.Sleep(3 * time.Second)
	//stat2 := getNodeStats(node)
	stat2, err := node.GetStats()
	if err != nil {
		return nil, err
	}
	stat1Idle := stat1.CpuStats.Idle + stat1.CpuStats.IOWait
	stat2Idle := stat2.CpuStats.Idle + stat2.CpuStats.IOWait
	stat1NonIdle := stat1.CpuStats.User + stat1.CpuStats.Nice + stat1.CpuStats.System + stat1.CpuStats.IRQ + stat1.CpuStats.SoftIRQ + stat1.CpuStats.Steal
	stat2NonIdle := stat2.CpuStats.User + stat2.CpuStats.Nice + stat2.CpuStats.System + stat2.CpuStats.IRQ + stat2.CpuStats.SoftIRQ + stat2.CpuStats.Steal
	stat1Total := stat1Idle + stat1NonIdle
	stat2Total := stat2Idle + stat2NonIdle
	total := stat2Total - stat1Total
	idle := stat2Idle - stat1Idle
	var cpuPercentUsage float64
	if total == 0 && idle == 0 {
		cpuPercentUsage = 0.00
	} else {
		cpuPercentUsage = (float64(total) - float64(idle)) / float64(total)
	}
	return &cpuPercentUsage, nil
}

func (e *Epvm) Score(t task.Task, nodes []*node.Node) map[string]float64 {
	nodeScores := make(map[string]float64)
	maxJobs := 4.0

	for _, node := range nodes {
		cpuUsage, err := calculateCpuUsage(node)
		if err != nil {
			log.Printf("error calculating CPU usage for node %s, skipping: %v", node.Name, err)
			continue
		}
		cpuLoad := calculateLoad(*cpuUsage, math.Pow(2, 0.8))

		memoryAllocated := float64(node.Stats.MemUsedKb()) + float64(node.MemoryAllocated)
		memoryPercentAllocated := memoryAllocated / float64(node.Memory)

		newMemPercent := (calculateLoad(memoryAllocated+float64(t.Memory/1000), float64(node.Memory)))
		memCost := math.Pow(LIEB, newMemPercent) + math.Pow(LIEB, (float64(node.TaskCount+1))/maxJobs) - math.Pow(LIEB, memoryPercentAllocated) - math.Pow(LIEB, float64(node.TaskCount)/float64(maxJobs))
		cpuCost := math.Pow(LIEB, cpuLoad) + math.Pow(LIEB, (float64(node.TaskCount+1))/maxJobs) - math.Pow(LIEB, cpuLoad) - math.Pow(LIEB, float64(node.TaskCount)/float64(maxJobs))

		nodeScores[node.Name] = memCost + cpuCost
	}
	return nodeScores
}

func (e *Epvm) Pick(scores map[string]float64, candidates []*node.Node) *node.Node {
	minCost := 0.00
	var bestNode *node.Node
	for idx, node := range candidates {
		if idx == 0 {
			minCost = scores[node.Name]
			bestNode = node
			continue
		}
		if scores[node.Name] < minCost {
			minCost = scores[node.Name]
			bestNode = node
		}
	}
	return bestNode
}
