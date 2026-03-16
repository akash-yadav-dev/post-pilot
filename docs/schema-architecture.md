Below is a **production-oriented relational schema** for a **Buffer-like social media scheduling platform** (similar to what you’re building). The schema is normalized, supports **multi-platform publishing, scheduling, worker queues, retries, and analytics**.

The system domains are:

```text
users
auth_accounts
social_accounts
posts
post_targets
jobs
job_runs
analytics_events
```

This structure works well with your **Go backend + Postgres + worker queue architecture**.

---

# 1. Users (core identity)

Stores the primary account in your system.

```sql
CREATE TABLE users (

  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  email VARCHAR(255) UNIQUE NOT NULL,
  name VARCHAR(255) NOT NULL,

  plan VARCHAR(50) NOT NULL DEFAULT 'free',

  credits INT NOT NULL DEFAULT 10
  CHECK (credits >= 0),

  is_active BOOLEAN DEFAULT TRUE,

  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()

);
```

Indexes:

```sql
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_plan ON users(plan);
```

---

# 2. Auth Accounts (login providers)

Handles authentication providers like:

* email/password
* Google
* GitHub

```sql
CREATE TABLE auth_accounts (

  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  user_id UUID NOT NULL
  REFERENCES users(id) ON DELETE CASCADE,

  provider VARCHAR(50) NOT NULL,

  provider_user_id VARCHAR(255) NOT NULL,

  password_hash TEXT,

  created_at TIMESTAMPTZ DEFAULT NOW(),

  UNIQUE(provider, provider_user_id)

);
```

Example:

| provider | provider_user_id |
| -------- | ---------------- |
| google   | 1123123123       |
| github   | 88923123         |

---

# 3. Social Accounts (platform connections)

Stores connected platforms:

* Twitter/X
* LinkedIn
* Mastodon
* Facebook

```sql
CREATE TABLE social_accounts (

  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  user_id UUID NOT NULL
  REFERENCES users(id) ON DELETE CASCADE,

  platform VARCHAR(50) NOT NULL,

  account_id VARCHAR(255) NOT NULL,

  account_name VARCHAR(255),

  access_token TEXT NOT NULL,

  refresh_token TEXT,

  token_expires_at TIMESTAMPTZ,

  metadata JSONB DEFAULT '{}'::jsonb,

  created_at TIMESTAMPTZ DEFAULT NOW(),

  UNIQUE(platform, account_id)

);
```

Example:

| platform | account   |
| -------- | --------- |
| twitter  | @user     |
| linkedin | John Page |

---

# 4. Posts (content created by users)

Stores the **base post content**.

```sql
CREATE TABLE posts (

  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  user_id UUID NOT NULL
  REFERENCES users(id) ON DELETE CASCADE,

  content TEXT NOT NULL,

  media_urls TEXT[],

  status VARCHAR(50) NOT NULL DEFAULT 'draft',

  scheduled_at TIMESTAMPTZ,

  created_at TIMESTAMPTZ DEFAULT NOW(),

  updated_at TIMESTAMPTZ DEFAULT NOW()

);
```

Status values:

```text
draft
scheduled
publishing
published
failed
```

Indexes:

```sql
CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_scheduled_at ON posts(scheduled_at);
```

---

# 5. Post Targets (multi-platform publishing)

One post can go to **multiple platforms**.

Example:

```text
post → twitter + linkedin
```

Schema:

```sql
CREATE TABLE post_targets (

  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  post_id UUID NOT NULL
  REFERENCES posts(id) ON DELETE CASCADE,

  social_account_id UUID NOT NULL
  REFERENCES social_accounts(id) ON DELETE CASCADE,

  status VARCHAR(50) DEFAULT 'pending',

  platform_post_id VARCHAR(255),

  published_at TIMESTAMPTZ

);
```

Example:

| post  | platform |
| ----- | -------- |
| post1 | twitter  |
| post1 | linkedin |

---

# 6. Jobs (scheduler queue)

Your worker system will process jobs from this table.

```sql
CREATE TABLE jobs (

  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  type VARCHAR(50) NOT NULL,

  payload JSONB NOT NULL,

  status VARCHAR(50) NOT NULL DEFAULT 'pending',

  run_at TIMESTAMPTZ NOT NULL,

  attempts INT DEFAULT 0,

  max_attempts INT DEFAULT 3,

  last_error TEXT,

  created_at TIMESTAMPTZ DEFAULT NOW()

);
```

Job types:

```text
publish_post
retry_publish
refresh_token
analytics_sync
```

Indexes:

```sql
CREATE INDEX idx_jobs_run_at ON jobs(run_at);
CREATE INDEX idx_jobs_status ON jobs(status);
```

---

# 7. Job Runs (worker execution history)

Tracks execution attempts.

```sql
CREATE TABLE job_runs (

  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  job_id UUID NOT NULL
  REFERENCES jobs(id) ON DELETE CASCADE,

  started_at TIMESTAMPTZ,

  finished_at TIMESTAMPTZ,

  status VARCHAR(50),

  error TEXT

);
```

This helps with:

* debugging
* monitoring
* retry logic

---

# 8. Analytics Events

Stores platform engagement data.

Example:

```text
likes
shares
comments
clicks
```

```sql
CREATE TABLE analytics_events (

  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  post_target_id UUID NOT NULL
  REFERENCES post_targets(id) ON DELETE CASCADE,

  event_type VARCHAR(50) NOT NULL,

  value INT DEFAULT 1,

  created_at TIMESTAMPTZ DEFAULT NOW()

);
```

Example events:

| type    |
| ------- |
| like    |
| share   |
| comment |
| click   |

---

# 9. Optional Table: Scheduled Queue Cache

If you want **fast worker scheduling**.

```sql
CREATE TABLE scheduler_locks (

  id SERIAL PRIMARY KEY,

  worker_id VARCHAR(255),

  locked_at TIMESTAMPTZ

);
```

Used to coordinate multiple workers.

---

# 10. Full Database Structure

Final system:

```text
users
│
├── auth_accounts
│
├── social_accounts
│
├── posts
│     │
│     └── post_targets
│
├── jobs
│     │
│     └── job_runs
│
└── analytics_events
```

---

# 11. Example Publishing Flow

User schedules post.

```text
User → create post
     → create post_targets
     → create publish_job
```

Worker flow:

```text
worker → fetch jobs
       → load post
       → load social_account
       → publish
       → update post_targets
```

---

# 12. Example Job Payload

```json
{
  "post_id": "uuid",
  "target_id": "uuid"
}
```

Worker:

```go
PublishPost(payload.PostID, payload.TargetID)
```

---

# 13. Key Design Advantages

This schema supports:

✔ multiple platforms
✔ multiple accounts per platform
✔ scheduled publishing
✔ worker retries
✔ analytics tracking
✔ queue workers

This architecture scales to **millions of scheduled posts**.

---

✅ This schema is **very close to what real scheduling SaaS tools use**.

---
