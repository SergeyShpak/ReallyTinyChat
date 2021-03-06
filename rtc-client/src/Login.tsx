import { Component } from 'react';
import * as React from 'react';
import * as Props from './props';
import RTCClient from './RTCClient/rtc-client';


const steps = {
  0: 'Login',
  1: 'Logging in...'
}

class Login extends Component<Props.ILoginProps, {
  login: string,
  disable: boolean,
  room: string,
  step: number
}> {

  constructor(props: Props.ILoginProps) {
    super(props);
    this.state = {
      disable: false,
      login: "SSH",
      room: "room",
      step: 0,
    };
    this.onLoginChange = this.onLoginChange.bind(this);
    this.onChatRoomChange = this.onChatRoomChange.bind(this);
    this.onLoggingIn = this.onLoggingIn.bind(this);
  }

  public render() {
    return (
      <div className="container login">
        <input
          className="input"
          type="text"
          value={this.state.login}
          onChange={this.onLoginChange}
          placeholder="Login"
        />
        <input
          className="input"
          type="text"
          value={this.state.room}
          onChange={this.onChatRoomChange}
          placeholder="Chat room name"
        />
        <input
          className="button"
          type="button"
          value={steps[this.state.step]}
          onClick={this.onLoggingIn}
          disabled={this.state.disable}
        />
      </div>
    );
  }

  private async onLoggingIn(e: React.MouseEvent<HTMLInputElement>) {
    const client = new RTCClient()
    client.SetOnRTCConnection(() => {
      this.props.onLogin(client)
    })
    client.SetOnRTCDataChannelOpen(this.props.onOpenDataChannel)
    client.SetOnServerError(this.props.onServerError)
    client.SetOnClose(this.props.onDataChannelClose)
    await client.Connect(this.state.login, this.state.room)
  }

  private async onChatRoomChange(e) {
    this.setState({ room: e.target.value });
  }

  private onLoginChange(e) {
    this.setState({ login: e.target.value });
  }
}

export default Login;