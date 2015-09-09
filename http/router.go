package http

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"
)

type RouteVariables map[string]string

type Router struct {
	mutex sync.RWMutex
	root  *route
}

const (
	kAnyPattern int = 1 << iota
	kRegexPattern
	kAbsolutePattern
)

type route struct {
	parent   *route
	child    *route
	lsibling *route
	rsibilng *route

	// "/{category}/{file}/{line}\d{1,}"
	// names are "category", "file" and "line"
	name string
	// may be "", in "any-pattern" matches anything
	pattern string
	// should use `regex` to match first, else fallback to use `pattern`
	regex *regexp.Regexp

	index    int
	priority int
	// min/max depth
	depths [2]int

	// handler, or may be uesed to store any data you want
	handler interface{}
}

func (r *route) String() string {
	h := "0"
	if r.handler != nil {
		h = fmt.Sprintf("%p", r.handler)
	}
	str := fmt.Sprintf("/%s{n:%s;[%d,%d,%d,%d];h:%s;}",
		r.pattern, r.name, r.index, r.priority, r.depths[0], r.depths[1], h)
	return str
}

func (r *route) isTail() bool      { return r.child == nil }
func (r *route) isAnyTail() bool   { return r.isTail() && r.priority == kAnyPattern }
func (r *route) isSlashTail() bool { return r.isTail() && r.name+r.pattern == "" }

func (r *route) match(part string, justRouteLevel bool, rvs RouteVariables) bool {
	result := false

	if justRouteLevel {
		for {
			_, tomatch, isRegex := splitRouteNameAndMatchPattern(part)
			if isRegex && r.regex == nil {
				result = false
				break
			}
			result = (r.pattern == tomatch)
			break
		}
	} else {
		switch r.priority {
		case kAnyPattern:
			// "Any pattern" only matches non-empty part.
			result = part != ""
		case kRegexPattern:
			result = r.regex != nil && r.regex.MatchString(part)
		case kAbsolutePattern:
			// http://stackoverflow.com/questions/7996919/should-url-be-case-sensitive
			// URLs in general are case-sensitive (with the exception of machine names).
			// There may be URLs, or parts of URLs, where case doesn't matter,
			// but identifying these may not be easy.
			// Users should always consider that URLs are case-sensitive.
			result = r.pattern == part
		}

		// Fill in the route variables map.
		if result && r.name != "" {
			rvs[r.name] = part
		}
	}

	return result
}

func (r *route) equal(other *route) bool {
	if r.index != other.index {
		return false
	}
	if r.pattern != other.pattern {
		return false
	}
	if r.name != other.name {
		return false
	}
	if r.priority != other.priority {
		return false
	}
	return true
}

func NewRouter() *Router { return &Router{root: &route{index: -1, pattern: ""}} }

// Just be responsible for mapping patterns to their specific handlers.
// Build a tree internally.
func (rt *Router) Handle(pattern string, handler interface{}) error {
	if isNil(handler) {
		return fmt.Errorf("nil handler")
	}

	// Make routes from the pattern.
	routes, err := makeRoutes(pattern, handler)
	if err != nil {
		return err
	}

	// Find if there already exists a same pattern.
	if rt.PatternExists(pattern) {
		return fmt.Errorf("pattern %q already exists", pattern)
	}

	rt.mutex.Lock()
	defer rt.mutex.Unlock()

	// Rebuild the tree.
	rt.rebuildRouteTree(routes)
	return nil
}

type matchResult struct {
	Path string
	// The handler which exactly matches the full path. Maybe nil.
	Handler interface{}
	// Non-nil handlers on the way in a reverse order. For example:
	// rt.Handle("/a/", handler_1)
	// rt.Handle("/a/b/", handler_2)
	// rt.Handle("/a/b/c", handler_3)
	// rt.Match("/a/b/c") will get result like this:
	// {
	//    Path: "/a/b/c",
	//    Handler: handler_3,
	//    HandlersOnTheWay: []{handler_2, handler_1},
	//    ...
	// }
	HandlersOnTheWay []RouteMatchItem
	// Route variables will be extracted from the route parts. For instance,
	// pattern "/{category}/{file}/{line}\\d+?" matches path "/golang/main.go/13",
	// and the route variables are:
	// map[string]string{
	//   "category":  "golang",
	//   "file":    "main.go",
	//   "line":    "13",
	// }
	RouteVars RouteVariables
}

