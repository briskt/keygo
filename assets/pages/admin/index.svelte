<script lang="ts">
  import {listUsers} from 'data/api/users'
  import {onMount} from 'svelte'

  let userCount

  onMount(async () => {
    userCount = (await listUsers()).length
  })
</script>

<h2>Home</h2>

Total number of users: {userCount || ''}
