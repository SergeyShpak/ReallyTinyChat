import RTCClient from './RTCClient/rtc-client';

export interface ILoginProps{
  onLogin: (login: RTCClient) => void
  onClientClose: (e: CloseEvent) => void
  onOpenDataChannel: () => void
  onDataChannelClose: (e: CloseEvent) => void
  onServerError: (code: number, hint: string) => void
}

export interface IContactsProps {
  contacts: string[]
}

export interface IChatProps {
  active: boolean
  client: RTCClient
}

export interface IChatAreaProps {
  active: boolean
  client: RTCClient
}

export interface IErrorBoundaryProps {
  onError: () => void
}