// Match the real path to a specified handler.
// You should give a normalized path. (eg. path.Clean(...))
func (rt *Router) Match(path string) matchResult {
	return rt.match(path, false)
}

func (rt *Router) PatternExists(pattern string) bool {
	mr := rt.match(pattern, true)
	return mr.Handler != nil
}

type RouteMatchItem struct {
	Path    string
	Handler interface{}
}

type traceTable map[int][]RouteMatchItem

func (t traceTable) tryInsert(depth, pos int, path string, handler interface{}) {
	if handler == nil {
		return
	}
	if _, exists := t[depth]; !exists {
		t[depth] = make([]RouteMatchItem, 2)
	}
	if t[depth][pos].Handler == nil {
		t[depth][pos].Path = path
		t[depth][pos].Handler = handler
	}
}

func (t traceTable) getHandlersOnTheWay() []RouteMatchItem {
	depths, handlers := make([]int, 0), make([]RouteMatchItem, 0)
	for depth := range t {
		depths = append(depths, depth)
	}
	if len(depths) == 0 {
		return handlers
	}
	sort.Ints(depths)
	for _, depth := range depths {
		for pos := 0; pos < 2; pos++ {
			if t[depth][pos].Handler != nil {
				handlers = append(handlers, RouteMatchItem{
					Path:    t[depth][pos].Path,
					Handler: t[depth][pos].Handler,
				})
			}
		}
	}
	// fmt.Printf("sorted handlers: %v\n", handlers)
	return handlers
}

func (rt *Router) match(path string, justRouteLevel bool) (result matchResult) {
	rt.mutex.RLock()
	defer rt.mutex.RUnlock()

	result.Path = path

	// The root.
	if path == "" || path == "/" {
		result.Handler = rt.root.handler
		return
	}

	// Filter out some dirty paths.
	// Path with multiple trailing "/" is a kind of them.
	if strings.HasSuffix(path, "//") {
		return
	}

	var (
		matched, backtrack bool
		theway             string
		tbltrace           traceTable     = make(traceTable)
		p                  *route         = rt.root
		parts              []string       = strings.Split(path[1:], "/")
		rvs                RouteVariables = make(RouteVariables)
	)

	result.RouteVars = rvs

	for i := 0; i < len(parts); i++ {
		if i < 0 {
			p = nil
			goto LAB_MATCH_END
		}

		theway = "/" + strings.Join(parts[:i+1], "/")

		// Here p shouldn't be nil.
		if backtrack {
			p = p.rsibilng
		} else {
			p = p.child
		}

		matched, backtrack = false, false

		if p == nil {
			goto LAB_MATCH_END
		}

		// Compares with the routes on the same level, excluding the most right oue.
		for ; p != nil && p.rsibilng != nil; p = p.rsibilng {

			// Em... Originally, I want to use this judgement statement to accelerate the matching
			// procedure. But now I want to traverse the whole tree to find the apporopriate "fallback"
			// handlers for each sub path (i.e. "/a", "/a/" of path "/a/b"), this statement should be commented.
			// if !justRouteLevel && (len(parts)-p.depths[0]-i < -1 || len(parts)-p.depths[1]-i > 1) {
			//  continue
			// }

			if p.match(parts[i], justRouteLevel, rvs) {
				matched = true
				break
			}
		}

		// Compares with the most right node.
		if !matched {
			matched = p.match(parts[i], justRouteLevel, rvs)
		}

		if !matched {
			// Backtrack.
			backtrack = true
			for i--; p != nil && p.rsibilng == nil; p = p.parent {
				if p == rt.root {
					break
				}
				i--
			}
		} else {
			if strings.HasSuffix(theway, "/") {
				tbltrace.tryInsert(i-1, 1, theway, p.handler)
			} else {
				tbltrace.tryInsert(i, 0, theway, p.handler)
				if p.child != nil && (p.child.isSlashTail() || p.child.isAnyTail()) {
					tbltrace.tryInsert(i, 1, theway+"/", p.child.handler)
				}
			}

		}
	}

LAB_MATCH_END:
	if p != nil {
		result.Handler = p.handler
	}
	result.HandlersOnTheWay = tbltrace.getHandlersOnTheWay()
	return
}

