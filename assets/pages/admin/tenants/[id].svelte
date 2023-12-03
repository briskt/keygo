<script lang="ts">
  import TenantUser from './_components/TenantUser.svelte'
  import {addTenantUser, getTenant, listTenants} from 'data/api/tenants'
  import type {Tenant} from 'data/types/tenant'
  import {localeTime} from 'helpers/time'
  import {Button, Dialog, Form, TextField} from '@silintl/ui-components'
  import {onMount} from 'svelte'

  export let id: string

  let newTenantUserEmail = ''
  let tenant = {} as Tenant
  let showAddTenantUserModal = false

  onMount(async () => {
    tenant = await getTenant(id)
  })

  const onClickAdd = () => {
    showAddTenantUserModal = true
  }

  const onAddTenantUserModalClosed = () => {
    showAddTenantUserModal = false
  }

  const onSubmitAddTenantUser = async () => {
    showAddTenantUserModal = false
    addTenantUser(id, newTenantUserEmail)
    newTenantUserEmail = ''
  }

  const onCancel = () => {
    showAddTenantUserModal = false
  }
</script>

<h2>Tenant</h2>

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

<Dialog.Alert
        open={showAddTenantUserModal}
        buttons={[]}
        defaultAction='cancel'
        title='Add Tenant'
        titleIcon='assignment_ind'
        on:closed={onAddTenantUserModalClosed}
>
  <Form on:submit={onSubmitAddTenantUser}>
    <p>
      <TextField maxlength="40" label="Email Address" bind:value={newTenantUserEmail} class="w-100" autofocus />
    </p>
    <div class="float-right form-button">
      <Button raised>Save</Button>
    </div>
    <div class="float-right form-button">
      <Button on:click={onCancel}>Cancel</Button>
    </div>
  </Form>
</Dialog.Alert>

<style>
  dd {
    margin: 0 0 1rem 0;
  }

  dt {
    font-weight: bold;
  }
</style>
