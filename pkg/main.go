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
	"github.com/proxy-wasm/proxy-wasm-go-sdk/proxywasm"
	"github.com/proxy-wasm/proxy-wasm-go-sdk/proxywasm/types"
)

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

	return types.ActionContinue
}

// OnHttpRequestBody implements types.HttpContext.
func (ctx *mcpGatewayContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {

    proxywasm.LogWarnf("OnHttpRequestBody()")

	if !endOfStream {
		return types.ActionPause
	}

	return types.ActionContinue
}

// OnHttpResponseHeaders implements types.HttpContext.
func (ctx *mcpGatewayContext) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {

	proxywasm.LogWarnf("OnHttpResponseHeaders()")

	return types.ActionContinue
}

// OnHttpResponseBody implements types.HttpContext.
func (ctx *mcpGatewayContext) OnHttpResponseBody(bodySize int, endOfStream bool) types.Action {

	proxywasm.LogWarnf("OnHttpResponseBody()")

	if !endOfStream {
		return types.ActionPause
	}

	return types.ActionContinue
}

