# Go Watch Party

Go Watch Party is a web application built with Go that enables users to upload videos and share invites for friends to watch in a private watch party. The platform synchronizes playback, ensuring all participants experience the video in real time.

## Features

- **Video Upload**: Users can upload videos to the platform.
- **Private Watch Parties**: Generate unique invitation links for friends.
- **Synchronized Playback**: Ensures all viewers watch the video at the same time.
- **Real-Time Chat**: Allows participants to communicate while watching.
- **Secure Authentication**: Only invited users can join a watch party.

## Tech Stack

- **Backend**: Go (Golang)
- **Frontend**: HTMX / Templ
- **Database**: PostgreSQL / Redis (for caching sessions)
- **WebSockets**: For real-time synchronization
- **Storage**: Local storage or cloud (e.g., AWS S3, Google Cloud Storage)


## Usage

1. Upload a video via the dashboard.
2. Generate a private watch party link.
3. Share the link with friends.
4. Watch and chat in sync!

---

### Future Enhancements
- Live streaming support
- Mobile app integration
- More authentication options