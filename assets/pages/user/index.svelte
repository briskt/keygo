<script lang="ts">
  import { user } from 'data/store/user'
  import {localeTime} from 'helpers/time'

  const { Email, FirstName, LastName, CreatedAt, LastLoginAt, Role } = $user
</script>

<h2>My Profile</h2>
<dl>
  <dt>Name</dt>
  <dd>{FirstName} {LastName}</dd>
  <dt>Email</dt>
  <dd>{Email}</dd>
  <dt>CreatedAt</dt>
  <dd>{localeTime(CreatedAt)}</dd>
  <dt>LastLoginAt</dt>
  <dd>{localeTime(LastLoginAt)}</dd>
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
