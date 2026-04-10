package tscli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jaxxstorm/tscli/pkg/apitype"
	tsapi "tailscale.com/client/tailscale/v2"
)

func ListDevicesJSON(ctx context.Context, client *tsapi.Client, allFields bool) (json.RawMessage, error) {
	path := "/tailnet/{tailnet}/devices"
	if allFields {
		path += "?fields=all"
	}
	return rawJSON(ctx, client, http.MethodGet, path, nil)
}

func GetDeviceJSON(ctx context.Context, client *tsapi.Client, deviceID string, allFields bool) (json.RawMessage, error) {
	path := fmt.Sprintf("/device/%s", url.PathEscape(deviceID))
	if allFields {
		path += "?fields=all"
	}
	return rawJSON(ctx, client, http.MethodGet, path, nil)
}

func ListDeviceRoutesJSON(ctx context.Context, client *tsapi.Client, deviceID string) (json.RawMessage, error) {
	path := fmt.Sprintf("/device/%s/routes", url.PathEscape(deviceID))
	return rawJSON(ctx, client, http.MethodGet, path, nil)
}

func SetDeviceRoutesJSON(ctx context.Context, client *tsapi.Client, deviceID string, request apitype.DeviceRoutesUpdateRequest) (json.RawMessage, error) {
	path := fmt.Sprintf("/device/%s/routes", url.PathEscape(deviceID))
	return rawJSON(ctx, client, http.MethodPost, path, request)
}

func GetTailnetSettingsJSON(ctx context.Context, client *tsapi.Client) (json.RawMessage, error) {
	return rawJSON(ctx, client, http.MethodGet, "/tailnet/{tailnet}/settings", nil)
}

func UpdateTailnetSettingsJSON(ctx context.Context, client *tsapi.Client, request apitype.UpdateTailnetSettingsRequest) (json.RawMessage, error) {
	return rawJSON(ctx, client, http.MethodPatch, "/tailnet/{tailnet}/settings", request)
}

func rawJSON(ctx context.Context, client *tsapi.Client, method, path string, body any) (json.RawMessage, error) {
	var raw json.RawMessage
	if _, err := Do(ctx, client, method, path, body, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}
