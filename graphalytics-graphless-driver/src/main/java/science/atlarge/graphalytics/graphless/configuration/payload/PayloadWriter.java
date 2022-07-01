package science.atlarge.graphalytics.graphless.configuration.payload;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardOpenOption;

/**
 * Configuration generator for Graphless.
 *
 * @author jmasic
 */
public class PayloadWriter {

    private static final int MAX_WORKERS = 400;

    private final Path payloadPath;
    private final String graphName;
    private final String extraArgs;
    private final String algorithm;
    private final long amountOfWork;

    public PayloadWriter(Path payloadPath, String graphName, String extraArgs, String algorithm, long amountOfWork) {
        this.payloadPath = payloadPath;
        this.graphName = graphName;
        this.extraArgs = extraArgs;
        this.algorithm = algorithm;
        this.amountOfWork = amountOfWork;
    }


    public final void writePayload(String runId) {
        try {
            String payloadContent = getPayloadContent(runId);
            Files.write(payloadPath, payloadContent.getBytes(), StandardOpenOption.TRUNCATE_EXISTING);
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }

    private String getPayloadContent(String runId) {
        // share of work that each worker takes
        long chunkSize = Math.max(30, amountOfWork / MAX_WORKERS + 1);

        return """
                {
                  "graphName": "%s",
                  "maxWorkers": %d,
                  "extraArgs": %s,
                  "chunkSize": %d,
                  "runId": "%s",
                  "testRun": false,
                  "algorithm": "%s",
                  "memoryClientConfig": {
                    "type": "Neo4j",
                    "db": {
                      "ip": "127.0.0.1",
                      "port": 7687,
                      "username": "neo4j",
                      "password": "n",
                      "shardsCount": 1
                    }
                  },
                  "messageClientConfig": {
                    "type": "Redis",
                    "db": {
                      "ip": "127.0.0.1",
                      "port": 6379,
                      "username": "",
                      "password": "",
                      "shardsCount": 1
                    }
                  },
                  "storageClientConfig": {
                    "type": "Local",
                    "storageConfig": {
                      "bucketName": "../simple-graphs/graphalytics/s3/",
                      "bucketKey": "graphFileKey",
                      "region": "n/a"
                    }
                  }
                }
                """.formatted(graphName, MAX_WORKERS, extraArgs, chunkSize, runId, algorithm);
    }
}
