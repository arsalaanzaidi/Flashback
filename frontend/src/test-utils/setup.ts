// @vitejs/plugin-react v2 injects a fast-refresh preamble check into every JSX
// module. In jsdom the preamble script never runs, so we set the flag here to
// prevent "can't detect preamble" errors while keeping fast refresh enabled in
// the dev build.
(window as any).__vite_plugin_react_preamble_installed__ = true

import '@testing-library/jest-dom'
