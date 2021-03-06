import * as React from 'react';
import './App.css';
import Chat from './Chat/Chat';
import Login from './Login';
import RTCClient from './RTCClient/rtc-client';

import 'bulma/css/bulma.css';
import logo from './logo.svg';


class App extends React.Component<{}, {
  client: RTCClient,
  isChatActive: boolean,
  state: number
}> {

  constructor(props: {}) {
    super(props)
    this.state = {
      client: null,
      isChatActive: false,
      state: 0,
    }
    this.onLogin = this.onLogin.bind(this)
    this.onOpenDataChannel = this.onOpenDataChannel.bind(this)
    this.onServerError = this.onServerError.bind(this)
    this.onDataChannelClose = this.onDataChannelClose.bind(this);
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
          <Login
            onLogin={this.onLogin}
            onClientClose={this.onError}
            onOpenDataChannel={this.onOpenDataChannel}
            onServerError={this.onServerError}
            onDataChannelClose={this.onDataChannelClose}
          /> :
          <Chat client={this.state.client} active={this.state.isChatActive}/>
        }
      </div>
    );
  }

  private onLogin(client: RTCClient): void {
    this.setState({state: 1, client})
  }

  private onDataChannelClose() {
    this.setState({isChatActive: false})
  }

  private onError() {
    console.log("Fired")
    this.setState({state: 0})
  }

  private onOpenDataChannel() {
    this.setState({isChatActive: true})
  }

  private onServerError(code: number, hint: string) {
    this.setState({state: 0})
    alert("ERROR: " + hint + " (" + code + ")")
  }
}

export default App;
