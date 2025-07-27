# 🛒 Tokoaku Backend

Tokoaku Backend is the REST API service powering the Tokoaku e-commerce platform. It is built with [Go Fiber](https://gofiber.io/), uses PostgreSQL for the main database, Redis for caching, Firebase Auth for authentication, and integrates with external services like Cloudinary, Resend SMTP, and Google Maps.

---

## 🚀 Features

- Firebase Authentication (Email/Password & Google)
- Role-based Access Control (Customer, Seller, Admin)
- Product, Variant, Discount, Cart, Order system
- Image storage via Cloudinary
- Email service via Resend
- Machine Learning integration (Sentiment, Summary, Forecast)

---

## 📦 Tech Stack

| Category         | Tech                          |
|------------------|-------------------------------|
| Language         | Go (Golang)                   |
| Framework        | Go Fiber                      |
| Database         | PostgreSQL (via Supabase)     |
| Cache / Session  | Redis (Upstash)               |
| Auth             | Firebase Authentication       |
| File Storage     | Cloudinary                    |
| Email SMTP       | Resend                        |
| ML API           | Flask (via HTTP)              |

---

## 📁 Environment Variables

Rename the .env.example file to .env, then manually fill in all the required values.

---

## 🔑 Note on PLATFORM_ID

`PLATFORM_ID` is a special `users.id` UUID value used internally to represent the platform’s official bank account. This account is required in the system for certain financial operations, such as receiving payments and calculating transfer fees.

To set this up:

- Make sure you have created a user (usually with role_id = 4 for “Platform”).
- Insert a bank account that references this user into the `bank_accounts` table, for example:

```sql
INSERT INTO bank_accounts (id, user_id, bank_id, account_number, account_name, is_active, created_at)
VALUES (
    gen_random_uuid(),
    '<PLATFORM_ID>',
    1, -- or any existing bank_id from the bank_list table
    '1234567890',
    'Platform Tokoaku',
    true,
    now()
);
``` 
---

## 🚀 Running the Project

1. **Clone the repository**

```bash
git clone https://github.com/winterheatherica/tokoaku-backend.git
cd tokoaku-backend
```

2. **Copy the .env.example file and rename it:**

```bash
cp .env.example .env
```

3. **Install dependencies**

```bash
go mod tidy
```

4. **Run database migration**

```bash
DB.AutoMigrate(...) // ← remove comment to enable migration
```

5. **Run seeders**

```bash
seed.RunAllSeeders(...) // ← remove comment to seed data
```

6. **Start the server**

```bash
go run ./cmd
``` 
