module github.com/jaxxstorm/tscli/tools

go 1.26.0

require (
	github.com/jaxxstorm/tscli v0.0.0
	gopkg.in/yaml.v3 v3.0.1
	tailscale.com/client/tailscale/v2 v2.10.1
)

require (
	github.com/tailscale/hujson v0.0.0-20260727124030-b80ff77dac4f // indirect
	golang.org/x/oauth2 v0.36.0 // indirect
)

replace github.com/jaxxstorm/tscli => ..
