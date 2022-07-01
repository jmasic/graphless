# Graphless driver for Graphalytics 


## Graphless-specific configuration
The class `science.atlarge.graphalytics.graphless.configuration.payload.PayloadWriter` contains some working defaults for a local setup.

These can be adapted as follows:
* `PayloadWriter#MAX_WORKERS`: maximum number of workers that can be active during the processing phase
* `PayloadWriter::getPayloadContent` has some default values for memory, message, and storage client configuration. More details on these configuration parameters can be found in the `README.md` file of Graphless A5


## Evaluation
For evaluating Graphless on Amazon AWS, run the following command from the Graphalytics Graphless driver folder:
```bash
$ ./init.sh $path_to_graph aws compile
```

Alternatively, for evaluating Graphless on Amazon AWS, run the following command from the Graphalytics Graphless driver folder:
```bash
$ ./init.sh $path_to_graph local compile
```
