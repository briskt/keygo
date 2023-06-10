<script lang="ts">
    import { getAuthStatus } from 'data/api/auth'
    import { onMount } from 'svelte'

    onMount(async () => {
        const status = await getAuthStatus()

        if (status.IsAuthenticated) {
            location.replace('/api/auth/logout')
        }
    })
</script>
