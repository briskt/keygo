import api from 'data/api'
import type {AuthStatus, Provider} from '../types/auth'

export const getAuthStatus = async (): Promise<AuthStatus> => {
  try {
    const response = await api.get('/api/auth')
    return response.json()
  } catch {
    const notAuth = <AuthStatus>{
      IsAuthenticated: false
    };
    return new Promise<AuthStatus>((resolve, reject) => {
      resolve(notAuth);
    });
  }
}

export const getLoginProviders = async (): Promise<Provider[]> => {
  const response = await api.get(`/api/auth/login`)
  return response.json()
}
