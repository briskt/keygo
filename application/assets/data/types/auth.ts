
export type AuthStatus = {
  IsAuthenticated: boolean
  Expiry: string // date
  UserID: string
}

export type Provider = {
  Key: string
  Name: string
  RedirectURL: string
}

export type Providers = Provider[]
