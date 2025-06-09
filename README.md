# GophKeeper Password Manager

GophKeeper is a client-server application that allows users to securely store logins, passwords, binary data, and other private information.

---

## Architecture Overview

* **Communication protocol:** gRPC
* **Authentication:** JWT tokens (HS256)
* **Server storage:** PostgreSQL
* **Client storage:** BoltDB (local cache)
* **Encryption:** End-to-end AES-256-GCM with Argon2id key derivation
* **Key storage:** OS-native keychain systems

---

## Client usage

The client is a CLI application supporting Windows, Linux, and macOS. Application configuration is stored in the user's home directory with OS-specific paths:

* Linux/macOS: `~/.config/goph-keeper-client/config.json`
* Windows: `%APPDATA%\goph-keeper-client\config.json`

CLI Commands (terminal user interface):

| Command                          | Description                             |
| -------------------------------- | --------------------------------------- |
| `goph-keeper-client version`     | Show client version and build date      |
| `goph-keeper-client register`    | Register a new user on the server       |
| `goph-keeper-client list`        | Retrieve and display all stored entries |
| `goph-keeper-client put`         | Add or update an entry                  |
| `goph-keeper-client delete {id}` | Delete an entry by its UUID             |

---

## Registration

```shell
goph-keeper-client register
```

1. Prompt for login (email or any string) and master password (no email verification is required).
2. Client sends registration request to server.
3. Server verifies login uniqueness and creates user record with UUID.
4. JWT token is generated and returned.
5. Client saves login and UUID to the config file (for Linux it will be `~/.config/goph-keeper-client/config.json`).
6. Derive an AES-256-GCM encryption key from the master password using Argon2id (salt = login + UUID).
7. Store derived key and JWT in the OS keychain.

*Note: The master password cannot be recovered or changed. Loss results in permanent data access loss.*

---

## Authentication

For any command except `version` and `register`:

1. Client checks the OS keychain for a stored encryption key and JWT.
2. If missing, prompt for login and master password (if login is in config, only prompt for password).
3. Send authentication request to server.
4. On success, derive encryption key, store key and JWT in keychain.
5. Execute the initial command.

---

## Synchronization

For any command except `version` and `register`:

1. Client checks local BoltDB cache for existing entries.
2. If missing, fetches all user data from server
3. Stores encrypted data in local BoltDB cache.

---

## Data Operations

### List Items

```shell
goph-keeper-client list
```

Displays all entries from the local BoltDB cache in JSON format.

### Add/Edit Items

```shell
goph-keeper-client put
```

1. Prompts for data type:

   * Login/password pairs
   * Text notes
   * Binary files
   * Bank card details
2. Prompts for any JSON-formatted metadata.
3. Generates UUID for new items.
4. Encrypt the data client-side and send to server.
5. Saves to server with composite primary key (UUID + user ID).
6. Updates local BoltDB cache (overwrites existing items)

### Delete Items

```shell
goph-keeper-client delete {id}
```

1. Sends delete request to server for specified ID
2. Removes item from local BoltDB cache
