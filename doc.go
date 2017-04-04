// Copyright 2016 Andrew O'Neill, Nordstrom

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
Package choices provides a library for simple a/b and multivariate
testing.

choices uses hashing to uniquely assign users to experiments and
decide the values the user is assigned. This allows us to quickly
assign a user to an experiment with out having to look up what they
were assigned previously. Most of the ideas in this package are based
off of Facebook's Planout.

In choices there are three main concepts. Experiments, Experiments, and
Params.  Experiments split traffic between the experiments they
contain. Experiments are the thing you are testing. Experiments are
made up of one or more Params.  Params have names and values. The
values can be either a uniform, equal weights between choices, or
weighted, user-specified weights for each choice.

In most cases you will want to create one namespace per experiment. If
you have experiments that might have interactions you can use
namespaces to split the traffic between them. For example, if you are
running a test on a banner and another test that takes over the whole
page, you will want to split the traffic between these two tests.
Another example, is if you want to run a test on a small percent of
traffic. The namespace will ensure that experiments won't overlap when
you evaluate the experiment.

Experiments contain the Params for a running experiment. When a caller
queries Experiments, they will first be hashed into a segment. We then
check if the segment is contained in that Experiment. If the segment
is contained in the experiment then the experiment will be evaluated.
An experiment will in turn evaluate each of it's Params.

Params are key-value pairs. They Options for value are Uniform choice
or Weighted choice. Uniform choices will be selected in a uniform
fashion. Weighted choices will be selected based on the proportions
supplied in the weights.
*/
package choices
