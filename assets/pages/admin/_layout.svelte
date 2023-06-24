<script lang="ts">
  import AdminTabs from './_components/AdminTabs.svelte'
  import { Page } from '@silintl/ui-components'
</script>

<Page>
  <h1>Admin</h1>
  <AdminTabs />
  <slot />
</Page>
