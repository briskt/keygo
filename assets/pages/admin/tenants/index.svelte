<script lang="ts">
  import {addTenant, listTenants} from 'data/api/tenants'
  import type {Tenant} from 'data/types/tenant'
  import {Button, Dialog, Form, TextField} from '@silintl/ui-components'
  import * as routes from 'helpers/routes'
  import {localeTime} from 'helpers/time'
  import {onMount} from 'svelte'

  let newTenantName = ''
  let tenants = [] as Tenant[]
  let showAddTenantModal = false

  onMount(async () => {
    tenants = await listTenants()
  })

  const onClickAdd = () => {
    showAddTenantModal = true
  }

  const onAddTenantModalClosed = () => {
    showAddTenantModal = false
  }

  const onSubmitAddTenant = async () => {
    showAddTenantModal = false
    addTenant(newTenantName)
    newTenantName = ''
    tenants = await listTenants()
  }

  const onCancel = () => {
    showAddTenantModal = false
  }
</script>

<h1>Tenants</h1>

<Button on:click={onClickAdd}>Add</Button>

<table>
  <tr>
    <th>Name</th>
    <th>Created At</th>
    <th>Updated At</th>
  </tr>
  {#each tenants as tenant (tenant.ID)}
    <tr>
      <td><a href="{routes.TENANTS + '/' + tenant.ID}">{tenant.Name}</a></td>
      <td>{localeTime(tenant.CreatedAt)}</td>
      <td>{localeTime(tenant.UpdatedAt)}</td>
    </tr>
  {/each}
</table>

<Dialog.Alert
        open={showAddTenantModal}
        buttons={[]}
        defaultAction='cancel'
        title='Add Tenant'
        titleIcon='assignment_ind'
        on:closed={onAddTenantModalClosed}
>
  <Form on:submit={onSubmitAddTenant}>
    <p>
      <TextField maxlength="40" label="Tenant Name" bind:value={newTenantName} class="w-100" autofocus />
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
  .float-left {
    float: left;
  }

  .float-right {
    float: right;
  }

  .form-button {
    margin: 0.5rem;
  }

  table {
    border-collapse: collapse;
  }
  table, th, td {
    border: 1px solid;
  }
  th,td {
    padding: 0.5rem;
  }
</style>
