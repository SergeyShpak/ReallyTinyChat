import EventedArray from './evented_array'
import * as Messages from './messages'


const server = "localhost:8080"

let clientInstance: RTCClient

export class RTCClient {
  private host: string
  private login: string = ""
  private messagesQueue: EventedArray
  private onRTCConnection: () => void
  private onRTCDataChannelOpen: () => void
  private onClose: (e: CloseEvent) => void
  private partner: string = ""
  private room: string = ""

  private rtcConn: RTCPeerConnection
  private rtcDataChannel: RTCDataChannel

  private serverWS: WebSocket

  constructor(host?: string) {
    this.host = host
    if (this.host === undefined || this.host === null) {
      this.host = server
    }
    const addr = "ws://" + this.host + "/conn"
    this.serverWS = new WebSocket(addr)
    clientInstance = this

    this.serverWS.onerror = this.onWsError
    this.serverWS.onmessage = this.onServerMessage
  }

  public async Connect(login: string, room: string) {
    this.login = login
    const helloMsg: Messages.IHello ={
      Login: login,
      Room: room,
    }
    const msg: Messages.IWsWrapper = {
      Message: JSON.stringify(helloMsg),
      Type: "HELLO"
    }
    await this.send(msg)
    this.createRTCConnection(null)
  }

  public SendOnDataChannel(msg: string) {
    // sending an empty string fails, seems not to be the deafault behavior
    if (!msg || msg.length === 0) {
      return
    }
    this.messagesQueue.push({from: this.login, msg})
    this.rtcDataChannel.send(msg)
  }

  public async SetDumpReceivedMessage(messagesQueue: EventedArray) {
    this.messagesQueue = messagesQueue
  }

  public SetOnClose(f: (e: CloseEvent) => void) {
    // this.rtcDataChannel.onclose = f
    this.onClose = f
  }

  public SetOnRTCConnection(f: () => void) {
    this.onRTCConnection = f
  }

  public SetOnRTCDataChannelOpen(f: () => void) {
    this.onRTCDataChannelOpen = f
  }

  public Partner(): string {
    return this.partner
  }

  private async sendOffer(partner: string) {
    this.partner = partner
    this.rtcDataChannel = this.rtcConn.createDataChannel("data")
    this.rtcDataChannel.onopen = this.onDataChannelOpen
    this.rtcDataChannel.onmessage = this.rtcHandleReceiveMessage
    const offer = await this.rtcConn.createOffer()
    await this.rtcConn.setLocalDescription(offer)
    const offerMsg: Messages.IOffer = {
      IsResponse: false,
      Login: this.login,
      Offer: JSON.stringify(offer),
      Partner: partner,
      Room: this.room,
    }
    const msg: Messages.IWsWrapper = {
      Message: JSON.stringify(offerMsg),
      Type: "OFFER"
    }
    await this.send(msg)
  }

  private async sendOfferResponse(offer: Messages.IOffer) {
    const rtcSessionDesc = JSON.parse(offer.Offer) as RTCSessionDescriptionInit
    await this.rtcConn.setRemoteDescription(rtcSessionDesc)
    const response = await this.rtcConn.createAnswer()
    await this.rtcConn.setLocalDescription(response)
    const offerResp: Messages.IOffer = {
      IsResponse: true,
      Login: this.login,
      Offer: JSON.stringify(response),
      Partner: this.partner,
      Room: this.room,
    }
    const msg: Messages.IWsWrapper = {
      Message: JSON.stringify(offerResp),
      Type: "OFFER"
    }
    await this.send(msg)
  }

  private async finalizeOffer(offer: Messages.IOffer) {
    const rtcSessionDescr = JSON.parse(offer.Offer) as RTCSessionDescriptionInit
    await this.rtcConn.setRemoteDescription(rtcSessionDescr)
  }

  private onIceCandidate(e: RTCPeerConnectionIceEvent) {
    if (!e.candidate) {
      return
    }
    clientInstance.sendIce(e.candidate)
  }

  private async sendIce(candidate: RTCIceCandidate) {
    const iceMsg: Messages.IIceCandidate = {
      Candidate: JSON.stringify(candidate),
      Partner: this.partner,
      Room: this.room,
    }
    const msg: Messages.IWsWrapper = {
      Message: JSON.stringify(iceMsg),
      Type: "ICE"
    }
    await this.send(msg)
  }

