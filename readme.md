Go Security 

Learning Golang by building a fully fresh SIEM/SOAR combo with built in event streaming rule detection and response.

Where we are now:
security.go is responsible for spinning things up and running them. 
Within that File we have the Configuration loading to start, this loads two configs:
- config.yaml - Responsible for the database configuration
- sources.yaml - Responsible for the file source configurations

Once the Configs are loaded we can load up the Database, right now we use Elastic, mainly because I wanted to store documents of unstructured type and also have the cabaility to use SQL queries if I needed those.
The way the system works though is we can add more databases quiet easily with the Database Interface.

Once the DB is initialized, we load up the indexer. The Open function opens the channel for new events. the Store function stores values that come through its channel into the DB.

So essentially you can send events from files to the indexer, it will parse them and store them in the Database.
We have querying and stuff already built too.

The next thing I'm building is an alert mechanism, so you can define alerts via YAML files which will execute on a schedule (probably) and send data (eventually) if their conditions trigger (AKA a SIEM alert).

Think thats all I got so far. (10/3/2024)

Things to do:
- Use AI to generate the YAML Alert files based on the DB values
- Use AI to generate the queries for the DB to gen the results
- Build testing for using the AI alert and have it write the file if approved


docker run -d --name elasticsearch --net elastic -p 9200:9200 -e "discovery.type=single-node" elasticsearch:8.9.0

https://www.elastic.co/guide/en/elasticsearch/reference/8.15/docker.html
https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/connecting.html
https://docs.go-blueprint.dev/blueprint-core/db-drivers/
https://go.dev/doc/tutorial/database-access