func (rt *Router) rebuildRouteTree(routes []*route) {
	parent := rt.root
	for _, node := range routes {
		rt.addChildRouteNode(parent, node)
		if node.index == -1 {
			rt.rearrange(rt.root)
		} else {
			rt.rearrange(parent)
		}
		parent = node
	}
}

// Rearrange the routes. Sort by priority and other policies, like pattern length
// and pattern alphabetic order.
// ## Why I rearrange the routes?
// For instance, pattern "/cpp/19911110/{file}\\w+?\.cxx" and "/cpp/{date}\\d{8}/born.cxx"
// both can match path "/cpp/19911110/born.txt", but the former has a higher priority,
// so the handler bound to which will be returned.
func (rt *Router) rearrange(parent *route) {
	if parent == nil {
		return
	}
	children := make([]*route, 0)
	for p := parent.child; p != nil; p = p.rsibilng {
		children = append(children, p)
	}

	if len(children) == 0 {
		return
	}

	sort.Sort(byRoutePriority(children))

	// Link them.
	for i := 0; i < len(children)-1; i++ {
		children[i].rsibilng = children[i+1]
	}
	for i := len(children) - 1; i > 0; i-- {
		children[i].lsibling = children[i-1]
	}
	// Head and tail.
	children[0].lsibling = nil
	children[len(children)-1].rsibilng = nil
	parent.child = children[0]
}

type byRoutePriority []*route

func (b byRoutePriority) Len() int      { return len(b) }
func (b byRoutePriority) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byRoutePriority) Less(i, j int) bool {
	if b[i].priority > b[j].priority {
		return true
	}
	if b[i].priority == b[j].priority {
		if len(b[i].pattern) < len(b[j].pattern) {
			return true
		}
		if len(b[i].pattern) == len(b[j].pattern) {
			return b[i].pattern < b[j].pattern
		}
	}
	return false
}

func (rt *Router) addChildRouteNode(parent, node *route) {
	// The root.
	if node.index == -1 {
		if rt.root.child != nil {
			rt.root.child.parent = node
		}
		node.child = rt.root.child
		rt.root = node
		return
	}

	node.parent = parent
	sib := parent.child

	// `node` is the first child of `parent`.
	if sib == nil {
		parent.child = node
		return
	}

	// Iterate on the same level to find if there's an equivalent node.
	for ; sib != nil && sib.rsibilng != nil; sib = sib.rsibilng {
		if sib.equal(node) {
			break
		}
	}

	// Take the place of the old one.
	if sib.equal(node) {
		node.child = sib.child

		// Fix issue: missing handler.
		// eg.
		// mux.Handle("/admin", handler_1)
		// mux.Handle("/admin/tag", handler_2)
		//
		//  /admin --> got a nil handler (be overwritten)
		if sib.handler != nil {
			if node.handler == nil {
				node.handler = sib.handler
			} else {
				panic("something went wrong")
			}
		}

		node.rsibilng = sib.rsibilng
		node.lsibling = sib.lsibling

		if node.depths[0] > sib.depths[0] {
			node.depths[0] = sib.depths[0]
		}
		if node.depths[1] < sib.depths[1] {
			node.depths[1] = sib.depths[1]
		}

		if sib.lsibling != nil {
			sib.lsibling.rsibilng = node
		}

		if sib.rsibilng != nil {
			sib.rsibilng.lsibling = node
		}

		if sib.child != nil {
			sib.child.parent = node
		}

		if sib == parent.child {
			parent.child = node
		}
		return
	}

	// Append.
	sib.rsibilng = node
	node.lsibling = sib
}

