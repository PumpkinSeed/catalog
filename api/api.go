package catalog

import "github.com/hashicorp/consul/api"

type Catalog interface {
	Register(req *api.CatalogRegistration, q *api.WriteOptions) (*api.WriteMeta, error)
	Deregister(dereg *api.CatalogDeregistration, q *api.WriteOptions) (*api.WriteMeta, error)
	Datacenters() ([]string, error)
	Nodes(q *api.QueryOptions) ([]*api.Node, *api.QueryMeta, error)
	Node(node string, q *api.QueryOptions) (*api.CatalogNode, *api.QueryMeta, error)
	Services(q *api.QueryOptions) (map[string][]string, *api.QueryMeta, error)
	Service(service, tag string, q *api.QueryOptions) ([]*api.CatalogService, *api.QueryMeta, error)
}

type catalog struct {
	addr string
}

func NewCatalog(addr string) Catalog {
	return &catalog{addr: addr}
}

func (c *catalog) Register(req *api.CatalogRegistration, q *api.WriteOptions) (*api.WriteMeta, error) {
	return nil, nil
}

func (c *catalog) Deregister(dereg *api.CatalogDeregistration, q *api.WriteOptions) (*api.WriteMeta, error) {
	return nil, nil
}

func (c *catalog) Datacenters() ([]string, error) {
	return nil, nil
}

func (c *catalog) Nodes(q *api.QueryOptions) ([]*api.Node, *api.QueryMeta, error) {
	return nil, nil, nil
}

func (c *catalog) Node(node string, q *api.QueryOptions) (*api.CatalogNode, *api.QueryMeta, error) {
	return nil, nil, nil
}

func (c *catalog) Services(q *api.QueryOptions) (map[string][]string, *api.QueryMeta, error) {
	return nil, nil, nil
}

func (c *catalog) Service(service, tag string, q *api.QueryOptions) ([]*api.CatalogService, *api.QueryMeta, error) {
	return nil, nil, nil
}
