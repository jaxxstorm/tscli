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
output: pretty # other options are: human, json or yaml
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

| API Area / Action                | Status | `tscli` Command |
| -------------------------------- | :----: | --------------- |
| **Devices**                      |        |                 |
| List tailnet devices             | :white_check_mark: | `tscli list devices` |
| Get a device                     | :white_check_mark: | `tscli get device --device <device>` |
| Delete a device                  | :white_check_mark: | `tscli delete device --device <device>` |
| Expire a device key              | :white_check_mark: | `tscli set expiry --device <device>` |
| List device routes               | :white_check_mark: | `tscli list routes --device <device>` |
| Set device routes                | :white_check_mark: | `tscli set routes --device <device> --route <cidr>` |
| Authorize / de-authorize device  | :white_check_mark: | `tscli set authorization --device <device> --approve=<bool>` |
| Set device name                  | :white_check_mark: | `tscli set name --device <device> --name <hostname>` |
| Set device tags                  | :white_check_mark: | `tscli set tags --device <device> --tag tag:<tag>` |
| Rotate device key                | :x: | — |
| Set device IPv4 address          | :white_check_mark: | `tscli set ip --device <device> --ip <ip>` |
| Get posture attributes           | :white_check_mark: | `tscli get posture --device <device>` |
| Set custom posture attributes    | :white_check_mark: | `tscli set posture --device <device> --key custom:x --value <v>` |
| Delete custom posture attributes | :white_check_mark: | `tscli delete posture --device <device> --key custom:x` |
| **Device Invites**               |        |                 |
| List device invites              | :white_check_mark: | `tscli list invites device --device <device>` |
| Create device invite             | :white_check_mark: | `tscli create invite device --device <device> --email <email>` |
| Get a device invite              | :x: | — |
| Delete a device invite           | :white_check_mark: | `tscli delete invite device --id <invite-id>` |
| Resend / accept device invite    | :x: | — |
| **User Invites**                 |        |                 |
| List user invites                | :white_check_mark: | `tscli list invites user [--state …]` |
| Create user invite               | :white_check_mark: | `tscli create invite user --email <email> [--role <role>]` |
| Get a user invite                | :x: | — |
| Delete a user invite             | :white_check_mark: | `tscli delete invite user --id <invite-id>` |
| Resend user invite               | :x: | — |
| **Logging**                      |        |                 |
| List configuration audit logs    | :white_check_mark: | `tscli list logs config --start <t> [--end <t>]` |
| List network flow logs           | :white_check_mark: | `tscli list logs network --start <t> [--end <t>]` |
| Log-streaming endpoints          | :x: | — |
| **DNS**                          |        |                 |
| List DNS nameservers             | :white_check_mark: | `tscli list nameservers` |
| Set DNS nameservers              | :white_check_mark: | `tscli set nameservers --nameserver <ip> …` |
| Get DNS preferences              | :white_check_mark: | `tscli get dns preferences` |
| Set DNS preferences              | :white_check_mark: | `tscli set dns preferences --magicdns=<bool>` |
| List DNS search paths            | :white_check_mark: | `tscli list dns searchpaths` |
| Set DNS search paths             | :white_check_mark: | `tscli set dns searchpaths --searchpath <domain> …` |
| Get split-DNS map                | :white_check_mark: | `tscli get dns split` |
| Patch split-DNS                  | :white_check_mark: | `tscli set dns split --domain <d>=<ip,ip> …` |
| Replace split-DNS                | :white_check_mark: | `tscli set dns split --replace --domain <d>=<ip>` |
| **Keys**                         |        |                 |
| List tailnet keys                | :white_check_mark: | `tscli list keys` |
| Create auth-key / OAuth client   | :white_check_mark: | `tscli create key --type authkey --oauthclient …` |
| Get key                          | :white_check_mark: | `tscli get key --key <id>` |
| Delete / revoke key              | :x: | — |
| **Policy File**                  |        |                 |
| Get policy file                  | :white_check_mark: | `tscli get policy [--json]` |
| Set policy file                  | :white_check_mark: | `tscli set policy --file <acl.hujson>` |
| Preview rule matches             | :white_check_mark: | `tscli get policy preview --type user\|ipport --value … [--current\|--file F]` |
| Validate / test policy           | :x: | — |
| **Posture Integrations**         |        |                 |
| List integrations                | :white_check_mark: | `tscli list posture-integrations` |
| Create integration               | :white_check_mark: | `tscli create posture-integration --provider <p> …` |
| Get integration                  | :white_check_mark: | `tscli get posture-integration --id <id>` |
| Update integration               | :white_check_mark: | `tscli set posture-integration --id <id> …` |
| Delete integration               | :x: | — |
| **Users**                        |        |                 |
| List users                       | :white_check_mark: | `tscli list users [--type …] [--role …]` |
| Get a user                       | :white_check_mark: | `tscli get user --user <id>` |
| Update user role                 | :white_check_mark: | `tscli set user-role --user <id> --role <role>` |
| Approve / suspend / restore user | :white_check_mark: | `tscli set user-access --user <id> --approve\|--suspend\|--restore` |
| Delete a user                    | :white_check_mark: | `tscli delete user --user <id>` |
| **Contacts**                     |        |                 |
| Get contacts                     | :white_check_mark: | `tscli get contacts` |
| Update contact                   | :white_check_mark: | `tscli set contacts --contact <id> --email <e@x>` |
| Resend verification              | :x: | — |
| **Webhooks**                     |        |                 |
| List webhooks                    | :white_check_mark: | `tscli list webhooks` |
| Create webhook                   | :white_check_mark: | `tscli create webhook --url <endpoint> …` |
| Get webhook                      | :white_check_mark: | `tscli get webhook --webhook <id>` |
| Update webhook                   | :white_check_mark: | `tscli set webhook --webhook <id> …` |
| Delete webhook                   | :white_check_mark: | `tscli delete webhook --webhook <id>` |
| Test / rotate webhook            | :white_check_mark: | `tscli rotate webhook --webhook <id>` |
| **Tailnet Settings**             |        |                 |
| Get tailnet settings             | :white_check_mark: | `tscli get settings` |
| Update tailnet settings          | :white_check_mark: | `tscli set settings --devices-approval …` |



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
