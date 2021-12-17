const express = require("express");
const path = require("path");
const request = require("request");
const queryString = require("querystring");
const cors = require("cors");
const cookieParser = require("cookie-parser");

const generateRandomString = require("./utils/generateRandomString");
const encodeBasic = require("./utils/encodeBasic");
const scopesArray = require("./utils/scopesArray");

const playlistMocks = require("./utils/mocks/playlist");

const { config } = require("./config");

const app = express();

// static files
app.use("/static", express.static(path.join(__dirname, "public")));

// view engine setup
app.set("views", path.join(__dirname, "views"));
app.set("view engine", "pug");

function getUserInfo(accessToken) {
  if (!accessToken) {
    return Promise.resolve(null);
  }

  const options = {
    url: "https://api.spotify.com/v1/me",
    headers: { Authorization: `Bearer ${accessToken}` },
    json: true,
  };

  return new Promise((resolve, reject) => {
    request.get(options, function (error, response, body) {
      if (error || response.statusCode !== 200) {
        reject(error);
      }

      resolve(body);
    });
  });
}

function getUserPlaylists(accessToken, userId) {
  if (!accessToken || !userId) {
    return Promise.resolve(null);
  }

  const options = {
    url: `https://api.spotify.com/v1/users/${userId}/playlists`,
    headers: { Authorization: `Bearer ${accessToken}` },
    json: true,
  };

  return new Promise((resolve, reject) => {
    request.get(options, function (error, response, body) {
      if (error || response.statusCode !== 200) {
        reject(error);
      }

      resolve(body);
    });
  });
}

app.use(cors());
app.use(cookieParser());

// routes
app.get("/", async function (req, res, next) {
  const { access_token } = req.cookies;
  try {
    const userInfo = await getUserInfo(access_token);
    res.render("playlists", {
      playlists: {
        userInfo,
        isHome: true,
        items: playlistMocks,
      },
    });
  } catch (err) {
    next(err);
  }
});

app.get("/playlists", async (req, res, next) => {
  const { access_token } = req.cookies;
  if (!access_token) {
    return res.redirect("/");
  }

  try {
    const userInfo = await getUserInfo(access_token);
    const userPlaylists = await getUserPlaylists();
    res.render("playlists", {
      playlists: {
        userInfo,
        userPlaylists,
      },
    });
  } catch (err) {
    next(err);
  }
});

app.get("/login", async (req, res) => {
  const state = generateRandomString(16);
  const query = queryString.stringify({
    // authorization code grant
    response_type: "code",
    client_id: config.spotifyClientId,
    scope: scopesArray.join(" "),
    redirect_uri: config.spotifyRedirectUrl,
    // state to avoid XSS
    state,
  });

  res.cookie("auth_state", state, { httpOnly: true });
  res.redirect(`https://accounts.spotify.com/authorize?${query}`);
});

app.get("/callback", async (req, res, next) => {
  const { code, state } = req.query;
  const { auth_state } = req.cookies;
  if (state === null || state !== auth_state) {
    next(new Error("The state doesn't match"));
  }
  res.clearCookie("auth_state");

  const authOptions = {
    url: "https://accounts.spotify.com/api/token",
    form: {
      code,
      redirect_uri: config.spotifyRedirectUrl,
      grant_type: "authorization_code",
    },
    headers: {
      Authorization: `Basic ${encodeBasic(
        config.spotifyClientId,
        config.spotifyClientSecret
      )}`,
    },
    json: true,
  };

  request.post(authOptions, function (err, res, body) {
    if (err || res.statusCode !== 200) {
      next(new Error("The token is invalid"));
    }

    res.cookie("acces_token", body.access_token, { httpOnly: true });
    res.redirect("/playlist");
  });
});

// server
const server = app.listen(3000, function () {
  console.log(`Listening http://localhost:${server.address().port}`);
});