  private async addIceCandidate(candidate: RTCIceCandidate) {
    this.rtcConn.addIceCandidate(candidate)
  }

  private createRTCConnection(config: any) {
    this.rtcConn = new RTCPeerConnection(config);
    this.rtcConn.ondatachannel = this.receiveChannel
    this.rtcConn.onicecandidate = this.onIceCandidate
    if (!!this.onRTCConnection) {
      this.onRTCConnection()
    }
  }

  private close(e?: CloseEvent) {
    if (!!this.rtcDataChannel) {
      this.rtcDataChannel.close()
    }
    if (!!this.rtcConn) {
      this.rtcConn.close()
    }
    if (!!this.serverWS) {
      this.serverWS.close()
    }
    this.onClose(e)
  }

  private receiveChannel(e: RTCDataChannelEvent) {
    clientInstance.rtcDataChannel = e.channel
    clientInstance.rtcDataChannel.onmessage = clientInstance.rtcHandleReceiveMessage
    clientInstance.onRTCDataChannelOpen()
  }

  private rtcHandleReceiveMessage(e: MessageEvent) {
    clientInstance.messagesQueue.push({from: clientInstance.partner, msg: e.data})
  }

  private async send(msg: Messages.IWsWrapper, ctr?: number, timeoutMs?: number): Promise<void> {
    const sleep = (ms: number) => {
      return new Promise(resolve => setTimeout(resolve, ms));
    }
    if (ctr === undefined || ctr === null || ctr < 0) {
      ctr = 10
    }
    if (timeoutMs === undefined || timeoutMs === null || timeoutMs < 0) {
      timeoutMs = 10
    }
    while (ctr !== 0) {
      if (this.serverWS.readyState === 1) {
        this.serverWS.send(JSON.stringify(msg));
        Promise.resolve()
        return
      }
      ctr--
      await sleep(timeoutMs)
    }
    Promise.reject("send timeout")
  }

  private async handleHelloOK(msg: Messages.IHelloOK) {
    this.login = msg.Login
    this.room = msg.Room
    const partners = msg.Partners.filter(p => p !== this.login)
    await Promise.all(partners.map(async (p) => {
      await this.sendOffer(p)
    }))
  }

  private async handleOffer(msg: Messages.IOffer) {
    if (msg.IsResponse) {
      await this.finalizeOffer(msg)
      return
    }
    this.partner = msg.Login
    await this.sendOfferResponse(msg)
    return
  }

  private async handleIce(msg: Messages.IIceCandidate) {
    await this.addIceCandidate(JSON.parse(msg.Candidate) as RTCIceCandidate)
  }

  private handleClose() {
    if (!!this.rtcDataChannel && !!this.onClose) {
      this.rtcDataChannel.close()
      this.onClose(null)
    }
  }

  private handleError(msg: Messages.IError) {
    console.log(msg)
    const closeEvent = new CloseEvent("ERROR code received", {code: 1002, reason: msg.Hint})
    this.close(closeEvent)
    return
  }

  private handleMessage(e: MessageEvent) {
    const msg = JSON.parse(e.data) as Messages.IWsWrapper
    switch(msg.Type) {
      case "HELLOOK":
        const helloResp = JSON.parse(msg.Message) as Messages.IHelloOK;
        this.handleHelloOK(helloResp);
        return
      case "OFFER":
        const offerPayload = JSON.parse(msg.Message) as Messages.IOffer
        this.handleOffer(offerPayload)
        return
      case "ICE":
        const candidatePayload = JSON.parse(msg.Message) as Messages.IIceCandidate
        this.handleIce(candidatePayload)
        return
      case "CLOSE":
        this.handleClose()
        return
      case "ERROR":
        const error = JSON.parse(msg.Message) as Messages.IError
        this.handleError(error)
        return
      default:
        throw Error("Bad message type: " + msg.Type)
    }
  }

  private onServerMessage(e: MessageEvent) {
    clientInstance.handleMessage(e)
  }

  private onDataChannelOpen(e) {
    if (!!clientInstance.onRTCDataChannelOpen) {
      clientInstance.onRTCDataChannelOpen()
    }
  }

  private onWsError(e: ErrorEvent) {
    console.log("an error occurred: ", e)
  }
}

export default RTCClient