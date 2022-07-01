# Graphless

## A5 - Loading based on work queue and value-reset
In this implementation, we provide a Graphless version which:
* supports backend configuration via JSON input
* split between Graph MaaS (`clients/memory`) and Message MaaS (`clients/message`)
* uses throttling in both `orchestrator_function.go` for the worker fan-out
* provides, as a loading mechanism, both a work queue and value reset (see "Loading based on value reset")
* fully supports local runs


### Configuration
Some JSON configuration is required to run Graphless.
For local runs, this goes inside `local_payload.json`.
For runs on Amazon AWS, this goes inside `main_payload.json`.
A typical configuration file has this format:
```json
{
  "graphName": <graphName:string>,
  "maxWorkers": <maxWorkers:int>,
  "extraArgs": <extraArgs:object>,
  "chunkSize": <chunkSize:int>,
  "runId": <runId:string>,
  "testRun": <testRun:boolean>,
  "algorithm": <algorithm:string>,
  "memoryClientConfig": {
    "type": <memory.type:string>,
    "db": {
      "ip": <memory.db.ip:string>,
      "port": <memory.db.port:int>,
      "username": <memory.db.username:string>,
      "password":  <memory.db.password:string>,
      "shardsCount":  <memory.db.shardsCount:int>
    }
  },
  "messageClientConfig": {
    "type": <message.type:string>,
    "db": {
      "ip": <message.db.ip:string>,
      "port": <message.db.port:int>,
      "username": <message.db.username:string>,
      "password":  <message.db.password:string>,
      "shardsCount":  <message.db.shardsCount:int>
    }
  },
  "storageClientConfig": {
    "type": <storage.type:string>,
    "storageConfig": {
      "bucketName": <storage.config.bucketName:string>,
      "bucketKey": <storage.config.bucketKey:string>,
      "region": <storage.config.region:string>
    }
  }
}
```

Where the keys have the following meaning:
* `<graphName:string>`: name of the graph being processed; example: `"dota-league"`
* `<maxWorkers:int>`: maximum number of workers used during the processing phase; example: `400`
* `<extraArgs:object>`: arguments needed for the execution of a specific algorithm; example (for BFS): `{"sourceVertex":1}`
* `<chunkSize:int>`: maximum number of vertices allocated to each worker for a superstep. Make sure `maxWorkers * chunkSize >= |V|`; example: `153`
* `<runId:string>`: identifier of a Graphless execution. Useful for identifying results; example: `"r837924"`
* `<testRun:boolean>`: deprecated; example: `false`
* `<algorithm:string>`: name of the algorithm being executed. Make sure that this algorithm is supported in `worker.go`; example: `"BFS"`
* `<memory.type:string>`: backend used for Graph MaaS. Check the switch in `memory/constants.go` for seeing the alternatives; example: `"Neo4j"`
* `<memory.db.ip:string>`: IP address at which Graph MaaS is reachable; example: `"127.0.0.1"`
* `<memory.db.port:int>`: port on which Graph MaaS listens; example: `7687`
* `<memory.db.username:string>`: username required to access Graph MaaS (optional); example: `"neo4j"`
* `<memory.db.password:string>`: password required to access Graph MaaS; example: `"n"`
* `<memory.db.shardsCount:int>`: count of shards that make up Graph MaaS; example: `1`
* `<message.type:string>`: backend used for Message MaaS. Check the switch in `message/constants.go` for seeing the alternatives; example: `"Redis"`
* `<message.db.ip:string>`: IP address at which Message MaaS is reachable; example: `"127.0.0.1"`
* `<message.db.port:int>`: port on which Message MaaS listens; example: `6379`
* `<message.db.username:string>`: username required to access Message MaaS (optional); example: `""`
* `<message.db.password:string>`: password required to access Graph MaaS; example: `"p4$$w0rD"`
* `<message.db.shardsCount:int>`: count of shards that make up Message MaaS; example: `1`
* `<storage.type:string>`: backend used for the storage layer. Check the switch in `storage/constants.go` for seeing the alternatives; example: `"Local"`
* `<storage.config.bucketName:string>`: name of the bucket on which graph data is stored. On a local file system, the folder in which graph data is located; example: `"~/graphs/example-directed/"`
* `<storage.config.bucketKey:string>`: name of the key prefix all graph data chunks; example:  `"graphFileKey"`
* `<storage.config.region:string>`: if running on cloud, region on which storage is located (optional); example: `"us-east2"`



### Loading based on value reset
This version contains both the loading solution based on **work queue** and **value reset**.
By default, the work queue implementation is used.
To switch to the value-reset solution, just:
* comment the current implementation of `main_function.loadGraphInMemory` and uncomment the part after `// NOTE: value-reset implementation`
* comment the current implementation of `neo4jClient.CreateVertices` and uncomment the part after `// NOTE: value-reset implementation`
* comment the current implementation of `neo4jClient.Clear` and uncomment the part after `// NOTE: value-reset implementation`

Please notice that you will need to load the graph at least once for the value-reset approach to work.
Further implementation to avoid this edge case is left as future work.

### MacOS installation
```bash
$ make clean-local build-local
```

### Linux installation
```bash
$ make clean build
```

### Amazon installation
* To deploy the cloud stack (Amazon Lambda) on AWS use the following:
    ```bash
    $ bash deploy.sh -t false
    ```
* Run on Amazon EC2 the following scripts depending on the technology used as MaaS:
    * neo4j_script.sh to install and configure Neo4j:
      ```bash
       $ ./neo4j_script.sh
      ```
    * redis_script.sh to install Redis:
      ```bash
       $ ./redis_script.sh 16 7379 redis
       $ ./redis_script.sh 16 7379 start
      ```

### Local runs
We run local experiments on a MacBook, with the MaaS components executing as Docker containers.
To execute the program locally, make sure to pass the `local` executing to the `main_function`.

We run Neo4j containers on Docker with:
```bash
$ docker run --rm --name neo4j \
        -p 7687:7687 -p 7474:7474 \
         --env NEO4J_AUTH=neo4j/n \
        neo4j:4.4.6
```

We run Redis containers on Docker with:
```bash
$ docker run --rm --name redis \
        -p 6379:6379 \
        redis:6.2.6-alpine3.15
```

To start Graphless locally run the following command:
```bash
$ ./bin-local/main_function/main_function local
```
