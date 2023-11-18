import type {Tenant, TenantCreate} from 'data/types/tenant'
import api from '../api'

export const addTenant = async (name: string): Promise<Tenant> => {
  const body: TenantCreate = {
    Name: name,
  }
  const response = await api.post('/api/tenants', body)
  return response.json()
}

export const listTenants = async (): Promise<Tenant[]> => {
  const response = await api.get('/api/tenants')
  return response.json()
}

export const getTenant = async (id: string): Promise<Tenant> => {
  const response = await api.get('/api/tenants/'+encodeURIComponent(id))
  return response.json()
}
