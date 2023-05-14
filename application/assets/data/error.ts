/**
 * Custom error type for non-ok API responses.
 *
 * See
 * [Custom error types](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Error#custom_error_types).
 */
export class ResponseError extends Error {
  response: Response
  constructor(response: Response, ...params: any[]) {
    super(...params)
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, ResponseError)
    }
    this.name = 'ResponseError'
    this.response = response
  }
}
