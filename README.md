# envoy-experiment

This Envoy experiment is to validate some of the recent fixes in proxy-wasm upstream are working as expected:
* Envoy: https://github.com/envoyproxy/envoy/pull/40213
* WASM sandbox: https://github.com/proxy-wasm/proxy-wasm-cpp-host/pull/434

In summary, a new flag `allow_on_headers_stop_iteration` was added to allow WASM filters to modify HTTP headers
based on the content of HTTP body, which is critical for handling MCP traffic. Previously, without this flag
attempt to modify headers would result in runtime error.

Using the head of the main branch from envoyproxy/envoy, we are able to verify that HTTP request headers can be 
now modified based on the content of HTTP body. Here are the steps to validate:

1. Pre-req

* Start a kind cluster locally
* Create a MCP server in the backend

```
kubectl apply -f deploy/mcp-gateway.yaml
```

2. Start Envoy proxy

```
kubectl apply -f deploy/envoy.yaml
```

and to exercise the proxy, we port-forward it locally:

```
kubectl port-forward service/envoy-proxy-service 8000:80
```

3. Test (normal)

To exercise the Envoy proxy to reach the backend MCP server:

```
curl -v -N -L -X POST http://localhost:8000/mcp/  \
-H "Content-Type: application/json" \
-H "Accept: application/json, text/event-stream" \
-d '
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2025-06-18",
    "capabilities": {
      "roots": {
        "listChanged": true
      },
      "sampling": {}
    },
    "clientInfo": {
      "name": "ExampleClient",
      "version": "1.0.0"
    }
  }
}
'
```

You should get a normal response like the following:

```
event: message
data: {"jsonrpc":"2.0","id":1,"result":{"protocolVersion":"2025-03-26","capabilities":{"experimental":{},"prompts":{"listChanged":true},"resources":{"subscribe":false,"listChanged":true},"tools":{"listChanged":true}},"serverInfo":{"name":"MCP Gateway","version":"1.9.4"}}}
```

4. Test (remove headers)

To test that we can modify headers based on the body, we add an additional field in the json payload `removeheaders`:

```
curl -v -N -L -X POST http://localhost:8000/mcp/  \
-H "Content-Type: application/json" \
-H "Accept: application/json, text/event-stream" \
-d '
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2025-06-18",
    "capabilities": {
      "roots": {
        "listChanged": true
      },
      "sampling": {}
    },
    "clientInfo": {
      "name": "ExampleClient",
      "version": "1.0.0"
    }
  },
  "removeacceptheader": true
}
'
```

You should get a 406 error this time, because we have removed the ACCEPT http header based on the body content.

```
* Host localhost:8000 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8000...
* Connected to localhost (::1) port 8000
> POST /mcp/ HTTP/1.1
> Host: localhost:8000
> User-Agent: curl/8.7.1
> Content-Type: application/json
> Accept: application/json, text/event-stream
> Content-Length: 333
> 
* upload completely sent off: 333 bytes
< HTTP/1.1 406 Not Acceptable
< date: Fri, 25 Jul 2025 13:15:00 GMT
< server: envoy
< content-type: application/json
< mcp-session-id: 2e64fbacbe4045ab91390e6bd90890ac
< content-length: 152
< x-envoy-upstream-service-time: 2
< 
{"jsonrpc":"2.0","id":"server-error","error":{"code":-32600,"message":"Not Acceptable: Client must accept both application/json and text/event-stream"}}* Connection #0 to host localhost left intact
```

5. MCP Inspector

One can start the MCP inspector locally and point to the gateway:

```
DANGEROUSLY_OMIT_AUTH=true npx @modelcontextprotocol/inspector
```

In the web browser, change transport to `Streamable HTTP` and use URL: `http://localhost:8000/mcp/`.

This should connect successfully. Other MCP gateway functions are not available in this repo as this project is to specifically 
to test a few upstream changes to make sure the WASM-specific problem is indeed fixed.
