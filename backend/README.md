# Transaction Matching Application

## User Management and Admin Roles

### Creating the First Admin User

To create the first admin user, you can use the included script:

```bash
cd backend
go run scripts/create_admin.go
```

This will create a default admin user with the following credentials:
- Email: admin@example.com
- Password: admin123

**Important**: Remember to change the admin password after first login!

### Managing User Roles via API

Once you have an admin user, you can manage other users and their roles through the API.

#### Creating a New User with Role

```bash
curl -X POST \
  http://localhost:8080/api/v1/users \
  -H 'Authorization: Bearer YOUR_JWT_TOKEN' \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "New User",
    "role": "preparer"
  }'
```

Valid roles are:
- `admin`: Full system access including user management
- `preparer`: Can upload and match transactions
- `approver`: Can review and approve matched transactions

#### Getting User Roles

```bash
curl -X GET \
  http://localhost:8080/api/v1/users/{user_id}/roles \
  -H 'Authorization: Bearer YOUR_JWT_TOKEN'
```

#### Updating User Roles

```bash
curl -X PUT \
  http://localhost:8080/api/v1/users/{user_id}/roles \
  -H 'Authorization: Bearer YOUR_JWT_TOKEN' \
  -H 'Content-Type: application/json' \
  -d '{
    "role": "admin",
    "operation": "add"
  }'
```

The operation can be either `add` or `remove`.

#### Granting Admin Role

There's a dedicated endpoint for granting admin privileges:

```bash
curl -X PUT \
  http://localhost:8080/api/v1/users/{user_id}/admin \
  -H 'Authorization: Bearer YOUR_JWT_TOKEN'
```

### User Management in the Frontend

The frontend application has an Admin section in the user profile page where users with admin privileges can:

1. View all users in the system
2. Assign or remove roles
3. Create new users with specific roles
4. Search and filter users

*Note: Only users with the admin role can access these features.* 