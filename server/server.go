package server

import (
	"context"
	"sort"
	"sync"

	"github.com/yun-jay/clerk-echo-oapi-middleware/api"
)

type server struct {
	sync.RWMutex
	lastID int64
	things map[int64]api.Thing
}

func NewServer() *server {
	return &server{
		lastID: 0,
		things: make(map[int64]api.Thing),
	}
}

// Ensure that we implement the server interface
var _ api.StrictServerInterface = (*server)(nil)

func (s *server) GetVersion(ctx context.Context, req api.GetVersionRequestObject) (api.GetVersionResponseObject, error) {

	swagger, err := api.GetSwagger()
	if err != nil {
		return nil, err
	}

	return api.GetVersion200JSONResponse{
		Message: &swagger.Info.Version,
	}, nil
}

func (s *server) ListThings(ctx context.Context, req api.ListThingsRequestObject) (api.ListThingsResponseObject, error) {
	// This handler will only be called when a valid JWT is presented for
	// access.
	s.RLock()

	thingKeys := make([]int64, 0, len(s.things))
	for key := range s.things {
		thingKeys = append(thingKeys, key)
	}
	sort.Sort(int64s(thingKeys))

	things := make([]api.ThingWithID, 0, len(s.things))

	for _, key := range thingKeys {
		thing := s.things[key]
		things = append(things, api.ThingWithID{
			Id:   key,
			Name: thing.Name,
		})
	}

	s.RUnlock()

	return api.ListThings200JSONResponse(things), nil
}

type int64s []int64

func (in int64s) Len() int {
	return len(in)
}

func (in int64s) Less(i, j int) bool {
	return in[i] < in[j]
}

func (in int64s) Swap(i, j int) {
	in[i], in[j] = in[j], in[i]
}

var _ sort.Interface = (int64s)(nil)

func (s *server) AddThing(ctx context.Context, req api.AddThingRequestObject) (api.AddThingResponseObject, error) {
	// This handler will only be called when the JWT is valid and the JWT contains
	// the scopes required.

	s.Lock()
	defer s.Unlock()

	thing := req.Body
	s.things[s.lastID] = *thing
	thingWithId := api.ThingWithID{
		Name: thing.Name,
		Id:   s.lastID,
	}
	s.lastID++

	return api.AddThing201JSONResponse(thingWithId), nil
}
