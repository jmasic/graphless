package science.atlarge.graphalytics.graphless_Aws;

import org.apache.commons.io.output.TeeOutputStream;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import science.atlarge.graphalytics.configuration.BuildInformation;
import science.atlarge.graphalytics.domain.algorithms.Algorithm;
import science.atlarge.graphalytics.domain.algorithms.BreadthFirstSearchParameters;
import science.atlarge.graphalytics.domain.algorithms.CommunityDetectionLPParameters;
import science.atlarge.graphalytics.domain.algorithms.PageRankParameters;
import science.atlarge.graphalytics.domain.algorithms.SingleSourceShortestPathsParameters;
import science.atlarge.graphalytics.domain.benchmark.BenchmarkRun;
import science.atlarge.graphalytics.domain.graph.FormattedGraph;
import science.atlarge.graphalytics.domain.graph.LoadedGraph;
import science.atlarge.graphalytics.execution.BenchmarkRunSetup;
import science.atlarge.graphalytics.execution.Platform;
import science.atlarge.graphalytics.execution.PlatformExecutionException;
import science.atlarge.graphalytics.execution.RunSpecification;
import science.atlarge.graphalytics.execution.RuntimeSetup;
import science.atlarge.graphalytics.graphless.algorithms.AlgorithmJob;
import science.atlarge.graphalytics.graphless.algorithms.BreadthFirstSearchJob;
import science.atlarge.graphalytics.graphless.algorithms.CommunityDetectionLPJob;
import science.atlarge.graphalytics.graphless.algorithms.LocalClusteringCoefficientJob;
import science.atlarge.graphalytics.graphless.algorithms.PageRankJob;
import science.atlarge.graphalytics.graphless.algorithms.SingleSourceShortestPathJob;
import science.atlarge.graphalytics.graphless.algorithms.WeaklyConnectedComponentsJob;
import science.atlarge.graphalytics.graphless.algorithms.params.GraphlessJobParams;
import science.atlarge.graphalytics.graphless.algorithms.params.GraphlessPlatform;
import science.atlarge.graphalytics.report.result.BenchmarkMetric;
import science.atlarge.graphalytics.report.result.BenchmarkMetrics;

import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.PrintStream;
import java.io.PrintWriter;
import java.math.BigDecimal;
import java.math.RoundingMode;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.List;
import java.util.Map;
import java.util.Properties;

import static java.nio.file.Files.readAllBytes;

/**
 * Implementation for running the Graphalytics benchmark against a instance of Graphless deployed on Amazon AWS.
 *
 * @author jmasic
 */
public class Graphless_AwsPlatform implements Platform {

    private static final Logger LOG = LogManager.getLogger();
    private static PrintStream sysOut;
    private static PrintStream sysErr;

    @Override
    public void verifySetup() {}

    @Override
    public LoadedGraph loadGraph(FormattedGraph formattedGraph) {
        // we assume the graph is already loaded in S3 (TODO: Upload the file at this step)
        return new LoadedGraph(
                formattedGraph,
                formattedGraph.getVertexFilePath(),
                formattedGraph.getEdgeFilePath()
        );
    }

