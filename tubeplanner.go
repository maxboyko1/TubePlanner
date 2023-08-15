package main

import (
	"container/heap"
	"fmt"
	"math"
	"os"
	"slices"
)

// Represents an "edge" in the transit graph, either a rail link or an interchange
type Link struct {
	endNode  *Node
	time     uint16
	linkType string
}

// Represents a "vertex" in the transit graph, with each existing combination
// of station name and line name being its own vertex
type Node struct {
	station   string
	line      string
	adj       []*Link
	totalTime uint16
	index     int
}

// Map of each station and line name combination to its corresponding Node
// pointer in the graph
type NodeMap map[string]map[string]*Node

// List of all nodes in the graph, min heap-ordered according to the shortest
// time taken to arrive there from user's chosen starting point (all necessary
// Go heap interface methods are implemented below)
type NodePriorityQueue []*Node

// Return number of nodes in the heap
func (npq NodePriorityQueue) Len() int {
	return len(npq)
}

// Return whether the total travel time to Node at index i is less than the
// total travel time to Node at index j
func (npq NodePriorityQueue) Less(i, j int) bool {
	return npq[i].totalTime < npq[j].totalTime
}

// Swap positions of Nodes at indices i and j in the heap
func (npq NodePriorityQueue) Swap(i, j int) {
	npq[i], npq[j] = npq[j], npq[i]
	npq[i].index = i
	npq[j].index = j
}

// Add a new Node to the end of the heap
func (npq *NodePriorityQueue) Push(x any) {
	n := len(*npq)
	node := x.(*Node)
	node.index = n
	*npq = append(*npq, node)
}

// Remove the minimum priority Node from the heap and return it
func (npq *NodePriorityQueue) Pop() any {
	old := *npq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.index = -1
	*npq = old[0 : n-1]
	return node
}

// Update the specified node with a new total travel time, then restore the heap ordering
func (npq *NodePriorityQueue) update(node *Node, newTotalTime uint16) {
	node.totalTime = newTotalTime
	heap.Fix(npq, node.index)
}

// Helper function for BuildTransitGraph() which adds a connection between two
// Nodes of the specified type and transit time to the graph
func AddConnection(npq *NodePriorityQueue, nodeMap NodeMap, connection any, lType string) {
	// Retrieve station/line names and transit time for the specified connection
	var stationA, lineA, stationB, lineB string
	var transitTime uint16
	switch conn := connection.(type) {
	case *RailLink:
		stationA, stationB = conn.fromStation, conn.toStation
		lineA, lineB = conn.line, conn.line
		transitTime = conn.transitTime
	case *Interchange:
		stationA, stationB = conn.fromStation, conn.toStation
		lineA, lineB = conn.fromLine, conn.toLine
		transitTime = conn.transitTime
	default:
		fmt.Fprintln(os.Stderr, "ERROR: Connection type must be RailLink or Interchange")
		os.Exit(1)
	}

	// Create a graph Node for the first station/line if it does not exist already,
	// with travel distance initialized to infinity
	var nodeAExists bool = false
	_, mapAExists := nodeMap[stationA]
	if mapAExists {
		_, nodeAExists = nodeMap[stationA][lineA]
	} else {
		nodeMap[stationA] = make(map[string]*Node)
	}
	if !nodeAExists {
		newNode := &Node{stationA, lineA, make([]*Link, 0), math.MaxUint16, 0}
		npq.Push(newNode)
		nodeMap[stationA][lineA] = newNode
	}

	// Create a graph Node for the second station/line, if it does not exist already,
	// with travel distance initialized to infinity
	var nodeBExists bool = false
	_, mapBExists := nodeMap[stationB]
	if mapBExists {
		_, nodeBExists = nodeMap[stationB][lineB]
	} else {
		nodeMap[stationB] = make(map[string]*Node)
	}
	if !nodeBExists {
		newNode := &Node{stationB, lineB, make([]*Link, 0), math.MaxUint16, 0}
		npq.Push(newNode)
		nodeMap[stationB][lineB] = newNode
	}

	// Add a link to node B to node A's adjacency list, and vice versa
	nodeA, nodeB := nodeMap[stationA][lineA], nodeMap[stationB][lineB]
	nodeA.adj = append(nodeA.adj, &Link{nodeB, transitTime, lType})
	nodeB.adj = append(nodeB.adj, &Link{nodeA, transitTime, lType})
}

// Retrieve the list of rail links and interchanges defined in transitdata.go
// and add each one as a connection in the transit graph
func BuildTransitGraph() (NodePriorityQueue, NodeMap) {
	railLinks, interchanges := GetRailLinks(), GetInterchanges()
	npq, nodeMap := make(NodePriorityQueue, 0), make(NodeMap)

	for _, rl := range railLinks {
		AddConnection(&npq, nodeMap, &rl, "rail")
	}
	for _, ic := range interchanges {
		if ic.fromStation == ic.toStation {
			AddConnection(&npq, nodeMap, &ic, "line interchange")
		} else {
			AddConnection(&npq, nodeMap, &ic, "station interchange")
		}
	}

	return npq, nodeMap
}

