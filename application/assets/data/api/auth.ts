import api from 'data/api'
import type {AuthStatus} from '../types/auth'
import { writable } from 'svelte/store'

export const authStatus = writable({} as AuthStatus)

export const getAuthStatus = async (): Promise<AuthStatus> => {
  try {
    const response = await api.get('/api/auth')
    const updatedAuthStatus = await response.json()
    authStatus.set(updatedAuthStatus)
    return updatedAuthStatus
  } catch {
    const notAuth = <AuthStatus>{
      IsAuthenticated: false
    };
    authStatus.set(notAuth)
    return new Promise<AuthStatus>((resolve, reject) => {
      resolve(notAuth);
    });
  }
}
