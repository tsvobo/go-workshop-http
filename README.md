# go-workshop-http

## Todo list application
Todo list application composed of client and server parts. Server is ready, though there are some missing parts in the client.

### Tasks to do
* TASK-1: Add json tags
* TASK-2: Add validation tags
* TASK-3: Implement task creation using POST request
* TASK-4: Implement task retrieval using GET request
* TASK-5.1: Add prometheus middleware and register request_duration_seconds histogramVec
* TASK-5.2: Implement RoundTripper using promhttp.RoundTripperFunc and record request duration
* TASK-6: Add tracing middleware (zipkinhttp)

### Misc
* To run a zipkin locally: `docker run  -p 9411:9411 openzipkin/zipkin`
