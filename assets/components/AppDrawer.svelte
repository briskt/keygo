<script lang="ts">
import AppFooter from './AppFooter.svelte'
import AppHeader from './AppHeader.svelte'
import { beforeUrlChange, url } from '@roxi/routify'
import { Drawer } from '@silintl/ui-components'
import { onMount } from 'svelte'

export let menuItems: any[]

let drawerEl = {} as any
let drawerWidth: string
let toggle = false
let currentUrl: string

onMount(() => {
  drawerEl = document.querySelector('.mdc-drawer')
  currentUrl = $url()
})

$: drawerWidth = `${drawerEl?.offsetWidth || 0}px`

$beforeUrlChange((event: CustomEvent, route: string, { url }: { url: string }) => {
  currentUrl = url
  return true
})
</script>

<style>
.logo {
  width: 10rem;
  display: block;
  margin: 0 auto;
}
:global(.drawer .mdc-drawer__content div a.mdc-deprecated-list-item) {
  margin: 16px 8px;
}
</style>

<Drawer
  {currentUrl}
  modal
  hideForPhonesOnly
  {toggle}
  {menuItems}
  title="Covered"
  class="drawer border-white"
>
  <a class="pointer" href="/" slot="header">
    <img class="logo" src="/logo.svg" alt="Logo" />
  </a>

  <AppHeader on:toggleDrawer={() => (toggle = !toggle)} />

  <slot />
  <AppFooter rightMargin={drawerWidth} />
</Drawer>
