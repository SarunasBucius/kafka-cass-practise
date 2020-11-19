# kafka-cass-practise

Purpose of application is to learn how to use kafka, cassandra and Rest to process data asynchronously.

MVP should be able to:
* Listen to requests
* Produce event to kafka from requests
* Consume event from kafka
* Insert data to cassandra from consumed kafka events

main() should contain 2 goroutines. One for listening to request, another for consuming events.

Simple example of an app:
Create route (for example POST /api/visited). Upon visiting, send event to kafka and return request.
Consume event to insert data to cassandra table (for example table 'visited' with column 'visited_at').
