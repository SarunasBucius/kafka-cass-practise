# kafka-cass-practise

Purpose of application is to learn how to use kafka, cassandra and Rest to process data asynchronously.

MVP should be able to:
* Listen to requests
* Produce event to kafka from requests
	* event should contain time at which request was processed
* Consume event from kafka
* Insert data to cassandra from consumed kafka events

main() should contain 2 goroutines. One for listening to request, another for consuming events.

Simple example of an app:
Create route (for example POST /api/visited). Upon visiting, send event to kafka, containing time request was processed, and return request.
Consume event to insert data to cassandra table (for example table 'visited' with column 'visited_at').

Increase use of kafka:
* Increase number of partitions to 2
* Use 2 consumer groups
	* Insert data to db, group.id=inserter
		* Use 2 consumers
	* Print day of the week, group.id=day
		* Use bulk consuming and parallelize returned events handling
* Use acks=1
* Change event to struct Visit with values visitedAt and ip
	* Use ip as event key