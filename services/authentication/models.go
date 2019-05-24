package authentication

type TokenIntrospectionResponse struct {
	Active          bool   `json:"active"`
	IssuedAt        int64  `json:"iat,omitempty"`
	ExpireAt        int64  `json:"exp,omitempty"`
	Subject         string `json:"sub,omitempty"` // user LDAP
	UserDomainRoles `json:""`
}

type UserDomainRoles struct {
	DomainRolesBinding  []*DomainRoles `json:"domain_roles_binding,omitempty"`
	DomainRolesComputed []*DomainRoles `json:"domain_roles_computed,omitempty"`
}

type DomainRoles struct {
	Domain            string             `json:"domain"`
	DomainDescription *string            `json:"domain_description"`
	Roles             []*Role            `json:"roles"`
	AggregatedRoles   []*ResourceBinding `json:"aggregated_roles"`
}

type ResourceBinding struct {
	ResourceName string   `json:"resource_name"`
	Verbs        []string `json:"verbs"`
}

type Role struct {
	Name             string             `json:"name"`
	ResourcesBinding []*ResourceBinding `json:"resources_bindings"`
}
