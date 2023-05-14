import api from 'data/api'
import type {AuthStatus, Provider} from '../types/auth'

const clientIDparam = 'clientID'
const tokenParam = 'token'

export const getToken = () => getClientID() + (localStorage.getItem(tokenParam) || '')

export const getAuthStatus = async (): Promise<AuthStatus> => {
  const response = await api.get('/api/auth')
  return response.json()
}

export const getLoginProviders = async (): Promise<Provider[]> => {
  const response = await api.get(`/api/auth/login?client_id=${getClientID()}`)
  return response.json()
}

init()
function init() {
  localStorage.getItem(clientIDparam) || localStorage.setItem(clientIDparam, makeRandomID())
}

const getClientID = () => localStorage.getItem(clientIDparam)

const makeRandomID = () => Math.random().toString(36).slice(2)

