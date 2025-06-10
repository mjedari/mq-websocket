## ToDo list for Websocket Project

1. Monitor all status
   codes. [here](https://kapeli.com/cheat_sheets/WebSocket_Status_Codes.docset/Contents/Resources/Documents/index#:~:text=1001%20indicates%20that%20an%20endpoint,navigated%20away%20from%20a%20page.)

2. Restructuring project
3. Error handling
4. Log kafka related logs into a seperated log
   file [here](https://github.com/confluentinc/librdkafka/issues/3468#issuecomment-1016289114)
5. Add this useful links to project readme
    * [Kafka Docker Quick Start](https://developer.confluent.io/quickstart/kafka-docker/)
    * [Kafka Go Client](https://docs.confluent.io/kafka-clients/go/current/overview.html#ak-go)
    * [Long running jobs](https://medium.com/codex/dealing-with-long-running-jobs-using-apache-kafka-192f053e1691)
    * [7 mistakes of using kafka](https://blog.softwaremill.com/7-mistakes-when-using-apache-kafka-44358cd9cd6)
    * [Configuration of kafka for confluent platform](https://docs.confluent.io/platform/current/installation/configuration/consumer-configs.html#max-partition-fetch-bytes)

### What to do for error:

* I can use a group of consumers (alt least 2 consumers to hand over partitions if the previous one is dead)

### Improvements RFC:

* [X] Give a robust nice structure before refactoring (2)
* [X] Use viper config package to make it more deployable (-)
* [X] Refactor public-channel to use private-channel structure (2)
* [X] Rate limiter for incoming requests (2)
* [X] Private Channel Authentication (2)
* [X] Retry pattern for redis (2)
* [X] Security improvements like, rate limiter, filtered messages and ...

* [ ] keep alive connection feature (ping/pong) (2)
* [ ] Metrics for prometheus (Grafana dashboard and integrate) (3)
* [ ] Circuit breaker (3)
* [ ] Refactor circuit breaker (5)
* [ ] Circuit breaker bug (investigating) (3)
*
* Centralized logging
* Add dedicate local and production config file
* Distributed tracing
* Unit testing after restructuring the project ( ? )
* investigating the potential problems of services (3)
* timeout cloud pattern for producing, consuming kafka
* use context to trance each request got the webserver and log its errors based of an unique identifier
  and use the id in metrics. Ex: I populate the request with uuid and pass it down to the modules, If I wrote
  response on connection calculate the time spend on so far. In gateway this would be more essential. We can use
  correlation id.

### Private Channel:

* We should add timeout mechanize to remove zombie connections through the life of web socket (2)
* We should unsubscribe user form private channel if logged out (2)

### Prometheus' metrics:
* get size of important objects like client, clientMap, hub, rooms
* per endpoint calls, latency  (implement with dimension and labels)
* number of online connections
* number of open rooms
* rate of writing messages (per room / per user / )
* size of allocated memory for rooms
* zombie connections
* number of goroutines
* connections living time
* amount of data is writing on connections (per user/ per endpoint)
* authentication failures
* connection codes monitoring (1001, 1006, ...) count, percentile per all connections to find out user behavior
* runtime configuration states
* memory usage (Alloc, TotalAlloc, Sys, HeapAlloc)
* cpu usage per goroutine and totally
* block IO (read/write)
* network usage (read/write)
* Garbage Collection (PauseTotalNs, NumGC)

### Open Source Project:

* [ ] Create a release branch
* [ ] Write test for critical points
* [ ] Create a UI for Kafka publish message
* [ ] Introduce project in LinkedIn and ask for feedback and contribution
* [ ] Add contribution and license file
* [ ] Checkout develop branch and decide what to do with it
* [ ] Use AI to write README.md
* [ ] Use AI to create a video for the project and how it works