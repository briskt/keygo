<script lang="ts">
  import { user } from 'data/store/user'

  const { Email, FirstName, LastName, CreatedAt, LastLoginAt, Role } = $user
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

<style>
  dd {
    margin: 0 0 1rem 0;
  }

  dt {
    font-weight: bold;
  }
</style>
