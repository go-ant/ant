// https://github.com/mikespook/gorbac
package rbac

import (
	"sync"
)

const (
	bufferSize = 16
)

// Assertion function supplies more fine-grained permission controls.
type AssertionFunc func(string, string, *Rbac) bool

// Export RBAC to a structure data
type Map map[string]RoleMap

// Return a RBAC structure.
func New() *Rbac {
	rbac := &Rbac{
		roles:   make(map[string]Role, bufferSize),
		factory: NewBaseRole,
	}
	return rbac
}

// RBAC
type Rbac struct {
	mutex   sync.RWMutex
	roles   map[string]Role
	factory RoleFactoryFunc
}

// Internal getRole
func (rbac *Rbac) getRole(name string) Role {
	role, ok := rbac.roles[name]
	if !ok {
		role = rbac.factory(rbac, name)
		rbac.roles[name] = role
	}
	return role
}

// Restore rbac from a map
func Restore(data Map) *Rbac {
	rbac := New()
	for role, value := range data {
		rbac.Add(role, value[PermissionKey], value[ParentKey])
	}
	return rbac
}

// Set a role with `name`. It has `permissions` and `parents`.
// If the role is not existing, a new one will be created.
// This function will cover role's orignal permissions and parents.
func (rbac *Rbac) Set(name string, permissions []string, parents []string) {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()
	role := rbac.getRole(name)
	role.Reset()
	for _, p := range permissions {
		role.AddPermission(p)
	}
	for _, pname := range parents {
		role.AddParent(pname)
	}
}

// Add a role with `name`. It has `permissions` and `parents`.
// If the role is not existing, a new one will be created.
// This function will add new permissions and parents to the role,
// and keep orignals.
func (rbac *Rbac) Add(name string, permissions []string, parents []string) {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()
	role := rbac.getRole(name)
	for _, p := range permissions {
		role.AddPermission(p)
	}
	for _, pname := range parents {
		role.AddParent(pname)
	}
}

// Remove a role.
func (rbac *Rbac) Remove(name string) {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()
	delete(rbac.roles, name)
}

// Return a role or nil if not exists.
func (rbac *Rbac) Get(name string) Role {
	rbac.mutex.RLock()
	defer rbac.mutex.RUnlock()
	role, ok := rbac.roles[name]
	if !ok {
		return nil
	}
	return role
}

// Test if the `name` has `permission` in the `assert` condition.
func (rbac *Rbac) IsGranted(name, permission string, assert AssertionFunc) bool {
	rbac.mutex.RLock()
	defer rbac.mutex.RUnlock()
	if assert != nil && !assert(name, permission, rbac) {
		return false
	}
	if role, ok := rbac.roles[name]; ok {
		return role.HasPermission(permission)
	}
	return false
}

// Dump RBAC
func (rbac *Rbac) Dump() Map {
	m := make(Map)
	for _, role := range rbac.roles {
		roleMap := RoleToMap(role)
		m[role.Name()] = roleMap
	}
	return m
}
