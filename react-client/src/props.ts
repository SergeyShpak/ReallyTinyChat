import * as Client from './client';

export interface ILoginProps{
  onLogin: (login: Client.WSClient) => void
}

export interface IContactsProps {
  contacts: string[]
}

export interface IChatProps {
  client: Client.WSClient
}

export interface IChatAreaProps {
  client: Client.WSClient
}