package postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"

	"github.com/aditwar-man/go-microservice-boilerplate/config"
)

const (
	maxOpenConns    = 60
	connMaxLifetime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

// Return new Postgresql db instance
func NewPsqlDB(c *config.Config) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s search_path=%s",
		c.Postgres.PostgresqlHost,
		c.Postgres.PostgresqlPort,
		c.Postgres.PostgresqlUser,
		c.Postgres.PostgresqlDbname,
		c.Postgres.PostgresqlPassword,
		c.Postgres.DefaultSchema,
	)

	db, err := sqlx.Connect(c.Postgres.PgDriver, dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func SeedRBACData(db *sql.DB) error {
	// Check if data already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM roles").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Data already seeded
	}

	// Seed default roles with hierarchy
	roles := []struct {
		name, description string
		parentName        *string
	}{
		{"administrator", "System Administrator with full access", nil},
		{"manager", "Manager with team oversight", stringPtr("administrator")},
		{"employee", "Regular employee user", stringPtr("manager")},
		{"viewer", "Read-only access user", stringPtr("employee")},
		{"guest", "Limited guest access", nil},
	}

	roleIDs := make(map[string]int)

	// Insert roles in order (parents first)
	for _, role := range roles {
		var parentID *int
		if role.parentName != nil {
			if pid, exists := roleIDs[*role.parentName]; exists {
				parentID = &pid
			}
		}

		var roleID int
		err := db.QueryRow("INSERT INTO roles (name, description, parent_role_id) VALUES ($1, $2, $3) RETURNING id",
			role.name, role.description, parentID).Scan(&roleID)
		if err != nil {
			return err
		}
		roleIDs[role.name] = roleID
	}

	// Seed default permissions
	permissions := []struct {
		name, description string
	}{
		{"create", "Create new records"},
		{"read", "Read/view records"},
		{"update", "Update existing records"},
		{"delete", "Delete records"},
		{"manage", "Full management access"},
		{"approve", "Approve actions/requests"},
		{"audit", "Access audit logs"},
		{"configure", "System configuration"},
	}

	permissionIDs := make(map[string]int)
	for _, perm := range permissions {
		var permID int
		err := db.QueryRow("INSERT INTO permissions (name, description) VALUES ($1, $2) RETURNING id",
			perm.name, perm.description).Scan(&permID)
		if err != nil {
			return err
		}
		permissionIDs[perm.name] = permID
	}

	// Seed default resources
	resources := []struct {
		name, description string
	}{
		{"users", "User management"},
		{"roles", "Role management"},
		{"permissions", "Permission management"},
		{"resources", "Resource management"},
		{"reports", "System reports"},
		{"settings", "System settings"},
		{"audit_logs", "Audit log access"},
		{"dashboard", "Dashboard access"},
	}

	resourceIDs := make(map[string]int)
	for _, res := range resources {
		var resID int
		err := db.QueryRow("INSERT INTO resources (name, description) VALUES ($1, $2) RETURNING id",
			res.name, res.description).Scan(&resID)
		if err != nil {
			return err
		}
		resourceIDs[res.name] = resID
	}

	// Seed default contexts
	contexts := []struct {
		name, description string
	}{
		{"global", "Global system context"},
		{"department", "Department-specific context"},
		{"project", "Project-specific context"},
		{"personal", "Personal user context"},
	}

	contextIDs := make(map[string]int)
	for _, ctx := range contexts {
		var ctxID int
		err := db.QueryRow("INSERT INTO context (name, description) VALUES ($1, $2) RETURNING id",
			ctx.name, ctx.description).Scan(&ctxID)
		if err != nil {
			return err
		}
		contextIDs[ctx.name] = ctxID
	}

	// Assign permissions to roles
	rolePermissions := []struct {
		roleName, permissionName, resourceName, contextName string
	}{
		// Administrator - full access
		{"administrator", "manage", "users", "global"},
		{"administrator", "manage", "roles", "global"},
		{"administrator", "manage", "permissions", "global"},
		{"administrator", "manage", "resources", "global"},
		{"administrator", "manage", "settings", "global"},
		{"administrator", "read", "audit_logs", "global"},
		{"administrator", "read", "reports", "global"},
		{"administrator", "read", "dashboard", "global"},

		// Manager - team management
		{"manager", "create", "users", "department"},
		{"manager", "read", "users", "department"},
		{"manager", "update", "users", "department"},
		{"manager", "read", "roles", "department"},
		{"manager", "approve", "reports", "department"},
		{"manager", "read", "dashboard", "department"},

		// Employee - basic access
		{"employee", "read", "users", "personal"},
		{"employee", "update", "users", "personal"},
		{"employee", "read", "dashboard", "personal"},
		{"employee", "create", "reports", "personal"},
		{"employee", "read", "reports", "personal"},

		// Viewer - read only
		{"viewer", "read", "dashboard", "personal"},
		{"viewer", "read", "reports", "personal"},

		// Guest - minimal access
		{"guest", "read", "dashboard", "global"},
	}

	for _, rp := range rolePermissions {
		roleID := roleIDs[rp.roleName]
		permissionID := permissionIDs[rp.permissionName]
		resourceID := resourceIDs[rp.resourceName]
		contextID := contextIDs[rp.contextName]

		_, err = db.Exec("INSERT INTO role_permissions (role_id, permission_id, resource_id, context_id) VALUES ($1, $2, $3, $4)",
			roleID, permissionID, resourceID, contextID)
		if err != nil {
			return err
		}
	}

	return nil
}

func stringPtr(s string) *string {
	return &s
}
