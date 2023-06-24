<script lang="ts">
  import UserTabs from './_components/UserTabs.svelte'
  import { user } from 'data/store/user'
  import { Page } from '@silintl/ui-components'

  $: userHasLoaded = !!$user.ID
</script>

<Page>
  {#if userHasLoaded}
    <h1>My Account</h1>
    <UserTabs />
    <slot />
  {:else}
    Loading...
  {/if}
</Page>
