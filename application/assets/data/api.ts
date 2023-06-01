import { getToken } from 'data/store/auth'
import { ResponseError } from './error'
import { setNotice } from '@silintl/ui-components'

type Method = 'delete' | 'get' | 'put' | 'post'

/**
 * @param method -- The HTTP method (e.g. 'get')
 * @param urlPath -- The URL path (e.g. '/api/something')
 * @param data -- (Optional:) The data to send in the body. If provided, it will be JSON encoded.
 * @param showError -- (Optional:) Whether to show the user the error message.
 * @throws {ResponseError}
 */
const call = async (method: Method, urlPath: string, data = null, showError = true): Promise<Response> => {
  const headers = {
    'Content-Type': 'application/json',
  }

  const response = await fetch(urlPath, {
    method,
    headers,
    body: data === null ? null : JSON.stringify(data),
  })
  if (!response.ok) {
    console.error(response)
    if (showError) {
      const responseData = await response.json()
      const errorMessage = responseData?.message || response.statusText
      setNotice(errorMessage)
    }
    throw new ResponseError(response, `${response.status} ${response.statusText}`)
  }
  return response
}

/**
 * @throws {ResponseError}
 */
const get = async (urlPath: string, showError = true): Promise<Response> => call('get', urlPath, null, showError)

/**
 * @throws {ResponseError}
 */
const post = async (urlPath: string, data: any = null): Promise<Response> => call('post', urlPath, data)

/**
 * @throws {ResponseError}
 */
const put = async (urlPath: string, data: any = null): Promise<Response> => call('put', urlPath, data)

/**
 * @throws {ResponseError}
 */
const remove = async (urlPath: string): Promise<Response> => call('delete', urlPath)

export default {
  get,
  post,
  put,
  remove,
}
