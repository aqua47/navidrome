import React, { useState } from 'react'
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Grid,
} from '@material-ui/core'
import {
  useDataProvider,
  useNotify,
  useRefresh,
  useTranslate,
} from 'react-admin'

const EditTagsDialog = ({ record, open, onClose }) => {
  const dataProvider = useDataProvider()
  const notify = useNotify()
  const refresh = useRefresh()
  const translate = useTranslate()

  const [values, setValues] = useState({
    title: record.title || '',
    artist: record.artist || '',
    album: record.album || '',
    albumArtist: record.albumArtist || '',
    genre: record.genre || '',
    year: record.year ? record.year.toString() : '',
    trackNumber: record.trackNumber ? record.trackNumber.toString() : '',
    disc: record.disc ? record.disc.toString() : '',
    bpm: record.bpm ? record.bpm.toString() : '',
    comment: record.comment || '',
  })

  const handleChange = (e) => {
    const { name, value } = e.target
    setValues((prev) => ({ ...prev, [name]: value }))
  }

  const handleSubmit = (e) => {
    e.preventDefault()
    const id = record.mediaFileId || record.id
    dataProvider
      .editSongTags(id, values)
      .then(() => {
        // Explicitly fetch the updated song to ensure the cache is fresh
        // before triggering a global refresh.
        return dataProvider.getOne('song', { id: id })
      })
      .then(() => {
        notify('notification.updated', { type: 'info', smart_count: 1 })
        refresh()
        onClose()
      })
      .catch((err) => {
        notify('notification.http_error', { type: 'warning' })
        console.error(err)
      })
  }

  return (
    <Dialog
      open={open}
      onClose={onClose}
      onClick={(e) => e.stopPropagation()}
      fullWidth
      maxWidth="sm"
    >
      <DialogTitle>{translate('resources.song.actions.editTags')}</DialogTitle>
      <form onSubmit={handleSubmit}>
        <DialogContent dividers>
          <Grid container spacing={2}>
            <Grid item xs={12}>
              <TextField
                name="title"
                value={values.title}
                onChange={handleChange}
                label={translate('resources.song.fields.title')}
                fullWidth
                variant="outlined"
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                name="artist"
                value={values.artist}
                onChange={handleChange}
                label={translate('resources.song.fields.artist')}
                fullWidth
                variant="outlined"
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                name="album"
                value={values.album}
                onChange={handleChange}
                label={translate('resources.song.fields.album')}
                fullWidth
                variant="outlined"
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                name="albumArtist"
                value={values.albumArtist}
                onChange={handleChange}
                label={translate('resources.song.fields.albumArtist')}
                fullWidth
                variant="outlined"
              />
            </Grid>
            <Grid item xs={6}>
              <TextField
                name="genre"
                value={values.genre}
                onChange={handleChange}
                label={translate('resources.song.fields.genre')}
                fullWidth
                variant="outlined"
              />
            </Grid>
            <Grid item xs={6}>
              <TextField
                name="year"
                value={values.year}
                onChange={handleChange}
                label={translate('resources.song.fields.year')}
                fullWidth
                variant="outlined"
              />
            </Grid>
            <Grid item xs={4}>
              <TextField
                name="trackNumber"
                value={values.trackNumber}
                onChange={handleChange}
                label={translate('resources.song.fields.trackNumber')}
                fullWidth
                variant="outlined"
              />
            </Grid>
            <Grid item xs={4}>
              <TextField
                name="disc"
                value={values.disc}
                onChange={handleChange}
                label="Disc #"
                fullWidth
                variant="outlined"
              />
            </Grid>
            <Grid item xs={4}>
              <TextField
                name="bpm"
                value={values.bpm}
                onChange={handleChange}
                label={translate('resources.song.fields.bpm')}
                fullWidth
                variant="outlined"
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                name="comment"
                value={values.comment}
                onChange={handleChange}
                label={translate('resources.song.fields.comment')}
                fullWidth
                multiline
                rows={2}
                variant="outlined"
              />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose}>{translate('ra.action.cancel')}</Button>
          <Button type="submit" variant="contained" color="primary">
            {translate('ra.action.save')}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  )
}

export default EditTagsDialog
