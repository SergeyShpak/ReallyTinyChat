import RTCClient from './RTCClient/rtc-client';

export interface ILoginProps{
  onLogin: (login: RTCClient) => void
  onClientClose: (e: CloseEvent) => void
}

export interface IContactsProps {
  contacts: string[]
}

export interface IChatProps {
  client: RTCClient
  closeChat: () => void
}

export interface IChatAreaProps {
  client: RTCClient
  closeChat: () => void
}

export interface IErrorBoundaryProps {
  onError: () => void
}