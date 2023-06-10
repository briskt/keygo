<script lang="ts">
    import {authStatus} from 'data/api/auth'

    $: $authStatus.IsAuthenticated && location.replace('/api/auth/logout')
</script>
