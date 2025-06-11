# GophKeeper Password Manager

GophKeeper is a client-server application that allows users to securely store logins, passwords, binary data, and other private information.

---

## Architecture Overview

- **Communication:** HTTP(S)
- **Authentication:** JWT (HS256)
- **Server Storage:** PostgreSQL
- **Client Storage:** BoltDB (local cache)
- **Security:**
  - End-to-end AES-256-GCM encryption
  - Client-side Argon2id key derivation
  - Irrecoverable master password (no reset mechanism)                          |

---

## Client Application

The client is a command-line interface (CLI) application built with BubbleTea TUI, compatible with Windows, Linux, and macOS. Configuration is stored in the user's home directory with OS-specific paths:

- **Linux/macOS:** `~/.config/goph-keeper-client/config.json`
- **Windows:** `%APPDATA%\goph-keeper-client\config.json`

Upon startup, the application checks for the configuration file. If it doesn’t exist, a new one is created with the default server address: `http://localhost:8080`.

**Workflow:**

1. **Initial Menu:**
   - `Login` / `Register` / `About`
2. **Authenticated Menu:**
   - `Show` / `Add` / `Sync` / `About`

---

## Operations

### 1. Registration (`POST /register`)

User registers with login (any string, no email verification) and master password.

**Server:**

  - Verifies the uniqueness of the `login`.
  - Creates a user record with a UUID and a password hash (using Argon2id).
  - Generates and returns a JWT token.

**Client:**

  - Derives an AES-256-GCM encryption key from the master password using Argon2id (salt = login + UUID).
  - Stores the key in memory.

*Note: The master password cannot be recovered or changed. Loss results in permanent data access loss.*

### 2. Authentication (`POST /login`)

**Server:** Verifies credentials against Argon2id hash

**Client:** Identical key derivation as registration

### 3. Synchronization

**Triggers:** Post-auth, every 10min (background), or manually

**Process:**

1. Check server availability via `GET /version`
  2. Pull server records via `GET /records` → merge into local BoltDB
  3. Push `pending` records via `PUT /records/{id}`
  4. Update statuses (`synced`/`conflict`)
  5. Delete `deleted` records via `DELETE /records/{id}`
  6. Completely remove records marked `deleted` locally
  7. Show message if has conflicts

### 4. Data Management

#### Show

Displays all records stored in the local BoltDB cache. User can select a record via the TUI to:
- View details
- Update
- Delete
- Resolve conflicts (if applicable)

#### Add/Update (`PUT /records/{id}`)
1. Select data type:
   - Credentials • Text • Binary • Bank card
2. Enter data + optional JSON metadata
3. Client:
   - Generates UUID (for new records)
   - Encrypts payload (AES-256-GCM)
   - Sends to server
4. Server stores record with a composite primary key (UUID + user ID)
5. Conflict handling triggers resolution UI

#### Delete (`DELETE /records/{id}`)

1. Mark `deleted` locally
2. Send delete request
3. Remove from cache on success

### 5. Conflict Resolution
**Detection:** Version mismatch during sync/update
**Resolution Workflow:**
1. Visual indicators:
   - Conflicted records are marked with `[CONFLICT]` flag in the `Show` list
   - Detailed conflict notification appears in the TUI interface
2. User selects conflicted record → "Resolve"
3. Side-by-side comparison:
   - Local version (client changes)
   - Remote version (server state)
4. User selects preferred version
5. Client:
   - Increments version number
   - Sends update via `PUT /records/{id}`
   - Updates status: `conflict` → `synced`

### 6. System Information (`About`)
Displays build metadata:
- Client version/date
- Server version/date
