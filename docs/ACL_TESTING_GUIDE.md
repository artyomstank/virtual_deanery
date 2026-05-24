# ACL System Testing Guide

## Setup Prerequisites

1. Ensure database is running and migrations are applied
2. Start the application: `make run` or `./cmd/api/main.go`
3. Default port: 8080

## Test Scenarios

### 1. User Registration and Role Assignment

#### Register Student (Default Role)
```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "student1",
    "email": "student1@school.com",
    "password": "password123"
  }'
```

Expected: User created with "student" role

#### Register Teacher
```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "teacher1",
    "email": "teacher1@school.com",
    "password": "password123",
    "role": "teacher"
  }'
```

Expected: User created with "teacher" role

#### Register Admin
```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin1",
    "email": "admin1@school.com",
    "password": "password123",
    "role": "admin"
  }'
```

Expected: User created with "admin" role

---

### 2. Authentication and Token Generation

#### Login as Student
```bash
STUDENT_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student1@school.com",
    "password": "password123"
  }' | jq -r '.access_token')

echo $STUDENT_TOKEN
```

#### Login as Admin
```bash
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin1@school.com",
    "password": "password123"
  }' | jq -r '.access_token')

echo $ADMIN_TOKEN
```

---

### 3. ACL Permission Tests

#### Student Accessing Grades (Should Succeed - Read Permission)
```bash
curl -X GET http://localhost:8080/api/v1/grades \
  -H "Authorization: Bearer $STUDENT_TOKEN"
```

Expected: HTTP 200 (if route is protected with read permission)

#### Student Trying to Modify Grades (Should Fail - No Write Permission)
```bash
curl -X PATCH http://localhost:8080/api/v1/grades/1 \
  -H "Authorization: Bearer $STUDENT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "score": 95
  }'
```

Expected: HTTP 403 Forbidden

#### Teacher Modifying Grades (Should Succeed - Write Permission)
```bash
TEACHER_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "teacher1@school.com",
    "password": "password123"
  }' | jq -r '.access_token')

curl -X PATCH http://localhost:8080/api/v1/grades/1 \
  -H "Authorization: Bearer $TEACHER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "score": 95
  }'
```

Expected: HTTP 200 OK

---

### 4. Admin Operations

#### Get All Roles
```bash
curl -X GET http://localhost:8080/api/v1/admin/roles \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

Expected: HTTP 200 with list of roles

#### Get ACL for Student Role
```bash
curl -X GET http://localhost:8080/api/v1/admin/acl/student \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

Expected: HTTP 200 with ACL entries for student role

#### Update ACL Entry - Grant Write Access to Students Resource
```bash
curl -X PATCH http://localhost:8080/api/v1/admin/acl \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role_id": 3,
    "resource_id": 1,
    "can_read": true,
    "can_write": true,
    "can_delete": false
  }'
```

Expected: HTTP 200 - ACL entry updated

---

### 5. Permission Denial Tests

#### Non-Admin Trying to Update ACL (Should Fail)
```bash
curl -X PATCH http://localhost:8080/api/v1/admin/acl \
  -H "Authorization: Bearer $STUDENT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role_id": 3,
    "resource_id": 1,
    "can_read": true,
    "can_write": true,
    "can_delete": false
  }'
```

Expected: HTTP 403 Forbidden

#### Accessing Protected Route Without Token
```bash
curl -X GET http://localhost:8080/api/v1/admin/roles
```

Expected: HTTP 401 Unauthorized

#### Accessing Protected Route with Invalid Token
```bash
curl -X GET http://localhost:8080/api/v1/admin/roles \
  -H "Authorization: Bearer invalid_token_here"
```

Expected: HTTP 401 Unauthorized

---

### 6. Cache Validation

#### Check Cache Update After ACL Change
1. Get ACL entries for a role (first call - from DB, cached)
2. Update an ACL entry
3. Get ACL entries again (should reflect changes)

```bash
# First call - loads from DB and caches
curl -s -X GET http://localhost:8080/api/v1/admin/acl/student \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq

# Update ACL
curl -X PATCH http://localhost:8080/api/v1/admin/acl \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role_id": 3,
    "resource_id": 1,
    "can_read": true,
    "can_write": true,
    "can_delete": true
  }'

# Second call - should show updated values
curl -s -X GET http://localhost:8080/api/v1/admin/acl/student \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq
```

Expected: can_delete changes from false to true

---

## Expected HTTP Status Codes

| Status | Meaning | Example |
|--------|---------|---------|
| 200 | Success | Login successful, resource accessed |
| 201 | Created | User registered successfully |
| 400 | Bad Request | Invalid input format |
| 401 | Unauthorized | Missing or invalid token |
| 403 | Forbidden | User lacks permission |
| 404 | Not Found | User/resource not found |
| 500 | Internal Error | Database or server error |

---

## Database Queries for Verification

### Check Roles
```sql
SELECT * FROM roles;
```

### Check User Roles
```sql
SELECT u.id, u.username, r.name FROM users u
JOIN user_roles ur ON u.id = ur.user_id
JOIN roles r ON ur.role_id = r.id;
```

### Check ACL Entries
```sql
SELECT 
  r.name as role, 
  res.name as resource,
  a.can_read, a.can_write, a.can_delete
FROM acl_entries a
JOIN roles r ON a.role_id = r.id
JOIN resources res ON a.resource_id = res.id
ORDER BY r.name, res.name;
```

### Check Specific User Permissions
```sql
SELECT 
  u.username,
  r.name as role,
  res.name as resource,
  a.can_read, a.can_write, a.can_delete
FROM users u
JOIN user_roles ur ON u.id = ur.user_id
JOIN roles r ON ur.role_id = r.id
JOIN acl_entries a ON r.id = a.role_id
JOIN resources res ON a.resource_id = res.id
WHERE u.username = 'student1'
ORDER BY res.name;
```

---

## Troubleshooting

### Issue: "role not found" error on registration
- Check that role exists in database: `SELECT * FROM roles;`
- Verify migration was applied successfully

### Issue: Permission denied when accessing admin routes
- Verify user is admin: Check `user_roles` table
- Check that user has "admin" role

### Issue: ACL changes not reflected immediately
- Cache is invalidated when ACL is updated
- If using multiple server instances, cache won't sync between them (implement Redis for distributed cache)

### Issue: Token validation errors
- Check JWT secret is consistent across application restarts
- Verify token hasn't expired (default: 24 hours)
