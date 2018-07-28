import * as React from 'react';
import './App.css';
import Chat from './Chat/Chat';
import Login from './Login';
import RTCClient from './RTCClient/rtc-client';
import StandBy from './StandBy';

import 'bulma/css/bulma.css';
import logo from './logo.svg';


class App extends React.Component<{}, {
  client: RTCClient,
  state: number
}> {

  constructor(props: {}) {
    super(props)
    this.state = {
      client: null,
      state: 0,
    }
    this.onLogin = this.onLogin.bind(this);
    this.onChatClose = this.onChatClose.bind(this);
    this.onError = this.onError.bind(this);
  }

  public render() {
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h1 className="App-title">Welcome to ReallyTinyChat</h1>
        </header>
        {this.state.state === 0 ?
          <Login onLogin={this.onLogin} onClientClose={this.onError}/> :
          this.state.state === 1 ?
          <StandBy />: <Chat client={this.state.client} closeChat={this.onChatClose}/>
        }
      </div>
    );
  }

  private onLogin(client: RTCClient) {
    this.setState({state: 1, client})
    const self = this
    const interval = setInterval(() => {
      if (self.state.client.State() === "open") {
        self.setState({state: 2})
        clearInterval(interval)
        return
      }
    }, 300)
  }

  private onChatClose() {
    this.setState({state: 0})
  }

  private onError() {
    console.log("Fired")
    this.setState({state: 0})
  }
}

export default App;
