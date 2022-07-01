package science.atlarge.graphalytics.graphless.algorithms.params;

import science.atlarge.graphalytics.domain.graph.LoadedGraph;
import science.atlarge.graphalytics.execution.RunSpecification;

import java.nio.file.Path;

/**
 * POJO holding job parameters for each graphless run.
 *
 * @author jmasic
 */
public class GraphlessJobParams {

    private final Path graphlessDirectory;
    private final LoadedGraph graph;
    private final GraphlessPlatform platform;
    private final String runId;


    public GraphlessJobParams(Path graphlessDirectory, LoadedGraph graph, GraphlessPlatform platform, RunSpecification runSpecification) {
        this.graphlessDirectory = graphlessDirectory;
        this.graph = graph;
        this.platform = platform;
        this.runId = runSpecification.getBenchmarkRun().getId();
    }

    public Path getGraphlessDirectory() {
        return graphlessDirectory;
    }

    public long getAmountOfWork() {
        return graph.getFormattedGraph().getNumberOfVertices();
    }

    public boolean isDirected() {
        return graph.getFormattedGraph().isDirected();
    }

    public GraphlessPlatform getPlatform() {
        return platform;
    }

    public String getRunId() {
        return runId;
    }
}
