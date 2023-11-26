<script lang="ts">
  import {getUser} from 'data/api/users'
  import type {User} from 'data/types/user'
  import {onMount} from 'svelte'

  export let id = '' as string

  let user = {} as User

  onMount(async () => {
    user = await getUser(id)
  })

</script>

<tr>
  <td>{user.Role}</td>
  <td>{user.Email}</td>
  <td>{user.FirstName}</td>
  <td>{user.LastName}</td>
</tr>
