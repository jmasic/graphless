package science.atlarge.graphalytics.graphless.algorithms;

import it.unimi.dsi.fastutil.longs.Long2DoubleMap;
import it.unimi.dsi.fastutil.longs.Long2DoubleOpenHashMap;
import it.unimi.dsi.fastutil.longs.Long2LongMap;
import it.unimi.dsi.fastutil.longs.Long2LongOpenHashMap;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import science.atlarge.graphalytics.graphless.algorithms.params.GraphlessJobParams;
import science.atlarge.graphalytics.graphless.configuration.payload.PayloadWriter;
import science.atlarge.graphalytics.graphless_Aws.Bucket;
import science.atlarge.graphalytics.util.ProcessUtil;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.Optional;
import java.util.concurrent.TimeUnit;

/**
 * Generic class holding shared logic for jobs running graph algorithms on Graphless.
 *
 * @author jmasic
 */
public abstract class AlgorithmJob {

    private static final Logger LOG = LogManager.getLogger();

    private static final String METADATA_FILE_SUFFIX = "-results-metadata.json";
    private static final String RESULTS_FILE_SUFFIX = "-results";

    public final Map<Long, ?> run(GraphlessJobParams jobParams) {
        LOG.info("Starting {} algorithm", getAlgorithmName());
        Path graphlessDirectory = jobParams.getGraphlessDirectory();

        final List<String> resultLines;
        final long loadingTime;
        final long executionTime;
        switch (jobParams.getPlatform()) {
            case AWS -> {
                writeInputPayload(graphlessDirectory, jobParams, "main_payload.json");
                runWithGraphlessOnAws(jobParams);
                waitForSeconds(30);
                resultLines = readAwsResult();
                loadingTime = extractLoadingTimeFromAwsResults();
                executionTime = extractExecutionTimeFromAwsResults();
            }
            case LOCAL -> {
                writeInputPayload(graphlessDirectory, jobParams, "local_payload.json");
                runWithGraphlessLocally(jobParams);
                resultLines = readLocalResult(jobParams);
                loadingTime = extractLoadingTimeFromLocalResults(jobParams);
                executionTime = extractExecutionTimeFromLocalResults(jobParams);
            }
            default -> throw new IllegalArgumentException("Don't know how to handle job params of class " + jobParams.getClass().getCanonicalName());
        }
        Map<Long, ?> result = parseResult(resultLines);

        LOG.info("Finished {} algorithm, with loading taking {} and execution taking {}",
                 getAlgorithmName(), loadingTime, executionTime);
        double executionTimeInSeconds = ((double) (executionTime / 1_000_000)) / 1000.0;
        final long makeSpan = loadingTime + (long) executionTimeInSeconds;
        LOG.info("Algorithm={};Execution_Time={};Execution_Time_Seconds={};Loading_Time={};Make_Span={};",
                 getAlgorithmName(), executionTime, executionTimeInSeconds, loadingTime, makeSpan);
        return result;
    }

    private void writeInputPayload(Path graphlessDirectory, GraphlessJobParams jobParams, String payloadFileName) {
        // main_payload.json
        Path payloadPath = Paths.get(graphlessDirectory.toAbsolutePath().toString(), payloadFileName);
        long amountOfWork = jobParams.getAmountOfWork();
        LOG.info("Amount of work: {}", amountOfWork);
        new PayloadWriter(
                payloadPath,
                "dota-league", // FIXME: Extract from somewhere, or should be cleaned up
                getExtraArgs(jobParams),
                getAlgorithmName(),
                amountOfWork
        ).writePayload(jobParams.getRunId());
    }


    private void runWithGraphlessLocally(GraphlessJobParams localJobParams) {
        executeShellCommand(
                localJobParams,
                new String[]{"./bin-local/main_function/main_function", "local"}
        );
    }
    private void runWithGraphlessOnAws(GraphlessJobParams awsJobParams) {
        executeShellCommand(
                awsJobParams,
                new String[]{"bash", "./start.sh", "--payload", "main_payload.json"}
        );
    }
    private void executeShellCommand(GraphlessJobParams jobParams, String[] command) {
        ProcessBuilder pb = new ProcessBuilder(command);
        pb.directory(jobParams.getGraphlessDirectory().toFile());
        LOG.info("Executing {} in directory '{}'...", getAlgorithmName(), jobParams.getGraphlessDirectory());
        try {
            final Process process = pb.start();
            ProcessUtil.monitorProcess(process, "---"); // FIXME: Take run id from somewhere
            int rc = process.waitFor();
            LOG.info("Process exited with {} for algorithm {}", rc, getAlgorithmName());
        } catch (IOException | InterruptedException e) {
            throw new RuntimeException(e);
        }
    }


    private List<String> readLocalResult(GraphlessJobParams localJobParams) {
        final String suffix = getAlgorithmName() + RESULTS_FILE_SUFFIX;
        return getLocalFileContentFromSuffix(localJobParams, suffix);
    }

    private long extractLoadingTimeFromLocalResults(GraphlessJobParams localJobParams) {
        final String suffix = getAlgorithmName() + METADATA_FILE_SUFFIX;
        List<String> fileLines = getLocalFileContentFromSuffix(localJobParams, suffix);
        String metadata = fileLines.get(0);
        return extractLoadingTimeFromMetadata(metadata);
    }

