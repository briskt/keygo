<script lang="ts">
  import {listUserTokens} from 'data/api/users'
  import {user} from 'data/store/user'
  import {onMount} from 'svelte'

  const { Email, FirstName, LastName, CreatedAt, LastLoginAt, Role } = $user

  let tokens = []

  onMount(async () => {
    tokens = await listUserTokens($user.ID)
  })
</script>

<h1>My Profile</h1>
<dl>
  <dt>Name</dt>
  <dd>{FirstName} {LastName}</dd>
  <dt>Email</dt>
  <dd>{Email}</dd>
  <dt>CreatedAt</dt>
  <dd>{new Date(CreatedAt).toLocaleString()}</dd>
  <dt>LastLoginAt</dt>
  <dd>{new Date(LastLoginAt).toLocaleString()}</dd>
  <dt>Role</dt>
  <dd>{Role}</dd>
</dl>

<h2>Tokens</h2>
<ul>
  {#each tokens as token (token.id)}
    <li>
      {JSON.stringify(token)}
    </li>
  {/each}
</ul>

<style>
  dd {
    margin: 0 0 1rem 0;
  }

  dt {
    font-weight: bold;
  }
</style>
