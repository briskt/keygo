<script lang="ts">
import { getAuthStatus, getLoginProviders } from 'data/api/auth'
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
        const providers = await getLoginProviders()

        // TODO: ask user what provider they want to use (or just use Auth0?)
        const google = providers.find((element) => element.Key === 'google')
        window.location = google.RedirectURL

    }
})
</script>

<svelte:head>
  <title>{$pageTitle}</title>
</svelte:head>

<Router {routes} />

<Snackbar />
