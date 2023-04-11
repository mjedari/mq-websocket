## ToDo list for Websocket Project

1. Monitor all status
codes. [here](https://kapeli.com/cheat_sheets/WebSocket_Status_Codes.docset/Contents/Resources/Documents/index#:~:text=1001%20indicates%20that%20an%20endpoint,navigated%20away%20from%20a%20page.)

2. Restruct project
3. Error handling
4. Log kafka related logs into a seperated log file [here](https://github.com/confluentinc/librdkafka/issues/3468#issuecomment-1016289114)
5. Add this useful links to project readme
   * [Kafka Docker Quick Start](https://developer.confluent.io/quickstart/kafka-docker/)
   * [Kafka Go Client](https://docs.confluent.io/kafka-clients/go/current/overview.html#ak-go)
   * [Long running jobs](https://medium.com/codex/dealing-with-long-running-jobs-using-apache-kafka-192f053e1691)
   * [7 mistakes of using kafka](https://blog.softwaremill.com/7-mistakes-when-using-apache-kafka-44358cd9cd6)



### What to do for error:
* I can use streaming instead of consumer
* I can use a group of consumers (alt least 2 consumers to hand over partitions if the previous one is dead)
* 