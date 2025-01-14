# gRPC

gRPC is a API(rpc) to exchange data, a error format and API to communicate for microservices, which is developed by Google and now became a part of Cloud Native.

Building an API is hard, we need to think about:

- **Payload Size**
  The amount of data sent over the network. Smaller payloads reduce bandwidth usage and improve performance.
  
- **Latency**
   The time it takes for a request to travel from the client to the server and back. Lower latency improves user experience.

- **Scalability**
  The ability of the system to handle increased load. Scalable APIs can serve more clients without degradation in performance.

- **Load Balancing**
  Distributing incoming requests across multiple servers to ensure no single server is overwhelmed.

- **Languages Interop**
  gRPC supports multiple programming languages (e.g., Java, Python, Go, C++, etc.), enabling seamless communication between services written in different languages. This is crucial for microservices architectures where different teams may use different tech stacks.
  
- **Auth, Monitoring, Log**
  - **Auth** gRPC provides built-in support for SSL/TLS and token-based authentication, ensuring secure communication between services.
  - **Monitoring** gRPC integrates with tools like Prometheus and OpenTelemetry to monitor service health, performance, and errors.
  - **Log** Structured logging can be implemented to track requests, responses, and errors, making debugging and auditing easier.
  
---

## Communicate

**gRPC Server <=(RPC)=> gRPC Stub Client(Any Language)**  
gRPC Server communicates bidirectionally with gRPC Stub Client via RPC. The client can be implemented in any programming language supported by gRPC (e.g., Java, Python, Go, C++, etc.).  

---

### why gRPC not JSON

JSON takes up more memory space and consumes more cpu resources than gRPC.
