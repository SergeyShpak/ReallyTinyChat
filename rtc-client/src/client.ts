import EventedArray from './evented_array'
import * as http from './http'


const server = "localhost:8080"
const client = new http.Client("http://localhost:8080")

interface IWsMessage {
  Type: string
  Message: string
}

interface IErrorMessage {
  Code: number
  Hint: string
}

interface IHelloMessage {
  Login: string
  Room: string
}

interface IHelloOKMessage {
  Login: string
  Room: string
}

interface IRoomInfoMessage {
  Connector: string,
  Connectee: string,
  Room: string,
}

interface IOffer {
  Login: string
  Offer: string
  IsResponse: boolean
  Room: string
}

interface IIceCandidate {
  Room: string
  Candidate: string
}

const config = null

let clientInstance: WSClient

export class WSClient {
  public candidates: number = 0

  private isOfferer: boolean = undefined
  private host: string
  private login: string = ""
  private messagesQueue: EventedArray
  private offerReceived: boolean = false
  private onClose: (e: CloseEvent) => void
  private partner: string = ""
  private room: string = ""

  private rtcConn: RTCPeerConnection
  private rtcDataChannel: RTCDataChannel

  private ws: WebSocket

  constructor(host?: string) {
    this.host = host
    if (this.host === undefined || this.host === null) {
      this.host = "localhost:8080"
    }
    const addr = "ws://" + this.host + "/conn"
    this.ws = new WebSocket(addr)
    clientInstance = this

    this.ws.onmessage = this.onServerMessage
  }

  public async Connect(login: string, room: string) {
    this.login = login
    const payload = JSON.stringify({
      Login: login,
      Room: room,
    })
    const helloMsg: IWsMessage = {
      Message: payload,
      Type: "HELLO"
    }
    await this.waitAndSend(helloMsg)
    this.rtcConn = new RTCPeerConnection(config);
    this.rtcConn.ondatachannel = this.receiveChannel
    this.rtcConn.onicecandidate = this.onIceCandidate
  }

  public async SendOffer() {
    this.rtcDataChannel = this.rtcConn.createDataChannel("data")
    this.rtcDataChannel.onmessage = this.rtcHandleReceiveMessage
    const offer = await this.rtcConn.createOffer()
    await this.rtcConn.setLocalDescription(offer)
    const payload = JSON.stringify({
      IsResponse: false,
      Login: this.login,
      Offer: JSON.stringify(offer),
      Room: this.room,
    })
    const msg: IWsMessage = {
      Message: payload,
      Type: "OFFER"
    }
    await this.waitAndSend(msg)
  }

  public async SendOfferResponse(offer: RTCSessionDescriptionInit) {
    await this.rtcConn.setRemoteDescription(offer)
    const response = await this.rtcConn.createAnswer()
    await this.rtcConn.setLocalDescription(response)
    const payload = JSON.stringify({
      IsResponse: true,
      Login: this.login,
      Offer: JSON.stringify(response),
      Room: this.room,
    })
    const msg: IWsMessage = {
      Message: payload,
      Type: "OFFER"
    }
    await this.waitAndSend(msg)
  }

  public State(): RTCDataChannelState {
    if (!this.rtcDataChannel) {
      return "connecting"
    }
    return this.rtcDataChannel.readyState
  }

  public async Wait() {
    await this.sleep(5000)
    console.log(this.rtcDataChannel)
  }

