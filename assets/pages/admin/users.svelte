<script lang="ts">
  import {listUsers} from 'data/api/users'
  import type {User} from 'data/types/user'
  import {localeTime} from 'helpers/time'
  import {onMount} from 'svelte'

  let users = [] as User[]

  onMount(async () => {
    users = await listUsers()
  })
</script>

<h2>Users</h2>

<table>
  <tr>
    <th>Role</th>
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
    <td>{user.Role}</td>
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
