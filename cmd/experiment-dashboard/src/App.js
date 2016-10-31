import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';
import { Namespaces } from './Namespaces.js';

class App extends Component {
  constructor() {
    super();
    this.state = {
      experiments: [],
    };
  }
  componentDidMount() {
    const headers = new Headers();
    headers.append('Accept', 'application/json');
    let req = new Request(`/api/v1/experiments`, {headers});
    fetch(req)
    .then(response => {
      response.json()
      .then(json => {
        this.setState({experiments: json});
      });
    });
  }
  render() {
    return (
      <div className="App">
        <div className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h2>Experiment Dashboard</h2>
        </div>
        <Namespaces namespaces={this.state.experiments} />
      </div>
    );
  }
}

export default App;
