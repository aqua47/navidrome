import React, { useState } from 'react';
import { Button, CircularProgress, Tooltip, Snackbar, Box } from '@material-ui/core';
import Alert from '@material-ui/lab/Alert'; 
import UploadIcon from '@material-ui/icons/Publish'; 

export const UploadButton = () => {
    const [uploading, setUploading] = useState(false);
    const [toast, setToast] = useState({
        open: false,
        message: '',
        severity: 'success',
    });

    const handleFileChange = async (event) => {
        const files = event.target.files;
        if (!files || files.length === 0) return;

        const file = files[0];
        const formData = new FormData();
        formData.append('file', file);

        setUploading(true);

        try {
            const response = await fetch('/api/song/upload', {
                method: 'POST',
                body: formData,
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText || 'Upload failed');
            }

            setToast({
                open: true,
                message: `"${file.name}" uploaded and scanned successfully!`,
                severity: 'success',
            });
        } catch (error) {
            setToast({
                open: true,
                message: `Error : ${error.message}`,
                severity: 'error',
            });
        } finally {
            setUploading(false);
            event.target.value = '';
        }
    };

    return (
        <Box display="inline-block">
            <input
                accept="audio/mp3,audio/*"
                style={{ display: 'none' }}
                id="navidrome-upload-input"
                type="file"
                onChange={handleFileChange}
                disabled={uploading}
            />
            <label htmlFor="navidrome-upload-input">
                <Tooltip title="Upload audio file">
                    <span>
                        <Button
                            variant="contained"
                            color="secondary"
                            component="span"
                            startIcon={uploading ? <CircularProgress size={20} color="inherit" /> : <UploadIcon />}
                            disabled={uploading}
                        >
                            {uploading ? 'Uploading...' : 'Uploader'}
                        </Button>
                    </span>
                </Tooltip>
            </label>

            <Snackbar 
                open={toast.open} 
                autoHideDuration={4000} 
                onClose={() => setToast({ ...toast, open: false })}
                anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
            >
                <Alert severity={toast.severity} onClose={() => setToast({ ...toast, open: false })} variant="filled">
                    {toast.message}
                </Alert>
            </Snackbar>
        </Box>
    );
};