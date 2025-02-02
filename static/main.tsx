import { render } from 'preact'
import { Router } from '~/router'
import I18n from 'i18n'

import './global.scss'

render(
  <I18n>
    <Router />
  </I18n>
  , document.getElementById('app')!)
