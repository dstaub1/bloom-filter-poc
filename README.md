# bloom-filter-poc

Epic: https://roadie.atlassian.net/browse/BE2-3430

Currently there are approximately 7.7M Kafka messages stored in the roadie-racecar-consumer-logs postgres table. This is a seasonal number and during peak season would be higher.
Messages are deleted weekly on a rolling basis. That is, a message stored in the postgres db is deleted after 7 days.
Eventually, we want to accommodate ~100x more volume. Estimate 700M - 1B Kafka messages at a given time. This is on the order of 10^8 - 10^9 messages during a week.

Coupled with this effort, we want to investigate consumer-specific delivery configurations. 
For example, our tracking points consumer will generally be ok if we miss a couple of messages, so we would set the configuration to at-most-once delivery.
For other consumers like the payment events consumer, we don't want to miss any messages, so we would set the configuration to at-least-once delivery.

Finally, we are evaluate the use of a Bloom filter coupled with Redis to alleviate some of the throughput issues we see in maintaining a massive Postgres table which is currently configured for 1 read and 1 write per Kafka message. As we scale, the limitations of the postgres db will become more apparent.