    private long extractExecutionTimeFromLocalResults(GraphlessJobParams localJobParams) {
        final String suffix = getAlgorithmName() + METADATA_FILE_SUFFIX;
        List<String> fileLines = getLocalFileContentFromSuffix(localJobParams, suffix);
        String metadata = fileLines.get(0);
        return extractExecutionTimeFromMetadata(metadata);
    }

    private List<String> getLocalFileContentFromSuffix(GraphlessJobParams localJobParams, String suffix) {
        Path resultFolder = getResultFolder(localJobParams);
        File[] matchingFiles = resultFolder.toFile().listFiles(
                (dir, name) -> name.startsWith("graphFileKey") && name.endsWith(suffix)
        );
        int matchingFilesCount = Objects.requireNonNull(matchingFiles).length;
        if (matchingFilesCount != 1) {
            throw new IllegalStateException("Only one result file should be there, found " + matchingFilesCount);
        }
        try {
            return Files.readAllLines(matchingFiles[0].toPath());
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }

    private static Path getResultFolder(GraphlessJobParams localJobParams) {
        return Paths.get(localJobParams.getGraphlessDirectory().toAbsolutePath().toString(),
                         "..",
                         "simple-graphs",
                         "graphalytics",
                         "s3");
    }


    private List<String> readAwsResult() {
        while (true) {
            final String suffix = getAlgorithmName() + RESULTS_FILE_SUFFIX;
            Optional<List<String>> maybeLines = getAwsFileContentFromSuffix(suffix);
            if (maybeLines.isEmpty()) {
                continue;
            }
            return maybeLines.get();
        }
    }

    private long extractLoadingTimeFromAwsResults() {
        while (true) {
            final String suffix = getAlgorithmName() + METADATA_FILE_SUFFIX;
            Optional<List<String>> maybeLines = getAwsFileContentFromSuffix(suffix);
            if (maybeLines.isEmpty()) {
                continue;
            }
            String metadata = maybeLines.get().get(0);
            return extractLoadingTimeFromMetadata(metadata);
        }
    }

    private long extractExecutionTimeFromAwsResults() {
        while (true) {
            final String suffix = getAlgorithmName() + METADATA_FILE_SUFFIX;
            Optional<List<String>> maybeLines = getAwsFileContentFromSuffix(suffix);
            if (maybeLines.isEmpty()) {
                continue;
            }
            String metadata = maybeLines.get().get(0);
            return extractExecutionTimeFromMetadata(metadata);
        }
    }

    private Optional<List<String>> getAwsFileContentFromSuffix(String suffix) {
        Bucket bucket = new Bucket();
        List<String> keys = bucket.list();
        List<String> matchingKeys = keys.stream()
                .filter(key -> key.startsWith("graphFileKey") && key.endsWith(suffix))
                .toList();

        int matchingFilesCount = matchingKeys.size();
        if (matchingFilesCount == 0) {
            waitForSeconds(5);
            return Optional.empty();
        }
        if (matchingFilesCount > 1) {
            throw new IllegalStateException("Only one result file should be there, found " + matchingFilesCount);
        }

        return Optional.of(bucket.getObjectLines(matchingKeys.get(0)));
    }

    private void waitForSeconds(int timeout) {
        try {
            TimeUnit.SECONDS.sleep(timeout);
        } catch (InterruptedException e) {
            throw new RuntimeException(e);
        }
    }

    private long extractLoadingTimeFromMetadata(String metadata) {
        return Long.parseLong(metadata.split("dataIngestionDuration\":")[1].split(",")[0]);
    }

    private long extractExecutionTimeFromMetadata(String metadata) {
        return Long.parseLong(metadata.split("executionDuration\":")[1].split(",")[0]);
    }


    private Map<Long, ?> parseResult(List<String> resultLines) {
        final var resultType = getResultType();
        switch (resultType) {
            case LONG_TO_LONG -> {
                Long2LongMap distances = new Long2LongOpenHashMap();
                for (String line : resultLines) {
                    String[] tokens = line.split(" ");
                    distances.put(Long.valueOf(tokens[0]), Long.valueOf(tokens[1]));
                }
                return distances;
            }
            case LONG_TO_DOUBLE -> {
                Long2DoubleMap distances = new Long2DoubleOpenHashMap();
                for (String line : resultLines) {
                    String[] tokens = line.split(" ");
                    distances.put(Long.valueOf(tokens[0]), parseDouble(tokens[1]));
                }
                return distances;
            }
            default -> throw new UnsupportedOperationException("Unsupported result type " + resultType);
        }
    }

    private Double parseDouble(String token) {
        if (token.equals("infinity")) {
            return Double.POSITIVE_INFINITY;
        }
        return Double.valueOf(token);
    }


    protected abstract String getExtraArgs(GraphlessJobParams jobParams);
    protected abstract String getAlgorithmName();
    protected abstract ResultType getResultType();

    protected enum ResultType {
        LONG_TO_DOUBLE,
        LONG_TO_LONG
    }
}
