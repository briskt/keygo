export type User = {
  AvatarURL: string
  CreatedAt: string //date
  Email: string
  FirstName: string
  ID: string
  LastLoginAt: string //date
  UpdatedAt: string //date
  LastName: string
  Role: string
  TenantID: string
}

export const isAdmin = (user: User) => user.Role == Admin

const Admin = 'Admin'
