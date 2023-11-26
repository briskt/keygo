<script lang="ts">
  import TenantUser from './_components/TenantUser.svelte'
  import {getTenant} from 'data/api/tenants'
  import type {Tenant} from 'data/types/tenant'
  import {localeTime} from 'helpers/time'
  import {Button, Dialog, Form, TextField} from '@silintl/ui-components'
  import {onMount} from 'svelte'

  export let id: string

  let tenant = {} as Tenant
  let showAddTenantUserModal = false

  onMount(async () => {
    tenant = await getTenant(id)
  })

  const onClickAdd = () => {
    showAddTenantUserModal = true
  }
</script>

<h1>Tenant</h1>

<Button on:click={onClickAdd}>Add User</Button>

<dl>
  <dt>Name</dt>
  <dd>{tenant.Name}</dd>
  <dt>CreatedAt</dt>
  <dd>{localeTime(tenant.CreatedAt)}</dd>
  <dt>UpdatedAt</dt>
  <dd>{localeTime(tenant.UpdatedAt)}</dd>
</dl>

<h2>Tenant Users</h2>

{#if tenant.UserIDs && tenant.UserIDs.length}
  <table>
    <tr>
      <th>Role</th>
      <th>Email</th>
      <th>First Name</th>
      <th>Last Name</th>
    </tr>
    {#each tenant.UserIDs as id (id)}
      <TenantUser {id} />
    {/each}
  </table>
{:else}
  <em>No users</em>
{/if}

<style>
  dd {
    margin: 0 0 1rem 0;
  }

  dt {
    font-weight: bold;
  }
</style>
