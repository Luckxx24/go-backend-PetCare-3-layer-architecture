🐾 PetCare Backend API
A robust backend system for pet boarding management. This project is built with Go and focuses on high-performance database interactions and secure access control using 3-Layer Architecture.

🏗️ Architecture: 3-Layer Pattern
The project follows a clean separation of concerns, making the codebase testable and scalable:

Delivery Layer (Handlers): Manages HTTP requests/responses, JSON parsing, and interacts with the middleware.

Service Layer (Business Logic): The core engine where roles are validated, booking rules are enforced, and automated notifications are triggered.

Repository Layer (Data Access): Powered by SQLC. It executes raw SQL queries with type-safety, ensuring maximum performance.

🛠️ Tech Stack
Language: Go (Golang)

Database: PostgreSQL

Query Builder: SQLC (Type-safe SQL)

Migration: Goose (Database version control)

Authentication: JWT (JSON Web Token)

Security: Role-Based Access Control (RBAC)

🔐 Authentication & Security
This project implements a secure authentication flow using JWT:

Token Generation: Upon login, the server generates a signed JWT containing the user_id and role.

Middleware Protection: Private routes are guarded by a middleware that:

Extracts the token from the Authorization: Bearer <token> header.

Validates the signature and expiration.

Injects the user_id and role into the Context, allowing the Service Layer to perform authorization checks.

RBAC (Role-Based Access Control): * Users: Can manage their own pets and bookings.

Staff: Can create Pet Status Logs and respond to chats.

Admin: Full access to system management.

🚀 Key Features
1. Secure Booking Management
End-to-end boarding process with ownership validation. Users can only view/chat about bookings they own.

2. Context-Aware Messaging
A specialized chat system grouped by booking_id.

Inbox: Aggregates the latest messages for each active transaction.

History: Paginated message retrieval to ensure high performance even with large chat logs.

3. Pet Status Logging
Staff can issue daily status reports (Food, Health, Grooming) including photos and notes to keep pet owners updated.

4. Automated Notifications
A system-driven notification engine that alerts users on new messages or pet status updates.
