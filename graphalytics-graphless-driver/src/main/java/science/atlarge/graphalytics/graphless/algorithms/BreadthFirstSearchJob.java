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

import science.atlarge.graphalytics.domain.algorithms.BreadthFirstSearchParameters;
import science.atlarge.graphalytics.graphless.algorithms.params.GraphlessJobParams;


public class BreadthFirstSearchJob extends AlgorithmJob {

	private final BreadthFirstSearchParameters parameters;

	public BreadthFirstSearchJob(BreadthFirstSearchParameters parameters) {
		this.parameters = parameters;
	}


	@Override
	protected String getExtraArgs(GraphlessJobParams jobParams) {
		return "{\"sourceVertex\":" + parameters.getSourceVertex() + "}";
	}

	@Override
	protected String getAlgorithmName() {
		return "BFS";
	}

	@Override
	protected ResultType getResultType() {
		return ResultType.LONG_TO_LONG;
	}
}