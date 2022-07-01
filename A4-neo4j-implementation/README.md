# Graphless

## A4 -  Improved Neo4j version 
In this implementation we modify the Neo4j version present in A3 implementation by providing:
* an improved Neo4j implementation as we use optimized queries for Neo4j (such as, `UNWIND`)
* barrier-based using during the loading phase


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

### Configuration
In this version of Graphless, differently from A5, the MaaS backend is chosen programmatically for local or AWS runs.
Therefore, to specify the backend used for each client you might need to do some code changes.

In `main_function.go`, `orchestrator_function.go`, and `worker_function.go`, set the following variables according to the backend you choose:
* `memoryClientType`: backend for MaaS
* `storageClientType`: backend for S3 bucket
