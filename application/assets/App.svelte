<script lang="ts">
import { getAuthStatus } from 'data/api/auth'
import {pageTitle} from 'data/store/page-title'
import { loadUser } from 'data/store/user'
import { routes } from '../.routify/routes'
import { Router } from '@roxi/routify'
import { Snackbar } from '@silintl/ui-components'
import { onMount } from 'svelte'

onMount(async () => {
    const status = await getAuthStatus()

    if (status.IsAuthenticated) {
        await loadUser(status.UserID)
    } else {
        location.replace('/api/auth/login')
    }
})
</script>

<svelte:head>
  <title>{$pageTitle}</title>
</svelte:head>

<Router {routes} />

<Snackbar />
