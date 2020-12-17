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
Create route (for example POST /api/visits). Upon visiting, send event to kafka, containing time request was processed, and return request.
Consume event to insert data to cassandra table (for example table 'visits' with column 'visited_at').

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

Increase use of cassandra db:
* Add GET /api/visits route
	* Returns JSON containing array with fields ip (string) and visited_at (array of strings)
	* Can be filtered using query parameters
		* By visited_at greater than (e.g. ?gt=2020-01)
		* By visited_at less than (e.g. ?lt=2020-01)
		* By visited_at between (e.g. ?gt=2020-01&lt=2020-02)
			* Value must be date yyyy-mm-dd, containing atleast year
		* By day of the week (e.g. ?day=Monday)
* Table visits contains columns:
	* ip - primary key
	* visited_at - cluster key
	* day - secondary index
* Add GET /api/visits/{ip}
	* Returns JSON containing ip and array of visited_at values
	* Supports same filters as /api/visits