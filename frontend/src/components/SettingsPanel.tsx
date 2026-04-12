import { useState, useEffect } from 'react'
import './SettingsPanel.css'
import { GetSettings, SaveSettings } from '../../wailsjs/go/main/App'
import { store } from '../../wailsjs/go/models'

interface SettingsPanelProps {
  open: boolean
  onClose: () => void
}

export function SettingsPanel({ open, onClose }: SettingsPanelProps) {
  const [settings, setSettings] = useState<store.Settings | null>(null)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    if (open) GetSettings().then(setSettings)
  }, [open])

  async function handleSave() {
    if (!settings) return
    setSaving(true)
    try { await SaveSettings(settings); onClose() }
    finally { setSaving(false) }
  }

  if (!open) return null

  return (
    <div className="settings-panel">
      <div className="settings-header">
        <span className="settings-title">Settings</span>
        <button className="settings-close" onClick={onClose}>✕</button>
      </div>

      {settings ? (
        <div className="settings-body">
          <div className="settings-field">
            <label className="settings-label">Retention</label>
            <div className="settings-row">
              <select
                className="settings-select"
                value={settings.retentionMode}
                onChange={e => setSettings({ ...settings, retentionMode: e.target.value })}
              >
                <option value="unlimited">Unlimited</option>
                <option value="count">Last N items</option>
                <option value="days">Last N days</option>
              </select>
              {settings.retentionMode !== 'unlimited' && (
                <input
                  type="number"
                  className="settings-number"
                  min={1}
                  value={settings.retentionValue}
                  onChange={e => setSettings({ ...settings, retentionValue: parseInt(e.target.value) || 1 })}
                />
              )}
            </div>
          </div>

          <div className="settings-field">
            <label className="settings-label">Launch at login</label>
            <input
              type="checkbox"
              className="settings-checkbox"
              checked={settings.launchAtLogin}
              onChange={e => setSettings({ ...settings, launchAtLogin: e.target.checked })}
            />
          </div>

          <div className="settings-field">
            <label className="settings-label">Global shortcut</label>
            <kbd className="settings-kbd">{settings.globalShortcut}</kbd>
          </div>
        </div>
      ) : (
        <div className="settings-loading">Loading…</div>
      )}

      <div className="settings-footer">
        <button className="btn-cancel" onClick={onClose}>Cancel</button>
        <button className="btn-save" onClick={handleSave} disabled={saving}>
          {saving ? 'Saving…' : 'Save'}
        </button>
      </div>
    </div>
  )
}
