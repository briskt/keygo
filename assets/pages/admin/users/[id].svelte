<script lang="ts">
  import {getUser, updateUser} from 'data/api/users'
  import type {User, UserUpdateInput} from 'data/types/user'
  import {goto} from '@roxi/routify'
  import {Button, Form, TextField} from '@silintl/ui-components'
  import {onMount} from 'svelte'

  export let id: string

  let user = {} as User

  const formData = {} as UserUpdateInput

  onMount(async () => {
    user = await getUser(id)
    formData.FirstName = user.FirstName
    formData.LastName = user.LastName
    formData.Email = user.Email
  })

  const onSubmit = (event: any) => {
    updateUser(id, formData)
    $goto('/admin/users')
  }

  const onCancel = (event: Event) => {
    event.preventDefault()
    $goto('/admin/users')
  }
</script>

<h2>Edit User</h2>

<Form on:submit={onSubmit}>
  <div class="my-1">
    <p>
      <TextField required label="First Name" bind:value={formData.FirstName} />
    </p>
    <p>
      <TextField required label="Last Name" bind:value={formData.LastName} />
    </p>
    <p>
      <TextField required label="Email" bind:value={formData.Email} />
    </p>
  </div>
  <div class="form-button">
    <Button raised>Save</Button>
  </div>
  <div class="form-button">
    <Button on:click={onCancel}>Cancel</Button>
  </div>
</Form>


<style>
  .form-button {
    float: right;
    margin: 0.5rem;
  }
</style>
