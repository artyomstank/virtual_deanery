# ACL (Access Control List) Implementation

## Overview

The ACL system provides role-based access control (RBAC) to the Virtual Deanery application. It controls which users can perform specific actions (read, write, delete) on protected resources.

## Architecture

### Core Components

1. **Entity Layer** (`internal/domain/entity/`)
   - `Role`: Represents user roles (admin, teacher, student, dean)
   - `Resource`: Represents protected resources (students, grades, schedule, etc.)
   - `ACLEntry`: Links roles to resources with permission flags (canRead, canWrite, canDelete)
   - `Action`: Enum for action types (read, write, delete)
   - `User.HasPermission()`: Method to check if user has permission

2. **Repository Layer** (`internal/domain/repository/`)
   - `RoleRepository`: Interface for role operations
   - `ACLRepository`: Interface for ACL operations
   - `UserRepository`: Interface for user operations (already exists)

3. **Service Layer**
   - `UserService` (`internal/domain/service/user_service.go`): Handles user registration, login, and permission checks
   - `ACLService` (`internal/domain/service/acl_service.go`): Interface for ACL management
   - Implementation (`internal/service/acl_service.go`): Concrete implementation with admin checks

4. **PostgreSQL Repositories** (`internal/repo/postgres/`)
   - `role_repo.go`: Implements RoleRepository
   - `acl_repo.go`: Implements ACLRepository
   - `user_repo.go`: User repository (existing)

5. **HTTP Layer**
   - **Middleware** (`internal/transport/http/middleware/acl.go`): ACL middleware for route protection
   - **DTOs** (`internal/transport/http/dto/acl_dto.go`): Request/response objects
   - **Handler** (`internal/transport/http/handler/acl_handler.go`): HTTP handlers for ACL management
   - **Router** (`internal/transport/http/router/router.go`): Route definitions with ACL middleware

## Database Schema

### Key Tables

1. **roles**
   - id (serial)
   - name (varchar, unique)
   - description (varchar)
   - created_at (timestamp)

2. **resources**
   - id (serial)
   - name (varchar, unique)
   - description (varchar)
   - created_at (timestamp)

3. **acl_entries**
   - id (serial)
   - role_id (FK to roles)
   - resource_id (FK to resources)
   - can_read (boolean)
   - can_write (boolean)
   - can_delete (boolean)
   - created_at (timestamp)
   - updated_at (timestamp)
   - unique constraint on (role_id, resource_id)

4. **user_roles**
   - user_id (FK to users, primary key)
   - role_id (FK to roles)
   - assigned_at (timestamp)

## API Endpoints

### Registration (Public)
```
POST /api/v1/users/register
Request:
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "securepassword123",
  "role": "student"  // optional, defaults to "student"
}
Response:
{
  "id": 1,
  "username": "john_doe",
  "email": "john@example.com",
  "role": "student",
  "is_active": true,
  "created_at": "2025-05-23T...",
  "updated_at": "2025-05-23T..."
}
```

### Login (Public)
```
POST /api/v1/users/login
Request:
{
  "email": "john@example.com",
  "password": "securepassword123"
}
Response:
{
  "access_token": "eyJhbGc...",
  "expires_at": "2025-05-24T...",
  "user": { /* UserResponse */ }
}
```

### Get User Profile (Protected)
```
GET /api/v1/users/:id
Headers: Authorization: Bearer {token}
Response: { /* UserResponse */ }
```

### Get All Roles (Admin Only)
```
GET /api/v1/admin/roles
Headers: Authorization: Bearer {token}
Response:
[
  { "id": 1, "name": "admin", "description": "Administrator..." },
  { "id": 2, "name": "teacher", "description": "..." },
  { "id": 3, "name": "student", "description": "..." },
  { "id": 4, "name": "dean", "description": "..." }
]
```

### Get ACL by Role (Admin Only)
```
GET /api/v1/admin/acl/:role
Headers: Authorization: Bearer {token}
Response:
[
  {
    "role_id": 3,
    "resource_id": 1,
    "role_name": "student",
    "resource_name": "students",
    "can_read": true,
    "can_write": false,
    "can_delete": false
  },
  ...
]
```

### Update ACL Entry (Admin Only)
```
PATCH /api/v1/admin/acl
Headers: Authorization: Bearer {token}
Request:
{
  "role_id": 3,
  "resource_id": 1,
  "can_read": true,
  "can_write": true,
  "can_delete": false
}
Response:
{
  "code": 200,
  "message": "ACL entry updated successfully"
}
```

## Default Roles and Permissions

### Admin
- **All Resources**: READ, WRITE, DELETE

### Teacher
- **grades**: READ, WRITE
- **schedule**: READ
- **students**: READ
- **teachers**: READ
- **profile**: READ, WRITE
- **reports**: READ

### Student
- **profile**: READ, WRITE
- **grades**: READ
- **schedule**: READ

### Dean
- Configurable through admin API

## Permission Checking Flow

1. User makes authenticated request (JWT token in Authorization header)
2. `AuthMiddleware` validates token and extracts `user_id` and `role`
3. `ACLMiddleware` (if configured on route) checks if user has permission for resource+action
4. `UserService.CheckPermission()` method:
   - Loads ACL entries for user's role from cache
   - Calls `User.HasPermission()` to check permission
   - Returns error if permission denied
5. Request proceeds to handler if permission granted

## Caching

ACL entries are cached in-memory per role:
- First access to a role's ACL entries triggers a database query
- Subsequent accesses use cached results
- Cache is invalidated when ACL entries are updated
- Cache is per-role for efficient lookups

## Usage Example

### 1. Register a User
```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@school.com",
    "password": "password123"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@school.com",
    "password": "password123"
  }'
```

### 3. Admin: Update ACL
```bash
curl -X PATCH http://localhost:8080/api/v1/admin/acl \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {admin_token}" \
  -d '{
    "role_id": 3,
    "resource_id": 2,
    "can_read": true,
    "can_write": true,
    "can_delete": false
  }'
```

## Integration with Routes

To protect a route with ACL checks:

```go
// In router.go
protected.Use(middleware.ACLMiddleware(
  userService,
  "students",           // resource name
  entity.ActionRead,    // action type
  logger,
))
protected.GET("/students", handler.GetStudents)
```

## Error Responses

### 401 Unauthorized
- Missing JWT token
- Invalid or expired token

### 403 Forbidden
- User lacks required permission for resource+action

### 400 Bad Request
- Invalid request format or validation error

### 500 Internal Server Error
- Database errors or unexpected failures

## Security Considerations

1. **Password Security**: Passwords are bcrypt-hashed with cost >= 12
2. **Token Security**: JWT tokens are signed with a secret key
3. **Permission Checks**: All admin operations verify the user is admin before proceeding
4. **Caching**: Cached ACL entries are invalidated on updates
5. **Error Messages**: Generic error messages prevent information leakage (e.g., "invalid credentials" instead of "user not found" or "wrong password")

## Future Enhancements

1. Add per-resource instance checks (e.g., can user edit only their own profile?)
2. Support attribute-based access control (ABAC)
3. Add audit logging for permission checks
4. Implement role hierarchies
5. Add temporary permission grants
6. Support resource ownership checks
