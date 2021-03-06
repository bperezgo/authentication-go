require("dotenv").config();

const config = {
    spotifyClientId: process.env.SPOTIFY_CLIENT_ID,
    spotifyClientSecret: process.env.SPOTIFY_CLIENT_SECRET,
    spotifyRedirectUrl: process.env.SPOTIFY_REDIRECT_URL,
}

module.exports = { config }