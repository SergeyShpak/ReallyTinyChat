export interface IWsWrapper {
  Type: string
  Message: string
}

export interface IHello {
  Login: string
  Room: string
}

export interface IHelloOK {
  Login: string
  Room: string
  Partners: string[]
}

export interface IOffer {
  Login: string
  Room: string
  Partner: string
  Offer: string
  IsResponse: boolean
}

export interface IIceCandidate {
  Room: string
  Partner: string
  Candidate: string
}

export interface IError {
  Code: number,
  Hint: string
}