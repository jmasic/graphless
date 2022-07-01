/*
 * Copyright 2015 Delft University of Technology
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package science.atlarge.graphalytics.graphless.algorithms;

import science.atlarge.graphalytics.graphless.algorithms.params.GraphlessJobParams;

public class LocalClusteringCoefficientJob extends AlgorithmJob {

	@Override
	protected String getExtraArgs(GraphlessJobParams jobParams) {
		return "{\"directed\": " + jobParams.isDirected() + "}";
	}

	@Override
	protected String getAlgorithmName() {
		return "LCC";
	}

	@Override
	protected ResultType getResultType() {
		return ResultType.LONG_TO_DOUBLE;
	}
}
