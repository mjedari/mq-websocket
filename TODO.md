## ToDo list for Websocket Project

1. Monitor all status
codes. [here](https://kapeli.com/cheat_sheets/WebSocket_Status_Codes.docset/Contents/Resources/Documents/index#:~:text=1001%20indicates%20that%20an%20endpoint,navigated%20away%20from%20a%20page.)

2. Restructuring project
3. Error handling
4. Log kafka related logs into a seperated log file [here](https://github.com/confluentinc/librdkafka/issues/3468#issuecomment-1016289114)
5. Add this useful links to project readme
   * [Kafka Docker Quick Start](https://developer.confluent.io/quickstart/kafka-docker/)
   * [Kafka Go Client](https://docs.confluent.io/kafka-clients/go/current/overview.html#ak-go)
   * [Long running jobs](https://medium.com/codex/dealing-with-long-running-jobs-using-apache-kafka-192f053e1691)
   * [7 mistakes of using kafka](https://blog.softwaremill.com/7-mistakes-when-using-apache-kafka-44358cd9cd6)
   * [Configuration of kafka for confluent platform](https://docs.confluent.io/platform/current/installation/configuration/consumer-configs.html#max-partition-fetch-bytes)



### What to do for error:
* I can use a group of consumers (alt least 2 consumers to hand over partitions if the previous one is dead)


### Improvements RFC:
* Give a robust nice structure before refactoring (2)
* Use viper config package to make it more deployable (-)
* Refactor public-channel to use private-channel structure (2)
* Security improvements like, rate limiter, filtered messages and ... 
* Using cloud patterns like retry pattern and request tracing (2 for retry)
* Metrics for prometheus (Grafana dashboard and integrate) (3)
* Centralized logging
* Distributed tracing
* Unit testing after restructuring the project ( ? )
* investigating the potential problems of services (3)

### Private Channel:
* We should add timeout mechanize to remove zombie connections through the life of web socket (2)
