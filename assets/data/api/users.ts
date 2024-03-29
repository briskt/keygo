import type {User, UserUpdateInput} from 'data/types/user'
import api from '../api'

export const getUser = async (id: string): Promise<User> => {
  const response = await api.get('/api/users/' + encodeURIComponent(id))
  return response.json()
}

export const listUsers = async (): Promise<User[]> => {
  const response = await api.get('/api/users')
  return response.json()
}

export const updateUser = async (id: string, input: UserUpdateInput): Promise<User> => {
  const response = await api.put('/api/users/' + encodeURIComponent(id), input)
  return response.json()
}
