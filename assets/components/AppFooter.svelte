<script lang="ts">
import { PRIVACY_POLICY, TERMS_OF_USE } from 'helpers/routes'

export let rightMargin: string
</script>

<style>
footer {
  height: 150px;
  display: flex;
  justify-content: center;
  align-items: end;
  margin-right: var(--right-margin);
}
a {
  padding-bottom: 2rem;
}
</style>

<footer style:--right-margin={rightMargin}>
  <a class="pr-1" href={TERMS_OF_USE}>Terms of Service</a>
  <a href={PRIVACY_POLICY}>Privacy Policy</a>
</footer>
