# Database Schema Documentation

This document describes the database schema for the light-social investing platform. All tables use snake_case naming. Data types and constraints are specified for clarity and best practices. PostgreSQL ENUM types are used where appropriate and are defined near their usage.

---

## users

| Column      | Type           | Constraints                | Description                |
|-------------|----------------|----------------------------|----------------------------|
| id          | SERIAL         | PRIMARY KEY                | User ID                    |
| email       | VARCHAR(255)   | NOT NULL, UNIQUE           | User email                 |
| created_at  | TIMESTAMP      | NOT NULL, DEFAULT now()    | Creation timestamp         |
| updated_at  | TIMESTAMP      | NOT NULL, DEFAULT now()    | Last update timestamp      |

---

-- ENUM for user_oauths
{% raw %}```sql
CREATE TYPE oauth_provider AS ENUM ('google', 'apple', 'x');
```{% endraw %}

## user_oauths

| Column        | Type           | Constraints                                | Description                |
|---------------|----------------|--------------------------------------------|----------------------------|
| user_id       | INTEGER        | REFERENCES users(id) NOT NULL               | User ID (FK)               |
| provider      | oauth_provider | NOT NULL                                   | OAuth provider             |
| access_token  | TEXT           | NOT NULL                                   | OAuth access token         |
| refresh_token | TEXT           |                                            | OAuth refresh token        |
| expires_at    | TIMESTAMP      |                                            | Token expiration           |
| created_at    | TIMESTAMP      | NOT NULL, DEFAULT now()                    | Creation timestamp         |
| updated_at    | TIMESTAMP      | NOT NULL, DEFAULT now()                    | Last update timestamp      |

**Primary Key:** (`user_id`, `provider`)

---

## user_email_verifications

| Column      | Type           | Constraints                | Description                |
|-------------|----------------|----------------------------|----------------------------|
| id          | SERIAL         | PRIMARY KEY                | Verification ID            |
| user_id     | INTEGER        | REFERENCES users(id)       | User ID (FK)               |
| email       | VARCHAR(255)   | NOT NULL                   | Email to verify            |
| token       | VARCHAR(255)   | NOT NULL, UNIQUE           | Verification token         |
| expires_at  | TIMESTAMP      | NOT NULL                   | Expiration timestamp       |
| created_at  | TIMESTAMP      | NOT NULL, DEFAULT now()    | Creation timestamp         |
| updated_at  | TIMESTAMP      | NOT NULL, DEFAULT now()    | Last update timestamp      |

---

## user_follower_counts

| Column      | Type           | Constraints                | Description                |
|-------------|----------------|----------------------------|----------------------------|
| user_id     | INTEGER        | PRIMARY KEY, REFERENCES users(id) | User ID (FK)      |
| count       | INTEGER        | NOT NULL, DEFAULT 0        | Follower count             |
| updated_at  | TIMESTAMP      | NOT NULL, DEFAULT now()    | Last update timestamp      |

---

-- ENUMs for dashboards
{% raw %}```sql
CREATE TYPE dashboard_status AS ENUM ('draft', 'active', 'paused', 'archived');
CREATE TYPE dashboard_publicity AS ENUM ('public', 'private', 'followers_only');
```{% endraw %}

## dashboards

| Column      | Type               | Constraints                | Description                |
|-------------|--------------------|----------------------------|----------------------------|
| id          | SERIAL             | PRIMARY KEY                | Dashboard ID               |
| user_id     | INTEGER            | REFERENCES users(id)       | Owner (FK)                 |
| title       | VARCHAR(255)       | NOT NULL                   | Dashboard title            |
| content     | JSONB              |                            | Dashboard content          |
| status      | dashboard_status   | NOT NULL, DEFAULT 'draft'  | Dashboard status           |
| publicity   | dashboard_publicity| NOT NULL, DEFAULT 'public' | Dashboard publicity        |
| created_at  | TIMESTAMP          | NOT NULL, DEFAULT now()    | Creation timestamp         |
| updated_at  | TIMESTAMP          | NOT NULL, DEFAULT now()    | Last update timestamp      |

---

## artifacts

| Column      | Type           | Constraints                | Description                |
|-------------|----------------|----------------------------|----------------------------|
| id          | SERIAL         | PRIMARY KEY                | Artifact ID                |
| user_id     | INTEGER        | REFERENCES users(id)       | Owner (FK)                 |
| title       | VARCHAR(255)   | NOT NULL                   | Artifact title             |
| content     | JSONB          |                            | Artifact content           |
| created_at  | TIMESTAMP      | NOT NULL, DEFAULT now()    | Creation timestamp         |
| updated_at  | TIMESTAMP      | NOT NULL, DEFAULT now()    | Last update timestamp      |

---

## dashboard_artifacts

| Column        | Type           | Constraints                | Description                |
|---------------|----------------|----------------------------|----------------------------|
| id            | SERIAL         | PRIMARY KEY                | Dashboard artifact ID       |
| dashboard_id  | INTEGER        | REFERENCES dashboards(id)  | Dashboard (FK)             |
| artifact_id   | INTEGER        | REFERENCES artifacts(id)   | Artifact (FK)              |
| created_at    | TIMESTAMP      | NOT NULL, DEFAULT now()    | Creation timestamp         |
| updated_at    | TIMESTAMP      | NOT NULL, DEFAULT now()    | Last update timestamp      |

---

## likes

| Column      | Type           | Constraints                | Description                |
|-------------|----------------|----------------------------|----------------------------|
| id          | SERIAL         | PRIMARY KEY                | Like ID                    |
| user_id     | INTEGER        | REFERENCES users(id)       | User (FK)                  |
| artifact_id | INTEGER        | REFERENCES artifacts(id)   | Artifact (FK)              |
| created_at  | TIMESTAMP      | NOT NULL, DEFAULT now()    | Creation timestamp         |
| updated_at  | TIMESTAMP      | NOT NULL, DEFAULT now()    | Last update timestamp      |

---

## follows

| Column           | Type           | Constraints                | Description                |
|------------------|----------------|----------------------------|----------------------------|
| id               | SERIAL         | PRIMARY KEY                | Follow ID                  |
| user_id          | INTEGER        | REFERENCES users(id)       | Follower (FK)              |
| followed_user_id | INTEGER        | REFERENCES users(id)       | Followed user (FK)         |
| created_at       | TIMESTAMP      | NOT NULL, DEFAULT now()    | Creation timestamp         |
| updated_at       | TIMESTAMP      | NOT NULL, DEFAULT now()    | Last update timestamp      |

---

## user_dashboard_feeds

| Column        | Type           | Constraints                | Description                |
|---------------|----------------|----------------------------|----------------------------|
| id            | SERIAL         | PRIMARY KEY                | Feed ID                    |
| user_id       | INTEGER        | REFERENCES users(id)       | User (FK)                  |
| dashboard_id  | INTEGER        | REFERENCES dashboards(id)  | Dashboard (FK)             |
| created_at    | TIMESTAMP      | NOT NULL, DEFAULT now()    | Creation timestamp         |
| updated_at    | TIMESTAMP      | NOT NULL, DEFAULT now()    | Last update timestamp      |

---