// Run a binary heap variation of Dijkstra's shortest paths algorithm on the
// completed transit graph to calculate the shortest possible trip between
// the provided start and end stations
func RunShortestPaths(npq *NodePriorityQueue, nodeMap NodeMap,
	start, dest string) ([]*Node, []string) {
	if start == dest {
		return nil, nil
	}
	nodePrev := make(map[*Node]*Node)
	linkPrev := make(map[*Node]*Link)
	// Initialize valid starting Nodes in graph (any transit line departing
	// from specified start station) with travel times of 0
	for _, node := range nodeMap[start] {
		npq.update(node, 0)
		nodePrev[node] = nil
		linkPrev[node] = nil
	}
	var curNode *Node = nil
	for len(*npq) > 0 {
		// Retrieve the Node of minimum established travel time from the heap
		curNode = heap.Pop(npq).(*Node)
		// If this Node represents the desired destination, we are done
		if curNode.station == dest {
			break
		}
		// For every node directly reachable from the current node, update the
		// travel time to that node if the path to it from the current node is
		// an improvement on its previously established travel time
		for _, link := range curNode.adj {
			altDistance := curNode.totalTime + link.time
			if altDistance < link.endNode.totalTime {
				link.endNode.totalTime = altDistance
				nodePrev[link.endNode] = curNode
				linkPrev[link.endNode] = link
				npq.update(link.endNode, altDistance)
			}
		}
	}
	// Construct the route from the start to ending Nodes by continually
	// following pointers to the previous node in the path until the start is
	// reached, tracking the type of the link at each step as well
	route, linkTypes := make([]*Node, 0), make([]string, 0)
	for linkPrev[curNode] != nil {
		route = append(route, curNode)
		linkTypes = append(linkTypes, linkPrev[curNode].linkType)
		curNode = nodePrev[curNode]
	}
	route = append(route, curNode)
	slices.Reverse(linkTypes)
	slices.Reverse(route)
	return route, linkTypes
}

// From the specified transit trip, as represented by the sequence of nodes
// visited as well as the types of connections between each, print a clear,
// readable series of directions for the user to follow to complete their trip
func PrintDirections(route []*Node, linkTypes []string) {
	if route == nil {
		fmt.Println("Already at destination!")
		return
	}
	fmt.Printf("1) Begin journey at %s station. (0 minutes)\n", route[0].station)
	var idx, step int
	for idx, step = 0, 2; idx < len(linkTypes); idx++ {
		switch linkTypes[idx] {
		case "rail":
			if idx == 0 || linkTypes[idx-1] != "rail" {
				fmt.Printf("%d) Travel on the %s line, through station stops:\n",
					step, route[idx+1].line)
				step++
			}
			fmt.Printf("- %s (%d minutes)\n", route[idx+1].station, route[idx+1].totalTime)
		case "line interchange":
			fmt.Printf("%d) Get off at %s and interchange to the %s line. (%d minutes)\n",
				step, route[idx+1].station, route[idx+1].line, route[idx+1].totalTime)
			step++
		case "station interchange":
			fmt.Printf("%d) From %s, interchange on foot to nearby %s station. (%d minutes)\n",
				step, route[idx].station, route[idx+1].station, route[idx+1].totalTime)
			step++
		default:
			fmt.Fprintf(os.Stderr, "ERROR: Invalid transit link type: %s\n", linkTypes[idx])
			os.Exit(1)
		}
	}
	fmt.Printf("%d) Reach destination at %s station. (%d minutes)\n",
		step, route[idx].station, route[idx].totalTime)
}

// Program that builds a graph to represent the London commuter transit map data
// specified in transitdata.go, computes the shortest possible trip (in minutes)
// between the user-provided start and end point stations, and prints to console
// a series of directions to follow to complete said trip
func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "USAGE: ./tubeplanner <start> <destination>")
		os.Exit(1)
	}
	graph, nodeMap := BuildTransitGraph()
	start, dest := os.Args[1], os.Args[2]
	if _, startExists := nodeMap[start]; !startExists {
		fmt.Fprintf(os.Stderr, "ERROR: %s is not a valid initial station\n", start)
		os.Exit(1)
	}
	if _, destExists := nodeMap[dest]; !destExists {
		fmt.Fprintf(os.Stderr, "ERROR: %s is not a valid destination\n", dest)
		os.Exit(1)
	}
	route, linkTypes := RunShortestPaths(&graph, nodeMap, start, dest)
	PrintDirections(route, linkTypes)
}
