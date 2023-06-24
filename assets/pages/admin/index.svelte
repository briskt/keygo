<script lang="ts">
  import {listUsers} from 'data/api/users'
  import {onMount} from 'svelte'

  let userCount

  onMount(async () => {
    userCount = (await listUsers()).length
  })
</script>

<h1>Home</h1>

Total number of users: {userCount || ''}
