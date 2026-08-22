package httpx

// Route registers a typed handler on the server.
func (s *Server) Route[I, O any](method, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return Route(s, method, path, handler, operationOptions...)
}

// Get registers a typed GET handler on the server.
func (s *Server) Get[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return Get(s, path, handler, operationOptions...)
}

// Post registers a typed POST handler on the server.
func (s *Server) Post[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return Post(s, path, handler, operationOptions...)
}

// Put registers a typed PUT handler on the server.
func (s *Server) Put[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return Put(s, path, handler, operationOptions...)
}

// Patch registers a typed PATCH handler on the server.
func (s *Server) Patch[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return Patch(s, path, handler, operationOptions...)
}

// Delete registers a typed DELETE handler on the server.
func (s *Server) Delete[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return Delete(s, path, handler, operationOptions...)
}

// Head registers a typed HEAD handler on the server.
func (s *Server) Head[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return Head(s, path, handler, operationOptions...)
}

// Options registers a typed OPTIONS handler on the server.
func (s *Server) Options[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return Options(s, path, handler, operationOptions...)
}

// MustRoute registers a typed handler and panics if registration fails.
func (s *Server) MustRoute[I, O any](method, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	MustRoute(s, method, path, handler, operationOptions...)
}

// MustGet registers a GET route and panics if registration fails.
func (s *Server) MustGet[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	MustGet(s, path, handler, operationOptions...)
}

// MustPost registers a POST route and panics if registration fails.
func (s *Server) MustPost[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	MustPost(s, path, handler, operationOptions...)
}

// MustPut registers a PUT route and panics if registration fails.
func (s *Server) MustPut[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	MustPut(s, path, handler, operationOptions...)
}

// MustPatch registers a PATCH route and panics if registration fails.
func (s *Server) MustPatch[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	MustPatch(s, path, handler, operationOptions...)
}

// MustDelete registers a DELETE route and panics if registration fails.
func (s *Server) MustDelete[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	MustDelete(s, path, handler, operationOptions...)
}

// MustHead registers a HEAD route and panics if registration fails.
func (s *Server) MustHead[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	MustHead(s, path, handler, operationOptions...)
}

// MustOptions registers an OPTIONS route and panics if registration fails.
func (s *Server) MustOptions[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	MustOptions(s, path, handler, operationOptions...)
}

// Route registers a typed handler under the group.
func (g *Group) Route[I, O any](method, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return GroupRoute(g, method, path, handler, operationOptions...)
}

// Get registers a typed GET handler under the group.
func (g *Group) Get[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return GroupGet(g, path, handler, operationOptions...)
}

// Post registers a typed POST handler under the group.
func (g *Group) Post[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return GroupPost(g, path, handler, operationOptions...)
}

// Put registers a typed PUT handler under the group.
func (g *Group) Put[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return GroupPut(g, path, handler, operationOptions...)
}

// Patch registers a typed PATCH handler under the group.
func (g *Group) Patch[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return GroupPatch(g, path, handler, operationOptions...)
}

// Delete registers a typed DELETE handler under the group.
func (g *Group) Delete[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) error {
	return GroupDelete(g, path, handler, operationOptions...)
}

// MustRoute registers a typed handler under the group and panics if registration fails.
func (g *Group) MustRoute[I, O any](method, path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	MustGroupRoute(g, method, path, handler, operationOptions...)
}

// MustGet registers a GET route under the group and panics if registration fails.
func (g *Group) MustGet[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	MustGroupGet(g, path, handler, operationOptions...)
}

// MustPost registers a POST route under the group and panics if registration fails.
func (g *Group) MustPost[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	MustGroupPost(g, path, handler, operationOptions...)
}

// MustPut registers a PUT route under the group and panics if registration fails.
func (g *Group) MustPut[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	MustGroupPut(g, path, handler, operationOptions...)
}

// MustPatch registers a PATCH route under the group and panics if registration fails.
func (g *Group) MustPatch[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	MustGroupPatch(g, path, handler, operationOptions...)
}

// MustDelete registers a DELETE route under the group and panics if registration fails.
func (g *Group) MustDelete[I, O any](path string, handler TypedHandler[I, O], operationOptions ...OperationOption) {
	MustGroupDelete(g, path, handler, operationOptions...)
}
