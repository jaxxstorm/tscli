# **tscli**

`tscli` is a fast, single-binary CLI for the [Tailscale HTTP API](https://tailscale.com/api).
From your terminal you can manage devices, users, auth keys, webhooks, posture integrations, tailnet-wide settings, and even hit raw endpoints when the SDK hasn’t caught up yet.

## ✨ Highlights

| Area                     | What you can do                                                                                             |
| ------------------------ | ----------------------------------------------------------------------------------------------------------- |
| **Devices**              | List, get, (de)authorize, rename, force IPv4, enable subnet routes, expire, set / delete posture attributes |
| **Keys**                 | List & get existing keys; create **auth-keys** _or_ **OAuth clients** (with full scope/tag validation)      |
| **Users**                | List (filter by type / role), get, suspend / restore / approve, manage invites                              |
| **Tailnet settings**     | Get & patch booleans + key-expiry with a single command (`tscli set settings …`)                            |
| **Policy file (ACL)**    | Fetch as raw HUJSON **or** canonical JSON                                                                   |
| **Webhooks**             | List, get, delete, **create** (generic / Slack) with subscription & provider validation                     |
| **Posture integrations** | List, get, create, patch existing integrations                                                              |
| **Invites**              | List / delete device- or user-invites                                                                       |
| **Contacts**             | Get & update contact emails                                                                                 |
| **Debug switch**         | `--debug` or `TSCLI_DEBUG=1` prints full HTTP requests / responses to stderr                                |
| **Config precedence**    | _flags_ → _env_ → `~/.tscli/.tscli.yaml` (or local `./.tscli.yaml`)                                         |

## 🔧 Install

### 🔧 Installation

#### macOS / Linux (Homebrew)

```bash
brew tap jaxxstorm/tap
brew install tscli          # upgrades via ‘brew upgrade’
```

#### Windows (Scoop)

```powershell
scoop bucket add jaxxstorm https://github.com/jaxxstorm/scoop-bucket.git
scoop install tscli
```

#### Nix

```bash
nix shell github:jaxxstorm/tscli
```

#### Manual download

Pre-built archives for **macOS, Linux, Windows (x86-64 / arm64)** are published on every release:

```bash
# example for Linux amd64
curl -sSfL \
  https://github.com/jaxxstorm/tscli/releases/latest/download/tscli_$(uname -s)_$(uname -m).tar.gz \
  | sudo tar -xz -C /usr/local/bin tscli
```

#### Go install (always builds from HEAD)

```bash
go install github.com/jaxxstorm/tscli@latest
```

After any method, confirm:

```bash
tscli --version
```

## ⚙️ Configuration

| Option            | Flag / Env var                          | YAML key  | Default |
| ----------------- | --------------------------------------- | --------- | ------- |
| Tailscale API key | `--api-key`, `-k` / `TAILSCALE_API_KEY` | `api-key` | —       |
| Tailnet name      | `--tailnet`, `-n` / `TAILSCALE_TAILNET` | `tailnet` | —       |

```yaml
# ~/.tscli/.tscli.yaml
api-key: tskey-abc123…
tailnet: example.com
format: pretty # other options are: human, json or yaml
```

## 🚀 Usage

```text
tscli <noun> <verb> [flags]
```

### Global flags

```
-k, --api-key string   Tailscale API key
-n, --tailnet string   Tailnet (default "-")
-d, --debug            Dump raw HTTP traffic to stderr
```

## 📜 Coverage

| API Area / Action                |       Status       | `tscli` Command                                             |
| -------------------------------- | :----------------: | ----------------------------------------------------------- |
| **Devices**                      |                    |                                                             |
| list devices                     | :white_check_mark: | `device list`                                               |
| get device                       | :white_check_mark: | `device get --device <id>`                                  |
| authorize / de-authorize device  | :white_check_mark: | `device authorize --device <id> [--approve=<bool>]`         |
| expire device key                | :white_check_mark: | `device expire --device <id>`                               |
| set device name                  | :white_check_mark: | `set name --device <id> --name <host>`                      |
| set device tags                  | :white_check_mark: | `set tags --device <id> --tag <tag>`                        |
| set device IPv4                  | :white_check_mark: | `set ip --device <id> --ip <addr>`                          |
| list subnet routes               | :white_check_mark: | `list routes --device <id>`                                 |
| set subnet routes                | :white_check_mark: | `set routes --device <id> --route <cidr> …`                 |
| **delete device**                | :white_check_mark: | `delete device --device <id>`                               |
| **Device-posture attributes**    |                    |                                                             |
| get attributes                   | :white_check_mark: | `get posture --device <id>`                                 |
| set attribute                    | :white_check_mark: | `set attribute --device <id> --key custom:x --value 42`     |
| delete attribute                 | :white_check_mark: | `delete attribute --device <id> --key custom:x`             |
| **Posture integrations**         |                    |                                                             |
| list integrations                | :white_check_mark: | `list posture-integrations`                                 |
| get integration                  | :white_check_mark: | `get posture-integration --id <id>`                         |
| create integration               | :white_check_mark: | `create posture-integration --provider …`                   |
| update integration               | :white_check_mark: | `set posture-integration --id <id> …`                       |
| delete integration               |        :x:         | —                                                           |
| **Auth / OAuth keys**            |                    |                                                             |
| list keys                        | :white_check_mark: | `list keys`                                                 |
| get key                          | :white_check_mark: | `get key --key <id>`                                        |
| create auth-key                  | :white_check_mark: | `create key --type authkey …`                               |
| create OAuth client              | :white_check_mark: | `create key --type oauthclient …`                           |
| delete / revoke key              |        :x:         | —                                                           |
| **Users**                        |                    |                                                             |
| list users                       | :white_check_mark: | `list users [--type …] [--role …]`                          |
| get user                         | :white_check_mark: | `get user --user <id>`                                      |
| approve / suspend / restore user | :white_check_mark: | `set user-access --user <id> --approve/--suspend/--restore` |
| delete user                      | :white_check_mark: | `delete user --user <id>`                                   |
| **Invites**                      |                    |                                                             |
| list user invites                | :white_check_mark: | `list invites user [--state …]`                             |
| list device invites              | :white_check_mark: | `list invites device --device <id>`                         |
| delete invite                    |        :x:         | —                                                           |
| get invite                       |        :x:         | —                                                           |
| **Contacts**                     |                    |                                                             |
| update contact                   | :white_check_mark: | `set contacts --contact <id> --email <e@x>`                 |
| list / create / delete contacts  |        :x:         | —                                                           |
| **DNS**                          |                    |                                                             |
| list nameservers                 | :white_check_mark: | `list nameservers`                                          |
| set nameservers                  | :white_check_mark: | `set nameservers --nameserver <ip> …`                       |
| advanced DNS settings            |        :x:         | —                                                           |
| **Policy file (ACL)**            |                    |                                                             |
| get policy file                  | :white_check_mark: | `get policy [--json]`                                       |
| set policy file                  | :white_check_mark: | `set policy --file <acl.hujson>`                            |
| preview rule matches             | :white_check_mark: | `get policy-preview --type … --value … [--file]`            |
| policy history / tests           |        :x:         | —                                                           |
| **Tailnet settings**             |                    |                                                             |
| get settings                     | :white_check_mark: | `get settings`                                              |
| update settings                  | :white_check_mark: | `set settings --devices-approval …`                         |
| **Webhooks**                     |                    |                                                             |
| list webhooks                    | :white_check_mark: | `list webhooks`                                             |
| get webhook                      | :white_check_mark: | `get webhook --id <id>`                                     |
| create webhook                   | :white_check_mark: | `create webhook --url <endpoint> --subscription …`          |
| update webhook                   |        :x:         | —                                                           |
| delete webhook                   | :white_check_mark: | `delete webhook --id <id>`                                  |
| rotate webhook secret            | :white_check_mark: | `rotate webhook --id <id>`                                  |

### Quick examples

```bash
# Approve a waiting device
tscli device authorize --device node-abc123 --approve

# Rotate an auth-key that expires in 30 days
tscli create key --description "CI" --expiry 720h | jq .key

# Create Slack webhook for device deletions
tscli create webhook \
  --url https://hooks.slack.com/services/T000/B000/XXXXX \
  --provider slack \
  --subscription nodeDeleted
```

## 🛠 Development

```bash
git clone https://github.com/jaxxstorm/tscli
cd tscli
TAILSCALE_API_KEY=tskey-… TAILSCALE_TAILNET=example.com go run ./cmd/tscli list devices
```

Tests & lint:

```bash
go test ./...
```

## 📄 License

MIT — see [`LICENSE`](./LICENSE).
