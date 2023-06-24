<script lang="ts">
  import {listUsers} from 'data/api/users'
  import type {User} from 'data/types/user'
  import {onMount} from 'svelte'

  let users = [] as User[]

  onMount(async () => {
    users = await listUsers()
  })
</script>

<h1>Users</h1>

<ul>
{#each users as user (user.ID)}
  <li>
    {user.FirstName}
  </li>
{/each}
</ul>
