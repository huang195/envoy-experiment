// Copyright 2020-2024 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/proxy-wasm/proxy-wasm-go-sdk/proxywasm"
	"github.com/proxy-wasm/proxy-wasm-go-sdk/proxywasm/types"
)

type RequestPayload struct {
	RemoveHeaders bool `json:"removeheaders"`
}

func main() {}
func init() {
	proxywasm.SetVMContext(&vmContext{})
}

// vmContext implements types.VMContext.
type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// NewPluginContext implements types.VMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {

	return &pluginContext{}
}

// pluginContext implements types.PluginContext.
type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
}

// NewHttpContext implements types.PluginContext.
func (ctx *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &mcpGatewayContext{}
}

// OnPluginStart implements types.PluginContext.
func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	proxywasm.LogWarnf("MCP Gateway started")
	return types.OnPluginStartStatusOK
}

// mcpGatewayContext implements types.HttpContext.
type mcpGatewayContext struct {
	// Embed the default root http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
}

// OnHttpRequestHeaders implements types.HttpContext.
func (ctx *mcpGatewayContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {

    proxywasm.LogWarnf("OnHttpRequestHeaders()")

	if !endOfStream {
        return types.ActionPause
    }

	return types.ActionContinue
}

// OnHttpRequestBody implements types.HttpContext.
func (ctx *mcpGatewayContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {

    proxywasm.LogWarnf("OnHttpRequestBody()")

	if !endOfStream {
		return types.ActionPause
	}

	payload, err := proxywasm.GetHttpRequestBody(0, bodySize)
	if err != nil {
		body := fmt.Sprintf("Cannot parse HTTP payload: %v", err)
		if err := proxywasm.SendHttpResponse(400, nil, []byte(body), -1); err != nil {
			proxywasm.LogWarnf("Failed to send HTTP response with error: %v", err)
		}
		return types.ActionContinue
	}

	proxywasm.LogWarnf("payload: %s", payload)

	var reqPayload RequestPayload
	if err := json.Unmarshal(payload, &reqPayload); err != nil {
		proxywasm.LogWarnf("Failed to parse JSON payload: %v", err)
		return types.ActionContinue
	}

	proxywasm.LogWarnf("removeHeaders: %v", reqPayload.RemoveHeaders)

	if reqPayload.RemoveHeaders {
		proxywasm.LogWarnf("removeHeaders is true, removing all request headers")
		
		headers, err := proxywasm.GetHttpRequestHeaders()
		if err != nil {
			proxywasm.LogWarnf("Failed to get request headers: %v", err)
			return types.ActionContinue
		}
		
		for _, header := range headers {
			headerName := header[0]
			if err := proxywasm.RemoveHttpRequestHeader(headerName); err != nil {
				proxywasm.LogWarnf("Failed to remove header %s: %v", headerName, err)
			}
		}
	}

	return types.ActionContinue
}
