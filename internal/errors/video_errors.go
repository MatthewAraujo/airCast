package errors

var (
	VideoNotFound = NewError(ERR_VIDEO_NOT_FOUND, "ERR_VIDEO_NOT_FOUND", "video not found")
	VideoInfo     = NewError(ERR_VIDEO_INFO, "ERR_VIDEO_INFO", "not able to get video info")
	VideoSeek     = NewError(ERR_VIDEO_SEEK, "ERR_VIDEO_SEEK", "error getting the start of the video")
	VideoSend     = NewError(ERR_VIDEO_SEND, "ERR_VIDEO_SEND", "error sending the video")
)
