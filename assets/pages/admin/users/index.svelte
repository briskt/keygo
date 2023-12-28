<script lang="ts">
  import {getTenant, listTenants} from 'data/api/tenants'
  import {listUsers} from 'data/api/users'
  import type {Tenant} from 'data/types/tenant'
  import type {User} from 'data/types/user'
  import {localeTime} from 'helpers/time'
  import {onMount} from 'svelte'

  let tenants = [] as Tenant[]
  let users = [] as User[]

  onMount(async () => {
    tenants = await listTenants()
    users = await listUsers()
  })

  const getTenantNameFromID = (id: string): string => {
    const tenant = tenants.find(t => t.ID === id)
    return tenant?.Name || '(none)'
  }
</script>

<h2>Users</h2>

<table>
  <tr>
    <th>Edit</th>
    <th>Role</th>
    <th>Tenant</th>
    <th>First Name</th>
    <th>Last Name</th>
    <th>Email</th>
    <th>Avatar URL</th>
    <th>Created At</th>
    <th>Updated At</th>
    <th>Last Login</th>
  </tr>
{#each users as user (user.ID)}
  <tr>
    <td><a href="/admin/users/{user.ID}">Edit</a></td>
    <td>{user.Role}</td>
    <td>{getTenantNameFromID(user.TenantID)}</td>
    <td>{user.FirstName}</td>
    <td>{user.LastName}</td>
    <td>{user.Email}</td>
    <td>{user.AvatarURL}</td>
    <td>{localeTime(user.CreatedAt)}</td>
    <td>{localeTime(user.UpdatedAt)}</td>
    <td>{localeTime(user.LastLoginAt)}</td>
  </tr>
{/each}

</table>
