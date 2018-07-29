import RTCClient from './RTCClient/rtc-client';

export interface ILoginProps{
  onLogin: (login: RTCClient) => void
  onClientClose: (e: CloseEvent) => void
  onOpenDataChannel: () => void
}

export interface IContactsProps {
  contacts: string[]
}

export interface IChatProps {
  active: boolean
  client: RTCClient
  closeChat: () => void
}

export interface IChatAreaProps {
  active: boolean
  client: RTCClient
  closeChat: () => void
}

export interface IErrorBoundaryProps {
  onError: () => void
}