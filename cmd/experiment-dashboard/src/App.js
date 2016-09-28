import React from 'react';
import logo from './logo.svg';
import './App.css';
import { Experiments } from './Experiments.js';

const App = props => {
  const data = [
    {
      name: "test",
      segments: "0123456789abcdef0123456789abcdef",
      params: [
        {
          name: "param1",
          values: [
            {name: "value1", weight: 1},
            {name: "value2", weight: 2},
            {name: "value3", weight: 3},
          ],
        },
        {
          name: "param2",
          values: [
            {name: "value1", weight: 3},
            {name: "value2", weight: 2},
            {name: "value3", weight: 1},
          ],
        },
      ],
    },
    {
      name: "test2",
      segments: "fedcba9876543210fedcba9876543210",
      params: [
        {
          name: "poop",
          values: [
            {name: "brown", weight: 1},
            {name: "brown-green", weight: 1},
            {name: "messy", weight: 1},
          ],
        },
      ],
    },
  ];
  return (
    <div className="App">
      <div className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <h2>Experiment Dashboard</h2>
      </div>
      <p className="App-intro">
        Do Stuff.
      </p>
      <Experiments experiments={data} />
    </div>
  );
}

export default App;
