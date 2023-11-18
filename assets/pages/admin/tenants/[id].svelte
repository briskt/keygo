<script lang="ts">
  import {getTenant} from 'data/api/tenants'
  import type {Tenant} from 'data/types/tenant'
  import {Button, Dialog, Form, TextField} from '@silintl/ui-components'
  import {onMount} from 'svelte'

  export let id: string

  let newTenantName = ''
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
  <dd>{new Date(tenant.CreatedAt).toLocaleString()}</dd>
  <dt>UpdatedAt</dt>
  <dd>{new Date(tenant.UpdatedAt).toLocaleString()}</dd>
</dl>



<style>
  dd {
    margin: 0 0 1rem 0;
  }

  dt {
    font-weight: bold;
  }
</style>
