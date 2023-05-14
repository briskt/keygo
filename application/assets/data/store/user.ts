import type { User } from 'data/types/user'
import api from '../api'
import { writable } from 'svelte/store'

export const user = writable({} as User)

/**
 * @throws {ResponseError}
 */
export async function loadUser(id: string): Promise<number> {
  const response = await api.get(`/api/users/${id}`)
  setUser(await response.json())
  return response.status
}

/**
 * Set the data in our `user` store. NOTE: Only to be called with data from an
 * API response, never with data directly from the end user.
 */
export const setUser = (updatedUserData: User) => user.set(updatedUserData)
