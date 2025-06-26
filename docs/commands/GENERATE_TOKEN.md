# 🔧 wnc generate token

Generate a basic authentication token to connect to the WNC.

## ✨ Features

- Generates Base64-encoded Basic Authentication tokens
- Secure password handling
- Compatible with all `wnc show` commands

## 📋 Syntax

```bash
wnc generate token [options...]
```

## ⚙️ Flags

| Flag         | Alias | Type   | Description                                      | Required |
| ------------ | ----- | ------ | ------------------------------------------------ | -------- |
| `--username` | `-u`  | string | Username to generate Basic Authentication header | Yes      |
| `--password` | `-p`  | string | Password to generate Basic Authentication header | Yes      |

## 📝 Usage

```bash
# Basic token generation
wnc generate token --username admin --password mypassword

# Store token in variable
TOKEN=$(wnc generate token --username admin --password mypassword)
echo "Generated token: $TOKEN"

# Use with abbreviated flags
wnc generate token -u admin -p mypassword
```

## 📤 Example Output

```text
❯ ./wnc generate token --username admin --password mypassword
YWRtaW46bXlwYXNzd29yZA==
```

## 📖 Related Commands

- [wnc show ap](SHOW_AP.md)
- [wnc show client](SHOW_CLIENT.md)
- [wnc show overview](SHOW_OVERVIEW.md)
