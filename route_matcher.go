package httpx

import (
	"strings"
	"sync"

	"github.com/samber/lo"
)

type routeMatcher struct {
	mu   sync.RWMutex
	root *routeMatcherNode
}

type routeMatcherNode struct {
	staticChildren map[string]*routeMatcherNode
	paramChild     *routeMatcherNode
	catchAllChild  *routeMatcherNode
	routes         []routeMatchEntry
	minSeq         uint64
}

type routeMatchEntry struct {
	seq   uint64
	route RouteInfo
}

func newRouteMatcher() *routeMatcher {
	return &routeMatcher{
		root: &routeMatcherNode{},
	}
}

func (m *routeMatcher) Add(path string, route RouteInfo, seq uint64) {
	if m == nil || seq == 0 {
		return
	}

	segments := splitRouteSegments(path)

	m.mu.Lock()
	defer m.mu.Unlock()

	root := m.ensureRootLocked()
	root.recordMinSeq(seq)
	node := lo.Reduce(segments, func(current *routeMatcherNode, segment string, _ int) *routeMatcherNode {
		next := current.ensureChild(segment)
		next.recordMinSeq(seq)
		return next
	}, root)

	node.routes = lo.Concat(node.routes, []routeMatchEntry{{
		seq:   seq,
		route: route,
	}})
}

func (m *routeMatcher) Match(path string) (RouteInfo, bool) {
	if m == nil {
		return RouteInfo{}, false
	}

	segments := splitRouteSegments(path)

	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.root == nil {
		return RouteInfo{}, false
	}

	matched, ok := m.root.match(segments, 0, pathHasTrailingSlash(path))
	if !ok {
		return RouteInfo{}, false
	}
	return matched.route, true
}

func (n *routeMatcherNode) ensureChild(segment string) *routeMatcherNode {
	if n == nil {
		return nil
	}
	if isCatchAllPathParameterSegment(segment) {
		if n.catchAllChild == nil {
			n.catchAllChild = &routeMatcherNode{}
		}
		return n.catchAllChild
	}
	if isPathParameterSegment(segment) {
		if n.paramChild == nil {
			n.paramChild = &routeMatcherNode{}
		}
		return n.paramChild
	}
	if n.staticChildren == nil {
		n.staticChildren = map[string]*routeMatcherNode{}
	}
	if n.staticChildren[segment] == nil {
		n.staticChildren[segment] = &routeMatcherNode{}
	}
	return n.staticChildren[segment]
}

func (m *routeMatcher) ensureRootLocked() *routeMatcherNode {
	if m.root == nil {
		m.root = &routeMatcherNode{}
	}
	return m.root
}

func (n *routeMatcherNode) match(segments []string, index int, trailingSlash bool) (routeMatchEntry, bool) {
	if n == nil {
		return routeMatchEntry{}, false
	}

	if index == len(segments) {
		return n.matchCurrent(trailingSlash)
	}

	return n.matchNext(segments, index, trailingSlash)
}

func (n *routeMatcherNode) matchCurrent(trailingSlash bool) (routeMatchEntry, bool) {
	if matched, ok := n.routeAtCurrentNode(); ok {
		return matched, true
	}
	if trailingSlash {
		return matchCatchAllRouteChild(n.catchAllChild)
	}
	return routeMatchEntry{}, false
}

func (n *routeMatcherNode) matchNext(segments []string, index int, trailingSlash bool) (routeMatchEntry, bool) {
	segment := segments[index]
	staticChild := n.staticChildren[segment]
	paramChild := n.paramChild
	for _, child := range orderedRouteChildren(staticChild, paramChild, n.catchAllChild) {
		if child.catchAll {
			if matched, ok := matchCatchAllRouteChild(child.node); ok {
				return matched, true
			}
			continue
		}
		if matched, ok := matchRouteChild(child.node, segments, index+1, trailingSlash); ok {
			return matched, true
		}
	}

	return routeMatchEntry{}, false
}

func (n *routeMatcherNode) routeAtCurrentNode() (routeMatchEntry, bool) {
	if len(n.routes) == 0 {
		return routeMatchEntry{}, false
	}
	return n.routes[0], true
}

func (n *routeMatcherNode) recordMinSeq(seq uint64) {
	if n == nil || seq == 0 {
		return
	}
	if n.minSeq == 0 || seq < n.minSeq {
		n.minSeq = seq
	}
}

type routeMatcherChild struct {
	node     *routeMatcherNode
	catchAll bool
}

func orderedRouteChildren(children ...*routeMatcherNode) []routeMatcherChild {
	ordered := orderedNonCatchAllRouteChildren(children)
	if len(children) > 2 && children[2] != nil {
		ordered = append(ordered, routeMatcherChild{
			node:     children[2],
			catchAll: true,
		})
	}
	return ordered
}

func orderedNonCatchAllRouteChildren(children []*routeMatcherNode) []routeMatcherChild {
	limit := min(len(children), 2)
	ordered := make([]routeMatcherChild, 0, len(children))
	for _, child := range children[:limit] {
		if child == nil {
			continue
		}
		ordered = insertRouteMatcherChild(ordered, routeMatcherChild{
			node: child,
		})
	}
	return ordered
}

func insertRouteMatcherChild(ordered []routeMatcherChild, candidate routeMatcherChild) []routeMatcherChild {
	insertAt := len(ordered)
	for i, existing := range ordered {
		if childSeqLess(candidate.node, existing.node) {
			insertAt = i
			break
		}
	}
	ordered = append(ordered, routeMatcherChild{})
	copy(ordered[insertAt+1:], ordered[insertAt:])
	ordered[insertAt] = candidate
	return ordered
}

func childSeqLess(left, right *routeMatcherNode) bool {
	switch {
	case right == nil:
		return true
	case left == nil:
		return false
	case left.minSeq == 0:
		return false
	case right.minSeq == 0:
		return true
	default:
		return left.minSeq < right.minSeq
	}
}

func matchRouteChild(node *routeMatcherNode, segments []string, index int, trailingSlash bool) (routeMatchEntry, bool) {
	if node == nil {
		return routeMatchEntry{}, false
	}
	return node.match(segments, index, trailingSlash)
}

func matchCatchAllRouteChild(node *routeMatcherNode) (routeMatchEntry, bool) {
	if node == nil {
		return routeMatchEntry{}, false
	}
	return node.routeAtCurrentNode()
}

func splitRouteSegments(path string) []string {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "/")
}

func isPathParameterSegment(segment string) bool {
	return strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}")
}

func isCatchAllPathParameterSegment(segment string) bool {
	return strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "...}")
}

func pathHasTrailingSlash(path string) bool {
	return path != "/" && strings.HasSuffix(path, "/")
}
