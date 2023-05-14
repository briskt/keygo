<script lang="ts">
import {refreshPageTitle} from 'data/store/page-title'
import {onMount, tick} from 'svelte'

onMount(async () => {
  await tick()
  refreshPageTitle()
})
</script>


<h1>
  Hello, world!
</h1>