func dumpRoute(buf *bytes.Buffer, r *route) {
	for p := r; p != nil; p = p.rsibilng {
		stack := make([]string, 0)
		for q := p.parent; q != nil; q = q.parent {

			if p.parent == q {
				if p.rsibilng == nil {
					stack = append(stack, "└── ")
				} else {
					stack = append(stack, "├── ")
				}
				continue
			}

			var z *route
			for z = p.parent; z != nil && z.parent != q; z = z.parent {
				// Ancestor bubble up.
			}
			if z != nil && z.rsibilng != nil {
				stack = append(stack, "│   ")
				continue
			}
			stack = append(stack, "    ")
		}

		// Pop stack.
		for i := len(stack) - 1; i >= 0; i-- {
			fmt.Fprintf(buf, "%s", stack[i])
		}
		fmt.Fprintf(buf, "%v\n", p)

		dumpRoute(buf, p.child)
	}

}

func (rt *Router) DumpTree() string {
	buf := bytes.NewBuffer(make([]byte, 0, 2048))
	dumpRoute(buf, rt.root)
	return buf.String()
}

func makeRoutes(pattern string, handler interface{}) ([]*route, error) {
	pattern = strings.TrimSpace(pattern)

	if pattern == "" {
		return nil, fmt.Errorf("empty pattern %q", pattern)
	}

	if pattern[0] != '/' {
		return nil, fmt.Errorf("pattern %q must start with %q", pattern, "/")
	}

	// The root.
	if pattern == "/" {
		return []*route{&route{index: -1, pattern: "", handler: handler}}, nil
	}

	parts := strings.Split(pattern[1:], "/")
	nameDuplicated := make(map[string]bool)
	routes := make([]*route, 0)
	depth := len(parts)

	for i, part := range parts {
		r := &route{index: i, depths: [2]int{depth - i, depth - i}, handler: nil}

		// Any part shouldn't be empty excluding the last one.
		if part == "" && i != (depth-1) {
			return nil, fmt.Errorf("empty route (or duplicated '/') in pattern %q", pattern)
		}

		var (
			isRegex      bool
			originalPart = part
		)

		r.name, part, isRegex = splitRouteNameAndMatchPattern(part)
		// Check duplicated name in the same pattern.
		if r.name != "" && nameDuplicated[r.name] {
			return nil, fmt.Errorf("duplicated route name %q in pattern %q",
				r.name, pattern)
		}
		nameDuplicated[r.name] = true

		if isRegex {
			if part == "^$" {
				return nil, fmt.Errorf("empty route (regex) %q in pattern %q", part, pattern)
			}
			if regex, err := regexp.Compile(part); err != nil {
				return nil, fmt.Errorf("unable to compile route %q in pattern %q, compile error: %s",
					originalPart, pattern, err.Error())
			} else {
				r.regex = regex
			}
		}

		r.pattern = part

		// Decide the priority.
		if r.pattern == "" && r.name != "" {
			r.priority = kAnyPattern
		} else if r.regex != nil {
			r.priority = kRegexPattern
		} else {
			r.priority = kAbsolutePattern
		}

		// Bind the handler to the tail.
		if i == depth-1 {
			r.handler = handler
		}
		routes = append(routes, r)
	}

	return routes, nil
}

// part: "<name>|{name}tomatch"
// `name` is enclosed in `<>` or `{}` as prefix in `part`, it might be empty("").
// `tomatch` is the remained content after name prefix, it also might be empty("").
// `isRegex` is true only when `tomatch` is not empty and `name` is enclosed by `{}`.
func splitRouteNameAndMatchPattern(part string) (name string, tomatch string, isRegex bool) {
	if part == "" {
		return
	}
	var closeBrace string
	if part[0] == '<' || part[0] == '{' {
		if part[0] == '<' {
			closeBrace = ">"
		} else {
			closeBrace = "}"
		}
		indexOfCloseBrace := strings.Index(part, closeBrace)
		if indexOfCloseBrace == -1 {
			closeBrace = ""
		} else {
			name = part[1:indexOfCloseBrace]
			part = part[indexOfCloseBrace+1:]
		}
	}

	for closeBrace == "}" {
		if part == "" {
			break
		}

		isRegex = true

		// Enclose the regular expression with ^ and $.
		if part[0] != '^' {
			part = "^" + part
		}
		if part[len(part)-1] != '$' {
			part = part + "$"
		}

		break
	}

	tomatch = part
	return
}

func isNil(i interface{}) bool {
	defer func() { recover() }()
	return i == nil || reflect.ValueOf(i).IsNil()
}
