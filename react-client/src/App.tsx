import * as React from 'react';
import './App.css';
import Chat from './Chat/Chat';
import * as Client from './client';
import Login from './Login';
import StandBy from './StandBy';

import 'bulma/css/bulma.css';
import logo from './logo.svg';


class App extends React.Component<{}, {
  client: Client.WSClient,
  state: number
}> {

  constructor(props: {}) {
    super(props)
    this.state = {
      client: null,
      state: 0,
    }
    this.onLogin = this.onLogin.bind(this);
  }

  public render() {
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h1 className="App-title">Welcome to ReallyTinyChat</h1>
        </header>
        {this.state.state === 0 ?
          <Login onLogin={this.onLogin}/> :
          this.state.state === 1 ?
          <StandBy />: <Chat client={this.state.client}/>
        }
        
      </div>
    );
  }

  private onLogin(client: Client.WSClient) {
    this.setState(s => ({state: 1, client}))
    const self = this
    const interval = setInterval(() => {
      if (self.state.client.State() === "open") {
        self.setState({state: 2})
        clearInterval(interval)
        console.log("Bye-bye")
        return
      }
    }, 300)
  }
}

export default App;
