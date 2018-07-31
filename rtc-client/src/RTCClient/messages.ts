import * as JWT from 'jsonwebtoken'

export interface IWsMessage {
  Type: string
  Login: string
  Room: string
  Token: string
}

export interface IServerMessage {
  Type: string
  Payload: string
}

export interface IHelloOK {
  Login: string
  Room: string
  Secret: string
  Partners: string[]
}

export interface IOffer {
  Login: string
  Partner: string
  Offer: string
  IsResponse: boolean
}

export interface IIceCandidate {
  Partner: string
  Candidate: string
}

export interface IError {
  Code: number
  Hint: string
}