  public async WaitForOpen(): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      if (this.rtcDataChannel.readyState === "open") {
        return resolve()
      }
    })
  }

  public Login(): string {
    return this.login
  }

  public async SetRemoteChannel(offer: RTCSessionDescriptionInit) {
    await this.rtcConn.setRemoteDescription(offer)
  }

  public SendOnDataChannel(msg: string) {
    // sending an empty string fails, seems not to be the deafault behavior
    if (!msg || msg.length === 0) {
      return
    }
    this.messagesQueue.push({from: this.login, msg})
    this.rtcDataChannel.send(msg)
  }

  public async FinalizeOffer(offer: RTCSessionDescriptionInit) {
    await this.rtcConn.setRemoteDescription(offer)
    await this.sleep(500)
    console.log("Finalize: ", this.rtcConn)
  }

  public async SendIce(candidate: RTCIceCandidate) {
    const payload = JSON.stringify({
      Candidate: JSON.stringify(candidate),
      Room: this.room,
    })
    const msg: IWsMessage = {
      Message: payload,
      Type: "ICE"
    }
    await this.waitAndSend(msg)
  }

  public async AddIceCandidate(candidate: RTCIceCandidate) {
    this.rtcConn.addIceCandidate(candidate)
    this.candidates += 1
  }

  public async SetDumpReceivedMessage(messagesQueue: EventedArray) {
    this.messagesQueue = messagesQueue
  }

  public SetOnClose(f: (e: CloseEvent) => void) {
    this.rtcDataChannel.onclose = f
    this.onClose = f
  }

  private send(msg: IWsMessage) {
    this.ws.send(JSON.stringify(msg));
  }

  private receiveChannel(e: RTCDataChannelEvent) {
    clientInstance.rtcDataChannel = e.channel
    clientInstance.rtcDataChannel.onmessage = clientInstance.rtcHandleReceiveMessage
  }

  private rtcHandleReceiveMessage(e: MessageEvent) {
    clientInstance.messagesQueue.push({from: clientInstance.partner, msg: e.data})
  }

  private onIceCandidate(e: RTCPeerConnectionIceEvent) {
    if (!e.candidate) {
      return
    }
    clientInstance.SendIce(e.candidate)
  }

  private async waitAndSend(msg: IWsMessage, ctr?: number, timeoutMs?: number): Promise<void> {
    if (ctr === undefined || ctr === null || ctr < 0) {
      ctr = 10
    }
    if (timeoutMs === undefined || timeoutMs === null || timeoutMs < 0) {
      timeoutMs = 10
    }
    while (ctr !== 0) {
      if (this.ws.readyState === 1) {
        this.send(msg)
        Promise.resolve()
        return
      }
      ctr--
      await this.sleep(timeoutMs)
    }
    Promise.reject("waitAndSend timeout")
  }

  private sleep(ms): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  private onServerMessage(e: MessageEvent) {
    const msg = JSON.parse(e.data) as IWsMessage
    switch(msg.Type) {
      case "HELLOOK":
        if (clientInstance.isOfferer !== undefined) {
          return
        }
        clientInstance.isOfferer = false
        const helloResp = JSON.parse(msg.Message) as IHelloOKMessage;
        clientInstance.room = helloResp.Room
        return
      case "ROOMINFO":
        const roomInfo = JSON.parse(msg.Message) as IRoomInfoMessage
        clientInstance.room = roomInfo.Room
        if (roomInfo.Connector === clientInstance.login) {
          clientInstance.partner = roomInfo.Connectee
          return
        }
        clientInstance.isOfferer = true
        clientInstance.partner = roomInfo.Connector
        clientInstance.SendOffer()
        return
      case "OFFER":
        if (clientInstance.offerReceived) {
          return
        }
        clientInstance.offerReceived = true
        const offerPayload = JSON.parse(msg.Message) as IOffer
        if (clientInstance.isOfferer && offerPayload.IsResponse) {
          const offer = JSON.parse(offerPayload.Offer) as RTCSessionDescriptionInit
          clientInstance.FinalizeOffer(offer).then(() =>
            clientInstance.Wait())
          return
        }
        if (!clientInstance.isOfferer && !offerPayload.IsResponse) {
          const offer = JSON.parse(offerPayload.Offer) as RTCSessionDescriptionInit
          clientInstance.SendOfferResponse(offer).then(() =>
            clientInstance.Wait()
          )
          return
        }
        return
      case "ICE":
        const candidatePayload = JSON.parse(msg.Message) as IIceCandidate
        clientInstance.AddIceCandidate(JSON.parse(candidatePayload.Candidate) as RTCIceCandidate)
        return
      case "CLOSE":
        if (!!this.rtcDataChannel && !!this.onClose) {
          this.rtcDataChannel.close()
          this.onClose(null)
        }
        break
      default:
        throw Error("Bad message type: " + msg.Type)
    }
  }
}