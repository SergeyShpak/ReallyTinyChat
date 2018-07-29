import * as React from 'react';
import { Component } from 'react';
import * as Props from '../props';
import EventedArray from '../RTCClient/evented_array';


class ChatArea extends Component<Props.IChatAreaProps, {
  chatMessages: string,
  cols: number,
  disabled: boolean,
  rows: number,

  message: string
  messagesQueue: EventedArray
}>{

  constructor(props: Props.IChatAreaProps) {
    super(props)
    this.onMessageChange = this.onMessageChange.bind(this)
    this.onSendClick = this.onSendClick.bind(this)
    // this.onCloseChat = this.onCloseChat.bind(this)

    const messagesQueue: EventedArray = new EventedArray(() => {
      const msg = this.state.messagesQueue.Stack.shift()
      this.setState({chatMessages: this.state.chatMessages + "\n" + msg.from + "> " + msg.msg})
    })
    this.props.client.SetDumpReceivedMessage(messagesQueue)
    this.state = {
      chatMessages: "",
      cols: 40,
      disabled: true,
      message: "",
      messagesQueue,
      rows: 15,
    }
  }

  public componentDidMount() {
    // this.props.client.SetOnClose(this.onCloseChat)
  }

  public render() {
    return (
      <div className="container">
        <textarea
          id="chat-area"
          name="chat-area" 
          disabled={true}
          rows={this.state.rows}
          cols={this.state.cols}
          value={
            this.props.active ?
            this.props.client.Partner() + " connected\n" + this.state.chatMessages :
            "Waiting for a partner"
          }
          wrap="hard"/>
        <div id="send-group">
          <input
            key="1234"
            className="input"
            disabled={!this.props.active}
            type="text"
            placeholder="Your message"
            value={this.state.message}
            onChange={this.onMessageChange}
          />
          <input
            className="button"
            disabled={!this.props.active}
            type="button"
            value="Send"
            onClick={this.onSendClick}
          />
        </div>
      </div>
    )
  }

  private onMessageChange(e) {
    this.setState({message: e.target.value})
  } 

  private onSendClick() {
    this.props.client.SendOnDataChannel(this.state.message)
    this.setState({message: ""})
  }
}

export default ChatArea