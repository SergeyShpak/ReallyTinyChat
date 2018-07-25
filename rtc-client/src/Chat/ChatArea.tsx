import * as React from 'react';
import { Component } from 'react';
import * as Client from '../client';
import * as Props from '../props';

import EventedArray from '../evented_array';


class ChatArea extends Component<Props.IChatAreaProps, {
  chatMessages: string,
  cols: number,
  disabled: boolean,
  rows: number,

  message: string
  messagesQueue: EventedArray
  client: Client.WSClient
}>{

  constructor(props: Props.IChatAreaProps) {
    super(props)
    this.onMessageChange = this.onMessageChange.bind(this)
    this.onSendClick = this.onSendClick.bind(this)
    this.onCloseChat = this.onCloseChat.bind(this)

    const messagesQueue: EventedArray = new EventedArray(() => {
      const msg = this.state.messagesQueue.Stack.shift()
      this.setState({chatMessages: this.state.chatMessages + "\n" + msg.from + "> " + msg.msg})
    })
    this.props.client.SetDumpReceivedMessage(messagesQueue)
    this.state = {
      chatMessages: "",
      client: this.props.client,
      cols: 40,
      disabled: true,
      message: "",
      messagesQueue,
      rows: 15,
    }
  }

  public componentDidMount() {
    this.state.client.SetOnClose(this.onCloseChat)
  }

  public render() {
    return (
      <div className="container">
        <textarea id="chat-area" name="chat-area" disabled={this.state.disabled} rows={this.state.rows} cols={this.state.cols} value={this.state.chatMessages} wrap="hard"/>
        <div id="send-group">
          <input
            key="1234"
            className="input"
            type="text"
            placeholder="Your message"
            value={this.state.message}
            onChange={this.onMessageChange}
          />
          <input
            className="button"
            type="button"
            value="Send"
            onClick={this.onSendClick}
          />
        </div>
      </div>
    )
  }

  private onCloseChat(e: CloseEvent) {
    this.props.closeChat()
  }

  private onMessageChange(e) {
    this.setState({message: e.target.value})
  } 

  private onSendClick() {
    this.state.client.SendOnDataChannel(this.state.message)
    this.setState({message: ""})
  }
}

export default ChatArea