    private Path resolveGraphlessDirectory() {
        final var properties = getProperties();
        final var graphlessDirectoryName = properties.getProperty("graphless.directory");
        try {
            return Paths.get(graphlessDirectoryName).toRealPath();
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }


    @Override
    public void prepare(RunSpecification runSpecification) {
        Bucket bucket = new Bucket();
        List<String> keys = bucket.list();
        LOG.info("Found {} objects in buckets", keys.size());
        for (String key : keys) {
            if (key.contains("results")) {
                bucket.delete(key);
                LOG.info("Deleted object '{}' in bucket", key);
            }
        }
    }


    @Override
    public void startup(RunSpecification runSpecification) {
        BenchmarkRunSetup benchmarkRunSetup = runSpecification.getBenchmarkRunSetup();
        startBenchmarkLogging(benchmarkRunSetup.getLogDir().resolve("platform").resolve("driver.logs"));
    }

    @SuppressWarnings("unchecked")
    @Override
    public void run(RunSpecification runSpecification) throws PlatformExecutionException {
        BenchmarkRun benchmarkRun = runSpecification.getBenchmarkRun();
        BenchmarkRunSetup benchmarkRunSetup = runSpecification.getBenchmarkRunSetup();
        RuntimeSetup runtimeSetup = runSpecification.getRuntimeSetup();

        Algorithm algorithm = benchmarkRun.getAlgorithm();
        Object parameters = benchmarkRun.getAlgorithmParameters();
        Map<Long, ?> output;

        LOG.info("Processing starts at: {}", System.currentTimeMillis());
        LoadedGraph graph = runtimeSetup.getLoadedGraph();
        GraphlessJobParams jobParams = new GraphlessJobParams(resolveGraphlessDirectory(), graph, GraphlessPlatform.AWS, runSpecification);
        AlgorithmJob job = switch (algorithm) {
            case BFS -> new BreadthFirstSearchJob((BreadthFirstSearchParameters) parameters);
            case CDLP -> new CommunityDetectionLPJob((CommunityDetectionLPParameters) parameters);
            case WCC -> new WeaklyConnectedComponentsJob();
            case PR -> new PageRankJob((PageRankParameters) parameters);
            case LCC -> new LocalClusteringCoefficientJob();
            case SSSP -> new SingleSourceShortestPathJob((SingleSourceShortestPathsParameters) parameters);
            default -> throw new PlatformExecutionException("Unsupported algorithm: " + algorithm);
        };
        output = job.run(jobParams);

        if (benchmarkRunSetup.isOutputRequired()) {
            try {
                String outputFile = benchmarkRunSetup.getOutputDir().resolve(benchmarkRun.getName()).toAbsolutePath().toString();
                writeOutput(outputFile, output);
            } catch(IOException e) {
                throw new PlatformExecutionException("An error while writing to output file", e);
            }
        }
        LOG.info("Processing ends at: " + System.currentTimeMillis());
    }

    @Override
    public BenchmarkMetrics finalize(RunSpecification runSpecification) {
        stopPlatformLogging();
        BenchmarkRunSetup benchmarkRunSetup = runSpecification.getBenchmarkRunSetup();
        Path path = benchmarkRunSetup.getLogDir().resolve("platform").resolve("driver.logs");
        final String logs;
        try {
            logs = new String(readAllBytes(path));
        } catch (IOException e) {
            e.printStackTrace();
            throw new IllegalStateException("Can't read file at " + path);
        }

        Long startTime = null;
        Long endTime = null;

        for (String line : logs.split("\n")) {
            try {
                if (line.contains("Processing starts at: ")) {
                    String[] lineParts = line.split("\\s+");
                    startTime = Long.parseLong(lineParts[lineParts.length - 1]);
                }

                if (line.contains("Processing ends at: ")) {
                    String[] lineParts = line.split("\\s+");
                    endTime = Long.parseLong(lineParts[lineParts.length - 1]);
                }
            } catch (Exception e) {
                LOG.error("Cannot parse line: {}", line, e);
            }

        }

        if(startTime != null && endTime != null) {
            BenchmarkMetrics metrics = new BenchmarkMetrics();
            long procTimeMS =  endTime - startTime;
            BigDecimal procTimeS = (new BigDecimal(procTimeMS)).divide(new BigDecimal(1000), 3, RoundingMode.CEILING);
            metrics.setProcessingTime(new BenchmarkMetric(procTimeS, "s"));

            return metrics;
        } else {
            return new BenchmarkMetrics();
        }
    }

    @Override
    public void terminate(RunSpecification runSpecification) {

    }

    @Override
    public void deleteGraph(LoadedGraph loadedGraph) {}


    private static void startBenchmarkLogging(Path fileName) {
        sysOut = System.out;
        sysErr = System.err;
        try {
            final File file = fileName.toFile();
            file.getParentFile().mkdirs();
            file.createNewFile();
            FileOutputStream fos = new FileOutputStream(file);
            TeeOutputStream bothStream =new TeeOutputStream(System.out, fos);
            PrintStream ps = new PrintStream(bothStream);
            System.setOut(ps);
            System.setErr(ps);
        } catch(Exception e) {
            e.printStackTrace();
            throw new IllegalArgumentException("cannot redirect to output file");
        }
    }

    public static void stopPlatformLogging() {
        System.setOut(sysOut);
        System.setErr(sysErr);
    }

    private void writeOutput(String path, Map<Long, ? extends Object> output) throws IOException {
        try (PrintWriter w = new PrintWriter(new FileOutputStream(path))) {
            for (Map.Entry<Long, ? extends Object> entry: output.entrySet()) {
                w.print(entry.getKey());
                w.print(" ");
                w.print(entry.getValue());
                w.println();
            }
        }
    }


    private static Properties getProperties() {
        final var buildInfoFile = "/project/build/platform.properties";
        try {
            return BuildInformation.loadBuildPropertiesFile(buildInfoFile);
        } catch (Exception e) {
            throw new IllegalStateException(String.format("Failed to load platform name from %s.", buildInfoFile));
        }
    }

    @Override
    public String getPlatformName() {
        return "graphless_local";
    }
}
