import type {User} from 'data/types/user'
import type {Token} from 'data/types/token'
import api from '../api'

export const viewUser = async (id: string): Promise<User> => {
  const response = await api.get('/api/users/' + encodeURIComponent(id))
  return response.json()
}

export const listUsers = async (): Promise<User[]> => {
  const response = await api.get('/api/users')
  return response.json()
}

export const listUserTokens = async (id: string): Promise<Token[]> => {
  const response = await api.get('/api/users/' + encodeURIComponent(id) + '/tokens')
  return response.json()
}
