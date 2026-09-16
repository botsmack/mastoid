package pkg

import (
	"testing"

	"github.com/mattn/go-mastodon"
)

func makeThreadWithMissingIntermediateStatus() *Thread {
	thread := &Thread{
		Nodes: map[mastodon.ID]*Node{},
	}

	root := thread.GetNode("root")
	root.Status = &mastodon.Status{
		ID: "root",
	}

	missing := thread.GetNode("missing")

	child := thread.GetNode("child")
	child.Status = &mastodon.Status{
		ID: "child",
	}

	root.Descendants.Set(missing.ID, missing)
	missing.Ancestors.Set(root.ID, root)

	missing.Descendants.Set(child.ID, child)
	child.Ancestors.Set(missing.ID, missing)

	return thread
}

func TestWalkDepthFirstSkipsMissingStatusAndContinuesTraversal(t *testing.T) {
	thread := makeThreadWithMissingIntermediateStatus()

	var visited []mastodon.ID

	err := thread.WalkDepthFirst(func(node *Node, depth int, siblingIdx int) error {
		visited = append(visited, node.Status.ID)
		return nil
	})
	if err != nil {
		t.Fatalf("WalkDepthFirst returned error: %v", err)
	}

	if len(visited) != 2 {
		t.Fatalf("expected 2 visited statuses, got %d: %v", len(visited), visited)
	}

	if visited[0] != "root" {
		t.Errorf("expected first status to be root, got %q", visited[0])
	}

	if visited[1] != "child" {
		t.Errorf("expected second status to be child, got %q", visited[1])
	}
}

func TestWalkBreadthFirstSkipsMissingStatusAndContinuesTraversal(t *testing.T) {
	thread := makeThreadWithMissingIntermediateStatus()

	var visited []mastodon.ID

	err := thread.WalkBreadthFirst(func(node *Node, depth int, siblingIdx int) error {
		visited = append(visited, node.Status.ID)
		return nil
	})
	if err != nil {
		t.Fatalf("WalkBreadthFirst returned error: %v", err)
	}

	if len(visited) != 2 {
		t.Fatalf("expected 2 visited statuses, got %d: %v", len(visited), visited)
	}

	if visited[0] != "root" {
		t.Errorf("expected first status to be root, got %q", visited[0])
	}

	if visited[1] != "child" {
		t.Errorf("expected second status to be child, got %q", visited[1])
	}
}
