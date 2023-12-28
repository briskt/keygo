import type {Tenant, TenantCreate, TenantUserCreate} from 'data/types/tenant'
import api from '../api'

// TODO: cache tenant list

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

export const addTenantUser = async (tenantID: string, email: string): Promise<Tenant> => {
  const body: TenantUserCreate = {
    Email: email,
  }
  const response = await api.post('/api/tenants/'+encodeURIComponent(tenantID)+'/users', body)
  return response.json()
}
