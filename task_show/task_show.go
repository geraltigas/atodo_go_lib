package task_show

import (
	"atodo_go/table"
	"fmt"
	"sort"
	"strconv"
)

type ShowEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type ShowNode struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Position Position `json:"position"`
	Status   string   `json:"status"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type ShowData struct {
	Nodes                []ShowNode `json:"nodes"`
	Edges                []ShowEdge `json:"edges"`
	NodeConnectedToStart []string   `json:"node_connected_to_start"`
	NodeConnectedToEnd   []string   `json:"node_connected_to_end"`
}

func inferenceStartAndEndNodes(nodes []ShowNode, relations []table.TaskRelation) ([]string, []string) {
	var nodeConnectedToStart []string
	var nodeConnectedToEnd []string

	sourceSet := make(map[int]bool)
	targetSet := make(map[int]bool)
	connectedMap := make(map[int]bool)

	for _, relation := range relations {
		sourceSet[relation.Source] = true
		targetSet[relation.Target] = true
	}

	intersection := make(map[int]bool)
	for k := range sourceSet {
		if targetSet[k] {
			intersection[k] = true
		}
	}

	sourceTarget := make(map[int]bool)
	for k := range sourceSet {
		if !intersection[k] {
			sourceTarget[k] = true
		}
	}

	targetSource := make(map[int]bool)
	for k := range targetSet {
		if !intersection[k] {
			targetSource[k] = true
		}
	}

	for k := range intersection {
		connectedMap[k] = true
	}

	for k := range sourceTarget {
		nodeConnectedToStart = append(nodeConnectedToStart, strconv.Itoa(k))
		connectedMap[k] = true
	}

	for k := range targetSource {
		nodeConnectedToEnd = append(nodeConnectedToEnd, strconv.Itoa(k))
		connectedMap[k] = true
	}

	for _, node := range nodes {
		nodeID, err := strconv.Atoi(node.ID)
		if err != nil {
			panic(err)
		}
		if _, ok := connectedMap[nodeID]; !ok {
			nodeConnectedToStart = append(nodeConnectedToStart, node.ID)
			nodeConnectedToEnd = append(nodeConnectedToEnd, node.ID)
		}
	}

	sort.Strings(nodeConnectedToStart)
	sort.Strings(nodeConnectedToEnd)

	return nodeConnectedToStart, nodeConnectedToEnd
}

func GetShowStack() ([]string, error) {
	nowViewingTaskID, err := table.GetNowViewingTask()
	if err != nil {
		return nil, err
	}
	if nowViewingTaskID == -1 {
		return nil, fmt.Errorf("no task is being viewed")
	}
	task, err := table.GetTaskByID(nowViewingTaskID)
	if err != nil {
		return nil, err
	}
	var stack []string
	stack = append(stack, task.Name)
	for task.ParentTask != -1 {
		task, err = table.GetTaskByID(task.ParentTask)
		if err != nil {
			return nil, err
		}
		stack = append(stack, task.Name)
	}
	return stack, nil
}

func GetShowData() (*ShowData, error) {
	nowViewingTaskID, err := table.GetNowViewingTask()
	if err != nil {
		return nil, err
	}
	if nowViewingTaskID == -1 {
		return nil, fmt.Errorf("no task is being viewed")
	}
	return GetShowDataByTaskID(nowViewingTaskID)
}

func GetShowDataByTaskID(id int) (*ShowData, error) {
	tasks, err := table.GetTasksByParentTask(id)
	if err != nil {
		return nil, err
	}
	showData := ShowData{}
	for _, task := range tasks {
		showData.Nodes = append(showData.Nodes, ShowNode{
			ID:   fmt.Sprintf("%d", task.ID),
			Name: task.Name,
			Position: Position{
				X: task.PositionX,
				Y: task.PositionY,
			},
		})
	}

	relations, err := table.GetRelationByParentTask(id)
	if err != nil {
		return nil, err
	}
	for _, relation := range relations {
		showData.Edges = append(showData.Edges, ShowEdge{
			Source: fmt.Sprintf("%d", relation.Source),
			Target: fmt.Sprintf("%d", relation.Target),
		})
	}

	nodeConnectedToStart, nodeConnectedToEnd := inferenceStartAndEndNodes(showData.Nodes, relations)
	showData.NodeConnectedToStart = nodeConnectedToStart
	showData.NodeConnectedToEnd = nodeConnectedToEnd
	return &showData, nil
}
