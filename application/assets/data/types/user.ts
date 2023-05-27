export type User = {
  AvatarURL: string
  CreatedAt: string //date
  Email: string
  FirstName: string
  ID: string
  UpdatedAt: string //date
  LastName: string
  Role: string
}

export const isAdmin = (user: User) => user.Role == Admin

const Admin = 'Admin'
