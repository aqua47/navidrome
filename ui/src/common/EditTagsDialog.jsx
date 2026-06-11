import React, { useState, useEffect } from 'react'
import PropTypes from 'prop-types'
import { useNotify, useTranslate, useDataProvider } from 'react-admin'
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
} from '@material-ui/core'

const EditTagsDialog = ({ open, record, onClose }) => {
  const translate = useTranslate()
  const notify = useNotify()
  const dataProvider = useDataProvider()
  const [tagsText, setTagsText] = useState('')
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    if (open && record) {
      const tags = record.tags || {}
      setTagsText(
        Object.entries(tags)
          .map(([key, value]) => `${key}: ${value.join(', ')}`)
          .join('\n'),
      )
    }
  }, [open, record])

  const parseTags = (text) => {
    return text
      .split('\n')
      .map((line) => line.trim())
      .filter(Boolean)
      .reduce((acc, line) => {
        const index = line.indexOf(':')
        if (index > 0) {
          const key = line.slice(0, index).trim()
          const value = line.slice(index + 1).trim()
          if (key) {
            acc[key] = value.split(',').map((item) => item.trim())
          }
        }
        return acc
      }, {})
  }

  const handleSave = async () => {
    const tags = parseTags(tagsText)
    setSaving(true)
    try {
      await dataProvider.editSongTags(record.mediaFileId || record.id, tags)
      notify('resources.song.notifications.tagsUpdated', { type: 'info' })
      onClose()
    } catch (error) {
      notify(
        translate('ra.notification.http_error') + ': ' + error.message,
        {
          type: 'warning',
          multiLine: true,
          duration: 0,
        },
      )
    } finally {
      setSaving(false)
    }
  }

  return (
    <Dialog open={open} onClose={onClose} fullWidth maxWidth="sm" aria-labelledby="edit-tags-dialog-title">
      <DialogTitle id="edit-tags-dialog-title">
        {translate('resources.song.actions.editTags')}
      </DialogTitle>
      <DialogContent>
        <TextField
          autoFocus
          margin="dense"
          label={translate('resources.song.fields.tags')}
          type="text"
          value={tagsText}
          onChange={(e) => setTagsText(e.target.value)}
          fullWidth
          multiline
          minRows={6}
          variant="outlined"
          helperText={translate('resources.song.helperTexts.editTags')}
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} disabled={saving}>
          {translate('ra.action.cancel')}
        </Button>
        <Button onClick={handleSave} color="primary" disabled={saving}>
          {translate('ra.action.save')}
        </Button>
      </DialogActions>
    </Dialog>
  )
}

EditTagsDialog.propTypes = {
  open: PropTypes.bool.isRequired,
  record: PropTypes.object,
  onClose: PropTypes.func.isRequired,
}

EditTagsDialog.defaultProps = {
  record: {},
}

export default EditTagsDialog
