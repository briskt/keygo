<script lang="ts">
  import {listTenants} from 'data/api/tenants'
  import {listUsers} from 'data/api/users'
  import {onMount} from 'svelte'

  let tenantCount
  let userCount

  onMount(async () => {
    tenantCount = (await listTenants()).length
    userCount = (await listUsers()).length
  })
</script>

<h2>Home</h2>

Total number of tenants: {tenantCount || ''} <br>
Total number of users: {userCount || ''} <br>
