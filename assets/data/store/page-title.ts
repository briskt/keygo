import { writable } from 'svelte/store'

const SITE_TITLE = 'Keygo'

/**
 * The current page title, shown in the browser's tab. It is (currently) added
 * via a `<svelte:head>` section in the `App.svelte` root level component, and
 * that is (in theory) the only place this store's value needs to be read.
 */
export const pageTitle = writable('')

const getH1Text = () => {
  const h1Tag = document.querySelector('h1')
  return h1Tag?.textContent?.trim() || ''
}

const getH2Text = () => {
  const h2Tag = document.querySelector('h2')
  return h2Tag?.textContent?.trim() || ''
}

const getSelectedTabText = () => {
  const selectedTabLabel = document.querySelector('.mdc-tab--active .mdc-tab__text-label')
  return selectedTabLabel?.textContent?.trim() || ''
}

const isNotEmpty = (text: string) => !!text

/**
 * Update the page title (shown in the browser tab) based on the `<h1>` text and
 * (if applicable) the tab bar's selected tab's text, with the site title added
 * to the end. If there is no `<h1>` text found, the first `<h2>` tag's text
 * will be used instead.
 *
 * If a page does not show the correct page title, ensure the file with the
 * `<Page>` component also has an `<h1>`, and if the `<h1>` or the active tab's
 * text depends on data fetched from an API, you may need to re-call this method
 * a tick after that data has loaded. Most pages won't need to call this at all,
 * since it is automatically called after each route change (see our usage of
 * Routify's `$afterPageLoad()` helper).
 *
 * Note that for pages which inherit their `<Page>` component, it is assumed
 * that the inherited _layout file also defines an `<h1>` tag, so anything in
 * those pages with an inherited `<Page>` should start with `<h2>` level
 * headings.
 */
export const refreshPageTitle = () => {
  const headerTagText = getH1Text() || getH2Text()
  const selectedTabText = getSelectedTabText()
  setPageTitle([headerTagText, selectedTabText, SITE_TITLE])
}

const setPageTitle = (titleSegments: string[]) => {
  const combinedTitle = titleSegments.filter(isNotEmpty).join(' | ')
  pageTitle.set(combinedTitle)
}
